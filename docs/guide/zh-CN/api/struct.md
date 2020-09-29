# 数据结构

IAM 系统数据结构。

## ObjectMeta

资源对象元数据，所有资源对象都具有此属性。注意：只有 `name` 是输入参数，其它全是输出参数。

| 参数名称  | 类型   | 必选 | 描述                     |
| --------- | ------ | ---- | ------------------------ |
| id        | uint64 | 否   | 资源 ID，唯一标识一个资源 |
| name      | String | 是   | 资源名称（输入参数）     |
| CreatedAt | String | 否   | 资源创建时间             |
| UpdatedAt | String     |   否   | 资源更新时间             |

## UserV2

查询用户列表接口中，返回的用户字段信息。

| 参数名称    | 类型                      | 描述               |
| ----------- | ------------------------- | ------------------ |
| metadata    | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname    | String                    | 昵称               |
| password    | String                    | 密码               |
| email       | String                    | 邮箱地址           |
| phone       | String                    | 电话号码           |
| totalPolicy | Uint64                    | 用户授权策略个数   |

## Secret

密钥信息。

| 参数名称    | 类型                                 | 描述                |
| ----------- | ------------------------------------ | ------------------- |
| metadata    | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| username    | String                               | 用户名              |
| secretID    | String                               | 密钥 ID              |
| secretKey   | String                               | 密钥 Key             |
| expires     | Int64                                | 过期时间            |
| description | String                               | 密钥描述            |

## Policy

IAM 授权策略字段信息。

| 参数名称 | 类型                                                   | 描述                |
| -------- | ------------------------------------------------------ | ------------------- |
| metadata | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| username | String                                                 | 用户名              |
| policy   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息              |

## ladon.DefaultPolicy

Ladon 授权策略定义。

| 参数名称    | 类型            | 描述           |
| ----------- | --------------- | -------------- |
| id          | String          | 授权策略唯一 ID |
| description | String          | 授权策略描述   |
| subjects    | Array of String | 主题列表       |
| effect      | String          | 效力           |
| resources   | Array of String | 资源列表       |
| actions     | Array of String | 操作列表       |
| conditions  | Object          | 生效条件       |
| meta        | String          | 元数据         |
