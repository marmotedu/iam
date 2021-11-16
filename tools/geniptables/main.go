// Copyright 2020 Linggei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/marmotedu/component-base/pkg/util/stringutil"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var (
	config     = pflag.StringP("config", "c", "access.yaml", "Access relationship configuration")
	hostType   = pflag.StringP("type", "t", "app", "Server type, suhc as: app, db")
	all        = pflag.BoolP("all", "a", false, "Generate the full iptables script")
	logTraffic = pflag.BoolP("log", "", false, "Log the traffic that matches each rule")
	deleteRule = pflag.BoolP("delete", "", false, "Delete access for a specified host")
	sshPort    = pflag.IntP("ssh-port", "", 22, "Target ssh port")
	cidr       = pflag.StringP("cidr", "", "", "Assumed intranet security")
	output     = pflag.StringP("output", "o", "", "output file name; default srcdir/<type>_string.go")
	help       = pflag.BoolP("help", "h", false, "Print this help message")
)

var head = `#!/usr/bin/env bash

#############################
#  SETUP
#############################

# Clear all rules
iptables -F

# Don't forward traffic
iptables -P FORWARD DROP 

# Allow outgoing traffic
iptables -P OUTPUT ACCEPT

# Allow established traffic
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT 

# Allow localhost traffic
iptables -A INPUT -i lo -j ACCEPT

# Allow keepalived vrrp protocol 
iptables -A INPUT -p vrrp -j ACCEPT

#############################
#  MANAGEMENT RULES
#############################

# Allow SSH (alternate port)
iptables -A INPUT -p tcp -s {{.SSHSource}} --dport {{.Port}} -j LOG --log-level 7 --log-prefix "Accept {{.Port}} alt-ssh"
iptables -A INPUT -p tcp -s {{.SSHSource}} --dport {{.Port}} -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    

# Allow nginx server access
iptables -A INPUT -p tcp -m multiport --dport 80,443 -j ACCEPT 
    
# Allow iam services
`

var tail = `
# Allow two types of ICMP
iptables -A INPUT -p icmp --icmp-type 8/0 -j LOG --log-level 7 --log-prefix "Accept Ping"
iptables -A INPUT -p icmp --icmp-type 8/0 -j ACCEPT
iptables -A INPUT -p icmp --icmp-type 8/0 -j LOG --log-level 7 --log-prefix "Accept Time Exceeded"
iptables -A INPUT -p icmp --icmp-type 11/0 -j ACCEPT

#############################
#  DEFAULT DENY
#############################

iptables -A INPUT -j LOG --log-level 7 --log-prefix "Default Deny"
iptables -A INPUT -j DROP 
`

var refreshDeny = `iptables -D INPUT -j LOG --log-level 7 --log-prefix "Default Deny" &>/dev/null
iptables -D INPUT -j DROP &>/dev/null
iptables -A INPUT -j LOG --log-level 7 --log-prefix "Default Deny"
iptables -A INPUT -j DROP
`

// Access defines hosts and ports need to accept by iptables.
type Access struct {
	SSHSource string   `yaml:"ssh-source"`
	Hosts     []string `yaml:"hosts"`
	Ports     []string `yaml:"ports"`
	DBPorts   []string `yaml:"dbports"`
}

// Generator generate a shell script contains iptables rules.
type Generator struct {
	access Access
	hosts  []string
	ports  []string
	filter []string
	buf    bytes.Buffer
	action string
	log    bool
}

// Jump defines jump information.
type Jump struct {
	// source ips can access sshd service
	SSHSource string

	// SSH Port allowed to access
	Port int
}

func main() {
	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Println(`Usage: geniptables [OPTIONS] [HOST]`)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()

		return
	}

	data, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalf("Read file error: %s", err.Error())
	}

	var access Access
	if err := yaml.Unmarshal(data, &access); err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	g := &Generator{
		access: access,
		action: "A",
		log:    *logTraffic,
	}
	if pflag.NArg() >= 1 {
		// if use cidr, tool will not allow generate rules for a specified host
		if *cidr != "" {
			pflag.Usage()

			return
		}
		g.filter = os.Args[1:]
	}
	if *all {
		jump := Jump{
			SSHSource: access.SSHSource,
			Port:      *sshPort,
		}

		tmpl, _ := template.New("jump").Parse(head)
		var buf bytes.Buffer
		_ = tmpl.Execute(&buf, jump)
		g.Printf(buf.String())
	}
	if *cidr != "" {
		g.hosts = []string{*cidr}
	} else {
		g.hosts = g.access.Hosts
	}
	if *deleteRule {
		g.action = "D"
	}

	switch *hostType {
	case "app":
		g.ports = g.access.Ports
	case "db":
		g.ports = g.access.DBPorts
	}
	g.generate()
	if *all {
		g.Printf(tail)
	} else {
		g.Printf(refreshDeny)
	}

	if *output != "" {
		if err := ioutil.WriteFile(*output, g.buf.Bytes(), 0o600); err != nil {
			log.Fatalf("writing output: %s", err)
		}

		return
	}

	fmt.Print(g.buf.String())
}

func (g *Generator) generate() {
	for _, h := range g.hosts {
		if len(g.filter) > 0 && !stringutil.StringIn(h, g.filter) {
			continue
		}
		for _, p := range g.ports {
			if g.log {
				g.Printf(
					"iptables -%s INPUT -p tcp -s %s --dport %s -j LOG --log-level 7 --log-prefix \"Accept %s access\"\n",
					g.action,
					h,
					p,
					p,
				)
			}
			g.Printf("iptables -%s INPUT -p tcp -s %s --dport %s -j ACCEPT\n", g.action, h, p)
		}
	}
}

// Printf like fmt.Printf, but add the string to g.buf.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}
