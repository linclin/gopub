# 发布系统配置说明

- sever端配置均放在 agent/server.json配置项内
- agent端配置为了保证代码简单性，均写在代码中

## 配置说明

使用json格式

```
Name:  #名称
DownDir： #下载位置(服务端要和客户端一致)
Net:
    MgntPort:  #管理端口，用于接收客户端的创建任务等Rest接口(默认配置为45003)
    DataPort:  #服务端数据下载端口(默认配置为45002)
    AgentMgntPort:  #Agent端的管理端口，用于接收Server下载的管理Rest接口（服务端使用，对应的agent的mgntPort）
    AgentDataPort:  #Agent端的数据下载端口（服务端使用，对应的agent的dataPort）
Auth:
    Username: 1234  #用于认证的用户名
    Password: 1234  #用于认证的密码
Control:
    Speed: 10  # 流量控制，单位为MBps
    CacheSize: 50 # 文件下载的内存缓存大小，单位为MB
    MaxActive: 10 # 并发的任务数
Server:  #是否为服务端(服务端使用必须为true)
```