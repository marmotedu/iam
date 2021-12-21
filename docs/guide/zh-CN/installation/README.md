# IAM 部署指南

为了让您更快的部署 iam，本文档提供清晰的：

## 架构说明

安装前可以先了解下 IAM 系统的软件架构：[架构说明](installation-architecture.md)

## 安装前检查

1. 服务器能访问外网

2. 操作系统：CentOS Linux 8.x (64-bit)

> 本安装脚本基于 CentOS 8.2 安装，建议你选择 CentOS 8.x 系统。其它Linux发行版、macOS也能安装，不过需要手动安装。

## 快速部署

分为以下 **2** 步骤：

1. 申请Linux服务器，并创建`going`用户
2. 一键部署 IAM 应用

### 1. 申请Linux服务器，并创建`going`用户

1. 创建普通用户（如果已有可不用创建）

一个项目由多个开发人员协作完成，为了节省企业成本，通常不会给每个开发人员都配备一台服务器。一般情况下，所有开发人员共用一个开发机，通过普通用户登录开发机进行开发，为了模拟真实的企业开发环境，本专栏也通过一个普通用户来进行项目的开发，创建方法如下：

```bash
# useradd going # 创建going用户，通过going用户登录开发机进行开发
# passwd going # 设置密码
Changing password for user going.
New password:
Retype new password:
passwd: all authentication tokens updated successfully.
```

这里假设我们设置 `going` 的密码是：`iam59!z$`

2. 添加sudoers

`root`用户的密码一般是由系统管理员维护，并定期更改。但普通用户可能要用到root的一些权限，不可能每次都向管理员询问密码。最常用的方法是，将普通用户加入到sudoers中，这样普通用户就可以通过sudo命令来暂时获取root的权限。执行如下命令添加：

```bash
# sed -i '/^root.*ALL=(ALL).*ALL/a\going\tALL=(ALL) \tALL' /etc/sudoers
```

### 2. 一键部署 IAM 应用

用新的用户名和密码，参考iam xshell session创建一个新的xshell session，并登录Linux服务器。执行如下命令：

```bash
$ export LINUX_PASSWORD='iam59!z$' # 重要！：这里要 export going 用户的密码
$ version=latest && curl https://marmotedu-1254073058.cos.ap-beijing.myqcloud.com/iam-release/${version}/iam.tar.gz | tar -xz -C /tmp/
$ cd /tmp/iam/ && ./scripts/install/install.sh iam::install::install
```

> 你也可以安装指定的版本，只需设置`version=$targetVersion`即可，例如：`version=v1.6.2`

通过以上方式安装好系统后，以下组件的密码均默认为 `iam59!z$`：
- MariaDB
- Redis
- MongoDB

### 3. 测试

通过步骤1、2你已经成功安装了IAM应用。接下来，你可以执行以下命令来测试IAM应用是否安装成功：

```bash
$ cd /tmp/iam/ && ./scripts/install/test.sh iam::test::test
```

如果运行结果如下图，则说明安装成功：

![测试结果](../../../images/iamtest运行结果.png)

## 快速卸载

```bash
$ export LINUX_PASSWORD='iam59!z$' # 重要！：这里要 export going 用户的密码
$ version=latest && curl https://marmotedu-1254073058.cos.ap-beijing.myqcloud.com/iam-release/${version}/iam.tar.gz | tar -xz -C /tmp/
$ cd /tmp/iam/ && ./scripts/install/install.sh iam::install::uninstall
```

## 更详细的安装方法

上面提供了一个快速部署方法，我还提供了一种更详细的安装方法，请参考：[手把手教你部署IAM系统](installation-procedures.md)
