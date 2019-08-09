## 熔断

#### go-kit 提供了三种熔断

- 1、 gobreaker
- 2、 handy
- 3、 hystrix-go

#### 客户端hystrix配置

- 1、Timeout 【请求超时的时间】
- 2、ErrorPercentThreshold【允许出现的错误比例】
- 3、SleepWindow【熔断开启多久尝试发起一次请求】
- 4、MaxConcurrentRequests【允许的最大并发请求数】
- 5、RequestVolumeThreshold 【波动期内的最小请求数，默认波动期10S】
