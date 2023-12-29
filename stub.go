package gomega_stub

import (
	"reflect"
	"sync"
)

type StubbableI interface {
	StubMethod(methodName string, fn interface{})
	MethodCalls(methodName string) [][]interface{}
	WasMethodCalledWith(methodName string, args ...interface{}) bool
	RecordMethodCall(methodName string, args ...interface{})
}

type Stubbed[SO any] struct {
	stubbedObj           SO
	callbackByMethodName map[string]interface{}
	callArgsByMethodName map[string][][]interface{}
	mu                   sync.Mutex // just in case tests run in parallel (is this overkill?)
}

func NewStubbed[SO any](objToStub SO) *Stubbed[SO] {
	return &Stubbed[SO]{
		stubbedObj:           objToStub,
		callbackByMethodName: make(map[string]interface{}),
		callArgsByMethodName: make(map[string][][]interface{}),
	}
}

func (s *Stubbed[SO]) StubMethod(methodName string, fn interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("fn must be a function")
	}
	s.callbackByMethodName[methodName] = fn
}

func (s *Stubbed[SO]) MethodCalls(methodName string) [][]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.callArgsByMethodName[methodName] == nil {
		return make([][]interface{}, 0)
	}
	return s.callArgsByMethodName[methodName]
}

func (s *Stubbed[SO]) WasMethodCalledWith(methodName string, args ...interface{}) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, callArgs := range s.MethodCalls(methodName) {
		if reflect.DeepEqual(callArgs, args) {
			return true
		}
	}
	return false
}

func (s *Stubbed[SO]) RecordMethodCall(methodName string, args ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.callArgsByMethodName[methodName] == nil {
		s.callArgsByMethodName[methodName] = make([][]interface{}, 0)
	}
	s.callArgsByMethodName[methodName] = append(s.callArgsByMethodName[methodName], args)
}
