# 3. Go 开发 IDE 安装和配置

编译环境准备完之后，你还需要一个代码编辑器才能开始 Go 项目开发。为了提高开发效率，你还需要将这个编辑器配置成 Go IDE。
目前，GoLand、VSCode 这些 IDE 都很优秀，但它们都是 Windows 系统下的 IDE。在 Linux 系统下我们可以选择将 Vim 配置成 Go IDE。熟练 Vim IDE 操作之后，开发效率不输 GoLand 和 VSCode。有多种方法可以配置一个Vim IDE，这里我选择使用 vim-go 将 Vim 配置成一个 Go IDE。vim-go 是社区比较受欢迎的 Vim Go 开发插件，可以用来方便的将一个 Vim 配置成 Vim IDE。
Vim IDE 的安装和配置分为以下两步。

1) 安装 vim-go

安装命令如下：

```bash
$ rm -f $HOME/.vim; mkdir -p ~/.vim/pack/plugins/start/
$ git clone --depth=1 https://github.com/fatih/vim-go.git ~/.vim/pack/plugins/start/vim-go
```

2) Go 工具安装

vim-go 会用到一些 Go 工具，比如在函数跳转时会用到 `guru`、`godef` 工具，在格式化时会用到 `goimports`，所以你也需要安装这些工具。安装方式如下：
执行 `vi /tmp/test.go`，然后输入 `:GoInstallBinaries` 安装 vim-go 需要的工具。
