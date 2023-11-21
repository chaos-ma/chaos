# server

```
builder 接口面临哪些问题需要解决：
1. 当grpc服务器只有一个的时候，如何处理？
2. 当grpc服务器有多个的时候，如何处理？
3. 当grpc服务器负载不同的时候如何处理？ 或者说： 当多个链接都处于ready状态的时候，应该如何选择？
4. 当链接失败的时候，如何处理？ 是否启动重试等等
grpc为了解决这些问题， 把链路分为不同的阶段：
balancer构建的阶段
子链接具体的连接阶段： 一个grpc服务器地址对应一个连接， 多个地址的时候就会有多个子连接
子连接的选择问题(picker接口完成)
balancer状态
链路创建， 删除，更新

负载均衡相关的接口：
1. Builder接口： 用于构建一个balancer接口实例 - 重点！！
2. SubConn接口： 主要负责具体的连接
3. Picker接口： 主要负责从众多的连接李，按照负载均衡算法选择一个连接供客户端使用 - 最重点！！
4. Balancer接口： 主要负责更新clientConn状态， 更新subConn状态 - 重点！！
5. ClientConn接口： 主要负责链路的维护， 包括创建一个子链路， 删除一个子链路， 更新ClientConn状态
```