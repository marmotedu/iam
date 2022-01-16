# 6. 安装和配置 IAM 系统

要想完成 IAM 系统的安装，你还需要安装和配置 iam-apiserver、iamctl、iam-authz-server、iam-pump、iam-watcher 和man page。下面先来看安装和配置的准备工作。

## 6.1 准备工作

在开始安装之前，需要先做一些准备工作，主要有以下 4 步：

1) 初始化 MariaDB 数据库，创建 `iam` 数据库。
2) 创建需要的目录。
3) 创建 CA 根证书和密钥。
4) 配置 hosts。

1. 初始化 MariaDB 数据库

安装完 MariaDB 数据库之后，还需要在 MariaDB 数据库中创建 IAM 系统需要的数据库、表和存储过程。创建 SQL 语句保存在`configs/iam.sql`文件中。具体的创建步骤如下：

1) 登录数据库并创建 `iam` 用户：

```bash
$ cd $IAM_ROOT
$ mysql -h127.0.0.1 -P3306 -uroot -p'iam59!z$' # 连接 MariaDB，-h 指定主机，-P 指定监听端口，-u 指定登录用户，-p 指定登录密码
MariaDB [(none)]> grant all on iam.* TO iam@127.0.0.1 identified by 'iam59!z$';
Query OK, 0 rows affected (0.000 sec)
MariaDB [(none)]> flush privileges;
Query OK, 0 rows affected (0.000 sec)
```

2) 用 `iam` 用户登录 MariaDB，执行`configs/iam.sql`文件，创建 `iam` 数据库：

```bash
$ mysql -h127.0.0.1 -P3306 -uiam -p'iam59!z$'
MariaDB [(none)]> source configs/iam.sql;
MariaDB [iam]> show databases;
+--------------------+
| Database             |
+--------------------+
| iam                   |
| information_schema |
| test                  |
+--------------------+
3 rows in set (0.000 sec)
```

上面的命令会创建 `iam` 数据库，也会在`iam`数据库中创建以下资源。

- **表：** `user` 是用户表，用来存放用户信息；`secret` 是密钥表，用来存放密钥信息；`policy` 是策略表，用来存放授权策略信息；`policy_audit` 是策略历史表，被删除的策略会被转存到该表。
- **admin 用户：** 在 `user` 表中，我们需要创建一个管理员用户，用户名是 `admin`，初始密码是 `Admin@2021`。
- **存储过程：** 删除用户时会自动删除该用户所属的密钥和策略信息。

2. 创建需要的目录

在安装和运行 IAM 系统的时候，我们需要将配置、二进制文件和数据文件存放到指定的目录。所以我们需要先创建好这些目录，创建命令如下：

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh # 加载环境文件中的环境配置
$ sudo mkdir -p ${IAM_DATA_DIR}/{iam-apiserver,iam-authz-server,iam-pump,iam-watcher} # 创建 Systemd WorkingDirectory 目录
$ sudo mkdir -p ${IAM_INSTALL_DIR}/bin #创建 IAM 系统安装目录
$ sudo mkdir -p ${IAM_CONFIG_DIR}/cert # 创建 IAM 系统配置文件存放目录
$ sudo mkdir -p ${IAM_LOG_DIR} # 创建 IAM 日志文件存放目录
```
IAM系统默认情况：
| 变量               | 默认值       | 作用                          |
| ------------------ | ------------ | ----------------------------- |
| ${IAM_DATA_DIR}    | /data/iam    | Systemd WorkingDirectory 目录 |
| ${IAM_INSTALL_DIR} | /opt/iam     | IAM 系统安装目录              |
| ${IAM_CONFIG_DIR}  | /etc/iam     | IAM 系统配置文件存放目录      |
| ${IAM_LOG_DIR}     | /var/log/iam | IAM 日志文件存放目录          |

注意`source scripts/install/environment.sh`时如果遇到`bash: XXXXXX: readonly variable`这类报错，说明变量已经被加载到当前Shell中了，可以忽略这类报错。

3. 创建 CA 根证书和密钥

为了确保安全，IAM 系统各组件需要使用 x509 证书对通信进行加密和认证。所以，这里我们需要先创建 CA 证书。CA 根证书是所有组件共享的，只需要创建一个 CA 证书，后续创建的所有证书都由它签名。

我们可以使用 CloudFlare 的 PKI 工具集 cfssl 来创建所有的证书。创建步骤如下：

1) 安装 cfssl 工具集。

我们可以直接安装 cfssl 已经编译好的二进制文件，cfssl 工具集中包含很多工具，这里我们需要安装 `cfssl`、`cfssljson`、`cfssl-certinfo`，功能如下：

- **cfssl：** 证书签发工具。
- **cfssljson：** 将 cfssl 生成的证书（json 格式）变为文件承载式证书。

可以通过以下命令来安装这2个工具：

```bash
$ cd $IAM_ROOT
$ ./scripts/install/install.sh iam::install::install_cfssl
```

2) 创建配置文件。

CA 配置文件用来配置根证书的使用场景 (profile) 和具体参数 (usage，过期时间、服务端认证、客户端认证、加密等)，可以在签名其他证书时用来指定特定场景：

```bash
$ cd $IAM_ROOT
$ tee ca-config.json << EOF
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "iam": {
        "usages": [
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ],
        "expiry": "876000h"
      }
    }
  }
}
EOF
```

上面的 JSON 配置中，一些字段的解释如下。
- **signing：** 表示该证书可用于签名其他证书（生成的 `ca.pem` 证书中 `CA=TRUE`）。
- **server auth：** 表示 client 可以用该证书对 server 提供的证书进行验证。
- **client auth：** 表示 server 可以用该证书对 client 提供的证书进行验证。
- **expiry：** 876000h，证书有效期设置为 100 年。

3) 创建证书签名请求文件。

通过以下命令来创建用于生成 CA 证书签名请求（CSR）的 JSON 配置文件：

```bash
$ cd $IAM_ROOT
$ tee ca-csr.json << EOF
{
  "CN": "iam-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "marmotedu",
      "OU": "iam"
    }
  ],
  "ca": {
    "expiry": "876000h"
  }
}
EOF
```

上面的 JSON 配置中，一些字段的解释如下。
- **C：** Country，国家。
- **ST：** State，省份。
- **L：** Locality (L) or City，城市。
- **CN：** Common Name，iam-apiserver 从证书中提取该字段作为请求的用户名 (User Name)，浏览器使用该字段验证网站是否合法。
- **O：** Organization，iam-apiserver 从证书中提取该字段作为请求用户所属的组 (Group)。
- **OU：** Organization Unit (或者 Company division )，部门/单位。

除此之外，还有两点需要注意：

- 不同证书 csr 文件的 `CN`、`C`、`ST`、`L`、`O`、`OU` 组合必须不同，否则可能出现 `PEER'S CERTIFICATE HAS AN INVALID SIGNATURE` 错误；
- 后续创建证书的 csr 文件时，`CN`、`OU`都不相同（`C`、`ST`、`L`、`O`相同），以达到区分的目的。

4) 创建 CA 证书和私钥。

可以通过 `cfssl gencert` 命令来创建CA 证书和私钥：

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ cfssl gencert -initca ca-csr.json | cfssljson -bare ca
$ ls ca*
ca-config.json  ca.csr  ca-csr.json  ca-key.pem  ca.pem
$ sudo mv ca* ${IAM_CONFIG_DIR}/cert # 需要将证书文件拷贝到指定文件夹下（分发证书），方便各组件引用
```

上述命令会创建运行 CA 所必需的文件 `ca-key.pem`（私钥）和 `ca.pem`（证书），还会生成 `ca.csr`（证书签名请求），用于交叉签名或重新签名。

创建完之后，我们可以通过 `cfssl certinfo` 命名查看 cert 和 csr 信息：

```bash
$ cfssl certinfo -cert ${IAM_CONFIG_DIR}/cert/ca.pem # 查看 cert(证书信息)
$ cfssl certinfo -csr ${IAM_CONFIG_DIR}/cert/ca.csr # 查看 CSR(证书签名请求)信息
```

4. 配置 hosts

IAM 通过域名访问 API 接口，因为这些域名没有注册过，还不能在互联网上被解析到，所以需要配置 hosts，通过hosts来被解析。配置命令如下：

```bash
$ sudo tee -a /etc/hosts <<EOF
127.0.0.1 iam.api.marmotedu.com
127.0.0.1 iam.authz.marmotedu.com
EOF
```
## 6.2 安装和配置 iam-apiserver

完成了准备工作之后，就可以安装 IAM 系统的各个组件了，通过以下 3 步来安装 iam-apiserver 服务。

1. 创建 iam-apiserver 证书和私钥

其他服务为了安全都是通过 HTTPS 协议访问 iam-apiserver，所以要先创建 iam-apiserver 的证书和私钥。具体步骤如下。

1) 创建证书签名请求。

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ tee iam-apiserver-csr.json <<EOF
{
  "CN": "iam-apiserver",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "marmotedu",
      "OU": "iam-apiserver"
    }
  ],
  "hosts": [
    "127.0.0.1",
    "localhost",
    "iam.api.marmotedu.com"
  ]
}
EOF
```

上述JSON配置中的 `hosts` 字段用来指定授权使用该证书的 IP 和域名列表，上面的 hosts 列出了 iam-apiserver 服务的 IP 和域名。

2) 生成证书和私钥。

通过 `cfssl gencert`来生成证书和私钥：

```bash
$ cfssl gencert -ca=${IAM_CONFIG_DIR}/cert/ca.pem \
  -ca-key=${IAM_CONFIG_DIR}/cert/ca-key.pem \
  -config=${IAM_CONFIG_DIR}/cert/ca-config.json \
  -profile=iam iam-apiserver-csr.json | cfssljson -bare iam-apiserver
$ sudo mv iam-apiserver*pem ${IAM_CONFIG_DIR}/cert # 将生成的证书和私钥文件拷贝到配置文件目录
```

2. 安装并运行 iam-apiserver

iam-apiserver 作为 IAM 系统的核心组件，需要第一个安装。具体安装步骤如下。

1) 安装 iam-apiserver 可执行二进制文件。

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ make build BINS=iam-apiserver
$ sudo cp _output/platforms/linux/amd64/iam-apiserver ${IAM_INSTALL_DIR}/bin
```
2) 生成并安装 iam-apiserver 的配置文件（iam-apiserver.yaml）。

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-apiserver.yaml > iam-apiserver.yaml
$ sudo mv iam-apiserver.yaml ${IAM_CONFIG_DIR}
```

3) 创建并安装 iam-apiserver systemd unit 文件。

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh init/iam-apiserver.service > iam-apiserver.service
$ sudo mv iam-apiserver.service /etc/systemd/system/
```

4) 启动 iam-apiserver 服务。

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable iam-apiserver
$ sudo systemctl restart iam-apiserver
$ systemctl status iam-apiserver # 查看 iam-apiserver 运行状态，如果输出中包含 active (running)字样说明 iam-apiserver 成功启动。
```

3. 测试 iam-apiserver 是否成功安装

测试 iam-apiserver 主要是测试用户、密钥、授权策略这3类 REST 资源的增删改查作。

首先，我们需要执行以下命令来获取访问 iam-apiserver 的 token：

```bash
$ curl -s -XPOST -H'Content-Type: application/json' -d'{"username":"admin","password":"Admin@2021"}' http://127.0.0.1:8080/login | jq -r .token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXBpLm1hcm1vdGVkdS5jb20iLCJleHAiOjE2MzA5NDExODMsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE2MzA4NTQ3ODMsInN1YiI6ImFkbWluIn0.TvraWpcjk2Izp2S98u4Tghrh20Qi2AeblUGhU0BnsTY
```

下面的操作都会使用curl命令来发送HTTP请求，在请求时，需要通过`-H'Authorization: Bearer $token'` 指定认证头信息。

为了操作方便，可以将获取到的token保存在token变量中：

```bash
$ export token='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXBpLm1hcm1vdGVkdS5jb20iLCJleHAiOjE2MzA5NDExODMsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE2MzA4NTQ3ODMsInN1YiI6ImFkbWluIn0.TvraWpcjk2Izp2S98u4Tghrh20Qi2AeblUGhU0BnsTY'
```

接下来，可以通过以下3步来完成iam-apiserver的测试。

1) 用户增删改查。

创建用户、列出用户、获取用户详细信息、修改用户、删除单个用户、批量删除用户的请求方法如下：

```bash
# 创建用户
$ curl -s -XPOST -H'Content-Type: application/json' -d'{"password":"User@2021","metadata":{"name":"colin"},"nickname":"colin","email":"colin@foxmail.com","phone":"1812884xxxx"}' http://127.0.0.1:8080/v1/users

# 列出用户
$ curl -s -XGET -H"Authorization: Bearer $token" 'http://127.0.0.1:8080/v1/users?offset=0&limit=10'

# 获取 colin 用户的详细信息
$ curl -s -XGET -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/users/colin

# 修改 colin 用户
$ curl -s -XPUT -H'Content-Type: application/json' -H"Authorization: Bearer $token" -d'{"nickname":"colin","email":"colin_modified@foxmail.com","phone":"1812884xxxx"}' http://127.0.0.1:8080/v1/users/colin

# 删除 colin 用户
$ curl -s -XDELETE -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/users/colin

# 批量删除用户
$ curl -s -XDELETE -H"Authorization: Bearer $token" 'http://127.0.0.1:8080/v1/users?name=colin&name=mark&name=john'
```

2) 密钥增删改查。

创建密钥、列出密钥、获取密钥详细信息、修改密钥、删除密钥的请求方法如下：

```bash
# 创建 secret0 密钥
$ curl -s -XPOST -H'Content-Type: application/json' -H"Authorization: Bearer $token" -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}' http://127.0.0.1:8080/v1/secrets

# 列出所有密钥
$ curl -s -XGET -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/secrets

# 获取 secret0 密钥的详细信息
$ curl -s -XGET -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/secrets/secret0

# 修改 secret0 密钥
$ curl -s -XPUT -H'Content-Type: application/json' -H"Authorization: Bearer $token" -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret(modified)"}' http://127.0.0.1:8080/v1/secrets/secret0

# 删除 secret0 密钥
$ curl -s -XDELETE -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/secrets/secret0
这里需要注意，因为密钥属于重要资源，被删除会导致所有的访问请求失败。为了安全，密钥不支持批量删除。
```

3) 授权策略增删改查。

创建策略、列出策略、获取策略详细信息、修改策略、删除策略的请求方法如下：

```bash
# 创建策略
$ curl -s -XPOST -H'Content-Type: application/json' -H"Authorization: Bearer $token" -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' http://127.0.0.1:8080/v1/policies

# 列出所有策略
$ curl -s -XGET -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/policies

# 获取 policy0 策略的详细信息
$ curl -s -XGET -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/policies/policy0

# 修改 policy0 策略
$ curl -s -XPUT -H'Content-Type: application/json' -H"Authorization: Bearer $token" -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' http://127.0.0.1:8080/v1/policies/policy0

# 删除 policy0 策略
$ curl -s -XDELETE -H"Authorization: Bearer $token" http://127.0.0.1:8080/v1/policies/policy0
```

## 6.3 安装 iamctl

上一小节，我们安装了 IAM 系统的 iam-apiserver 服务。但是想要访问 iam 服务，我们还需要安装客户端工具 iamctl。具体来说，我们可以通过 3 步完成 iamctl 的安装和配置。

1. 创建 iamctl 证书和私钥

iamctl 使用 HTTPS 协议与 iam-apiserver 进行通信，iam-apiserver 对 iamctl 请求包含的证书进行认证和授权。iamctl还用于对 IAM 系统的访问和管理，所以这里创建具有最高权限的 admin 证书。

1) 创建证书签名请求。
下面创建的该证书只会被 iamctl 当做 client 证书使用，所以 `hosts` 字段为空。创建命令如下：

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ cat > admin-csr.json <<EOF
{
  "CN": "admin",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "marmotedu",
      "OU": "iamctl"
    }
  ],
  "hosts": []
}
EOF
```

2) 生成证书和私钥。

```bash
$ cfssl gencert -ca=${IAM_CONFIG_DIR}/cert/ca.pem \
  -ca-key=${IAM_CONFIG_DIR}/cert/ca-key.pem \
  -config=${IAM_CONFIG_DIR}/cert/ca-config.json \
  -profile=iam admin-csr.json | cfssljson -bare admin
$ mkdir -p $(dirname ${CONFIG_USER_CLIENT_CERTIFICATE}) $(dirname ${CONFIG_USER_CLIENT_KEY}) # 创建客户端证书存放的目录
$ mv admin.pem ${CONFIG_USER_CLIENT_CERTIFICATE} # 安装 TLS 的客户端证书
$ mv admin-key.pem ${CONFIG_USER_CLIENT_KEY} # 安装 TLS 的客户端私钥文件
```

2. 安装 iamctl

iamctl 是 IAM 系统的客户端工具，其安装位置和 iam-apiserver、iam-authz-server、iam-pump、iam-watcher位置不同，为了能够在 shell 下直接运行 iamctl 命令，我们需要将 iamctl 安装到`$HOME/bin`下，同时将 iamctl 的配置存放到默认加载的$HOME/.iam目录下。具体来说，可以通过以下2步来完成。

1) 安装 iamctl 可执行程序。

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ make build BINS=iamctl
$ cp _output/platforms/linux/amd64/iamctl $HOME/bin
```
2) 生成并安装 iamctl 的配置文件（`iamctl.yaml`）。

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iamctl.yaml > iamctl.yaml
$ mkdir -p $HOME/.iam
$ mv iamctl.yaml $HOME/.iam
```

因为 iamctl 是一个客户端工具，可能会在多台机器上运行。为了简化部署 iamctl 工具的复杂度，我们可以把 `iamctl.yaml`配置文件中跟 CA 认证相关的 CA 文件内容用 base64 加密后，放置在 `iamctl.yaml`配置文件中，这样我们就不需要再部署这些CA文件。具体的加密方法为：把 `iamctl.yaml`文件中的配置项`client-certificate`、`client-key`、`certificate-authority` 分别用如下配置项替换 `client-certificate-data`、`client-key-data`、`certificate-authority-data`。这些配置项的值可以通过对 CA 文件使用 base64 加密获得。假如，`certificate-authority` 值为`/etc/iam/cert/ca.pem`，则 `certificate-authority-data` 的值为 `cat "/etc/iam/cert/ca.pem" | base64 | tr -d '\r\n'`，其他`xxx-data`变量的值类似。

3. 测试 iamctl 是否成功安装

执行 `iamctl user list` 可以列出预创建的 `admin` 用户，如图2-3所示。 

![用户列表](../../../images/iamctl_user_list.png)

## 6.4 安装和配置 iam-authz-server

接下来，我们需要安装另外一个核心组件：iam-authz-server，可以通过以下 3 步来安装。

1.创建 iam-authz-server 证书和私钥

1) 创建证书签名请求。

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ tee iam-authz-server-csr.json <<EOF
{
  "CN": "iam-authz-server",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "marmotedu",
      "OU": "iam-authz-server"
    }
  ],
  "hosts": [
    "127.0.0.1",
    "localhost",
    "iam.authz.marmotedu.com"
  ]
}
EOF
```

iam-authz-server-csr.json配置中的 `hosts` 字段指定授权使用该证书的 IP 和域名列表，上面的hosts列出了 iam-authz-server 服务的 IP 和域名。

2) 生成证书和私钥。

```bash
$ cfssl gencert -ca=${IAM_CONFIG_DIR}/cert/ca.pem \
  -ca-key=${IAM_CONFIG_DIR}/cert/ca-key.pem \
  -config=${IAM_CONFIG_DIR}/cert/ca-config.json \
  -profile=iam iam-authz-server-csr.json | cfssljson -bare iam-authz-server
$ sudo mv iam-authz-server*pem ${IAM_CONFIG_DIR}/cert # 将生成的证书和私钥文件拷贝到配置文件目录
```

2. 安装并运行 iam-authz-server

安装并运行 iam-authz-server步骤和安装并运行 iam-apiserver 步骤基本一样，也需要 4 步。

1) 安装 iam-authz-server 可执行程序。

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ make build BINS=iam-authz-server
$ sudo cp _output/platforms/linux/amd64/iam-authz-server ${IAM_INSTALL_DIR}/bin
```

2) 生成并安装 iam-authz-server 的配置文件（iam-authz-server.yaml）。

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-authz-server.yaml > iam-authz-server.yaml
$ sudo mv iam-authz-server.yaml ${IAM_CONFIG_DIR}
```

3) 创建并安装 iam-authz-server systemd unit 文件。

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh init/iam-authz-server.service > iam-authz-server.service
$ sudo mv iam-authz-server.service /etc/systemd/system/
```

4) 启动 iam-authz-server 服务。

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable iam-authz-server
$ sudo systemctl restart iam-authz-server
$ systemctl status iam-authz-server # 查看 iam-authz-server 运行状态，如果输出中包含 active (running)字样说明 iam-authz-server 成功启动。
```

3. 测试 iam-authz-server 是否成功安装

1) 重新登录系统，并获取访问令牌。

```bash
$ token=`curl -s -XPOST -H'Content-Type: application/json' -d'{"username":"admin","password":"Admin@2021"}' http://127.0.0.1:8080/login | jq -r .token`
```

2) 创建授权策略。

```bash
$ curl -s -XPOST -H"Content-Type: application/json" -H"Authorization: Bearer $token" -d'{"metadata":{"name":"authztest"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' http://127.0.0.1:8080/v1/policies
```

3) 创建密钥，并从命令的输出中提取`secretID`和`secretKey`。

```bash
$ curl -s -XPOST -H"Content-Type: application/json" -H"Authorization: Bearer $token" -d'{"metadata":{"name":"authztest"},"expires":0,"description":"admin secret"}' http://127.0.0.1:8080/v1/secrets
{"metadata":{"id":23,"name":"authztest","createdAt":"2021-04-08T07:24:50.071671422+08:00","updatedAt":"2021-04-08T07:24:50.071671422+08:00"},"username":"admin","secretID":"ZuxvXNfG08BdEMqkTaP41L2DLArlE6Jpqoox","secretKey":"7Sfa5EfAPIwcTLGCfSvqLf0zZGCjF3l8","expires":0,"description":"admin secret"}
```

4) 生成访问 iam-authz-server 的 token。

iamctl 提供了 `jwt sigin` 命令，可以根据 `secretID` 和 `secretKey` 签发 Token，方便你使用。

```bash
$ iamctl jwt sign ZuxvXNfG08BdEMqkTaP41L2DLArlE6Jpqoox 7Sfa5EfAPIwcTLGCfSvqLf0zZGCjF3l8 # iamctl jwt sign $secretID $secretKey
eyJhbGciOiJIUzI1NiIsImtpZCI6Ilp1eHZYTmZHMDhCZEVNcWtUYVA0MUwyRExBcmxFNkpwcW9veCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXV0aHoubWFybW90ZWR1LmNvbSIsImV4cCI6MTYxNzg0NTE5NSwiaWF0IjoxNjE3ODM3OTk1LCJpc3MiOiJpYW1jdGwiLCJuYmYiOjE2MTc4Mzc5OTV9.za9yLM7lHVabPAlVQLCqXEaf8sTU6sodAsMXnmpXjMQ
```

如果你开发过程中有些重复性的操作，为了方便使用，也可以将这些操作以iamctl子命令的方式集成到iamctl命令行中。

5) 测试资源授权是否通过。

我们可以通过请求 `/v1/authz` API接口来完成资源授权：

```bash
$ curl -s -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6Ilp1eHZYTmZHMDhCZEVNcWtUYVA0MUwyRExBcmxFNkpwcW9veCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXV0aHoubWFybW90ZWR1LmNvbSIsImV4cCI6MTYxNzg0NTE5NSwiaWF0IjoxNjE3ODM3OTk1LCJpc3MiOiJpYW1jdGwiLCJuYmYiOjE2MTc4Mzc5OTV9.za9yLM7lHVabPAlVQLCqXEaf8sTU6sodAsMXnmpXjMQ' -d'{"subject":"users:maria","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"192.168.0.5"}}' http://127.0.0.1:9090/v1/authz
{"allowed":true}
```

如果授权通过会返回：`{"allowed":true}`。

## 6.5 安装和配置 iam-pump

安装 iam-pump 步骤和安装 iam-apiserver、iam-authz-server 步骤基本一样，可以通过5步来安装。

1. 安装 iam-pump 可执行程序

可以通过执行以下命令来安装iam-pump可执行程序：

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ make build BINS=iam-pump
$ sudo cp _output/platforms/linux/amd64/iam-pump ${IAM_INSTALL_DIR}/bin
```

2. 生成并安装 iam-pump 的配置文件

iam-pump的配置文件为iam-pump.yaml，生成并安装命令如下：

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-pump.yaml > iam-pump.yaml
$ sudo mv iam-pump.yaml ${IAM_CONFIG_DIR}
```

3. 创建并安装systemd unit 文件

iam-pump的systemd uint文件为iam-pump.service，生成并安装命令如下：

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh init/iam-pump.service > iam-pump.service
$ sudo mv iam-pump.service /etc/systemd/system/
```

4. 启动 iam-pump 服务

可以通过执行以下命令来启动iam-pump服务：

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable iam-pump
$ sudo systemctl restart iam-pump
$ systemctl status iam-pump # 查看 iam-pump 运行状态，如果输出中包含 active (running)字样说明 iam-pump 成功启动。
```

5. 测试 iam-pump 是否成功安装

可以通过执行以下命令来测试iam-pump服务是否安装成功：

```bash
$ curl http://127.0.0.1:7070/healthz
{"status": "ok"}
```

经过以上 5 个步骤，如果返回 `{"status": "ok"}` 就说明 iam-pump 服务健康。

## 6.6 安装和配置 iam-watcher

安装 iam-watcher步骤和安装 iam-apiserver、iam-authz-server、iam-pump步骤基本一样，可以通过5步来安装。

1. 安装 iam-watcher可执行程序

可以通过执行以下命令来安装iam-watcher可执行程序：

```bash
$ cd $IAM_ROOT
$ source scripts/install/environment.sh
$ make build BINS=iam-watcher
$ sudo cp _output/platforms/linux/amd64/iam-watcher ${IAM_INSTALL_DIR}/bin
```

2. 生成并安装 iam-watcher的配置文件

iam-watcher的配置文件为iam-watcher.yaml，生成并安装命令如下：

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh configs/iam-watcher.yaml > iam-watcher.yaml
$ sudo mv iam-watcher.yaml ${IAM_CONFIG_DIR}
```

3. 创建并安装systemd unit 文件

iam-watcher的systemd uint文件为iam-watcher.service，生成并安装命令如下：

```bash
$ ./scripts/genconfig.sh scripts/install/environment.sh init/iam-watcher.service > iam-watcher.service
$ sudo mv iam-watcher.service /etc/systemd/system/
```

4. 启动 iam-watcher服务

可以通过执行以下命令来启动iam-watcher服务：

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable iam-watcher
$ sudo systemctl restart iam-watcher
$ systemctl status iam-watcher # 查看 iam-watcher运行状态，如果输出中包含 active (running)字样说明 iam-watcher成功启动。
```

5. 测试 iam-watcher是否成功安装

可以通过执行以下命令来测试iam-watcher服务是否安装成功：

```bash
$ curl http://127.0.0.1:5050/healthz
{"status": "ok"}
```

经过以上 5 个步骤，如果返回 `{"status": "ok"}` 就说明 iam-watcher服务健康。

## 6.7 安装 man 文件

IAM 系统通过组合调用包：`github.com/cpuguy83/go-md2man/v2/md2man` 和 `github.com/spf13/cobra` 提供的函数生成了各个组件的 man1 文件，具体来说，需要 3 步实现。

1) 生成各个组件的 man1 文件。


```bash
$ cd $IAM_ROOT
$ ./scripts/update-generated-docs.sh
```

2) 安装生成的 man1 文件。

```bash
$ sudo cp docs/man/man1/* /usr/share/man/man1/
```

3) 检查是否成功安装 man1 文件。

```bash
$ man iam-apiserver
```

执行 man iam-apiserver 命令后，会弹出 man 文档界面，如下图所示。

![man page](../../../images/iam-apiserver-man1.png)
 
至此，IAM 系统所有组件都已经安装成功了，你可以通过 `iamctl version` 查看客户端和服务端的版本，命令如下：

```bash
$ iamctl version -o yaml
clientVersion:
  buildDate: "2021-10-15T07:47:54Z"
  compiler: gc
  gitCommit: 34444662f5e09449caca1f6b318c7e0896d7d4d5
  gitTreeState: dirty
  gitVersion: v1.6.2
  goVersion: go1.17.2
  platform: linux/amd64
serverVersion:
  buildDate: "2021-10-15T07:40:26Z"
  compiler: gc
  gitCommit: 34444662f5e09449caca1f6b318c7e0896d7d4d5
  gitTreeState: dirty
  gitVersion: v1.6.2
  goVersion: go1.17.2
  platform: linux/amd64
```

## 6.8 测试IAM系统是否安装成功

最后，我们可以执行以下命令来测试整个IAM系统是否被成功安装：

```bash
$ ./scripts/install/test.sh iam::test::test
```

如果安装成功，最后会输出一行：`congratulations, install iam application successfully!`。如下图所示。

![iam::test::test运行结果](../../../images/iamtest运行结果.png)
