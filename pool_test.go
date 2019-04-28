package mqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"reflect"
	"testing"
)

type PoolTestSuite struct {
	suite.Suite
	Pool *pool
}

func (suite *PoolTestSuite) SetupTest() {
	suite.Pool = &pool{
		maxWorkers:  1,
		requestChan: make(chan interface{}),
		stopChan:    make(chan int),
	}
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(PoolTestSuite))
}

func (suite *PoolTestSuite) TestNewPool() {
	ass := assert.New(suite.T())
	args := []interface{}{1}
	actualPool := newPool(func(some int) error { return nil }, 100, args)
	ass.IsType(&pool{}, actualPool)
	ass.NotNil(actualPool.getRequestChan())
}

func (suite *PoolTestSuite) TestSetRequestChan() {
	wantChan := make(chan interface{})
	suite.Pool.setRequestChan(wantChan)
	assert.Equal(suite.T(), wantChan, suite.Pool.requestChan)
}

func (suite *PoolTestSuite) TestGetRequestChan() {
	assert.Equal(suite.T(), suite.Pool.requestChan, suite.Pool.getRequestChan())
}

func (suite *PoolTestSuite) TestStart() {
	ch := make(chan bool)
	go func() {
		suite.Pool.stopChan <- 1
		ch <- true
	}()
	err := suite.Pool.start()
	assert.Nil(suite.T(), err)
	<-ch
}

func (suite *PoolTestSuite) TestCallbackCommunication() {
	ass := assert.New(suite.T())
	wantArg := "test val"
	testCallback := func(val string) error {
		//manually stops pool at the end of callback
		defer func() {
			suite.Pool.stopChan <- 1
		}()
		ass.Equal(val, wantArg)
		return nil
	}
	suite.Pool.callback = reflect.ValueOf(testCallback)

	ass.Nil(suite.Pool.start())
	suite.Pool.requestChan <- wantArg
	<-suite.Pool.stopChan
}

func (suite *PoolTestSuite) TestCallbackCommunicationHandlePanic() {
	ass := assert.New(suite.T())
	wantArg := "test val"
	testCallback := func(val string) error {
		//manually stops pool at the end of callback
		defer func() {
			suite.Pool.stopChan <- 1
		}()
		log.Panic("test panic")
		return nil
	}
	suite.Pool.callback = reflect.ValueOf(testCallback)

	ass.Nil(suite.Pool.start())
	suite.Pool.requestChan <- wantArg
	<-suite.Pool.stopChan
}

func (suite *PoolTestSuite) TestCallbackCommunicationHandleError() {
	ass := assert.New(suite.T())
	wantArg := "test val"
	testCallback := func(val string) error {
		//manually stops pool at the end of callback
		defer func() {
			suite.Pool.stopChan <- 1
		}()
		return fmt.Errorf("test error")
	}
	suite.Pool.callback = reflect.ValueOf(testCallback)

	ass.Nil(suite.Pool.start())
	suite.Pool.requestChan <- wantArg
	<-suite.Pool.stopChan
}

func (suite *PoolTestSuite) TestStop() {
	_ = suite.Pool.start()
	suite.Pool.stop()
}
