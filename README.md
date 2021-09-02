# redis
基于net/http封装的http库

# 功能
 - 简化Server实例初始化流程，基于配置自动对Server进行初始化且启动
 - 简化Client实例初始化流程，指定请求参数及响应结构体，基于配置自动发起网络请求并将响应结果反序列化为指定结构体
 - Client支持连接池、多实例等场景
 - Client的数据序列化支持json、pb、gob、thrift，请求方式支持get、post、json\_post、file\_post、bytes\_post


# HttpServer 基本使用
 - toml 配置文件
    ```
    [rest.server.default]
        addr = ":8080"
        read_timeout = 1000
        write_timeout = 1000
        dial_timeout = 1000
    ```

 - 使用方法
	```golang
    // 初始化config包，参考config模块
    code...

    // 创建Gin对象并配置路由
    server := gin.New()
    code...

    // 如下方式可以直接使用http server实例
    if err := rest.StartServer("default", server); err != nil {
        fmt.Printf("[ERROR] Build rest service error: %s\n", err)
        os.Exit(-1)
    }

    // 停止http Server
    rest.StopServer("default")
    ```

# HttpClient 基本使用
 - toml 配置文件
    ```
    [rest.client.ping]
        url = "http://127.0.0.1/ping"
        codec = "json"
        method = "post"
        fail_retry = 1
        fail_retry_interval = 300
        dial_timeout = 600
        dial_disable_keep_alive = false
        pool_max_idle = 100
        pool_max_idle_per_host = 100
        pool_idle_time = 90000
    ```

 - 使用方法
	```golang
    // 初始化config包，参考config模块
    code...

    params := map[string]string{"seconds": "10"}
    body := &Body{...}
    resp := &Resp{}
    _, err := http_cli.Call("ping", params, body, resp)
    if err == nil {
        fmt.Println(resp)
    }
    ```
