# 飞鹅博客部署
## 一、准备工作
+ 服务器或者云主机（如果没有，个人PC装个虚拟机也可以用）

## 一、后端部署
### 1.1 本部署文档环境说明
+ OS：ubuntu 22.04
+ MySql 8.0+
+ CPU制造商为Intel，指令集为x86_64，名称为amd64

> 说明：以下命令，如果使用root账户执行，则无需添加sudo，若是非root账户执行，需要添加sudo，本示例使用非root账户执行

### 1.2 安装辅助工具

```shell
# apt 包管理工具依赖更新
sudo apt-get update
# vim是文本工具 wget是下载flygoose二进制文件使用的
sudo apt-get install vim wget curl -y
```

### 1.3 安装和配置MySql（如果你已经有了，请跳过本步骤）
```shell
# 安装mysql
sudo apt-get install mysql-server -y
# 安装完成后登录
sudo mysql -u root -p
```
因为MySql安装完成后没有密码，当出现以下内容时直接按回车即可
```shell
Enter password:
```
登录成功就会出现以下
```shell
mysql>
```
接下来在`mysql>`输入
```sql
-- 选中mysql这个数据库
use mysql;
-- 创建用户名为flygoose密码为flygoose的用户 这个账户可以管理任何一个数据库
create user 'flygoose'@'%' identified by 'flygoose';
-- 授权flygoose这个账户登录
grant all privileges on *.* to 'flygoose'@'%' with grant option;
-- 刷新权限，使上面的配置生效
flush privileges;
-- 退出 mysql>
exit
```
退出后就自动回到了系统终端。

### 1.4 配置和启动飞鹅后端服务
先下载最新版本的飞鹅二进制包
```shell
# 我将二进制文件放在了/opt下
cd /opt
# 我使用的是linux系统，Intel的CPU，指令集是x86_64 下载的文件是名字中带linux-amd64的
# wget后面的地址是从github仓库中拿到的 页面地址是：https://github.com/helloworld-Co/flygoose-api/releases
# 打开页面后选择最新版本 并且选择服务器对应的二进制文件即可
sudo wget https://github.com/helloworld-Co/flygoose-api/releases/download/tag-2.0-rc1/flygoose-api-linux-amd64-2.0-rc1
```
下载完成后整理配置文件
```shell
touch flygoose-config.yml
vim flygoose-config.yml
```
下面开始整理配置文件(若是有与我的配置有出入的，以你的为准，比如有些同学已经有了自己的mysql，用户名、密码、ip等都以你的为准)
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
  # 数据库名称 如果你跳过了安装MySql那一步，这里需要自己创建一下这个数据库，名字可以自己定
  name: db_flygoose
  # 用户名
  user: flygoose
  # 密码
  password: flygoose
```
vim打开`flygoose-config.yml`后，按键盘上（此时键盘应处于英文状态，中文状态是无法出现效果的）的`i`进入插入数据模式，左下角会出现 `Insert`字样
然后把上面修改好的配置粘贴进去。然后按`ESC`键，就退出编辑状态，然后输入`:wq`，此时左下角会出现`:wq`，然后按回车，就保存并退出了。
接下来给下载的二进制文件赋执行权限并启动
```shell
# 加执行权限 flygoose-api-linux-amd64-2.0-rc1应替换成你下载的二进制文件的名字
sudo chmod +x flygoose-api-linux-amd64-2.0-rc1
# 执行 flygoose-api-linux-amd64-2.0-rc1应替换成你下载的二进制文件的名字
sudo nohup /opt/flygoose-api-linux-amd64-2.0-rc1 -c /opt/flygoose-config.yml & 
```
### 1.5 验证
使用curl在服务器
```shell
curl localhost:29090/api/health
```
或者使用本地浏览器输入`你的服务器ip:29090/api/health`或者POSTMAN使用GET请求`你的服务器ip:29090/api/health`,有下面数据返回即成功。
```json
{"code":1,"data":null,"message":"success"}
```



至此已完成基本的部署

如果你有服务器或者云主机并且具备基本的运维能力，可以继续完成以下配置。

注意：

1. 域名需要一定的RMB
2. CA证书有些平台免费，有些平台付费
3. 服务器或者云主机需要一定的RMB
4. 域名解析的话云服务商大多都有免费的解析可以白嫖

### 1.6 Nginx转发

安装Nginx

配置Nginx转发

### 1.7 域名解析

如果有域名，可以开启HTTPS并解析到你的服务器

### 1.8 Nginx配置HTTPS证书

将你的域名CA证书配置到Nginx

## 附录1 打包命令参考

```shell
# windows Intel芯片
GOOS=windows GOARCH=amd64 go build -o flygoose-api-2.0-win-amd64.exe
# windows 支持arm指令集芯片
GOOS=windows GOARCH=arm64 go build -o flygoose-api-2.0-win-arm64.exe
# linux Intel芯片
GOOS=linux GOARCH=amd64 go build -o flygoose-api-2.0-linux-amd64
# linux 支持arm指令集芯片，如海思麒麟
GOOS=linux GOARCH=arm64 go build -o flygoose-api-2.0-linux-arm64
# mac os 苹果芯片
GOOS=darwin GOARCH=arm64 go build -o flygoose-api-2.0-darwin-arm64
# mac os Intel芯片
GOOS=darwin GOARCH=amd64 go build -o flygoose-api-2.0-darwin-amd64
```

