package test

import (
        "log"
        "time"
        "testing"

        "github.com/nsqio/go-nsq"
        "strconv"
)

// 2个Producer  1个Consumer
// produce1() 发布publish "x","y" 到 topic "test"
// produce2() 发布publish "z" 到 topic "test"
// consumer1() 订阅subscribe  channel "sensor01"  of topic "test"
func TestNSQ(t *testing.T) {
        go consumer1()
        go produce1()
        go produce2()
        time.Sleep(20 * time.Second)
}

// 生产者1
func produce1() {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:4150"
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        // 发布消息

        if err := producer.Publish("test", []byte("x")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
        if err := producer.Publish("test", []byte("y")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }
}

// 生产者2
func produce2() {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:4152"
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        // 发布消息

        if err := producer.Publish("test", []byte("z")); err != nil {
                log.Fatal("publish error: " + err.Error())
        }

}

// 消费者
func consumer1() {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor01", cfg)
        if err != nil {
                log.Fatal(err)
        }
        // 设置消息处理函数
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C1")
                        return nil
                }))
        // 连接到单例nsqd
        //if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
        //        log.Fatal(err, " C1")
        //}

        // 连接到多个nsqd
        if err := consumer.ConnectToNSQDs([]string{"127.0.0.1:4150","127.0.0.1:4152"}); err != nil {
                log.Fatal(err, " C1")
        }
        <-consumer.StopChan
}

func consumer2() {
        cfg := nsq.NewConfig()
        consumer, err := nsq.NewConsumer("test", "sensor02", cfg)
        if err != nil {
                log.Fatal(err)
        }
        // 设置消息处理函数
        consumer.AddHandler(nsq.HandlerFunc(
                func(message *nsq.Message) error {
                        log.Println(string(message.Body) + " C2")
                        return nil
                }))
        // 连接到单例nsqd
        //if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
        //        log.Fatal(err, " C2")
        //}
        // 连接到多个nsqd
        if err := consumer.ConnectToNSQDs([]string{"127.0.0.1:4150", "127.0.0.1:4152"}); err != nil {
                log.Fatal(err, " C2")
        }
        <-consumer.StopChan
}

func producerNo(prodno int) {
        cfg := nsq.NewConfig()
        nsqdAddr := "127.0.0.1:" + strconv.Itoa(4150+prodno*2)
        producer, err := nsq.NewProducer(nsqdAddr, cfg)
        if err != nil {
                log.Fatal(err)
        }
        // 发布消息
        count := 0
        for {
                if err := producer.Publish("test", []byte("test message "+strconv.Itoa(prodno))); err != nil {
                        log.Fatal("publish error: " + err.Error())
                }
                time.Sleep(time.Duration(2/(1+prodno)) * time.Second)
                count++
                if prodno == 0 && count > 50 {
                        break
                }
        }
}
