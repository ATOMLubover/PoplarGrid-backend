# 后端实现架构说明

## 使用方式

### 启动各个服务器

直接在项目根目录使用如下命令即可启动服务器。  

```bash
bash serve_gateway.bash
bash serve_api.bash
```

这两个脚本的使用方式为 `serve_${SERVER_NAME}.bash [build|run]`，无参数默认编译（包含更新 Swagger 文档） + 运行，`build` 参数将只更新，`run` 参数将运行已存在的可执行文件。  

> 脚本已经经过特别处理，在任意工作目录下启动均不影响效果。  

#### 补充：扩展服务器启动脚本

`serve_gateway.bash` 等底层均是调用的 `serve.bash TARGET_DIR_NAME [build|run]`，因此可以轻松扩展服务器并设置新的启动脚本。  

---

## 配置文件

每个服务由自己的单独的 yaml 配置文件设定参数。  
使用 yaml 的原因是 json 不支持注释。  

---

## 文件架构方式

除了入口 main.go 文件，所有源码均在 ./internal/ 文件夹下，根据是否会共享重用，分布在 ./internal/shared 和 ./internal/$SERVER_NAME/ 下。  

对于提供 HTTP 服务的服务器，其源码在 ./internal/$SERVER_NAME/ 下的文件布局为简化版的 MVC 布局，因为我没有学习过 Spring Boot 框架，可能细化方面不够标准。  

Update Server 是不直接面向客户端开放接口的非 HTTP 服务器，其源码按照功能模块区分，同样位于 ./internal/$SERVER_NAME/ 之下。  

---

## HTTP/WS 框架

整体使用 Iris 框架搭建，采用 MVC 三层经典架构搭建。  
但是减少了 interface 的使用，为了快速适应开发初期的需求变更，直接使用 struct 对接的方式。

### 编码细节

---

## ORM 框架

使用了经典 GORM 框架。所有模型均基于 GORM 标签开始设计。
