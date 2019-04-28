package mqb

//go:generate mockgen -source=pool.go -destination=mock_pool_test.go -package=mqb PoolInterface

import (
	"fmt"
	"reflect"
	"sync"
)

//poolInterface interface
type poolInterface interface {
	start() error
	stop()
	setRequestChan(chan interface{})
	getRequestChan() chan interface{}
}

type pool struct {
	sync.Mutex
	callback    reflect.Value
	args        []reflect.Value
	maxWorkers  int
	requestChan chan interface{}
	stopChan    chan int
}

func newPool(callback interface{}, maxWorkers int, args []interface{}) (p poolInterface) {
	pool := &pool{
		maxWorkers:  maxWorkers,
		callback:    reflect.ValueOf(callback),
		requestChan: make(chan interface{}),
		stopChan:    make(chan int, 1),
	}
	if len(args) != 0 {
		for _, arg := range args {
			pool.args = append(pool.args, reflect.ValueOf(arg))
		}
	}
	return pool
}

func (pool *pool) setRequestChan(requestChan chan interface{}) {
	pool.requestChan = requestChan
}

func (pool *pool) getRequestChan() chan interface{} {
	return pool.requestChan
}

func (pool *pool) start() (err error) {
	sema := make(chan int, pool.maxWorkers)
	wg := sync.WaitGroup{}
	go func() {
		for {
			select {
			case task := <-pool.requestChan:
				sema <- 1
				wg.Add(1)
				go func() {
					defer func() {
						wg.Done()
						<-sema
					}()
					func() {
						defer func() {
							if r := recover(); r != nil {
								Logger.Error(fmt.Sprintf("panic was recovered=[%v]", r))
							}
						}()
						e := pool.callback.Call(append([]reflect.Value{reflect.ValueOf(task)}, pool.args...))
						if e != nil && e[0].Interface() != nil {
							Logger.Error(e[0].Interface().(error).Error())
						}
					}()
				}()
			case <-pool.stopChan:
				wg.Wait()
				pool.stopChan <- 1
				return
			}
		}
	}()
	return
}

func (pool *pool) stop() {
	pool.stopChan <- 1
	<-pool.stopChan
}
