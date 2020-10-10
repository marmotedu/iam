
<a name="v0.1.2"></a>
## [v0.1.2](https://github.com/marmotedu/iam/compare/v0.1.1...v0.1.2) (2020-10-10)

### Bug Fixes

* **apiserver:** set check url to 127.0.0.1 when bind-address is 0.0.0.0
* **apiserver:** fix compile error

### Code Refactoring

* **apiserver:** remove middleware and add more header to cors
* **apiserver:** change the position of fs := cmd.Flags()
* **apiserver:** change to cobra functions which Run with error
* **pkg:** remove default middlewares and rewrite wrktest.sh

