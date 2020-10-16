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

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/marmotedu/component-base/pkg/util/stringutil"
)

var (
	config     = pflag.StringP("config", "c", "access.yaml", "Access relationship configuration")
	hostType   = pflag.StringP("type", "t", "app", "Server type, suhc as: app, db")
	all        = pflag.BoolP("all", "a", false, "Generate the full iptables script")
	logTraffic = pflag.BoolP("log", "", false, "Log the traffic that matches each rule")
	cidr       = pflag.StringP("cidr", "", "10.0.4.0/24", "Only allow to login other internal host from the specified CIDR")
	jumpServer = pflag.StringP("jump-server", "", "", "Jump server used to login other internal host")
	sshPort    = pflag.IntP("ssh-port", "", 30022, "Target ssh port")
	output     = pflag.StringP("output", "o", "", "output file name; default srcdir/<type>_string.go")
	help       = pflag.BoolP("help", "h", false, "Print this help message")
)

var head string = `#!/usr/bin/env bash

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

#############################
#  MANAGEMENT RULES
#############################

# Allow SSH (alternate port)
iptables -A INPUT -p tcp -s {{.Server}} --dport {{.Port}} -j LOG --log-level 7 --log-prefix "Accept {{.Port}} alt-ssh"
iptables -A INPUT -p tcp -s {{.Server}} --dport {{.Port}} -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    
    
# Allow iam services
`

var tail string = `
# Allow two types of ICMP
iptables -A INPUT -p icmp --icmp-type 8/0 -j LOG --log-level 7 --log-prefix "Accept Ping"
iptables -A INPUT -p icmp --icmp-type 8/0 -j ACCEPT
iptables -A INPUT -p icmp --icmp-type 8/0 -j LOG --log-level 7 --log-prefix "Accept Time Exceeded"
iptables -A INPUT -p icmp --icmp-type 11/0 -j ACCEPT

#############################
#  DEFAULT DENY
#############################

iptables -A INPUT -d -j LOG --log-level 7 --log-prefix "Default Deny"
iptables -A INPUT -j DROP 
`

// Access defines hosts and ports need to accept by iptables.
type Access struct {
	Hosts   []string `yaml:"hosts"`
	Ports   []string `yaml:"ports"`
	DBPorts []string `yaml:"dbports"`
}

// Generator generate a shell script contains iptables rules.
type Generator struct {
	access Access
	ports  []string
	filter []string
	buf    bytes.Buffer
	log    bool
}

// Jump defines jump information.
type Jump struct {
	// Jump Server IP address
	Server string

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
		log:    *logTraffic,
	}
	if pflag.NArg() >= 1 {
		g.filter = os.Args[1:]
	}
	if *all {
		loginServer := *cidr
		if *jumpServer != "" {
			loginServer = *jumpServer
		}

		jump := Jump{
			Server: loginServer,
			Port:   *sshPort,
		}

		tmpl, _ := template.New("jump").Parse(head)
		var buf bytes.Buffer
		_ = tmpl.Execute(&buf, jump)
		g.Printf(buf.String())
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
	}

	if *output != "" {
		if err := ioutil.WriteFile(*output, g.buf.Bytes(), 0755); err != nil {
			log.Fatalf("writing output: %s", err)
		}

		return
	}

	fmt.Print(g.buf.String())
}

func (g *Generator) generate() {
	for _, h := range g.access.Hosts {
		if len(g.filter) > 0 && !stringutil.StringIn(h, g.filter) {
			continue
		}
		for _, p := range g.ports {
			if g.log {
				g.Printf("iptables -A INPUT -p tcp -s %s --dport %s -j LOG --log-level 7 --log-prefix \"Accept %s access\"\n", h, p, p)
			}
			g.Printf("iptables -A INPUT -p tcp -s %s --dport %s -j ACCEPT\n", h, p)
		}
	}
}

// Printf like fmt.Printf, but add the string to g.buf.
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}
