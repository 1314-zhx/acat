# ACAT 实验室纳新网站源码(mvp版本)
> 后续开发者要了解该项目，从本文档开始


> 前言：
> 本项目后端使用`Go`语言编写，并且源码已在`github`上开源。
> 第一版源码在`github.com/1314-zhx/acat`，
> 

> 
> 项目介绍：
> 本项目是使用`Go+gin框架`搭建的一个简单的招生项目web网站。

# 开发人员名单

- 第一期开发人员：张皓翔[后端Go组] 缑梦洁[前端] 张心怡[前端] `有疑问可以咨询`
- 第二期开发人员：(请后续开发者在此处填写，并删除这段话)

## 目录
- [项目环境](#项目环境)
- [组件安装](#组件安装)
- [项目结构](#项目结构)
- [项目部署](#项目部署)
- [压力测试](#压力测试)
- [快速开始](#快速开始)
- [API接口](#API接口)
- [拓展须知](#拓展须知)
- [后续优化建议](#后续优化建议)

## 项目环境
> 开发环境

go      version     1.24.0

gin     version     1.11.0 

mysql   version     9.3.0

redis   version     5.0.14.1

docker  version     27.5.1

nginx   version     1.28.0

## 组件安装

> 1.zap 日志安装
```
go get go.uber.org/zap@v1.27.0
```

> 2.Lumberjack  日志轮转库安装

```
go get gopkg.in/natefinch/lumberjack.v2@latest
```

> 3.swag 安装

```
go install github.com/swaggo/swag/v2/cmd/swag@latest
```

> 4.gorm 安装

```
go get -u gorm.io/gorm
```

> 5.viper 安装

```
go get github.com/spf13/viper
```
## 项目结构
```
ACAT-
    -conf           // 配置文件
    -controller     // 控制函数
    -dao            // 数据库操作
    -docs           // 文档
    -logger         // 日志文件
    -logic          // 逻辑处理
    -model          // 模型
    -router         // 路由
    -setting        // 环境配置
    -util           // 工具包
```

## 项目部署
> 本地部署，用于测试


详情见 `Dockerfile` 文件

本地构建 `docker` 镜像
```
docker build -t acat .
```
运行容器
```
docker run -d -p 9090:9090 --name acat acat
```

> 远程部署
通过 scp + docker load

`docker` 镜像推送至服务器，由 `nginx` 暴露端口。本地监听 `9090` 端口

- 步骤 1：本地保存镜像为 tar 文件
```
docker save acat -o acat.tar
```
- 步骤 2：上传到服务器
```
scp acat.tar root@47.92.xxx.xxx:/root/
```
- 步骤 3：在服务器上加载并运行
```
ssh root@47.92.xxx.xxx

docker load -i /root/acat.tar

docker run -d -p 9090:9090 --name acat --restart unless-stopped acat
```

## 压力测试

本地压力测试，使用 `wrk` ，也可以使用 `ab` ，但本人使用的是 `wrk`。

因为 `Windows` 无法使用 `wrk` ，需要 `liunx` 环境下

在 `ubuntu` 终端输入 `Windows` 在 `liunx` 下的文件路径(以我的开发环境为例，我的本地路径是D://ACAT)

```
cd /mnt/d/acat 
go run .
```
在另一 `ubuntu` 终端输入
```
wrk -t10 -c1000 -d5 http://127.0.0.1:9090
```
以根路由为例，后续测试自行更改参数和路径

## 快速开始

>本地快速开始，无须使用docker
> 温馨提示，国内下载较慢可以配置GOPROXY，用.cn域名，.io有时也比较慢
```
go env -w GOPROXY=https://goproxy.cn,direct
```
任意文件夹下使用

```
git clone https://github.com/1314-zhx/acat.git
```

将源码克隆本地

在主文件夹下运行

```
cd acat
go run .
```

## API 接口

nil

## 拓展须知

> 后续开发者须知

若后续版本更新，或要修改源码，进行业务拓展。需要做到如下几点，以保证项目的延续

1.修改完后，必须修改`README.md`文档，根目录下的，和`docs`两处。要求修改
项目环境，组件安装，等。

2.更新`github`克隆的地址，为你的地址。

3.开发者根据服务器IP和域名更新外部部署的命令路径。

4.每期开发人员须在`开发人员名单`中留名，以保证项目知识迁移，新老人员对接。
以 `xxx[小组名]` 格式。

**协作原则
“改代码，必改文档；换环境，必更新指引”**

## 后续优化建议

我有点懒，并且因为经验不足导致在刚开始设计数据库模型的时候，很多东西没有考虑到
，比如同一个东西在甲表叫Id，在另一个表叫Slot_Id，或者叫ID，导致我有时候我都看不懂
我在干什么，并且写完了才发现我好像没做日志跟踪。

所以后续要优化的话，可以把表结构优化一下，让一个东西在所有表中名字含义尽量相同。可以增加日志追踪，
我的建议是加一个AdminLogger结构体，保存操作者ID，干了什么，什么时候干的。

并且可以把前端做的再优美一点，因为前端学姐比较高冷，向她们请教问题，比见明星都难，所以这一版
的前端比较简陋。后续如果有前端大佬可以加个轮播图之类的，多加点照片啥的，好看就行。

还有就是我其实在admin模块做的比较简陋，没有日志分级，导致全部都是Info级别，后面也可以优化一下，
我不想搞了。

哦，再说一点，其实我没用swag，虽然有docs这个目录，但其实我没写swag注释，所以其实那个接口页面什么也没有，
本来刚起步时是想用的，但是写到后面，干劲没了就没加，感觉太烦人了，后面想加也可以。

还有在注册的时候，为什么用邮箱发送验证码，这个确实没办法，本来想的是手机号验证码和邮箱验证码
都有的，但是咨询阿里云客服时，人家说最近把个人短信权限取消了，就是说现在个人不能
调用阿里云API发送短信不信了，只能由公司的注册商标去申请，并且邮箱还是免费的，索性就
把手机短信方案取消了。但是我用的QQ邮箱，这个可能会被限流，所以日后可以用专业的发送邮箱的工具。
