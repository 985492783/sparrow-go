# Sparrow-go

以go实现核心和控制台，数据面可接入本地文件、配置中心或者k8s configMap中实现云原生

* [快速开始](docs/start/quick_start.md)

## 启动

config.toml中开启自选功能，目前有三种：

* Console：控制台，默认开，端口9800
* Switcher：开关平台，默认开，端口9854
* Logger：日志平台，默认关

> 当前版本未提供，预计v1.1提供

如果无控制台，用户可以通过`cli`操控数据，例如`switcher-cli get com.sparrow.Config#isEnabled -n dev`

## 数据面

| 支持计划          | v1.0 | v1.1 | ... |
|---------------|------|------|-----|
| LocalFile     | [x]  |      |     |
| ConfigCenter  |      | [x]  |     |
| k8s ConfigMap |      |      |     |

## 客户端

Java-SDK： 
* sparrow-client
* sparrow-spring-boot-starter
  * sparrow-switcher-spring-boot-starter
  * sparrow-logger-spring-boot-starter