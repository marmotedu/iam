# mysqldump命令使用指南

参数：

- `--no-data`: 只导出表结构不导出数据
- `--routines`: 导出存储过程和自定义函数

## 1. 导出所有数据库

```bash
mysqldump -uroot -proot --databases iam > /tmp/iam.sql
```

## 2. 导出iam数据库的所有数据

```bash
mysqldump -uroot -proot --databases iam > /tmp/iam.sql
```

## 3. 导出初始化iam数据库的sql语句

```bash
mysqldump -hxxx.xx.xx.xxx -uiam --databases iam -p'iam59!z$' iam --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > /tmp/iam.sql
```
