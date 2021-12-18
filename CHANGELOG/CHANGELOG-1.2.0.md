
<a name="v1.2.0"></a>
## [v1.2.0](https://github.com/marmotedu/iam/compare/v1.1.0...v1.2.0) (2021-12-18)

### Bug Fixes

* use the same key type for context
* fix install script cannot clone expected version
* not add global flagset if options is nil
* fix no usage and help template set for cmd when app options is nil
* **pump:** fix iam-pump exit where get no data from redis bug
* **watcher:** add missing fields in user table

### Code Refactoring

* improve code reuse
* improve graceful shutdown for authzserver to reduce data loss
* code about apiserver
* only publish redis message when request is successful
* optimize the code
* change param type of `NewAutoStrategy` to AuthStrategy
* optimize func name `addNamedCmdTemplate` to `addCmdTemplate`
* optimize the code
* optimize code
* return `User already exist` instead of `Database error`
* **authzserver:** retry when list policy and secret failed

### Features

* add /etc/iam as the configuration file query path
* add iam-watcher service to do periodic works
* **pump:** add distributed lock for iam-pump

