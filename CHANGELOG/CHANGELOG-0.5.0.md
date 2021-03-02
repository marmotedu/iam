
<a name="v0.5.0"></a>
## [v0.5.0](https://github.com/marmotedu/iam/compare/v0.4.0...v0.5.0) (2021-03-02)

### Bug Fixes

* fix compile error
* fix the wrong information link in command long description
* **authzserver:** fix context bug, cancel context in Run function

### Code Refactoring

* optimize variable name Store to store
* change code architecture according to go  clean arch
* change the way to create mysql db instance
* add missing doc.go and the generate file
* add context.Context parameter to some functions
* optimize log output
* **authzserver:** optimize log output
* **makefile:** change tools install method

### Features

* support graceful shutdown
* add graceful shutdown
* **pump:** add graceful stop for pump

