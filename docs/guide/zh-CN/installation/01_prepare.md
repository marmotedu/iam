# Linux 服务器基本配置

1. 登陆Linux系统

假设我们使用 **going** 用户作为实战用户，使用Xshell / SecureCRT等工具登陆Linux系统，推荐使用Xshell。

2. 配置 `$HOME/.bashrc`

```bash
$ tee $HOME/.bashrc << 'EOF'
# .bashrc    
    
# User specific aliases and functions    
    
alias rm='rm -i'    
alias cp='cp -i'    
alias mv='mv -i'    
    
# Source global definitions    
if [ -f /etc/bashrc ]; then    
        . /etc/bashrc    
fi    
    
# User specific environment    
# Basic envs    
export LANG="en_US.UTF-8" # 设置系统语言为 en_US.UTF-8，避免终端出现中文乱码    
export PS1='[\u@dev \W]\$ ' # 默认的 PS1 设置会展示全部的路径，为了防止过长，这里只展示："用户名@dev 最后的目录名"    
export WORKSPACE="$HOME/workspace" # 设置工作目录    
export PATH=$HOME/bin:$PATH # 将 $HOME/bin 目录加入到 PATH 变量中    
EOF

$ mkdir -p $HOME/workspace
$ source $HOME/.bashrc
```

3. 安装依赖包

```bash
$ sudo yum -y install make autoconf automake cmake perl-CPAN libcurl-devel libtool gcc gcc-c++ glibc-headers zlib-    devel git-lfs telnet ctags lrzsz
```

4. 安装 Git

```bash
$ cd /tmp                         
$ wget https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.30.2.tar.gz
$ tar -xvzf git-2.30.2.tar.gz                                                
$ cd git-2.30.2/                    
$ ./configure                       
$ make                              
$ sudo make install
```
5. 配置 Git
```bash
$ cat << 'EOF' >> $HOME/.bashrc
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF     

$ git config --global user.name "Lingfei Kong"    # 用户名改成自己的
$ git config --global user.email "colin404@foxmail.com"    # 邮箱改成自己的
$ git config --global credential.helper store    # 设置 Git，保存用户名和密码
$ git config --global core.longpaths true # 解决 Git 中 'Filename too long' 的错误
$ git config --global core.quotepath off
$ git config --global url."https://github.com.cnpmjs.org/".insteadOf "https://github.com/"
$ git lfs install --skip-repo
$ source $HOME/.bashrc
```
