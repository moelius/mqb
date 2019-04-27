package main

import (
	"context"
	"fmt"
	"github.com/moelius/mqb"
	"log"
	"sync"
	"time"
)

const workersCount = 100
const messagesCount = 1000

var onResponseCallback = func(task *mqb.AmqpTask, corrID string, results chan []byte) error {
	if corrID != task.Delivery.CorrelationId {
		return nil
	}
	results <- task.Delivery.Body
	return nil
}

func RPCall(ctx context.Context, server *mqb.Server, request *mqb.AmqpRequest, results chan []byte) ([]byte, error) {
	if err := server.Publish(request); err != nil {
		return nil, err
	}
	select {
	case r := <-results:
		return r, nil
	case <-ctx.Done():
		return nil, ctx.Err()

	}
}

func main() {
	server, err := mqb.NewServer(context.Background(), mqb.URL("amqp://"))
	if err != nil {
		log.Fatal(err)
	}

	server.WithLogger(mqb.NewLoggerLog(mqb.DebugLevel))

	defer func() {
		if err := server.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	producerOptions := &mqb.AmqpProducerOptions{Exchange: "rpc"}
	if err := server.NewProducer(producerOptions); err != nil {
		log.Fatal(err)
	}

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	allStart := time.Now()
	for workerNum := 1; workerNum < workersCount+1; workerNum++ {
		wg.Add(1)
		go func(w int) {
			replyTo := fmt.Sprintf("client-%d", time.Now().UnixNano())
			consumerOptions := &mqb.AmqpConsumerOptions{
				AutoAck:        true,
				ChannelOptions: &mqb.AmqpChannelOptions{Exchange: "rpc", ExchangeType: mqb.AmqpExchangeDirect},
				BindingOptions: &mqb.AmqpBindingOptions{Key: replyTo, Queue: &mqb.AmqpQueueOptions{Name: replyTo, AutoDelete: true, Exclusive: true}},
			}

			results := make(chan []byte)
			corrID := fmt.Sprintf("client-%d", w)
			_, err = server.NewConsumer(context.Background(), consumerOptions, onResponseCallback, mqb.CallbackMaxWorkers(1), mqb.CallbackArg(corrID), mqb.CallbackArg(results))
			if err != nil {
				log.Fatal(err)
			}

			for m := 1; m < messagesCount+1; m++ {
				start := time.Now()
				request := &mqb.AmqpRequest{Exchange: "rpc", Key: "request", CorrelationID: corrID, ReplyTo: replyTo, Body: []byte(fmt.Sprintf("%d_%d", w, 1))}
				ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
				res, err := RPCall(ctx, server, request, results)
				if err != nil {
					mqb.Logger.Warning(err.Error())
					continue
				}
				mqb.Logger.Info(fmt.Sprintf("time spent: %s, clientN: %d, messageN:%d, result: %s", time.Since(start), w, m, string(res)))
			}
			wg.Done()
		}(workerNum)
	}
	wg.Wait()

	allTimeSpent := time.Since(allStart)
	fmt.Println(fmt.Sprintf("all time spent: %s, rps: %fs", allTimeSpent, workersCount*messagesCount/allTimeSpent.Seconds()))
}
