# IAM - 身份识别与访问管理系统

IAM = Identity and Access Management

IAM 是一个基于 Go 语言开发的身份识别与访问管理系统，用于对资源访问进行授权。同时也具有如下能力：

1. 配合极客时间专栏 **《[Go 语言项目开发实战](https://time.geekbang.org/column/intro/100079601)》**，讲解如何用 Go 做企业级应用的开发，是该项目的理论课程，包含了项目各个知识点和构建思路的讲解，也会包含我的一线研发经验和建议。

   目录请参考：[《Go 语言项目开发实战》课程目录](./docs/guide/zh-CN/geekbang/geekbang_course_catalog.md)


2. 作为一个开发脚手架，供开发者克隆后二次开发，快速构建自己的应用。

IAM 项目会长期维护、定期更新，**欢迎兄弟们 Star & Contributing**

## 功能特性

本项目用到了Go企业开发的大部分核心技能点，见下图：

![技术思维导图](./docs/images/技术思维导图.png)

更多请参考：[marmotedu/gocollect](https://github.com/marmotedu/gocollect)

## 软件架构

![IAM架构](./docs/images/IAM架构.png)

架构解析见：[IAM 架构 & 能力说明](docs/guide/zh-CN/installation/installation-architecture.md)

## 快速开始

### 依赖检查

**Minimum Requirements**

- Hardware
  - 2 GB of Memory
  - 50 GB of Disk Space
- 操作系统：CentOS Linux 8.2 (64-bit)
- 正常访问外网

 **需求检查 & 依赖安装** 

 请参考：[](docs/guide/zh-CN/installation/installation-requirement.md)

### 构建

1. 代码包下载

```
$ git clone https://github.com/marmotedu/iam
```

2. 编译

```bash
$ cd iam
$ make
```

### 运行

```bash
./scripts/install/install.sh iam::install::install_iam    
```

## 使用指南

[IAM Documentation](docs/guide/zh-CN)

## 如何贡献

欢迎贡献代码，贡献流程可以参考 [developer's documentation](docs/devel/zh-CN/development.md)。

## 社区

You are encouraged to communicate most things via [GitHub issues](https://github.com/marmotedu/iam/issues/new/choose) or pull requests.

## 关于作者

- Lingfei Kong <colin404@foxmail.com>

为了方便交流，我建了微信群，可以加我 **微信：nightskong**，拉你入群，方便交流。

## 谁在用

如果你有项目在使用iam系统模板，也欢迎联系作者，加入使用案例。

## 许可证

IAM is licensed under the MIT. See [LICENSE](LICENSE) for the full license text.
