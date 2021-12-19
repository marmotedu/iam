#  IAM 架构 & 能力说明

## 总体架构

IAM 架构如下图所示：

![IAM架构](../../../images/IAM架构.png)

## 架构说明

IAM 采用前后端分离的软件架构，基于Go语言开发，包含多个组件，共同组成一个认证与授权系统。

## 模块说明

- **iam-apiserver**：核心组件，通过 RESTful API 完成用户、密钥和授权策略的增删改查。
- **iam-authz-server**：授权服务，从 iam-apiserver 拉取密钥和授权策略，并缓存在内存中，用户通过请求 iam-authz-server 提供的 /v1/authz 接口来完成资源的授权。/v1/authz 接口会查询缓存的授权策略，根据这些策略决定授权是否通过。iam-authz-server 也会将授权日志上报的 Redis 中。
- **iam-pump**：从 redis 中拉取缓存的授权日志，分析后存入 mongo 数据库中。
- **iam-watcher**：分布式作业服务，间隔一定时间查询MariaDB数据库，执行一些业务逻辑处理，例如：从policy_audit表中删除超过指定天数的授权策略、禁用超过指定天数还没有登录过的用户。
- **marmotedu-sdk-go**：IAM 的 golang sdk，参考了 kubernetes 的 client-go，封装了 iam-apiserver 和 iam-authz-server 的所有 RESTful API，方便用户调用。
- **iamctl**：IAM 的客户端，参考了 kubernetes 的客户端工具 kubectl，通过 marmotedu-sdk-go 访问 iam-apiserver 和 iam-authz-server。iamctl 封装了 iam-apiserver 的所有 RESTful API，还封装了其它功能。用户可以通过命令行的方式访问 iam-apiserver。
- **redis**：缓存数据库，用来缓存密钥和授权策略，降低访问延时。同时也会缓存授权日志，作为运营系统的数据来源。
- **mysql**：持久性存储用户、密钥和授权策略信息。
- **mongo**：存储授权日志，供后期运营系统展示和分析。

上图中，灰色部分也是 IAM 项目需要的组件：
- **app**：第三方应用，是 IAM 的使用方，通过 RESTful API 或者 marmotedu-sdk-go 调用 iam-authz-server 提供的 /v1/authz 接口完成对资源的授权。
- **iam-webconsole**：IAM 的前端，通过 RESTful API 调用 iam-apiserver 实现用户、密钥和策略的增删改查。
- **iam-operating-system**：IAM 运营系统，可以用来展示运营数据或者对 IAM 进行运营类管理，比如提供上帝视角查看所有用户的资源，调整某个用户下密钥的最大个数等。
- **Loadbalance**：负载均衡器，可以是 Nginx、Haproxy 或者 API 网关，后端挂载多个 iam-apiserver 和 iam-authz-server 实例，实现 iam-apiserver 和 iam-authz-server 组件的高可用。

## 能力说明

- **RESTful 资源管理** IAM 支持对user、secret、policy资源进行CRUD管理。
- **资源授权** 可以对资源访问进行授权。
- **授权日志处理** 支持对授权日志进行处理并展示。
- **命令行工具** 通过iamctl命令行工具，可以很方便的进行各类操作。
- **分布式作业** iam-watcher为IAM项目的分布式作业服务，可以实现异步任务，并插件化的添加新的任务类型。
