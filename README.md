NSQ是一个基于Go语言的分布式实时消息平台，它基于MIT开源协议发布，由bitly公司开源出来的一款简单易用的消息中间件。   

NSQ可用于大规模系统中的实时消息服务，并且每天能够处理数亿级别的消息，其设计目标是为在分布式环境下运行的去中心化服务提供一个强大的基础架构。   

NSQ具有分布式、去中心化的拓扑结构，该结构具有无单点故障、故障容错、高可用性以及能够保证消息的可靠传递的特征。NSQ非常容易配置和部署，且具有最大的灵活性，支持众多消息协议。   

# NSQ组件
 NSQ 由 3 个守护进程组成:

- ***nsqd*** 是接收、队列和传送消息到客户端的守护进程。

- ***nsqlookupd*** 是管理的拓扑信息，并提供了最终一致发现服务的守护进程。

- ***nsqadmin*** 是一个 Web UI 来实时监控集群(和执行各种管理任务)。 


# NSQ架构  

![NSQ](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/NSQ.png)  

* ## ***producer***   消息的生产者   
  +  ***```producer```*** 通过 **```HTTP API```** 将消息发布到 ***```nsqd```*** 的指定 **```topic```** ，一般有 **```pub/mpub```** 两种方式， **```pub```** 发布一个消息， **```mpub```** 一个往返发布多个消息。  
  +  ***```producer```*** 也可以通过 **```nsqd客户端```** 的 **```TCP接口```** 将消息发布给 ***```nsqd```*** 的指定 **```topic```** 。
  +  当生产者 ***```producer```*** 初次发布带 **```topic```** 的消息给 ***```nsqd```*** 时,如果 **```topic```** 不存在，则会在 ***```nsqd```*** 中创建 **```topic```** 。
 
* ##  ***channels***   
  + 当生产者每次发布消息的时候,消息会采用多播的方式被拷贝到各个channel中,channel起到队列的作用。
  + 

* ## 概述  
  1. NSQ推荐通过他们相应的nsqd实例使用协同定位发布者，这意味着即使面对网络分区，消息也会被保存在本地，直到它们被一个消费者读取。更重要的是，发布者不必去发现其他的nsqd节点，他们总是可以向本地实例发布消息。
  2. 首先，一个发布者向它的本地nsqd发送消息，要做到这点，首先要先打开一个连接，然后发送一个包含topic和消息主体的发布命令，在这种情况下，我们将消息发布到事件topic上以分散到我们不同的worker中。
事件topic会复制这些消息并且在每一个连接topic的channel上进行排队，在我们的案例中，有三个channel，它们其中之一作为档案channel。消费者会获取这些消息并且上传到S3。
  3. 每个channel的消息都会进行排队，直到一个worker把他们消费，如果此队列超出了内存限制，消息将会被写入到磁盘中。Nsqd节点首先会向nsqlookup广播他们的位置信息，一旦它们注册成功，worker将会从nsqlookup服务器节点上发现所有包含事件topic的nsqd节点。
  4. 然后每个worker向每个nsqd主机进行订阅操作，用于表明worker已经准备好接受消息了。这里我们不需要一个完整的连通图，但我们必须要保证每个单独的nsqd实例拥有足够的消费者去消费它们的消息，否则channel会被队列堆着。


# NSQD详解
![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd.gif)   

* topic: topic是nsq的消息发布的逻辑关键词。当程序初次发布带topic的消息时,如果topic不存在,则会被创建。
* channels: 当生产者每次发布消息的时候,消息会采用多播的方式被拷贝到各个channel中,channel起到队列的作用。
* messages: 数据流的形式。


# NSQ工具
* nsq_pubsub
* nsq_stat
* nsq_tail
* nsq_to_file
* nsq_to_http
* nsq_to_nsq
* to_nsq

# 参考文献

[[1]GoDoc of nsq [EB/OL]](https://godoc.org/github.com/bitly/nsq)  
[[2]NSQ v1.0.0-compat DESIGN [EB/OL]](https://nsq.io/overview/design.html)  
[[3]消息中间件NSQ深入与实践 [EB/OL]](https://www.jianshu.com/p/2b01403e3443)  
[[4]NSQ消息队列详解 [EB/OL]](http://xiaoh.me/2018/02/13/nsq-summary/)  
[[5]NSQ使用总结 [EB/OL]](http://lihaoquan.me/2016/6/20/using-nsq.html)



