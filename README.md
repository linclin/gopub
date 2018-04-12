gopub（基于vue.js element框架+golang beego框架开发）是一个基于运维场景设计的企业级运维发布系统。配置简单、功能完善、界面流畅、开箱即用！支持git、jenkins版本管理，支持各种web代码发布，一键完成Golang,nodejs,PHP，Python，JAVA等代码的发布、回滚操作。

我们运维团队前期使用walle web部署系统进行发布操作,在此也感谢walle团队贡献的优秀开源项目.walle的web体验比较好,本次开源的gopub前台完全模仿walle前台,使用vue.js element框架重写.

gopub已持续运行近两年时间,在我们预发布和生产环境完成37000+次稳定部署,支持单项目100+台服务器部署110秒左右,支持2G+CDN静态资源发布传输30秒完成.
![统计](docs/images/index.png)

## 代码地址
* [开源中国码云 https://gitee.com/dev-ops/gopub](https://gitee.com/dev-ops/gopub)
* [Github https://github.com/linclin/gopub](https://github.com/linclin/gopub) 

## 使用框架
* [Element](http://element-cn.eleme.io/#/zh-CN)
* [Beego](https://beego.me/)
* [httprouter](https://github.com/julienschmidt/httprouter) 
* [Taipei-Torrent](https://github.com/jackpal/Taipei-Torrent) 

## 功能特性
* Docker&k8s支持：Docker镜像仅60M,kubernetes编排文件一键部署运行
* 部署简便：go二进制部署,无需安装运行环境.
* gitlab发布支持：配置每个项目git地址,自动获取分支,commit选择并自动拉取代码
* jenkins发布支持：支持jenkins编译包一键发布
* ssh执行命令/传输文件：使用golang内置ssh库高效执行命令/传输文件
* BT支持：大文件和大批量机器文件传输使用BT协议支持
* 多项目部署:支持多项目多任务并行,内置[grpool协程池](https://github.com/linclin/grpool)支持并发操作命令和传输文件
* 全web化操作：web配置项目,一键发布,一键快速回滚
* API支持：提供所有配置和发布操作API,便于对接其他系统  [API使用example](api_example/example.go)
* 部署钩子：支持部署前准备任务,代码检出后处理任务,同步后更新软链前置任务,发布完毕后收尾任务4种钩子函数脚本执行

## Docker 快速启动
``` shell
#使用dockerhub镜像启动,连接外部数据库
sudo docker run --name gopub -e MYSQL_HOST=x.x.x.x -e MYSQL_PORT=3306  -e MYSQL_USER=root -e MYSQL_PASS=123456 -e MYSQL_DB=walle -p 8192:8192  --restart always  -d   lc13579443/gopub:latest 
```
### Docker 镜像制作
``` shell
# 使用multi-stage(多阶段构建)需要docker 17.05+版本支持
sudo docker build -t  gopub .
sudo docker run --name gopub -e MYSQL_HOST=x.x.x.x  -e MYSQL_PORT=3306  -e MYSQL_USER=root -e MYSQL_PASS=123456 -e MYSQL_DB=walle -p 8192:8192  --restart always  -d  gopub:latest 

```
### Kubernetes 快速部署
``` shell 
# apiVersion: apps/v1需要kubernetes 1.9.0+版本支持
kubectl apply -f  gopub-kubernetes.yml

```

## 源码编译安装
### 编译环境
- golang >= 1.8+ 
- nodejs >= 4.0.0（编译过程中需要可以连公网下载依赖包）

### 源码下载

``` shell
# 克隆项目
git clone https://gitee.com/dev-ops/gopub.git

# 编译前端,npm较慢可使用cnpm

cd vue-gopub
npm install
npm run build

#修改配置 数据库配置文件在 src/conf/app.conf

#编译,control需要给可执行权限,并修改go安装目录 export GOROOT=/usr/local/go
./control build

#执行数据库初始化
./control init

#启动服务 启动成功后 可访问 127.0.0.1:8192 用户名:admin 密码:123456
./control start

#停止服务
./control stop

#重启服务
./control restart
```

### 快速使用
#### 下载项目[二进制包](https://gitee.com/dev-ops/gopub/attach_files/download?i=127803&u=http%3A%2F%2Ffiles.git.oschina.net%2Fgroup1%2FM00%2F03%2F4D%2FPaAvDFrLHHaAOuz_AJ9X2n198H45982.gz%3Ftoken%3D1ab1bd5c19af447d6024db8ce7054df1%26ts%3D1523330628%26attname%3Dgopub-1.0.1.tar.gz)，无需安装go环境和node环境
``` shell
#给control和src/gopub给可执行权限

#执行数据库初始化
./control init

#启动服务 启动成功后 可访问 127.0.0.1:8192 用户名:admin 密码:123456
./control start

#停止服务
./control stop

#重启服务
./control restart
```
## 配置ssh-key信任
前提条件:gopub运行用户(如root)ssh-key必须加入目标机器的{remote_user}用户ssh-key信任列表

``` shell

#添加机器信任
su {local_user} && ssh-copy-id -i ~/.ssh/id_rsa.pub remote_user@remote_server

#need remote_user's password
#免密码登录需要远程机器权限满足以下三个条件：
/home/{remote_user} 755
~/.ssh 700
~/.ssh/authorized_keys 644 或 600
```


## Getting started
### 1. 项目配置

![项目配置](docs/images/project.png)

* 项目名称：xxx.example.com   （项目命名一定要规范并唯一）

* 项目环境：现在只用到预发布环境和线上环境。

* 地址：支持gitlab,jenkins,file三种发布方式.

 选用Git在地址栏里面填入git地址，https方式需在地址中加入账号密码,ssh方式需保证gopub所在服务器有代码拉取权限.我们一般在gitlab创建一个public用户,将gopub所在服务器key加入public用户deploy-keys设置,并将public用户授权可拉取所有gitlab项目.

 选用jenkins需要录入jenkins对于的job地址和账号密码,


#### 宿主机
* 代码检出库：/data/www/xxx (名称需要唯一)
* 排除文件：默认不填写,可填写.git(tar打包忽略.git目录)等其他需要打包忽略文件

#### 目标机器
* 用户：www  (目标机执行操作用户)
* webroot：/data/htdocs/shell_php7 (目标机代码发布目录,软链目录)
* 发布版本库：/data/htdocs/backup/shell_php7 (目标机代码备份目录,实体目录,* * * webroot软链到该目录下的对应发布目录)
* 版本保留数：20 (发布版本库中保留多少个发布历史)
* 机器列表：一行一个IP  （复制粘贴ip的时候注意特殊字符）

#### 高级任务
前面两个任务的执行是在管理机上，后面两个任务的执行是在目标机器上

* 代码检出前任务：视情况而定（默认为空）
* 代码检出后任务： 需要composer的项目需要添加：cd {WORKSPACE} && rm -rf composer.lock vendor && composer install --optimize-autoloader --no-dev -vvv --ignore-platform-reqs ，否则为空
* 同步完目标机后任务：视情况而定（默认为空）
* 更改版本软链后任务：视情况而定（默认为空）

### 2. 创建上线单
![创建上线单](docs/images/pub1.png)
* gitlab上线单
![git配置](docs/images/pub2-git.png)
* jenkins上线单
![jenkins配置](docs/images/pub2-jenkins.png)

### 3. 部署操作 
![选择上线单](docs/images/pub3.png)
![部署](docs/images/pub4.png)


## 开发团队
* [林超](https://github.com/linclin)
* 高传泽
* 金阳
* 赵连启
* 张群烽
 
## 下个版本计划
* 与jenkins接口对接,支持发布包可下拉选择
* 支持选择蓝鲸CMDB3.0业务模块发布,避免维护IP列表
* 现有的Docker镜像基于centos打包,镜像超过1.4G,下个版本使用alpine作为基础镜像,减少镜像大小.并支持kubernetes编排

## 技术支持 
联系我们，技术交流QQ群：

![wechat](docs/images/qq.png)