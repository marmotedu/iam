# 授权策略相关接口

## 1. 创建授权策略

### 1.1 接口描述

创建授权策略。

### 1.2 请求方法

POST /v1/policies

### 1.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                                                   | 描述                |
| -------- | ---- | ------------------------------------------------------ | ------------------- |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| policy   | 是   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息   |

### 1.4 输出参数

| 参数名称 | 类型                                                   | 描述                |
| -------- | ------------------------------------------------------ | ------------------- |
| metadata | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| policy   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息   |

### 1.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "policy"
  },
  "policy": {
    "description": "One policy to rule them all.",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "conditions": {
      "remoteIPAddress": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    }
  }
}' http://marmotedu.io:8080/v1/policies
```
**输出示例**

```json
{
  "metadata": {
    "id": 41,
    "name": "policy",
    "createdAt": "2020-09-23T11:42:36.94274418+08:00",
    "updatedAt": "2020-09-23T11:42:36.94274418+08:00"
  },
  "username": "admin",
  "policy": {
    "id": "",
    "description": "One policy to rule them all.",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "conditions": {
      "remoteIPAddress": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    },
    "meta": null
  }
}
```

## 2. 批量删除授权策略

### 2.1 接口描述

批量删除授权策略。

### 2.2 请求方法

DELETE /v1/policies

### 2.3 输入参数

**Query 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（授权策略名） |

### 2.4 输出参数

Null

### 2.5 请求示例

**输入示例**

```bash
curl -XDELETE -H'Content-Type: application/json' -H'Authorization: Bearer $Token' http://marmotedu.io:8080/v1/policies?name=policy&name=sdk
```

**输出示例**

```json
null
```

## 3. 删除授权策略

### 3.1 接口描述

删除授权策略。

### 3.2 请求方法

DELETE /v1/policies/:name

### 3.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（授权策略名） |

### 3.4 输出参数

Null

### 3.5 请求示例

**输入示例**

```bash
curl -XDELETE -H'Content-Type: application/json' -H'Authorization: Bearer $Token' http://marmotedu.io:8080/v1/policies/policy
```

**输出示例**

```json
null
```

## 4. 修改授权策略属性

### 4.1 接口描述

修改授权策略属性。

### 4.2 请求方法

PUT /v1/policies/:name

### 4.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                                                   | 描述                |
| -------- | ---- | ------------------------------------------------------ | ------------------- |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| policy   | 是   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息   |

### 4.4 输出参数

| 参数名称 | 类型                                                   | 描述                |
| -------- | ------------------------------------------------------ | ------------------- |
| metadata | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| policy   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息   |

### 4.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "policy"
  },
  "policy": {
    "description": "One policy to rule them all.(modify)",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "conditions": {
      "remoteIPAddress": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    }
  }
}' http://marmotedu.io:8080/v1/policies
```
**输出示例**

```json
 {
  "metadata": {
    "id": 42,
    "name": "policy",
    "createdAt": "2020-09-23T11:45:16+08:00",
    "updatedAt": "2020-09-23T11:46:11.309424642+08:00"
  },
  "username": "admin",
  "policy": {
    "id": "",
    "description": "One policy to rule them all.(modify)",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "conditions": {
      "remoteIPAddress": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    },
    "meta": null
  }
}
```

## 5. 查询授权策略信息

### 5.1 接口描述

查询授权策略信息。

### 5.2 请求方法

GET /v1/policies/:name

### 5.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（授权策略名） |

### 5.4 输出参数

| 参数名称 | 类型                                                   | 描述                |
| -------- | ------------------------------------------------------ | ------------------- |
| metadata | [ObjectMeta](./struct.md#ObjectMeta)                   | REST 资源的功能属性 |
| policy   | [ladon.DefaultPolicy](./struct.md#ladon.DefaultPolicy) | Ladon 授权策略信息   |

### 5.5 请求示例

**输入示例**

```bash
curl -XGET -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/policies/policy
```

**输出示例**

```json
{
  "metadata": {
    "id": 42,
    "name": "policy",
    "createdAt": "2020-09-23T11:45:16+08:00",
    "updatedAt": "2020-09-23T11:46:11+08:00"
  },
  "username": "admin",
  "policy": {
    "id": "",
    "description": "One policy to rule them all.(modify)",
    "subjects": [
      "users:<peter|ken>",
      "users:maria",
      "groups:admins"
    ],
    "effect": "allow",
    "resources": [
      "resources:articles:<.*>",
      "resources:printer"
    ],
    "actions": [
      "delete",
      "<create|update>"
    ],
    "conditions": {
      "remoteIPAddress": {
        "type": "CIDRCondition",
        "options": {
          "cidr": "192.168.0.1/16"
        }
      }
    },
    "meta": null
  }
}
```

## 6. 查询授权策略列表

### 6.1 接口描述

查询授权策略列表。

### 6.2 请求方法

GET /v1/policies

### 6.3 输入参数

**Query 参数**

| 参数名称      | 必选 | 类型   | 描述                                                           |
| ------------- | ---- | ------ | -------------------------------------------------------------- |
| fieldSelector | 否   | String | 字段选择器，格式为 `name=policy,description=admin`,当前只支持 name 字段过滤 |

### 6.4 输出参数

| 参数名称   | 类型     | 描述               |
| ---------- | -------- | ------------------ |
| totalCount | Uint64     | 资源总个数         |
| items      | Array of [Policy](./struct.md#Policy) | 符合条件的授权策略列表 |

### 6.5 请求示例

**输入示例**

```bash
curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/policies?offset=0&limit=10&fieldSelector=name=policy
```

**输出示例**

```json
{
  "totalCount": 1,
  "items": [
    {
      "metadata": {
        "id": 42,
        "name": "policy",
        "createdAt": "2020-09-23T11:45:16+08:00",
        "updatedAt": "2020-09-23T11:46:11+08:00"
      },
      "username": "admin",
      "policy": {
        "id": "",
        "description": "One policy to rule them all.(modify)",
        "subjects": [
          "users:<peter|ken>",
          "users:maria",
          "groups:admins"
        ],
        "effect": "allow",
        "resources": [
          "resources:articles:<.*>",
          "resources:printer"
        ],
        "actions": [
          "delete",
          "<create|update>"
        ],
        "conditions": {
          "remoteIPAddress": {
            "type": "CIDRCondition",
            "options": {
              "cidr": "192.168.0.1/16"
            }
          }
        },
        "meta": null
      }
    }
  ]
}
```
