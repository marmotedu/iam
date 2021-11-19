
<a name="v1.2.0"></a>
## [v1.2.0](https://github.com/marmotedu/iam/compare/v1.1.0...v1.2.0) (2021-11-19)

### Bug Fixes

* fix install script cannot clone expected version
* not add global flagset if options is nil
* fix no usage and help template set for cmd when app options is nil
* **pump:** fix iam-pump exit where get no data from redis bug
* **watcher:** add missing fields in user table

### Code Refactoring

* return `User already exist` instead of `Database error`

### Features

* add iam-watcher service to do periodic works

