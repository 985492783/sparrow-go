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

| 支持计划          | v1.0               | v1.1               | ...                  |
|---------------|--------------------|--------------------|----------------------|
| LocalFile     | :white_check_mark: |                    |                      |
| ConfigCenter  |                    | :white_check_mark: |                      |
| k8s ConfigMap |                    |                    | :white_large_square: |


## 客户端

Java-SDK：

* sparrow-client
* sparrow-spring-boot-starter
    * sparrow-switcher-spring-boot-starter
    * sparrow-logger-spring-boot-starter

## Sparrow-Cli
相对控制台来说，cli更轻量，其原理就是一个go sdk+cobra

## 鉴权

> 用户权限通过配置Auth.Permits a:b的形式，多个权限用英文逗号隔开即可，禁用权限加前缀-  
> 客户端与控制台共用一套Auth，所以分配权限时请注意

鉴权开关：authEnabled，默认关闭防止用户poc时产生疑问  
鉴权规则如下：

* `*:*`用户拥有所有权限
* ` `无权限
* `a:b`用户拥有a应用的b权限
* `a:*`用户拥有a应用的所有权限

权限列表

<table>
<tr>
<td>应用</td>
<td>功能</td>
<td>权限标识</td>
<td>备注</td>
</tr>
<tr> 
<td rowspan="4">开关平台Switcher</td>
<td>查询</td>
<td>switcher:list</td>
<td></td>
</tr>

<tr> 
<td>注册并初始化｜注销</td>
<td>switcher:register</td>
<td></td>
</tr>

<tr> 
<td>控制台修改值</td>
<td>switcher:update</td>
<td></td>
</tr>

<tr>
<td></td>
<td></td>
<td></td>
</tr>

<tr> 
<td rowspan="4">控制台</td>
<td>查询配置</td>
<td>console:config_list</td>
<td></td>
</tr>

</table>

**最佳实践**  
样例一：java sdk接入sparrow switcher，但是不希望透出修改能力给客户端
```toml
authEnbaled = true
[Auth."user1"]
password="pass1"
#拥有switcher所有能力，除了更新
permits="switcher:*,-switcher:update"
```

样例二：sparrow-cli查询修改switcher配置
```toml
authEnbaled = true
[Auth."user1"]
password="pass1"
#拥有switcher所有能力，除了更新
permits="switcher:*"
```