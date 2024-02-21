# flygoose

！！！TODO: 需要改成2.0-rc版本的部署文档  下面这个还没改！！！

## 一、依赖

+ go 1.19+（需要自己安装），建议proxy修改为国内的地址，不然会被墙，参考[https://goproxy.cn/](https://goproxy.cn/)
+ mysql 8.0+ 或 postgresql 12+（需要自己安装）

## 二、架构

![架构](architecture.png)

## 三、部署文档

### 3.1 本地启动

1. Gland打开项目，在go.mod文件夹右键，选 `Go Mod Tidy`，会自动下载好依赖

2. 创建数据库`db_flygoose`。

   注意: Mysql8.0+ 字符集使用`utf8mb4`，排序规则使用 `utf8mb4_0900_ai_ci`，Postgresql创建数据库字符集为`utf8`

#### 3.1.1 Admin启动

环境配置：flygoose/cmd/admin/admin-config.yaml

```yml
# 服务端口
http:
  port: 29091
database:
  # 数据库类型 取值mysql或postgresql
  driver: mysql
  # 数据库IP
  host: 127.0.0.1
  # 数据库端口
  port: 3306
  # 数据库名称
  name: db_flygoose
  # 用户名
  user: flygoose
  # 密码
  password: flygoose
```

启动配置在`flygoose/cmd/admin/main.go`，下面这行：

```go
configPath := flag.String("c", "cmd/admin/admin-config.yaml", "指定配置文件路径")
```

#### 3.1.2 Flygoose启动

dev环境配置：flygoose/cmd/flygoose/flygoose-config.yaml

```yaml
# 服务端口
http:
  port: 29090
database:
  # 数据库类型 取值mysql或postgresql
  driver: mysql
  # 数据库IP
  host: 127.0.0.1
  # 数据库端口
  port: 3306
  # 数据库名称
  name: db_flygoose
  # 用户名
  user: flygoose
  # 密码
  password: flygoose
```

启动配置在`flygoose/cmd/flygoose/main.go`，下面这行：

```go
configPath := flag.String("c", "cmd/flygoose/flygoose-config.yaml", "指定配置文件路径")
```

### 3.2 打包部署

1. 下载项目
2. 修改配置文件 `cmd/admin/admin-config.yaml `和 `cmd/flygoose/flygoose-config.yaml`
3. 打包(本示例假设服务器系统为linux，CPU型号是x86_64)

```shell
# 进入flygoose根目录执行以下脚本
go mod tidy
# 打包flygoose脚本 ----start
cd cmd/flygoose
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0
go build main.go -o flygoose
# 打包flygoose脚本 ----end
# 打包admin脚本 ----start
cd ../admin
set GOARCH=amd64
set GOOS=linux
set CGO_ENABLED=0
go build main.go -o flygoose-admin
# 打包admin脚本 ----end
```

执行完成后会在`cmd/flygoose`生成一个名为`flygoose`可执行文件，在`cmd/admin`下回生成一个名为`flygoose-admin`

4. 添加权限并启动

```shell
# 加执行权限
chomod +x flygoose
chomod +x flygoose-admin
# 启动
nohup ./flygoose -c ./flygoose-config.yaml &
nohup ./flygoose-admin -c ./admin-config.yaml &
```

5. 测试（以下测试端口号请换成自己配置的）

```log
# 服务器上curl 测试  输出 {"code":1,"data":null,"message":"success"} 即认为成功
curl 127.0.0.1:29090/v1/health
curl 127.0.0.1:29091/v8/health
# 本地postman 测试 输出 {"code":1,"data":null,"message":"success"} 即认为成功
GET http://你的服务器ip:29090/v1/health
GET http://你的服务器ip:29091/v8/health
```

### 3.3 二进制文件部署

1. 下载二进制文件，需要下载`flygoose`和`flygoose-admin`  地址：`https://github.com/helloworld-Co/flygoose-api/releases`
2. 修改配置并保存文件

```yaml
# 服务端口
http:
  port: 29090
database:
  # 数据库类型 取值mysql或postgresql
  driver: mysql
  # 数据库IP
  host: 127.0.0.1
  # 数据库端口
  port: 5432
  # 数据库名称
  name: db_ee
  # 用户名
  user: postgres
  # 密码
  password: postgres
```

以上是`flygoose`的配置文件信息，内容修改成自己的配置，保存为`flygoose-config.yaml`

```yaml
# 服务端口
http:
  port: 29091
database:
  # 数据库类型 取值mysql或postgresql
  driver: postgresql
  # 数据库IP
  host: 127.0.0.1
  # 数据库端口
  port: 5432
  # 数据库名称
  name: db_ee
  # 用户名
  user: postgres
  # 密码
  password: postgres
```

以上是`flygoose-admin`admin的配置文件信息，内容修改成自己的配置，保存为`admin-config.yaml`

3. 加执行权限并启动

```
# 加执行权限
chomod +x flygoose
chomod +x flygoose-admin
# 启动
nohup ./flygoose -c flygoose-config.yaml &
nohup ./flygoose-admin -c admin-config.yaml &
```

> ！！！特别注意！！！： flygoose和flygoose-admin请放在同一个文件夹
> 因为两个服务共用部署文件夹下的/static目录，如果不放在同一个目录会导致flygoose-admin上传的图片flygoose无法访问到
> 如果想扩展可以考虑以下思路：
> 1. flygoose提供图片上传和访问接口，flygoose-admin上传图片时会把请求转发到flygoose，上传和访问都用flygoose接口
> 2. flygoose前端项目部署用到了nginx，可以放在nginx下的静态文件目录中，利用nginx提供图片访问功能
> 3. 自己部署或者使用云服务，例如fastDFS，MinIO，各类图床，阿里OSS

5. 测试（以下端口号请换成自己配置的）

```log
# 服务器上 curl 测试  输出 {"code":1,"data":null,"message":"success"} 即认为成功
curl 127.0.0.1:29090/v1/health
curl 127.0.0.1:29091/v8/health
# 本地 postman 测试 输出 {"code":1,"data":null,"message":"success"} 即认为成功
GET http://你的服务器ip:29090/v1/health
GET http://你的服务器ip:29091/v8/health
```

## Docker 方式运行
以下为快速运行飞鹅后端的容器部署步骤方式，生产环境请把数据库持久化高可用部署
### install postgresql
```
docker run --rm --name flygoose-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres  -p 5432:5432 postgres
```

### 创建db
```
docker exec -it -u postgres flygoose-postgres psql -c "create database db_ee"
```

### 构建flygoose镜像
```
sh build.sh
```

### 运行飞鹅后端接口服务
```
docker run  --net=host  -it --rm flygoose:v1
```

### 运行flygoose admin api
```
docker run  --net=host  -it --rm flygoose:v1 sh -c "/apps/admin/admin -c /apps/admin/admin-config.yaml"
```

### 测试验证
```
http://127.0.0.1:29090/v1/health  # flygoose api
http://127.0.0.1:29091/v8/health  # flygoose admin api
```