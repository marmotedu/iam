#  部署环境要求

## 硬件要求

> **特别注意**：
>
> 1. 安装的时候，至少需要**一个 Installer 节点**和**一个作为 Global 集群的 master 节点**共**两个节点**。
>
>    v1.3.0 之后的版本可直接使用 All-In-One 的安装模式，此时 Installer 节点也可以作为 Global 集群的节点。但注意：此时 Installer 的节点配置要以 Global 集群的节点配置为准，否则 Installer 节点配置太低很容易安装失败。另外该功能还不是很成熟，为避免安装失败，尽量**将 Installer 节点和 Global 节点分开始用**
>
> 2. **Installer 节点**：是单独的用作安装的节点，不能作为 Global 集群的节点使用。因为在安装 Global 集群时，需要多次重启 docker，此时如果 Global 集群里面有 Installer 节点，重启 docker 会中断 Global 集群的安装。该节点需要一台**系统盘 100G** 的机器，系统盘要保证剩余 **50GB 可用的空间**。
>
>    **v1.3.0 之后 Installer 节点支持作为 Global 集群的节点使用，但注意此时 Installer 节点配置以 Global 集群的节点为准**
>
> 3. **Global 集群**：至少需要一台 **8核16G内存，100G系统盘**的机器。
>
> 4. **业务集群**：业务集群是在部署完 Global 集群之后再添加的。


* **最小化部署硬件配置：**

    <table>
        <tr>
            <td><strong>安装/业务集群</strong></td>
            <td><strong>节点/集群 </td>
            <td><strong>CPU 核数 </td>
            <td><strong>内存</td>
            <td><strong>系统盘</td>
            <td><strong>数量</td>
        </tr>
        <tr>
            <td rowspan="2">安装</td>
            <td>Installer 节点</td>
            <td>1</td>
            <td>2G</td>
            <td>100G</td>
            <td>1</td>
        </tr>
        <tr>
            <td>Global 集群</td>
            <td>8</td>
            <td>16G</td>
            <td>100G</td>
            <td>1</td>
        </tr>
        <tr>
            <td rowspan="2">业务集群</td>
            <td>Master & ETCD</td>
            <td>4</td>
            <td>8G</td>
            <td>100G</td>
            <td>1</td>
        </tr>
        <tr>
            <td>Node</td>
            <td>8</td>
            <td>16G</td>
            <td>100G</td>
            <td>3</td>
        </tr>
    </table>



* **推荐硬件配置：**

    <table>
        <tr>
            <td><strong>安装/业务集群</strong></td>
            <td><strong>节点/集群 </td>
            <td><strong>CPU 核数 </td>
            <td><strong>内存</td>
            <td><strong>系统盘</td>
            <td><strong>数量</td>
        </tr>
        <tr>
            <td rowspan="2">安装</td>
            <td>Installer 节点</td>
            <td>1</td>
            <td>2G</td>
            <td>100G</td>
            <td>1</td>
        </tr>
        <tr>
            <td>Global 集群</td>
            <td>8</td>
            <td>16G</td>
            <td>100G SSD</td>
            <td>3</td>
        </tr>
        <tr>
            <td rowspan="2">业务集群</td>
            <td>Master & ETCD</td>
            <td>16</td>
            <td>32G</td>
            <td>300G SSD</td>
            <td>3</td>
        </tr>
        <tr>
            <td>Node</td>
            <td>16</td>
            <td>32G</td>
            <td>系统盘：100G<br>数据盘：300G （/var/lib/docker） </td>
            <td>>3</td>
        </tr>
    </table>
    
    > 注意：上表中的**数据盘**（/var/lib/docker）表示的是 docker 相关信息在主机中存储的位置，即**容器数据盘**，包括 docker 的镜像、容器、日志（如果容器的日志文件所在路径没有挂载 volume，日志文件会被写入容器可写层，落盘到容器数据盘里）等文件。建议给此路径挂盘，避免与系统盘混用，避免因容器、镜像、日志等 docker 相关信息导致磁盘压力过大。

## 软件要求

> **注意，以下要求针对集群中的所有节点**

| 需求项                          | 具体要求                                                     | 命令参考<br>（以 CentOS 7.6为例）                            |
| ------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 操作系统                        | Ubuntu 16.04/18.04 LTS (64-bit) <br>CentOS Linux 7.6 (64-bit)<br>Tencent Linux 2.2 | `cat /etc/redhat-release`                                    |
| kernel 版本                     | >= Kernel 3.10.0-957.10.1.el7.x86_64                         | `uname -sr`                                                  |
| ssh<br />sudo<br />yum<br />CLI | 确保<br> Installer 节点及其容器、<br>Global 集群节点及其容器、<br>业务集群节点及其容器、<br>之间能够 ssh 互联；<br />确保每个节点都有基础工具 | `1. 确保在添加所有节点时，IP 和密码输入正确。`<br/>`2. 确保每个节点都有 sudo 或 root 权限`<br />`3. 如果是 CentOS，确保拥有 yum；其他操作系统类似，确保拥有包管理器`<br />`4. 确保拥有命令行工具` |
| Swap                            | 关闭。<br />如果不满足，系统会有一定几率出现 io 飙升，造成 docker 卡死。kubelet 会启动失败(可以设置 kubelet 启动参数 --fail-swap-on 为 false 关闭 swap 检查) | `sudo swapoff -a`<br/>`sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab`<br/>`# 注意：如果 /etc/fstab有挂载 swap，必须要注释掉，不然重新开机时又会重新挂载 swap` |
| 防火墙                          | 关闭。<br />或者至少要放通22、80、8080、443、6443、2379、2380、10250-10255、31138 端口 | `可通过以下关闭防火墙`<br />`systemctl stop firewalld && systemctl disable firewalld`<br />`或者通过以下命令放通指定端口，例如只放通80端口`<br />`firewall-cmd --zone=public --add-port=80/tcp --permanent` |
| SELinux                         | 关闭。<br />Kubernetes 官方要求，否则 kubelet 挂载目录时可能报错 `Permission denied` | `setenforce 0` <br/>`sed -i "s/SELINUX=enforcing/SELINUX=disabled/g" /etc/selinux/config` |
| 时区                            | 所有服务器时区必须统一，建议设置为 `Asia/Shanghai`           | `timedatectl set-timezone Asia/Shanghai`                     |
| 时间同步                        | ETCD 集群各机器需要时间同步，可以利用 chrony 用于系统时间同步；所有服务器要求时间必须同步，误差不得超过 2 秒 | `yum install -y chronyd` <br/>`systemctl enable chronyd && systemctl start chronyd` |
| 路由检查                        | 有些设备可能会默认配置一些路由，这些路由可能与 TKEStack 冲突，建议删除这些路由并做相关配置 | `ip link delete docker0`<br/>`ip link add name docker0 type bridge`<br/>`ip addr add dev docker0 172.17.0.1/16` |
| docker 检查                     | 有些设备可能会默认安装 docker，该 docker 版本可能与 TKEStack 不一致，建议在安装 TKEStack 之前删除docker | `yum remove docker-ce containerd docker-ce-cli -y`           |