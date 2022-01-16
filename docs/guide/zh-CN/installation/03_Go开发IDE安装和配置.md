# 3. Go 开发 IDE 安装和配置

编译环境准备完之后，你还需要一个代码编辑器才能开始 Go 项目开发。为了提高开发效率，你还需要将这个编辑器配置成 Go IDE。

目前，GoLand、VSCode 这些 IDE 都很优秀，但它们都是 Windows 系统下的 IDE。在 Linux 系统下我们可以选择将 Vim 配置成 Go IDE。熟练 Vim IDE 操作之后，开发效率不输 GoLand 和 VSCode。有多种方法可以配置一个Vim IDE，但当前比较受欢迎的是通过 SpaceVim 将 Vim 配置成一个 Go IDE。

SpaceVim 是一个社区驱动的、模块化的 Vim IDE，它以模块的方式组织管理插件以及相关配置，为不同的语言开发量身定制了相关的开发模块。SpaceVim提供代码自动补全、语法检查、格式化、调试、REPL 等特性。开发者只需要载入相关语言的模块就能得到一个开箱即用的 Vim IDE。

Vim 可以选择 NeoVim。NeoVim 是基于 Vim 的一个 fork 分支，它主要解决了 Vim8 之前版本中的异步执行、开发模式等问题，对 Vim 的兼容性很好，同时对 Vim 的代码进行了大量的清理和重构，去掉了对老旧系统的支持，添加了新的特性。

虽然 Vim8 后来也新增了异步执行等特性，两者在使用层面差异不大，但是 NeoVim 开发更激进，新特性更多，架构也相对更合理，所以我选择了 Neo Vim，你也可以根据个人爱好来选择。Vim IDE 的安装和配置分为以下五步。

1) 安装 NeoVim。

你可以直接执行 `pip3` 和 `yum` 命令安装，安装方法如下：

centos:
```bash
$ sudo pip3 install pynvim
$ sudo yum -y install neovim
```
ubuntu:
```bash
$ sudo apt update
$ sudo apt install -y neovim
$ sudo apt install python3-pip
$ sudo pip3 install neovim jedi python-language-server --user
```



2) 配置 `$HOME/.bashrc`。

配置 `nvim` 的别名为 `vi`，这样，当我们执行 `vi` 时，Linux系统就会默认调用 `nvim`。配置 EDITOR 环境变量可以使一些工具，例如 Git 默认使用 `nvim`。配置方法如下：

```bash
$ tee -a $HOME/.bashrc <<'EOF'
# Configure for nvim
export EDITOR=nvim # 默认的编辑器（git 会用到）
alias vi="nvim"
EOF
```

3) 检查 `nvim` 是否安装成功。

你可以通过查看 NeoVim 的版本来确认是否安装成功，如果成功输出版本号，说明 NeoVim 安装成功。

```bash
$ bash
$ vi --version # 输出 NVIM v0.3.8 说明安装成功
NVIM v0.3.8
Build type: RelWithDebInfo
...
```

4) 离线安装 SpaceVim。

安装 SpaceVim 步骤稍微有点复杂，为了简化你的安装，我将安装和配置 SpaceVim 的步骤做成了一个离线安装包 marmotVim 。marmotVim 可以进行 SpaceVim 的安装、卸载、打包等操作，安装步骤如下：

```bash
$ cd /tmp
$ wget https://marmotedu-1254073058.cos.ap-beijing.myqcloud.com/tools/marmotVim.tar.gz
$ tar -xvzf marmotVim.tar.gz
$ cd marmotVim
$ ./marmotVimCtl install
```

SpaceVim 配置文件为：`$HOME/.SpaceVim.d/init.toml`和`$HOME/.SpaceVim.d/autoload/custom_init.vim`，你可以根据需要自行配置（配置文件中有配置说明）：

- **init.toml：** SpaceVim 的配置文件。
- **custom_init.vim：** 用户自定义的配置文件，兼容 vimrc。

> **提示：** 如果离线安装遇到问题，你也可以直接参考SpaceVim的官方文档来安装：https://spacevim.org/quick-start-guide/。

SpaceVim Go IDE 常用操作的按键映射如下表所示：

| 按键                 | 功能描述                                                                    |
| ------------------- | --------------------------------------------------------------------------- |
| F2                  | 显示函数、变量、结构体等列表                                                |
| F3                  | 显示当前目录下文件列表                                                      |
| gd/ctrl + ]/<Enter> | `:GoDef`，跳转到光标所在标识符的声明或者定义的位置                            |
| Ctrl+I              | `:GoDefPop`，跳转到跳转堆栈的上一个位置                                       |
| Ctrl+O              | 回上一次位置                                                                |
| Shift+K             | `:GoDoc`，在新 Vim 窗口中显示光标处 word 或者给定 word 的 Go 文档             |
| Shift+L             | `:GoIfErr`，生成 `if err != nil { return ... }`示例代码                         |
| Shift+T             | `:GoDefType`，跳转到光标所在标识符的类型定义的位置                            |
| Shift+M             | `:GoInfo`，显示光标所在的标识符的信息，比如显示函数的声明信息，变量的数据类型 |
| Crl + N             | 自动补全时下一个补全项                                                      |
| Ctrl + P            | 自动补全时上一个补全项                                                      |


5) Go 工具安装。

SpaceVim 会用到一些 Go 工具，比如在函数跳转时会用到 `guru`、`godef` 工具，在格式化时会用到 goimports，所以你也需要安装这些工具。安装方法有两种：
- Vim 底线命令安装：`vi test.go`，然后执行：`:GoInstallBinaries`安装。
- 复制工具：直接将笔者整理好的工具文件复制到`$GOPATH/bin` 目录下：

```bash
$ cd /tmp
$ wget https://marmotedu-1254073058.cos.ap-beijing.myqcloud.com/tools/gotools-for-spacevim.tgz
$ mkdir -p $GOPATH/bin
$ tar -xvzf gotools-for-spacevim.tgz -C $GOPATH/bin
```
