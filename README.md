NSQ是一个基于Go语言的分布式实时消息平台，它基于MIT开源协议发布，由bitly公司开源出来的一款简单易用的消息中间件。   

NSQ可用于大规模系统中的实时消息服务，并且每天能够处理数亿级别的消息，其设计目标是为在分布式环境下运行的去中心化服务提供一个强大的基础架构。   

NSQ具有分布式、去中心化的拓扑结构，该结构具有无单点故障、故障容错、高可用性以及能够保证消息的可靠传递的特征。NSQ非常容易配置和部署，且具有最大的灵活性，支持众多消息协议。   

# NSQ组件
 NSQ 由 3 个守护进程组成:

- ***nsqd*** 是接收、队列和传送消息到客户端的守护进程。

- ***nsqlookupd*** 是管理的拓扑信息，并提供了最终一致发现服务的守护进程。

- ***nsqadmin*** 是一个 Web UI 来实时监控集群(和执行各种管理任务)。 

# NSQ工具
* nsq_pubsub
* nsq_stat
* nsq_tail
* nsq_to_file
* nsq_to_http
* nsq_to_nsq
* to_nsq

# NSQ架构  

![NSQ](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/NSQ.png)  

# NSQ详解
![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd.gif)



