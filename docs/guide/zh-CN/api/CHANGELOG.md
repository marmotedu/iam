# 更新历史

## 第1次发布

发布时间：2020-09-30 17:32:12

本次发布包含了以下内容：

**新增接口:**

| 接口名称                                                      | 接口功能     |
| ------------------------------------------------------------- | ------------ |
| [POST /v1/users](./user.md#)                                  | 创建用户     |
| [DELETE /v1/users](./user.md#批量删除用户)                    | 批量删除用户 |
| [DELETE /v1/users/:name](./user.md#删除用户)                  | 删除用户     |
| [PUT /v1/users/:name/change_password](./user.md#修改用户密码) | 修改用户密码 |
| [PUT /v1/users/:name/](./user.md#修改用户属性)                | 修改用户属性 |
| [GET /v1/users/:name](./user.md#查询用户信息)                 | 查询用户信息 |
| [GET /v1/users](./user.md#查询用户列表)                       | 查询用户列表 |
| [POST /v1/secrets](./secret.md#创建密钥)           | 创建密钥     |
| [DELETE /v1/secrets/:name](./secret.md#删除密钥)   | 删除密钥     |
| [PUT /v1/secrets/:name/](./secret.md#修改密钥属性) | 修改密钥属性 |
| [GET /v1/secrets/:name](./secret.md#查询密钥信息)  | 查询密钥信息 |
| [GET /v1/secrets](./secret.md#查询密钥列表)        | 查询密钥列表 |
| [POST /v1/policies](./policy.md#创建授权策略)           | 创建授权策略     |
| [DELETE /v1/policies](./policy.md#批量删除授权策略)     | 批量删除授权策略 |
| [DELETE /v1/policies/:name](./policy.md#删除授权策略)   | 删除授权策略     |
| [PUT /v1/policies/:name/](./policy.md#修改授权策略属性) | 修改授权策略属性 |
| [GET /v1/policies/:name](./policy.md#查询授权策略信息)  | 查询授权策略信息 |
| [GET /v1/policies](./policy.md#查询授权策略列表)        | 查询授权策略列表 |

**新增数据结构:**

- [ObjectMeta](./struct.md#ObjectMeta)
- [UserV2](./struct.md#UserV2)
- [Secret](./struct.md#Secret)
- [Policy](./struct.md#Policy)
- [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy)
