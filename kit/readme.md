## gokit + grpc

#### 介绍
go-kit 是一个微服务的开发工具集，微服务系统中的大多数常见问题，因此，使用者可以将精力集中在业务逻辑上。

grpc 缺乏服务治理的功能，我们可以通过 go-kit 结合 grpc 来实现我们的完整需求。go-kit 抽象的 endpoint 设计让我们可以很容易包装其它微服务框架使用的协议。

go-kit 提供以下功能：

- 1、Circuit breaker（熔断器）
- 2、Rate limiter（限流器）
- 3、Logging（日志）
- 4、Metrics（Prometheus 统计）
- 5、Request tracing（请求跟踪）
- 6、Service discovery and load balancing（服务发现和负载均衡）


#### go-kit [TransportServer]
一个 Transport 的 Server 必须要拥有 endPoint、decodeRequestFunc、encodeResponseFunc

- 1、 endPoint 一个端点代表一个 RPC，也就是我们服务接口中的一个函数
- 2、 decodeRequestFunc 请求参数解码函数
- 3、 encodeResponseFunc 返回参数编码函数
