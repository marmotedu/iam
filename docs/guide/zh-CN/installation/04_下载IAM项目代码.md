# 4. 下载 IAM 项目代码

因为 IAM 的安装脚本存放在 iam 代码仓库中，安装需要的二进制文件也需要通过编译iam 源码来获得，所以在安装之前，你需要先下载 iam 源码：

```bash
$ mkdir -p $WORKSPACE/golang/src/github.com/marmotedu
$ cd $WORKSPACE/golang/src/github.com/marmotedu
$ git clone -b v1.6.2 --depth=1 https://github.com/marmotedu/iam
```

其中，`marmotedu` 和 `marmotedu/iam` 目录存放了本实战项目的代码。在学习的过程中，你需要频繁的访问这 2 个目录，为了方便访问，你可以追加以下 2 个环境变量和 2 个 alias 到`$HOME/.bashrc`文件中：

```bash
$ tee -a $HOME/.bashrc << 'EOF'
# Alias for quick access
export GOWORK="$WORKSPACE/golang/src"
export IAM_ROOT="$GOWORK/github.com/marmotedu/iam"
alias mm="cd $GOWORK/github.com/marmotedu"
alias i="cd $GOWORK/github.com/marmotedu/iam"
EOF
$ bash
```

之后，你就可以先通过执行 alias 命令 `mm` 访问 `$GOWORK/github.com/marmotedu` 目录；通过执行 alias 命令 `i` 访问 `$GOWORK/github.com/marmotedu/iam` 目录。我也建议你将常用操作配置成 alias，方便以后操作。

在安装配置IAM系统之前需要你执行以下命令export `going`用户的密码，这里假设密码是 `iam59!z$`：

```bash
export LINUX_PASSWORD='iam59!z$'
```

在项目开发中，像密码、密钥Key这类敏感信息，一般不会直接硬编码在系统中，而是通过环境变量的方式来使用。现网应用的配置文件是存放在一个安全的网络环境中，并且有访问授权流程，比较安全，这种配置文件中是可以配置密码等敏感信息的。

