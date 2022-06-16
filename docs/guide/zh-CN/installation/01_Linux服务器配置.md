# 1. Linux 服务器配置

你可以通过 Xshell 等工具登录 Linux 服务器，登录后需要对服务器做一些简单但必要的配置，包括创建普通用户、添加 sudoers、配置 `$HOME/.bashrc` 文件，具体步骤如下。

## 1.1 创建 `going` 普通用户

1) 用 `root` 用户登录 Linux 系统，并创建普通用户。

一般来说，一个项目会由多个开发人员协作完成，为了节省企业成本，公司不会给每个开发人员都配备一台服务器，而是让所有开发人员共用一个开发机，通过普通用户登录开发机进行开发。为了模拟真实的企业开发环境，我们也通过一个普通用户进行项目开发。创建普通用户方法如下：

```bash
# useradd going # 创建 going 用户，通过 going 用户登录开发机进行开发
# passwd going # 设置 going 用户的登录密码
Changing password for user going.
New password:
Retype new password:
passwd: all authentication tokens updated successfully.
```
不仅如此，使用普通用户登录和操作开发机也可以保证系统的安全性，这是一个比较好的习惯，所以你在日常开发中也要尽量避免使用 root 用户。

2) 添加 sudoers。

我们知道很多时候，普通用户在进行一些操作时也需要 root 权限，但 `root` 用户的密码一般是由系统管理员维护并定期更改的，每次都向管理员询问密码又很麻烦。因此，我建议你将普通用户加入到 sudoers 中，这样普通用户就可以通过 sudo 命令来暂时获取 root 权限。具体来说，你可以执行以下命令将 `going` 用户添加到 sudoers 中：

centos:

```bash
# sed -i '/^root.*ALL=(ALL).*ALL/a\going\tALL=(ALL) \tALL' /etc/sudoers
```

ubuntu:

```bash
# sed -i '/^root.*ALL=(ALL:ALL).*ALL/a\going\tALL=(ALL) \tALL' /etc/sudoers
```

3) 替换 CentOS 8.4 系统中自带的 Yum 源

由于 Red Hat 提前宣布 CentOS 8 于 2021 年 12 月 31 日停止维护，官方的 Yum 源已不可使用，所以需要切换官方的 Yum 源，这里选择阿里提供的 Yum 源。切换命令如下：

```bash
# mv /etc/yum.repos.d /etc/yum.repos.d.bak # 先备份原有的 Yum 源
# mkdir /etc/yum.repos.d 
# wget -O /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-vault-8.5.2111.repo
# yum clean all && yum makecache
```


## 1.2 `going` 用户 Shell 环境设置

1) 登录 Linux 系统

假设我们使用 **going** 用户作为实战用户，使用 Xshell/SecureCRT 等工具登录 Linux 系统，推荐使用 Xshell。


2) 配置 `$HOME/.bashrc` 文件。

我们登录新服务器后的第一步就是配置 `$HOME/.bashrc` 文件，以使 Linux 登录 Shell 更加易用，例如：配置 `LANG` 解决中文乱码；配置 `PS1` 可以使命令行提示符显示更简介。配置后的内容如下：

```bash
# .bashrc

# User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Source global definitions
if [ -f /etc/bashrc ]; then
    . /etc/bashrc
fi

if [ ! -d $HOME/workspace ]; then
    mkdir -p $HOME/workspace
fi

# User specific environment
# Basic envs
export LANG="en_US.UTF-8" # 设置系统语言为 en_US.UTF-8，避免终端出现中文乱码
export PS1='[\u@dev \W]\$ ' # 默认的 PS1 设置会展示全部的路径，为了防止过长，这里只展示："用户名@dev 最后的目录名"
export WORKSPACE="$HOME/workspace" # 设置工作目录
export PATH=$HOME/bin:$PATH # 将 $HOME/bin 目录加入到 PATH 变量中

# Default entry folder
cd $WORKSPACE # 登录系统，默认进入 workspace 目录

# User specific aliases and functions
```

有一点需要注意，在 `export PATH` 时，最好把 `$PATH` 放到最后，因为添加到 `$HOME/bin` 目录中的命令是期望被优先搜索并使用的。

`$HOME/.bashrc` 文件会自动创建工作目录 `workspace`，所有的工作都可以在这个目录下展开。这样做可以带来以下几点好处：

- 可以使我们的 `$HOME` 目录保持整洁，便于以后的文件查找和分类。
- 如果哪一天 `/` 分区空间不足，可以将整个 `workspace` 目录 `mv` 到另一个分区中，并在`/`分区中保留软连接，例如：`/home/going/workspace -> /data/workspace/`。
- 如果哪天想备份所有的工作文件，可以直接备份 `workspace`。

配置好 `$HOME/.bashrc` 文件后，我们就可以执行 `bash` 命令将配置加载到当前 Shell 中了。

至此，我们就完成了 Linux 开发机的初步配置。

## 1.3 依赖安装和配置

在 Linux 系统上安装 IAM 应用会依赖一些 RPM 包和工具，有些是直接依赖，有些是间接依赖。为了避免后续的操作出现依赖错误，例如因为包不存在而导致的编译、命令执行错误等，本节会预先安装和配置这些依赖包和工具。依赖安装和配置的具体步骤如下：

1) 安装依赖。

你可以在 CentOS 系统上通过 `yum` 命令来安装需要的依赖工具，安装命令如下：

```bash
$ sudo yum -y install make autoconf automake cmake perl-CPAN libcurl-devel libtool gcc gcc-c++ glibc-headers zlib-devel git-lfs telnet lrzsz jq expat-devel openssl-devel
```

如果系统提示 `Package xxx is already installed.`，说明`xxx`包在系统中已经被安装过，你可以忽略该类报错提示。

你可以在Ubuntu 系统上通过 `apt` 命令来安装需要的依赖工具，安装命令如下：

```bash
$ sudo apt-get update 
$ sudo apt-get install build-essential
$ sudo apt-get install dh-autoreconf libcurl4-gnutls-dev libexpat1-dev gettext libz-dev libssl-dev
$ sudo apt install libcurl4-openssl-dev
```


2) 安装 Git。

因为安装 IAM 应用、执行 `go get` 命令、安装 protobuf 工具等都会通过 Git 来下载安装包，所以我们还需要安装 Git。由于低版本的 Git 不支持`--unshallow`参数，而 `go get` 在安装 Go 包时会用到 `git fetch --unshallow` 命令，因此我们要确保安装一个高版本的 Git，具体的安装方法如下：

```bash
$ cd /tmp
$ wget --no-check-certificate https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.36.1.tar.gz
$ tar -xvzf git-2.36.1.tar.gz
$ cd git-2.36.1/
$ ./configure
$ make
$ sudo make install
$ git --version          # 输出 git 版本号，说明安装成功
git version 2.36.1
```

安装好Git后，还需要把 Git 的二进制目录添加到 `PATH` 路径中，不然 Git 可能会因为找不到一些命令导致 Git 报错。你可以执行以下命令来添加目录：

```bash
$ tee -a $HOME/.bashrc <<'EOF'
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF
```
3) 配置 Git。

你可以直接执行以下命令来配置 Git：

```bash
$ git config --global user.name "Lingfei Kong"    # 用户名改成自己的
$ git config --global user.email "colin404@foxmail.com"    # 邮箱改成自己的
$ git config --global credential.helper store    # 设置 git，保存用户名和密码
$ git config --global core.longpaths true # 解决 Git 中 'Filename too long' 的错误
```

除了按照上述步骤配置 Git 之外，这里还有以下两点需要注意。

- 在 Git 中，我们把非 ASCII 字符叫作 Unusual 字符。这类字符在 Git 输出到终端的时候默认是用 8 进制转义字符输出的（以防乱码），但现在的终端多数都支持直接显示非 ASCII 字符，所以我们可以关闭掉这个特性，具体的命令如下：

```bash
$ git config --global core.quotepath off
```

- GitHub 限制最大只能克隆 100M 的单个文件，为了能够克隆大于 100M 的文件，我们还需要安装 Git Large File Storage，安装方式如下：

```bash
$ git lfs install --skip-repo
```
