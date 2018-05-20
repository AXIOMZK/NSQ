> NSQ是一个基于Go语言的分布式实时消息平台，它基于MIT开源协议发布，由bitly公司开源出来的一款简单易用的消息中间件。

> NSQ可用于大规模系统中的实时消息服务，并且每天能够处理数亿级别的消息，其设计目标是为在分布式环境下运行的去中心化服务提供一个强大的基础架构。

> NSQ具有分布式、去中心化的拓扑结构，该结构具有无单点故障、故障容错、高可用性以及能够保证消息的可靠传递的特征。NSQ非常容易配置和部署，且具有最大的灵活性，支持众多消息协议。  

* [NSQ组件](#nsq%E7%BB%84%E4%BB%B6)
* [NSQ架构](#nsq%E6%9E%B6%E6%9E%84)
  * [<em><strong>topic</strong></em>  消息的逻辑关键词](#topic--%E6%B6%88%E6%81%AF%E7%9A%84%E9%80%BB%E8%BE%91%E5%85%B3%E9%94%AE%E8%AF%8D)
  * [<em><strong>producer</strong></em>   消息的生产者/发布者](#producer---%E6%B6%88%E6%81%AF%E7%9A%84%E7%94%9F%E4%BA%A7%E8%80%85%E5%8F%91%E5%B8%83%E8%80%85)
  * [<em><strong>channel</strong></em>  消息传递的通道](#channel--%E6%B6%88%E6%81%AF%E4%BC%A0%E9%80%92%E7%9A%84%E9%80%9A%E9%81%93)
  * [<em><strong>consumer</strong></em>  消息的消费者](#consumer--%E6%B6%88%E6%81%AF%E7%9A%84%E6%B6%88%E8%B4%B9%E8%80%85)
  * [概述](#%E6%A6%82%E8%BF%B0)
* [NSQD](#nsqd)
  * [<em><strong>nsqd</strong></em> 基本结构](#nsqd-%E5%9F%BA%E6%9C%AC%E7%BB%93%E6%9E%84)
  * [<em><strong>nsqd</strong></em> 详细流程图](#nsqd-%E8%AF%A6%E7%BB%86%E6%B5%81%E7%A8%8B%E5%9B%BE)
  * [<em><strong>nsqd</strong></em> 源码概述](#nsqd-%E6%BA%90%E7%A0%81%E6%A6%82%E8%BF%B0)
* [NSQ使用](#nsq%E4%BD%BF%E7%94%A8)
  * [测试1](#%E6%B5%8B%E8%AF%951)
  * [测试2](#%E6%B5%8B%E8%AF%952)
* [参考文献](#%E5%8F%82%E8%80%83%E6%96%87%E7%8C%AE)

# NSQ组件
 NSQ 由 3 个守护进程组成:

- ***nsqd*** 是接收、队列和传送消息到客户端的守护进程。

- ***nsqlookupd*** 是管理的拓扑信息，并提供了最终一致发现服务的守护进程。

- ***nsqadmin*** 是一个 Web UI 来实时监控集群(和执行各种管理任务)。 


# NSQ架构  

![NSQ](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/NSQ.png)

* ##  ***topic***  消息的逻辑关键词
  +  **```topic```** 是 **```NSQ```** 消息发布的 **```逻辑关键词```** ，可以理解为人为定义的一种消息类型。当程序初次发布带 **```topic```** 的消息时,如果 **```topic```** 不存在,则会在 ***```nsqd```*** 中创建。
  

* ## ***producer***   消息的生产者/发布者   
  +  ***```producer```*** 通过 **```HTTP API```** 将消息发布到 ***```nsqd```*** 的指定 **```topic```** ，一般有 **```pub/mpub```** 两种方式， **```pub```** 发布一个消息， **```mpub```** 一个往返发布多个消息。  
  +  ***```producer```*** 也可以通过 **```nsqd客户端```** 的 **```TCP接口```** 将消息发布给 ***```nsqd```*** 的指定 ***```topic```*** 。
  +  当生产者 ***```producer```*** 初次发布带 ***```topic```*** 的消息给 ***```nsqd```*** 时,如果 ***```topic```*** 不存在，则会在 ***```nsqd```*** 中创建 ***```topic```*** 。
  

* ##  ***channel***  消息传递的通道
  + 当生产者每次发布消息的时候,消息会采用多播的方式被拷贝到各个 ***```channel```*** 中, ***```channel```*** 起到队列的作用。
  +  ***```channel```*** 与 ***```consumer(消费者)```*** 相关，是消费者之间的负载均衡,消费者通过这个特殊的channel读取消息。
  + 在 ***```consumer```*** 想单独获取某个 **```topic```** 的消息时，可以 ***```subscribe(订阅)```*** 一个自己单独命名的 ***```nsqd```*** 中还不存在的 ***```channel```*** , ***```nsqd```*** 会为这个 ***```consumer```*** 创建其命名的 ***```channel```*** 。 
  + ***```Channel```*** 会将消息进行排列，如果没有 ***```consumer```*** 读取消息，消息首先会在内存中排队，当量太大时就会被保存到磁盘中。可以在配置中配置具体参数。
  + 一个 ***```channel```*** 一般会有多个 ***```consumer```*** 连接。假设所有已连接的 ***```consumer```*** 处于准备接收消息的状态，每个消息将被传递到一个随机的 ***```consumer```*** 。
  + Go语言中的channel是表达队列的一种自然方式，因此一个NSQ的topic/channel，其核心就是一个存放消息指针的Go-channel缓冲区。缓冲区的大小由  --mem-queue-size 配置参数确定。

  
* ## ***consumer***  消息的消费者 
  +   ***consumer*** 通过 **TCP** ***```subscribe```*** 自己需要的 ***```channel```*** 
  +   ***```topic```*** 和 ***channel*** 都没有预先配置。 ***```topic```*** 由第一次发布消息到命名 ***```topic```*** 的 ***producer*** 创建 ```或``` 第一次通过 ***```subscribe```*** 订阅一个命名 ***```topic```*** 的 ***consumer*** 来创建。 ***channel*** 被 ***consumer*** 第一次 ***```subscribe```*** 订阅到指定的 ***channel*** 创建。
  +  多个 ***consumer***  ***```subscribe```*** 一个 ***channel*** ，假设所有已连接的客户端处于准备接收消息的状态，每个消息将被传递到一个 **```随机```** 的 ***consumer*** 。
  +   NSQ 支持延时消息， ***consumer*** 在配置的延时时间后才能接受相关消息。
  +   Channel在 ***consumer*** 退出后并不会删除，这点需要特别注意。

* ## 概述  
  1. NSQ推荐通过 ***```nsqd```*** 实例使用协同定位 ***producer*** ，这意味着即使面对网络分区，消息也会被保存在本地，直到它们被一个 ***consumer*** 读取。更重要的是， ***producer*** 不必去发现其他的 ***```nsqd```*** 节点，他们总是可以向本地 ***```nsqd```*** 实例发布消息。
  2. 一个 ***producer*** 向它的本地 ***```nsqd```*** 发送消息，要做到这点，首先要先打开一个连接( NSQ 提供 ```HTTP API``` 和 ```TCP 客户端``` 等2种方式连接到 ***```nsqd```*** )，然后发送一个包含 ***```topic```*** 和消息主体的发布命令(pub/mpub/publish)，在这种情况下，我们将消息发布到 ***```topic```*** 上，消息会采用多播的方式被拷贝到各个 ***```channel```*** 中, 然后通过多个 ***```channel```*** 以分散到我们不同需求的 ***consumer*** 中。  
  ![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd/nsqd.gif)
  3. ***```channel```*** 起到队列的作用。 多个 ***producer*** 产生的 ***```topic```*** 消息在每一个连接 ***```topic```*** 的 ***```channel```*** 上进行排队。
  4. 每个 ***```channel```*** 的消息都会进行排队，直到一个 ***```consumer```*** 把他们消费，如果此队列超出了内存限制，消息将会被写入到磁盘中。 ***```nsqd```*** 节点首先会向 ```nsqlookup``` 广播他们的位置信息，一旦它们注册成功， ***```consumer```*** 将会从 ```nsqlookup``` 服务器节点上发现所有包含事件 ***```topic```*** 的 ***```nsqd```*** 节点。
  5. 每个 ***```consumer```*** 向每个 ***```nsqd```*** 主机进行订阅操作，用于表明 ***```consumer```*** 已经准备好接受消息了。这里我们不需要一个完整的连通图，但我们必须要保证每个单独的 ***```nsqd```*** 实例拥有足够的消费者去消费它们的消息，否则 ***```channel```*** 会被队列堆着。


# NSQD  


* ##  ***```nsqd```*** 基本结构  
![SVG](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd/SVG&WG.png)
>  + 利用svc框架来启动服务, Run 时, 先后调用svc框架的 Init 和 Start 方法 ，然后开始不断监听退出的信号量, 最后调用 svc框架的Stop 方法来退出。
>  + svc框架的Start方法从本地文件读取数据初始化topic和channel，然后调用功能入口Main方法。Main方法利用waitGroup框架来启动4个服务线程，至此启动完毕。
>  + WaitGroup来自sync包，用于线程同步，单从字面意思理解，wait等待的意思，group组、团队的意思，WaitGroup就是等待一组服务执行完成后才会继续向下执行，涉及到WG个数的操作都使用原子操作来保证线程安全。  
  

* ##   ***```nsqd```*** 详细流程图  
![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd/nsqdflow.png)
   
* ##  ***```nsqd```*** 源码概述  
![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/nsqd/nsqd.png)
> +  ***```nsqd```*** 服务开启时启动 ***``` TCP```*** 服务供客户端连接，启动 ***```HTTP```*** 服务，提供 ***```HTTP API```***    
  
```go
//nsqd/nsqd.go:238
tcpServer := &tcpServer{ctx: ctx}
n.waitGroup.Wrap(func() {
	protocol.TCPServer(n.tcpListener, tcpServer, n.logf)
})
httpServer := newHTTPServer(ctx, false, n.getOpts().TLSRequired == TLSRequired)
n.waitGroup.Wrap(func() {
	http_api.Serve(n.httpListener, httpServer, "HTTP", n.logf)
})
```
> +  ***``` TCP```*** 接收到客户端的请求后，创建protocol实例并调用nsqd/tcp.go中IOLoop()方法  
 
```go
//nsqd/tcp.go:31
var prot protocol.Protocol
switch protocolMagic {
case "  V2":
	prot = &protocolV2{ctx: p.ctx}
default:
	protocol.SendFramedResponse(clientConn, frameTypeError, []byte("E_BAD_PROTOCOL"))
	clientConn.Close()
	p.ctx.nsqd.logf(LOG_ERROR, "client(%s) bad protocol magic '%s'",
		clientConn.RemoteAddr(), protocolMagic)
	return
}
err = prot.IOLoop(clientConn)
```

> +  protocol的IOLoop接收客户端的请求，根据命令的不同做相应处理。同时nsqd/protocol_v2.go中IOLoop会起一个goroutine运行messagePump()，该函数从该client订阅的channel中读取消息并发送给client(```consumer```)  

  
```go
//nsqd/protocol_v2.go:41
func (p *protocolV2) IOLoop(conn net.Conn) error {
	...

	clientID := atomic.AddInt64(&p.ctx.nsqd.clientIDSequence, 1)
	client := newClientV2(clientID, conn, p.ctx)

	// synchronize the startup of messagePump in order
	// to guarantee that it gets a chance to initialize
	// goroutine local state derived from client attributes
	// and avoid a potential race with IDENTIFY (where a client
	// could have changed or disabled said attributes)
	messagePumpStartedChan := make(chan bool)
	go p.messagePump(client, messagePumpStartedChan)
	<-messagePumpStartedChan

	...
}
```
```go
//nsqd/protocol_v2.go:200
func (p *protocolV2) messagePump(client *clientV2, startedChan chan bool) {
	...
	var memoryMsgChan chan *Message
	var backendMsgChan chan []byte
	var subChannel *Channel
	...

	select {
			...
		case b := <-backendMsgChan:
			...
			msg, err := decodeMessage(b)
			...
			subChannel.StartInFlightTimeout(msg, client.ID, msgTimeout)
			client.SendingMessage()
			err = p.SendMessage(client, msg)
			...
		case msg := <-memoryMsgChan:
			...
			subChannel.StartInFlightTimeout(msg, client.ID, msgTimeout)
			client.SendingMessage()
			err = p.SendMessage(client, msg)
			...
		}
	...
```

> +  然后我们看下memoryMsgChan，backendMsgChan是如何产生的。我们知道producer通过TCP或HTTP来发布消息。我们重点看下TCP时的处理过程。首先protocol的IOLoop会根据producer的不同请求做相应处理，Exec方法判断请求的参数，调用不同的方法。

```go
//nsqd/protocol_v2.go:165
func (p *protocolV2) Exec(client *clientV2, params [][]byte) ([]byte, error) {
	...
	switch {
	case bytes.Equal(params[0], []byte("FIN")):
		return p.FIN(client, params)
	case bytes.Equal(params[0], []byte("RDY")):
		return p.RDY(client, params)
	case bytes.Equal(params[0], []byte("REQ")):
		return p.REQ(client, params)
	case bytes.Equal(params[0], []byte("PUB")):
		return p.PUB(client, params)
	case bytes.Equal(params[0], []byte("MPUB")):
		return p.MPUB(client, params)
	case bytes.Equal(params[0], []byte("DPUB")):
		return p.DPUB(client, params)
	case bytes.Equal(params[0], []byte("NOP")):
		return p.NOP(client, params)
	case bytes.Equal(params[0], []byte("TOUCH")):
		return p.TOUCH(client, params)
	case bytes.Equal(params[0], []byte("SUB")):
		return p.SUB(client, params)
	case bytes.Equal(params[0], []byte("CLS")):
		return p.CLS(client, params)
	case bytes.Equal(params[0], []byte("AUTH")):
		return p.AUTH(client, params)
	}
	return nil, protocol.NewFatalClientErr(nil, "E_INVALID", fmt.Sprintf("invalid command %s", params[0]))
}
```

> +  我们重点看下“PUB”时的运行过程。调用了p.pub(client, params)。从TCP中读到messageBody，然后处理后调用topic.PutMessage(msg)发送给topic。
> +  topic.PutMessage（）首先对topic加一个锁，通过t.put(m)方法将消息m发送memoryMsgChan中，然后释放锁。如果memoryMsgChan满了，申请一个buff，把消息写到Backend，后期被backendMsgChan接收


```go
//nsqd/protocol_v2.go:757
func (p *protocolV2) PUB(client *clientV2, params [][]byte) ([]byte, error) {
	...
	topic := p.ctx.nsqd.GetTopic(topicName)
	msg := NewMessage(topic.GenerateID(), messageBody)
	err = topic.PutMessage(msg)
    ...
}

//nsqd/topic.go:197
func (t *Topic) put(m *Message) error {
	select {
	case t.memoryMsgChan <- m:
	default:
		b := bufferPoolGet()
		err := writeMessageToBackend(b, m, t.backend)
		bufferPoolPut(b)
		t.ctx.nsqd.SetHealth(err)
		if err != nil {
			t.ctx.nsqd.logf(LOG_ERROR,
				"TOPIC(%s) ERROR: failed to write message to backend - %s",
				t.name, err)
			return err
		}
	}
	return nil
}
```


> +  了解了消息的产生，现在看下消息的传递。在nsqd/topic.go中有一个NewTopic()。其中又调用了messagePump()，注意这个和上面IOLoop的messagePump()不一样。

```go
//nsqd/topic.go:44
func NewTopic(topicName string, ctx *context, deleteCallback func(*Topic)) *Topic {
	t := &Topic{
		name:              topicName,
		channelMap:        make(map[string]*Channel),
		memoryMsgChan:     make(chan *Message, ctx.nsqd.getOpts().MemQueueSize),
		exitChan:          make(chan int),
		channelUpdateChan: make(chan int),
		ctx:               ctx,
		pauseChan:         make(chan bool),
		deleteCallback:    deleteCallback,
		idFactory:         NewGUIDFactory(ctx.nsqd.getOpts().ID),
	}

	...

	t.waitGroup.Wrap(func() { t.messagePump() })

	t.ctx.nsqd.Notify(t)

	return t
}
```

> +  看下t.messagePump()。topic的messagePump函数会不断从memoryMsgChan/backend队列中读消息，并将消息每个复制一遍，发送给topic下的所有channel。

```go
//nsqd/topic.go:220
func (t *Topic) messagePump() {
	...
	for {
		select {
		case msg = <-memoryMsgChan:
		case buf = <-backendChan:
			msg, err = decodeMessage(buf)
			...
		case <-t.channelUpdateChan:
			...
		case pause := <-t.pauseChan:
			...
		case <-t.exitChan:
			goto exit
		}

		for i, channel := range chans {
			chanMsg := msg
			// copy the message because each channel
			// needs a unique instance but...
			// fastpath to avoid copy if its the first channel
			// (the topic already created the first copy)
			if i > 0 {
				chanMsg = NewMessage(msg.ID, msg.Body)
				chanMsg.Timestamp = msg.Timestamp
				chanMsg.deferred = msg.deferred
			}
			if chanMsg.deferred != 0 {
				channel.PutMessageDeferred(chanMsg, chanMsg.deferred)
				continue
			}
			err := channel.PutMessage(chanMsg)
			...
			}
		}
	}
	...
}
```

> +  channel的PutMessage方法和topic类似的，也是调用了put,首先写入memoryMsgChan，满了写入backend。
> + 最后由一开始介绍的protocol实例的messagePump方法从memoryMsgChan或backendMsgChan读取消息并通过p.SendMessage(client, msg)发送到客户端 ，消息写入client.Writer。

```go
//nsqd/channel.go:220
func (c *Channel) put(m *Message) error {
	select {
	case c.memoryMsgChan <- m:
	default:
		b := bufferPoolGet()
		err := writeMessageToBackend(b, m, c.backend)
		bufferPoolPut(b)
		c.ctx.nsqd.SetHealth(err)
		if err != nil {
			c.ctx.nsqd.logf(LOG_ERROR, "CHANNEL(%s): failed to write message to backend - %s",
				c.name, err)
			return err
		}
	}
	return nil
}

//nsqd/protocol_v2.go:258
		select {
       		...
		case b := <-backendMsgChan:
			...
			subChannel.StartInFlightTimeout(msg, client.ID, msgTimeout)
			client.SendingMessage()
			err = p.SendMessage(client, msg)
			...
		case msg := <-memoryMsgChan:
			...
			subChannel.StartInFlightTimeout(msg, client.ID, msgTimeout)
			client.SendingMessage()
			err = p.SendMessage(client, msg)
			...
        }
```


# NSQ使用

> 首先启动 ```nsdlookupd```

```shell
nsqlookupd
```

* 客户端通过查询 ```nsdlookupd``` 来发现指定topic的生产者，并且 ```nsqd``` 节点广播 ```topic``` 和通道 ```channel``` 信息
* 该服务运行后有两个端口：TCP 接口，```nsqd``` 用它来广播；HTTP 接口，客户端用它来发现和管理。
* 在生产环境中，为了高可用，最好部署三个nsqlookupd服务。

> 先创建 ```nsqd``` 的数据路径

```shell
mkdir /tmp/nsqdata1 /tmp/nsqdata2
```


> 运行两个测试的 ```nsqd``` 实例

```shell
nsqd --lookupd-tcp-address=127.0.0.1:4160 -broadcast-address=127.0.0.1 -tcp-address=127.0.0.1:4150 -http-address=0.0.0.0:4151 -data-path=/tmp/nsqdata1

nsqd --lookupd-tcp-address=127.0.0.1:4160 -broadcast-address=127.0.0.1 -tcp-address=127.0.0.1:4152 -http-address=0.0.0.0:4153 -data-path=/tmp/nsqdata2
```

*  ```nsqd``` 可以独立运行，不过通常它是由 ```nsdlookupd``` 实例所在集群配置的(它在这能声明 ```topics``` 和 ```channels``` ，以便大家能找到)
* 服务启动后有两个端口：一个给客户端(TCP)，另一个是 HTTP API。还能够开启HTTPS。
* 同一台服务器启动多个 ```nsqd``` ，要注意端口和数据路径必须不同，包括： ```–lookupd-tcp-address``` 、  ```-tcp-address``` 、 ```–data-path``` 。
* 删除```topic```、```channel```需要```HTTP API```调用。

> 启动 ```nsqadmin``` 前端Web监控

```shell
nsqadmin --lookupd-http-address=localhost:4161
```
*  ```nsqadmin``` 是一套 ```WEB UI``` ，用来汇集集群的实时统计，并执行不同的管理任务。
* 运行后，能够通过4171端口查看并管理 ```topic``` 和 ```channel``` 。
*  ```nsqadmin``` 通常只需要运行一个。

## 测试1
> * 2个Producer  1个Consumer
> * produce1() 发布publish "x","y" 到 topic "test"
> * produce2() 发布publish "z" 到 topic "test"
> * consumer1() 订阅subscribe  channel "sensor01"  of topic "test"

```go
package test

import (
        "log"
        "time"
        "testing"
        "strconv"

        "github.com/nsqio/go-nsq"
)

func TestNSQ1(t *testing.T) {
       NSQDsAddrs := []string{"127.0.0.1:4150", "127.0.0.1:4152"}
       go consumer1(NSQDsAddrs)
       go produce1()
       go produce2()
       time.Sleep(30 * time.Second)
}

func produce1() {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:4150"
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        if err := producer.Publish("test", []byte("x")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
        if err := producer.Publish("test", []byte("y")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
}

func produce2() {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:4152"
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        if err := producer.Publish("test", []byte("z")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
}

func consumer1(NSQDsAddrs []string) {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor01", cfg)
        if err != nil {
                log.Fatal(err)
        }
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C1")
                        return nil
                }))
        if err := consumer.ConnectToNSQDs(NSQDsAddrs); err != nil {
                log.Fatal(err, " C1")
        }
        <-consumer.StopChan
}
```

**测试结果**

![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/test/test1.png)
>  ```x,y,z``` 都被 ```consumer1``` 接收了。注意到接收时间， ```x,y``` 几乎同时被接收，它们都由 ```producer1``` 发布，而 ```z``` 由 ```producer2``` 发布，中间间隔10秒。测试了很多次都是10秒,偶尔是15秒或20秒。查看了ConnectToNSQDs()

```go
// ConnectToNSQDs takes multiple nsqd addresses to connect directly to.
//
// It is recommended to use ConnectToNSQLookupd so that topics are discovered
// automatically.  This method is useful when you want to connect to local instance.
```
> Consumer每隔x秒，向nsqlookud进行http轮询，用来更新自己的nsqd地址目录,当一个producer的channel一直没有数据时，则会轮询到下一个producer

## 测试2
> * 1个Producer  3个Consumer
> * produce3() 发布publish "x","y","z" 到 topic "test"
> * consumer1() 订阅subscribe  channel "sensor01"  of topic "test"
> * consumer2() 订阅subscribe  channel "sensor01"  of topic "test"
> * consumer3() 订阅subscribe  channel "sensor02"  of topic "test"

```go
package test

import (
        "log"
        "time"
        "testing"
        "strconv"

        "github.com/nsqio/go-nsq"
)

func TestNSQ2(t *testing.T) {
        NSQDsAddrs := []string{"127.0.0.1:4150"}
        go consumer1(NSQDsAddrs)
        go consumer2(NSQDsAddrs)
        go consumer3(NSQDsAddrs)
        go produce3()
        time.Sleep(5 * time.Second)
}

func produce3() {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:4150"
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        if err := producer.Publish("test", []byte("x")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
        if err := producer.Publish("test", []byte("y")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
        if err := producer.Publish("test", []byte("z")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
}

func consumer1(NSQDsAddrs []string) {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor01", cfg)
        if err != nil {
                log.Fatal(err)
        }
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C1")
                        return nil
                }))
        if err := consumer.ConnectToNSQDs(NSQDsAddrs); err != nil {
                log.Fatal(err, " C1")
        }
        <-consumer.StopChan
}

func consumer2(NSQDsAddrs []string) {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor01", cfg)
        if err != nil {
                log.Fatal(err)
        }
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C2")
                        return nil
                }))
        if err := consumer.ConnectToNSQDs(NSQDsAddrs); err != nil {
                log.Fatal(err, " C2")
        }
        <-consumer.StopChan
}

func consumer3(NSQDsAddrs []string) {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor02", cfg)
        if err != nil {
                log.Fatal(err)
        }
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C3")
                        return nil
                }))
        if err := consumer.ConnectToNSQDs(NSQDsAddrs); err != nil {
               log.Fatal(err, " C3")
        }
        <-consumer.StopChan
}
```

**测试结果**

![nsqd](https://github.com/VeniVidiViciVK/NSQ/raw/master/docs/test/test2.png)
> * ```consumer1``` 接收到了 ```y```
> * ```consumer2``` 接收到了 ```x,z```
> * ```consumer3``` 接收到了 ```x,y,z```

  * ***```channel```*** ***```sensor01```*** 中的消息被随机的分到了 ```consumer1``` 和 ```consumer2```
  * ```consumer3``` 单独占有 ***```channel```*** ***```sensor02```*** ，接收了其中的所有消息



# 参考文献

[[1]GoDoc of nsq [EB/OL]](https://godoc.org/github.com/bitly/nsq)  
[[2]NSQ v1.0.0-compat DESIGN [EB/OL]](https://nsq.io/overview/design.html)  
[[3]A Journey Into NSQ [EB/OL]](https://blog.gopheracademy.com/advent-2013/day-22-a-journey-into-nsq/)
