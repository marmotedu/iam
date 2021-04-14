# IAM 系统接口文档规范

接口文档拆分为以下几个 Markdown 文件，并存放在目录 `docs/guide/zh-CN/api` 中：
- [README.md](./README.md)：API 接口介绍文档，会分类介绍 IAM 支持的 API 接口，并会存放相关 API 接口文档的链接，方便开发者查看。
- [CHANGELOG.md](./CHANGELOG.md)：API 接口文档变更历史，方便进行历史回溯，也可以使调用者决定是否进行功能更新和版本更新。
- [generic.md](./generic.md)：通用说明，用来说明通用的请求参数、返回参数、认证方法和请求方法等。
- [struct.md](./struct.md)：数据结构，用来列出接口文档中使用的数据结构。这些数据结构可能被多个 API 接口使用，会在 user.md、secret.md、policy.md 文件中被引用。
- [user.md](./user.md)、[secret.md](./secret.md)、[policy.md](./policy.md)：API接口文档，相同 REST 资源的接口会存放在一个文件中，以 REST 资源名命名文档名。
- [error_code_generated.md](./error_code_generated.md)：错误码描述，通过程序自动生成。

`user.md` 文件记录了用户相关的接口，每个接口按顺序排列包含如下 5 部分：
- 接口描述：描述接口实现了什么功能。
- 请求方法：接口的请求方法，格式为：`HTTP方法 请求路径`，例如 `POST /v1/users`。在 **通用说明** 中的 **请求方法**部分，会说明接口的请求协议和请求地址。
- 输入参数：接口的输入字段，又分为：Header 参数、Query 参数、Body 参数、Path 参数。每个字段通过：**参数名称**、**必选**、**类型** 和 **描述** 4 个属性来描述。如果参数有限制或者默认值，可以在描述部分注明。
- 输出参数：接口的返回字段，每个字段通过 **参数名称**、**类型** 和 **描述** 3 个属性来描述。
- 请求示例：一个真实的 API 接口请求和返回示例。
