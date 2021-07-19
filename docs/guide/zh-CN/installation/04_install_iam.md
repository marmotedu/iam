# 安装并配置IAM系统

### 前置条件

**确保按如下顺序执行了：**    
1. [Linux 服务器基本配置](./01_prepare.md)    
2. [Go 开发环境配置](./02_install_go_env.md) 
3. [安装并配置数据库](./03_install_storage.md)

### 安装并配置 IAM 系统
1. 准备工作

```bash
cd /tmp/iam && ./scripts/install/install.sh iam::install::prepare_iam
```

### 1. 安装并配置 iam-apiserver

```bash
cd /tmp/iam && ./scripts/install/iam-apiserver.sh iam::apiserver::install
```

### 2. 安装并配置 iam-authz-server

```bash
cd /tmp/iam && ./scripts/install/iam-authz-server.sh iam::authzserver::install
```

### 3. 安装并配置 iam-pump

```bash
cd /tmp/iam && ./scripts/install/iam-pump.sh iam::pump::install
```

### 4. 安装并配置 iamctl

```bash
cd /tmp/iam && ./scripts/install/iamctl.sh iam::iamctl::install
```

### 5. 安装 man 文件

```bash
cd /tmp/iam && ./scripts/install/man.sh iam::man::install
```

### 6. 测试

```bash
cd /tmp/iam && ./scripts/install/test.sh iam::test::test
```
