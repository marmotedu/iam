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
- -o, --output: iptables规则输出的脚本文件名
- -h, --help: 打印帮助信息


## 使用示例

**注意：**

iptables规则中需要配置登录当前节点的规则，有2种配置方法：

- 允许某个网段登录SSH端口（默认方式，且默认CIDR为：10.0.4.0/24）：不失安全性，又比较灵活，仅内网机器可登录。
- 允许某个IP登录SSH端口：安全性更高，但是如果更换跳板机需要刷新所有节点的iptables规则

本文采用第一个方法，且geniptables工具默认的网段和SSH端口都符合本书的配置，所以以下命令中，都没有明确指定，也即采用默认值：

- CIDR: 10.0.4.0/24
- SSH Port: 30022

如果读者需要变更，可分别指定：`--cidr`和`--ssh-port`参数。如果想只允许某个IP登录服务器节点，可指定`--jump-server`。当指定了`--jump-server`后`--cidr`参数失效。

具体iptables规则配置如下：

1. 进入到iam源码根目录

2. 配置accesss.yaml

```yaml
# IAM应用节点列表（来源IP）
hosts:
  - 10.0.4.20
  - 10.0.4.21
# 来源IP可以访问的应用端口列表（iam-apiserver, iam-authz-server, iam-pump对外暴露的的端口）
ports:
  - 38080
  - 38443
  - 39090
  - 39443
  - 37070
# 来源IP可以访问的数据库端口列表（Redis, MariaDB, MongoDB）
dbports:
  - 33306
  - 36379
  - 57017
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

#############################
#  MANAGEMENT RULES
#############################

# Allow SSH (alternate port)
iptables -A INPUT -p tcp --dport 30022 -j LOG --log-level 7 --log-prefix "Accept 30022 alt-ssh"
iptables -A INPUT -p tcp --dport 30022 -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    
    
# Allow iam services
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 38080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 38443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 39090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 39443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 37070 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 38080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 38443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 39090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 39443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 37070 -j ACCEPT

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

#############################
#  MANAGEMENT RULES
#############################

# Allow SSH (alternate port)
iptables -A INPUT -p tcp --dport 30022 -j LOG --log-level 7 --log-prefix "Accept 30022 alt-ssh"
iptables -A INPUT -p tcp --dport 30022 -j ACCEPT 

#############################    
#  ACCESS RULES    
#############################    
    
# Allow iam services
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 33306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 36379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.20 --dport 57017 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 33306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 36379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.21 --dport 57017 -j ACCEPT

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
```

5. 新增应用节点

编辑access.yaml在host列表下新增10.0.4.22节点IP。编辑后内容如下：

```yaml
# IAM应用节点列表（来源IP）
hosts:
  - 10.0.4.20
  - 10.0.4.21
  - 10.0.4.22
# 来源IP可以访问的应用端口列表（iam-apiserver, iam-authz-server, iam-pump对外暴露的的端口）
ports:
  - 38080
  - 38443
  - 39090
  - 39443
  - 37070
# 来源IP可以访问的数据库端口列表（Redis, MariaDB, MongoDB）
dbports:
  - 33306
  - 36379
  - 57017
```

新增应用节点时，已有应用节点和新增的应用节点都需要更新iptables规则:

- 新增应用节点：参考步骤 3
- 已有应用节点：参考步骤 6和步骤 7

6. 应用节点新增iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t app 10.0.4.22
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 38080 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 38443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 39090 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 39443 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 37070 -j ACCEPT
```

7. 数据库节点新增iptables规则

```bash
$ go run tools/geniptables/main.go -c access.yaml -t db 10.0.4.22
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 33306 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 36379 -j ACCEPT
iptables -A INPUT -p tcp -s 10.0.4.22 --dport 57017 -j ACCEPT
```
