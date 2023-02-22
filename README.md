# 数据库助手/dbhelper （MYSQL）

> 一个极简的数据库辅助助手，可以帮助你创建gorm/model，迁移数据表表结构或者迁移一个数据库（MYSQL）表数据到另一个数据库中

# 安装/install

如果你的本地有`go 1.18`以上的环境。你可以使用以下命令直接安装。

```
go install github.com/purerun/dbhelper@latest
```

执行后你可以运行`dbhelper`查看是否安装成功

```
$ dbhelper 
数据库助手/database helper

Usage:
  dbhelper [command]

Available Commands:
  completion        Generate the autocompletion script for the specified shell
  help              Help about any command
  init              初始化配置文件
  make:model        从db创建gorm
  migrate:table     导表助手
  migrate:tabledata 导表助手2

Flags:
  -h, --help   help for dbhelper

Use "dbhelper [command] --help" for more information about a command.

```

# 使用

当你初次使用 `dbhelper` 的时候你可以按照以下步骤

```shell
mkdir tmp
cd tmp
dbhelper init
```

`dbhelper init` 执行完毕后会在当前目录生成一个`.env` 文件，如果存在则不会再次生成。

## 使用 dbhelper 进行数据库表迁移 命令：`dbhelper migrate:table`

`.env` 中有两个变量和表结构迁移相关，分别是 `ORIGIN_DATABASE_URL` 原数据库 `TARGET_DATABASE_URL` 目标数据库。
变量配置格式如下

```
root:password@tcp(127.0.0.1:3306)/thh_database?charset=utf8mb4&parseTime=True&loc=Local
[账号]:[密码]@tcp([ip]:[端口])/[数据库名]?charset=utf8mb4&parseTime=True&loc=Local
```

配置好后即可执行迁移指令 `dbhelper migrate:table`

## 使用 dbhelper 进行数据库表内容迁移 命令：`dbhelper migrate:tabledata`

配置与表结构迁移相同，需注意保持两个数据库的表相同，如果出现表结构不同，或 mysql 版本不同都可能出现一些不符合预期的结果。

配置好后即可执行迁移指令 `dbhelper migrate:tabledata`
