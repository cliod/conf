# Conf

## 简介

`Conf`是一款用于加载用户配置文件的工具，兼容`yaml`, `json`和`props`格式

## 使用

```go
// 全局配置使用(默认在项目根目录下的conf/app.yaml文件，只加载一次)
v := conf.GetInt("app.version")
n := conf.GetString("app.name")
fmt.Println(v, n)

// 自定义配置(可变长参数，自定义配置)
c := conf.New("custom-conf-dir", "custom-conf-filename", "custom-file-type", "<StoreVariable>")
ver := c.Get("app.version")
fmt.Println(ver)
```
