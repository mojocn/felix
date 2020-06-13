# ginbro(gin and gorm's brother) 详解

[![](/images/ginbro_coverage.jpg)](https://github.com/libragen/felix)

## 安装felix

```bash
git clone https://github.com/libragen/felix
cd felix
go mod download

go install
echo "添加 GOBIN 到 PATH环境变量"

echo "或者"

go get github.com/libragen/felix

echo "go build && ./felix -h"

```

## What is Ginbro

- Gin脚手架工具:因为工作中非常多次的使用mysql数据库 + gin + GORM 开发RESTful API程序,所以开发一个Go语言的RESTful APIs的脚手架工具
- Ginbro代码来源:Ginrbo的代码迭代自[github.com/dejavuzhou/ginbro](https://github.com/dejavuzhou/ginbro)
- SPA二进制化工具:vuejs全家桶代码二进制化成go代码,编译的时候变成二进制,运行的时候直接加载到内存中,同时和gin API在一个域名下不需要再nginx中配置rewrite或者跨域,加快API访问速度


## 功能一:Gin+GORM_SQL RESTful 脚手架工具

### 工作原理

1. 通过cobra 获取命令行参数
2. 使用sql参数连接数据库
3. 获取数据库表的名称和字段类型等数据库
4. 数据库边的表名和字段信息,转换成 [Swagger doc 规范](https://swagger.io/specification/)字段 和 GORM 模型字段
5. 使用标准库 [`text/template`](https://golang.google.cn/pkg/text/template/) 生成swagger.yaml, GORM 模型文件, GIN handler 文件 ...
6. 使用 `go fmt ./...` 格式化代码
7. 使用标准库`archive/zip`打包`*.go config.toml ...`代码,提供zip文件下载(命令行模式没有)

### 支持数据库大多数SQL数据库
- mysql
- SQLite
- postgreSQL
- mssql(TODO:: sqlserver)

### ginbro 生成app代码包含功能简介

- 每一张数据库表生成一个RESTful规范的资源(`GET-pagination/POST/GET-one/PATCH/DELETE`)
- 支持API-json数据分页-和总数分页缓存,减少全表扫描
- 支持golang-内存单机缓存
- 支持`gin autotls`
- 前端代码和API公用一个服务,减少跨域OPTIONS的请求时间和配置时间,同时完美支持前后端分离
- 开箱支持jwt-token认证和Bearer Token 路由中间件
- 开箱即用的logrus数据库
- 开箱即用的viper配置文件
- 开箱即用的swagger API 文档
- 开箱即用的定时任务系统

### 项目演示地址

#### [felix sshw 网页UI演示地址 用户名和密码都是admin](http://felix.mojotv.cn/#/)
#### [生成swagger API交互文档地址 http://ginbro.mojotv.cn/swagger/](http://ginbro.mojotv.cn/swagger/)
#### [msql生成go代码地址](https://github.com/dejavuzhou/ginbro-son)
#### [bili命令行演示视频地址](https://www.bilibili.com/video/av36804258/)


### 命令行参数详解

```bash
[root@ericzhou felix]# felix ginbro -h
generate a RESTful APIs app with gin and gorm for gophers

Usage:
  felix ginbro [flags]

示例:
felix ginbro -a dev.wordpress.com:3306 -P go_package_name -n db_name -u db_username -p 'my_db_password' -d '~/thisDir'

Flags:
      --authColumn string   使用bcrypt方式加密的用户表密码字段名称 (default "password")
      --authTable string    认知登陆用户表名称 (default "users")
  -a, --dbAddr string       数据库连接的地址 (default "127.0.0.1:3306")
  -c, --dbChar string       数据库字符集 (default "utf8")
  -n, --dbName string       数据库名称
  -p, --dbPassword string   数据库密码 (default "password")
  -t, --dbType string       数据库类型: mysql/postgres/mssql/sqlite (default "mysql")
  -u, --dbUser string       数据库用户名 (default "root")
  -d, --dir string          golang代码输出的目录,默认是当前目录 (default ".")
  -h, --help                帮助
  -l, --listen string       生成go app 接口监听的地址 (default "127.0.0.1:5555")
      --pkg string          生成go app 包名称(go version > 1.12) 生成go.mod文件, eg: ginbroSon

[root@ericzhou felix]# 
```


### web界面

对于那些喜欢使用命令行的同学,你们可以选择使用web界面来操作

```bash
git clone https://github.com/libragen/felix
cd felix
go mod download

go install
echo "添加 GOBIN 到 PATH环境变量"

echo "go build && ./felix -h"

echo 打开Web界面

felix sshw -h

felix sshw

echo "三秒钟之后会自动帮助你打开浏览器,如果如果你使用的windows或者mac系统"

```

#### 1.登陆界面

默认用户名和密码都是 `admin`

![](/images/ginrbo_00.png)

#### 2.填写数据库连接信息

![](/images/ginrbo_01.png)

#### 3.配置app用户认证的表和字段

![](/images/ginrbo_02.png)

#### 4.配置app 包名称,导出目录和监听地址
![](/images/ginrbo_03.png)

#### 5.生成go代码
![](/images/ginrbo_04.png)

#### 6.下载代码或cd者到指定目录
![](/images/ginrbo_05.png)


## 功能二:前端代码二进化,通过gin中间件整合到API服务

### 工作原理
1. 遍历编译好的前端代码目录
2. 使用`archive/zip`写入到`bytes.buffer`中
3. 格式化输出层 字符串常量的 go文件中
4. 创建gin中间件,加载字符串处理,解析出文件
5. 中间件path如果命中文件,这http 输出文件,否在交给下一个handler

### 参数说明
```bash

$ felix ginbin -h
示例: felix ginbin -s dist -p staticbin
Usage:
  felix ginbin [flags]

Flags:
  -c, --comment string   代码注释说明.
  -d, --dest string      出输go代码到目录. (default ".")
  -f, --force            是否覆盖输出. (default true)
  -h, --help             帮助
  -m, --mtime            是否修改文件时间戳.
  -p, --package string   输出的包名称. (default "felixbin")
  -s, --src string       前端静态文件的目录地址. (default "dist")
  -t, --tags string      go 语言的标签.
  -z, --zip              是否zip压缩.

```

### 使用说明:生成的二进制化go文件

vuejs/dist 使用 `felix ginbin` 生成的go文件
[https://github.com/libragen/felix/blob/master/staticbin/gin_static.go](https://github.com/libragen/felix/tree/master/staticbin)

gin 路由应用二进制化的前端代码中间件如下:

`import "github.com/libragen/felix/staticbin" //导入felix ginbin 生成的二进制化包`

[https://github.com/libragen/felix/blob/master/ssh2ws/ssh2ws.go](https://github.com/libragen/felix/blob/master/ssh2ws/ssh2ws.go)

````bash
	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20

	//sever static file in http's root path
	binStaticMiddleware, err := staticbin.NewGinStaticBinMiddleware("/")
	if err != nil {
		return err
	}
	r.Use(binStaticMiddleware)
````

## 引用和代码仓库

### [dejavuzhou/felix Golang 工具集](https://github.com/libragen/felix)
### [felix ginbro 命令逻辑代码目录](https://github.com/libragen/felix/tree/master/ginbro)
### [前端代码二进制化成gin中间件代码](https://github.com/libragen/felix/blob/master/ginbro/ginstatic.go)
### 文章来源 [MojoTech](https://tech.mojotv.cn/2019/05/22/golang-felix-ginbro)
