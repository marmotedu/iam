# git-chglog 使用指南

使用git-chglog需要配置：
1. CHANGELOG模板
2. git-chglog配置

## 安装

```bash
$ go get github.com/git-chglog/git-chglog/cmd/git-chglog
```

## 使用

```bash
$ git-chglog --init
```

选项：

- What is the URL of your repository?: https://github.com/marmotedu/iam
- What is your favorite style?: github
- Choose the format of your favorite commit message: <type>(<scope>): <subject> -- feat(core): Add new feature
- What is your favorite template style?: standard
- Do you include Merge Commit in CHANGELOG?: n
- Do you include Revert Commit in CHANGELOG?: y
- In which directory do you output configuration files and templates?: .chglog

```bash
$ git-chglog -o CHANGELOG/CHANGELOG-0.1.md
```

**其它使用方法：**

```bash
$ git-chglog

  If <tag query> is not specified, it corresponds to all tags.
  This is the simplest example.

$ git-chglog 1.0.0..2.0.0

  The above is a command to generate CHANGELOG including commit of 1.0.0 to 2.0.0.

$ git-chglog 1.0.0

  The above is a command to generate CHANGELOG including commit of only 1.0.0.

$ git-chglog $(git describe --tags $(git rev-list --tags --max-count=1))

  The above is a command to generate CHANGELOG with the commit included in the latest tag.

$ git-chglog --output CHANGELOG.md

  The above is a command to output to CHANGELOG.md instead of standard output.

$ git-chglog --config custom/dir/config.yml

  The above is a command that uses a configuration file placed other than ".chglog/config.yml".
```


