package main

import (
	"context"
	"fmt"
	"github.com/moelius/mqb"
	"log"
	"strconv"
	"strings"
	"time"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func RPCallback(task *mqb.AmqpTask) error {
	//Do some job
	args := strings.Split(string(task.Delivery.Body), "_")
	n, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	result := fmt.Sprintf("response from client_%s: %d", args[0], fib(n))
	//Define response request for client
	response := &mqb.AmqpRequest{
		Exchange:      "rpc",
		Key:           task.Delivery.ReplyTo,
		CorrelationID: task.Delivery.CorrelationId,
		Body:          []byte(result),
	}
	//publish response to client
	return task.Broker.Publish(response)
}

func main() {
	server, err := mqb.NewServer(context.Background(), mqb.URL("amqp://"), mqb.MaxReconnects(3), mqb.ReconnectTimeout(3*time.Second))
	if err!=nil{
		log.Fatal(err)
	}
	server.WithLogger(mqb.NewLoggerLog(mqb.InfoLevel))
	consumerOptions := &mqb.AmqpConsumerOptions{
		AutoAck:        true,
		ChannelOptions: &mqb.AmqpChannelOptions{Exchange: "rpc", ExchangeType: mqb.AmqpExchangeDirect},
		BindingOptions: &mqb.AmqpBindingOptions{Key: "request", Queue: &mqb.AmqpQueueOptions{Name: "rpc.request"}},
		QosOptions:     &mqb.AmqpQosOptions{PrefetchCount: 500},
	}

	_, err = server.NewConsumer(context.Background(), consumerOptions, RPCallback, mqb.CallbackMaxWorkers(100))
	if err != nil {
		log.Fatal(err)
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
