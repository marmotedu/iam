# 用户相关接口

## 1. 创建用户

### 1.1 接口描述

创建用户。

### 1.2 请求方法

POST /v1/users

### 1.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                      | 描述               |
| -------- | ---- | ------------------------- | ------------------ |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname | 是   | String                    | 昵称               |
| password | 是   | String                    | 密码               |
| email    | 是   | String                    | 邮箱地址           |
| phone    | 否   | String                    | 电话号码           |

### 1.4 输出参数

| 参数名称 | 类型                      | 描述               |
| -------- | ------------------------- | ------------------ |
| metadata | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname | String                    | 昵称               |
| password | String                    | 密码               |
| email    | String                    | 邮箱地址           |
| phone    | String                    | 电话号码           |

### 1.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "foo"
  },
  "nickname": "foo",
  "password": "Foo@2020",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}' http://marmotedu.io:8080/v1/users
```
**输出示例**

```json
 {
  "metadata": {
    "name": "foo",
    "id": 31,
    "createdAt": "2020-09-23T00:27:23.432346108+08:00",
    "updatedAt": "2020-09-23T00:27:23.432346108+08:00"
  },
  "nickname": "foo",
  "password": "$2a$10$5M4m97yo4fZAHPwcRQdr1e0NaX7qMYKRIv0xePDtI8bk0ZGLN9X/6",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}
```

## 2. 批量删除用户

### 2.1 接口描述

批量删除用户。

### 2.2 请求方法

DELETE /v1/users

### 2.3 输入参数

**Query 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（用户名） |

### 2.4 输出参数

Null

### 2.5 请求示例

**输入示例**

```bash
curl -XDELETE -H'Content-Type: application/json' -H'Authorization: Bearer $Token' http://marmotedu.io:8080/v1/users?name=foo&name=fooo
```

**输出示例**

```json
null
```

## 3. 删除用户

### 3.1 接口描述

删除用户。

### 3.2 请求方法

DELETE /v1/users/:name

### 3.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（用户名） |

### 3.4 输出参数

Null

### 3.5 请求示例

**输入示例**

```bash
curl -XDELETE -H'Content-Type: application/json' -H'Authorization: Bearer $Token' http://marmotedu.io:8080/v1/users/foo
```

**输出示例**

```json
null
```

## 4. 修改密码

### 4.1 接口描述

修改用户密码。

### 4.2 请求方法

PUT /v1/users/:name/change_password

### 4.3 输入参数

**Body 参数**

| 参数名称    | 必选 | 类型   | 描述   |
| ----------- | ---- | ------ | ------ |
| oldPassword | 是   | String | 旧密码 |
| newPassword | 是   | String | 新密码 |

### 4.4 输出参数

Null

### 4.5 请求示例

**输入示例**

```bash
curl -XPUT -H'Content-Type: application/json' -H'Authorization: Basic $(echo -n 'foo:Foo@2020'|base64)' -d'{
  "oldPassword": "Foo@2020",
  "newPassword": "Foo@2021"
}' http://marmotedu.io:8080/v1/users/foo/change-password
```

**输出示例**

```json
null
```

## 5. 修改用户属性

### 5.1 接口描述

修改用户属性。

### 5.2 请求方法

PUT /v1/users/:name

### 5.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                      | 描述               |
| -------- | ---- | ------------------------- | ------------------ |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname | 是   | String                    | 昵称               |
| password | 是   | String                    | 密码               |
| email    | 是   | String                    | 邮箱地址           |
| phone    | 否   | String                    | 电话号码           |

### 5.4 输出参数

| 参数名称 | 类型                      | 描述               |
| -------- | ------------------------- | ------------------ |
| metadata | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname | String                    | 昵称               |
| password | String                    | 密码               |
| email    | String                    | 邮箱地址           |
| phone    | String                    | 电话号码           |

### 5.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "foo"
  },
  "nickname": "foo1",
  "password": "Foo@2020",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}' http://marmotedu.io:8080/v1/users
```
**输出示例**

```json
 {
  "metadata": {
    "name": "foo",
    "id": 31,
    "createdAt": "2020-09-23T00:27:23.432346108+08:00",
    "updatedAt": "2020-09-23T00:27:23.432346108+08:00"
  },
  "nickname": "foo1",
  "password": "$2a$10$5M4m97yo4fZAHPwcRQdr1e0NaX7qMYKRIv0xePDtI8bk0ZGLN9X/6",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}
```

## 6. 查询用户信息

### 6.1 接口描述

查询用户信息。

### 6.2 请求方法

GET /v1/users/:name

### 6.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（用户名） |

### 6.4 输出参数

| 参数名称 | 类型                      | 描述               |
| -------- | ------------------------- | ------------------ |
| metadata | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| nickname | String                    | 昵称               |
| password | String                    | 密码               |
| email    | String                    | 邮箱地址           |
| phone    | String                    | 电话号码           |

### 6.5 请求示例

**输入示例**

```bash
curl -XGET -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/users/foo
```

**输出示例**

```json
{
  "metadata": {
    "id": 35,
    "name": "foo",
    "createdAt": "2020-09-23T07:33:14+08:00",
    "updatedAt": "2020-09-23T07:53:09+08:00"
  },
  "nickname": "foo1",
  "password": "$2a$10$nJ0edVsVnmpVXPSm93g9SuwQjbdzL.ZgjQO3wdaMEgJ85ilX5bSK2",
  "email": "foo@foxmail.com",
  "phone": "1812884xxxx"
}
```

## 7. 查询用户列表

### 7.1 接口描述

查询用户列表。

### 7.2 请求方法

GET /v1/users

### 7.3 输入参数

**Query 参数**

| 参数名称      | 必选 | 类型   | 描述                                                           |
| ------------- | ---- | ------ | -------------------------------------------------------------- |
| fieldSelector | 否   | String | 字段选择器，格式为 `name=foo,phone=181`,当前只支持 name 字段过滤 |

### 7.4 输出参数

| 参数名称   | 类型     | 描述               |
| ---------- | -------- | ------------------ |
| totalCount | Uint64     | 资源总个数         |
| items      | Array of [UserV2](./struct.md#UserV2) | 符合条件的用户列表 |

### 7.5 请求示例

**输入示例**

```bash
curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/users?offset=0&limit=10&fieldSelector=name=foo
```

**输出示例**

```json
{
  "totalCount": 1,
  "items": [
    {
      "metadata": {
        "id": 35,
        "name": "foo",
        "createdAt": "2020-09-23T07:33:14+08:00",
        "updatedAt": "2020-09-23T07:53:09+08:00"
      },
      "nickname": "foo1",
      "password": "",
      "email": "foo@foxmail.com",
      "phone": "1812884xxxx",
      "totalPolicy": 0
    }
  ]
}
```
