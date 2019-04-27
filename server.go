package mqb

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var Logger LoggerInterface

//CallbackOption callback option func
type CallbackOption func(options *CallbackOptions)

//CallbackOptions struct
type CallbackOptions struct {
	MaxWorkers int
	Args       []interface{}
}

//CallbackMaxWorkers callback max workers option
func CallbackMaxWorkers(max int) CallbackOption {
	return func(o *CallbackOptions) { o.MaxWorkers = max }
}

//CallbackArg callback static argument, this value must be in callback func
func CallbackArg(arg interface{}) CallbackOption {
	return func(o *CallbackOptions) {
		o.Args = append(o.Args, arg)
	}
}

var defaultCallbackOptions = CallbackOptions{MaxWorkers: 1}

//Server mqb server struct
type Server struct {
	broker  BrokerInterface
	osSig   chan os.Signal
	ctx     context.Context
	started bool
}

//NewServer make new mqb server
func NewServer(ctx context.Context, brokerOptions ...BrokerOption) (server *Server, err error) {
	opts := defaultBrokerOptions
	for _, o := range brokerOptions {
		o(&opts)
	}
	server = new(Server)

	server.WithLogger(NewLoggerLog(InfoLevel))

	server.ctx = ctx
	if strings.HasPrefix(opts.url, "amqp://") {
		server.broker = newAmqpBroker()
	} else {
		return nil, fmt.Errorf("unknown transport url: [%s]", opts.url)
	}
	if err = configure(server.broker, opts); err != nil {
		return
	}
	server.osSig = make(chan os.Signal, 1)
	return
}

//WithLogger set up custom logger interface. logger must implement LoggerInterface
func (server *Server) WithLogger(logger LoggerInterface) {
	Logger = logger
}

//NewConsumer make new amqpConsumer
func (server *Server) NewConsumer(ctx context.Context, consumerOptions interface{}, callback interface{}, callbackOptions ...CallbackOption) (consumer ConsumerInterface, err error) {
	opts := defaultCallbackOptions
	for _, o := range callbackOptions {
		o(&opts)
	}
	broker := server.getBroker()
	consumer, err = newConsumer(ctx, broker, consumerOptions, newPool(callback, opts.MaxWorkers, opts.Args))
	if err != nil || !server.started {
		return
	}
	if err = consumer.getPool().start(); err != nil {
		return
	}
	if err = startConsumer(server.getBroker(), consumer.getName()); err != nil {
		return nil, err
	}
	return
}

//NewProducer make new producer
func (server *Server) NewProducer(producerOptions interface{}) (err error) {
	return newProducer(server.getBroker(), producerOptions)
}

//ListenAndServe serve and listen for os signals to stop server
func (server *Server) ListenAndServe() (err error) {
	if err = server.Serve(); err != nil {
		return
	}
	server.signalHandler()
	return
}

//Serve starts server
func (server *Server) Serve() error {
	if err := serve(server.broker); err != nil {
		return err
	}
	Logger.Info("server started")
	server.started = true
	return nil
}

//Stop stop server
func (server *Server) Stop() error {
	err := stop(server.broker)
	if err == nil {
		server.started = false
	}
	return err
}

//Publish publish request with broker
func (server *Server) Publish(requestInterface interface{}) error {
	return publish(server.getBroker(), requestInterface)
}

//GetSignal returns os signal for graceful server shutdown
func (server *Server) GetSignal() chan os.Signal {
	signal.Notify(server.osSig, syscall.SIGINT, syscall.SIGTERM)
	return server.osSig
}

func (server *Server) signalHandler() {
	sig := server.GetSignal()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-sig:
				if err := server.Stop(); err != nil {
					Logger.Error(fmt.Sprintf("server was not gracefully stopped with error=[%s]", err))
					wg.Done()
					return
				}
				Logger.Info("server was gracefully stopped")
				wg.Done()
				return
			case err := <-server.broker.getNotify():
				if err != nil {
					Logger.Error(fmt.Sprintf("server down with transport error=[%s]", err))
				}
				wg.Done()
				return
			case <-server.ctx.Done():
				if err := server.Stop(); err != nil {
					Logger.Error(fmt.Sprintf("server was not gracefully stopped with error=[%s]", err))
					wg.Done()
					return
				}
				Logger.Info(fmt.Sprintf("server was gracefully stopped with [%s]", server.ctx.Err()))
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

func (server *Server) getBroker() BrokerInterface {
	return server.broker
}
