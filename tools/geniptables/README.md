# geniptables

geniptables工具用来为来源IP和端口生成iptables规则

## geniptables工具使用方法

命令格式；

```
geniptables [OPTIONS] [HOST]
```

参数说明
- -c, --config: 需要生成iptables的来源IP和访问的目标端口配置
- -t, --type: app：生成应用节点的iptables规则，db:生成数据库节点iptables规则。
- -a, --all: 生成完整的iptables命令，一般在第一次初始化节点iptables时需要
- --log: 记录符合iptables规则的网络包访问
- --ssh-port: SSH端口号
- -o, --output: iptables规则输出的脚本文件名
- -h, --help: 打印帮助信息


## 使用示例

1. 进入到iam源码根目录

2. 配置accesss.yaml

```yaml
# 允许登录SSH节点的来源IP，可以是固定IP(例如10.0.4.2)，也可以是个网段，0.0.0.0/0代表不限制来源IP
ssh-source: 10.0.4.0/24 

# IAM应用节点列表（来源IP）
hosts:
  - 10.0.4.20
  - 10.0.4.21

# 来源IP可以访问的应用端口列表（iam-apiserver, iam-authz-server, iam-pump对外暴露的的端口）
ports:
  - 8080
  - 8443
  - 9090
  - 9443
  - 7070
  - 5050

# 来源IP可以访问的数据库端口列表（Redis, MariaDB, MongoDB）
dbports:
  - 3306
  - 6379
  - 27017
```

3. 第一次初始化应用节点的iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t app -a -o firewall.sh
```

生成的firewall.sh文件内容如下：

```bash
#!/usr/bin/env bash

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
iptables -A INPUT -p tcp -s 10.0.4.0/24 --dport 22 -j LOG --log-level 7 --log-prefix "Accept 22 alt-ssh"
iptables -A INPUT -p tcp -s 10.0.4.0/24 --dport 22 -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    

# Allow nginx server access
iptables -A INPUT -p tcp -m multiport --dport 80,443 -j ACCEPT 
    
# Allow iam services
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 8080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 8443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 9090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 9443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 7070 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 5050 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 8080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 8443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 9090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 9443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 7070 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 5050 -j ACCEPT

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
```
4. 第一次初始化数据库节点的iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t db -a -o firewall.sh
```

生成的firewall.sh文件内容如下：

```bash
#!/usr/bin/env bash

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
iptables -A INPUT -p tcp -s 10.0.4.0/24 --dport 22 -j LOG --log-level 7 --log-prefix "Accept 22 alt-ssh"
iptables -A INPUT -p tcp -s 10.0.4.0/24 --dport 22 -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    

# Allow nginx server access
iptables -A INPUT -p tcp -m multiport --dport 80,443 -j ACCEPT 
    
# Allow iam services
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 3306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 6379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 27017 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 3306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 6379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 27017 -j ACCEPT

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
```

5. 新增应用节点

编辑access.yaml在hosts列表下新增10.0.4.22节点IP。编辑后内容如下：

```yaml
# 允许登录SSH节点的来源IP，可以是固定IP(例如10.0.4.2)，也可以是个网段，0.0.0.0/0代表不限制来源IP
ssh-source: 10.0.4.0/24 

# IAM应用节点列表（来源IP）
hosts:
  - 10.0.4.20
  - 10.0.4.21
  - 10.0.4.22

# 来源IP可以访问的应用端口列表（iam-apiserver, iam-authz-server, iam-pump对外暴露的的端口）
ports:
  - 8080
  - 8443
  - 9090
  - 9443
  - 7070
  - 5050

# 来源IP可以访问的数据库端口列表（Redis, MariaDB, MongoDB）
dbports:
  - 3306
  - 6379
  - 27017
```

新增应用节点时，已有应用节点和新增的应用节点都需要更新iptables规则:

- 新增应用节点：参考步骤 3
- 已有应用节点：参考步骤 6和步骤 7

6. 应用节点新增iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t app 10.0.4.22
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 8080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 8443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 9090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 9443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 7070 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 5050 -j ACCEPT
```

7. 数据库节点新增iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t db 10.0.4.22
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 3306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 6379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 27017 -j ACCEPT
```


## 其他iptables命令

### 1. 删除规则：

```bash
$ iptables -L INPUT --line-numbers
$ iptables -D INPUT <line-number>
```
