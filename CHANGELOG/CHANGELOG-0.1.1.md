
<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/marmotedu/iam/compare/v0.1.0...v0.1.1) (2020-10-06)

### Bug Fixes

* **apiserver:** set check url to 127.0.0.1 when bind-address is 0.0.0.0
* **apiserver:** fix compile error

### Code Refactoring

* **apiserver:** remove middleware and add more header to cors
* **apiserver:** change the position of fs := cmd.Flags()
* **apiserver:** change to cobra functions which Run with error

