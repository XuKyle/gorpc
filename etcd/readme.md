## 服务注册，发现 etcd

##### https://github.com/etcd-io/etcd/releases
##### https://segmentfault.com/a/1190000014045625


客户端流程
- 1、 连接注册中心
- 2、 获取提供的服务
- 3、 监听服务目录变化，目录变化更新本地缓存
- 4、 创建负载均衡器
- 5、 获取请求的 endPoint