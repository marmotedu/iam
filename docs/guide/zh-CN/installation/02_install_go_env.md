# Go 开发环境配置

### 前置条件

**确保执行了：**[Linux服务器基本配置](./01_prepare.md)

### Go 开发环境配置

1. 下载 iam 源码包，里面包含安装脚本

```bash
$ git clone --depth=1 https://github.com/marmotedu/iam.git /tmp/iam
```

2. Go 编译环境安装和配置

```bash
$ cd /tmp/iam/ && ./scripts/install/install.sh iam::install::go
$ source $HOME/.bashrc # iam::install::go 会修改 $HOME/.bashrc文件，这里需要重新加载到当前SHELL
```

3. Go 开发 IDE 安装和配置

```bash
$ cd /tmp/iam && ./scripts/install/install.sh iam::install::vim_ide
```
