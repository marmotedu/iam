# 安装并配置数据库

### 前置条件

**确保按如下顺序执行了：**

1. [Linux 服务器基本配置](./01_prepare.md)
2. [Go 开发环境配置](./02_install_go_env.md)

### 安装并配置数据库

1. 安装并配置MariaDB

```bash
$ cd /tmp/iam && ./scripts/install/mongodb.sh iam::mongodb::install
```

2. 安装并配置Redis

```bash
$ cd /tmp/iam && ./scripts/install/redis.sh iam::redis::install
```

3. 安装并配置MongoDB

```bash
$ cd /tmp/iam && ./scripts/install/mongodb.sh iam::mongodb::install
```

