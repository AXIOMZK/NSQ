 NSQ 由 3 个守护进程组成:


- ***nsqd*** 是接收、队列和传送消息到客户端的守护进程。

- ***nsqlookupd*** 是管理的拓扑信息，并提供了最终一致发现服务的守护进程。

- ***nsqadmin*** 是一个 Web UI 来实时监控集群(和执行各种管理任务)。

我们来看一下它的架构   

![NSQ](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/NSQ架构-6.pdf)   



