# 部署架构

## 总体架构

![](https://github.com/marmotedu/iam/blob/master/docs/images/%E9%83%A8%E7%BD%B2%E6%9E%B6%E6%9E%84v1.png?raw=true)

> `iam-pump` 和 `mongo` 2 个组件正在开发中

## 架构说明

架构为了能够尽可能的用到常用的 golang 开发知识点，采用了 RESTful 和 grpc 2 种通信协议，采用了 SQL 和 NoSQL 数据库，同时大量借鉴了 `kubernetes` 和 `tkestack/tke` 优秀的设计理念。


## 模块说明

- **iam-apiserver:** iam 核心组件，用来进行用户、密钥和授权策略管理
- **iam-authz-server:** 授权服务器，从 **iam-apiserver** 拉取密钥和授权策略，根据匹配的策略进行授权
- **iamctl:** iam 系统的客户端，类似于 `kubectl`，通过 `marmotedu-sdk-go` 访问 `iam-apiserver`
- **marmotedu-sdk-go:**  iam 系统的 golang sdk，类似于 `client-go`
- **redis:** redis 缓存， 用来存储授权审计信息，供 `iam-pump` 后期进行数据分析
- **mysql:** 持久性存储用户、密钥和授权策略
- **iam-pump:** 从 redis 里面拉取授权审计数据，分析后存入 mongo
- **mongo:** 授权审计数据，供后期运营展示和分析
