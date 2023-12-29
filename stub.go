package gomega_stub

import (
	"fmt"
	"reflect"
	"sync"
)

type StubbedI interface {
	StubMethod(methodName string, fn interface{})
	IsStubbed(methodName string) bool

	RestoreBehavior(methodName string)

	MethodCalls(methodName string) [][]interface{}
	WasMethodCalledWith(methodName string, args ...interface{}) bool

	Call(methodName string, args ...interface{}) []interface{}
}

//	Stubbed is a struct that provides concrete implementations
//	for the methods of the StubbedI interface. This struct is
//	intended to be embedded in a stub struct "wrapper".
//
//	The final hierarchy should look like this:
//		StubWrapper -> Stubbed -> StubbedStruct;
//		[where each struct embeds the next]
//
//	I realize that creating and formatting a wrapper to implement
//	this struct may be a painful and delicate process, so I (plan on)
//	providing a generator that will create the wrapper and the stubbed
//	methods for you. The generator would be run as a separate step in
//	the build process.
//
//	* See (TODO - insert file name of example) for an example of usage.
//
//	* See (TODO - insert file name of generated file) for the generated file.
//
//	* See (TODO - insert file name of generator here) for the generator source code.
type Stubbed[SO any] struct {
	wrapper              interface{} // the struct that wraps the stubbed object
	stubbedObj           *SO         // the struct being stubbed
	callbackByMethodName map[string]interface{}
	callArgsByMethodName map[string][][]interface{}
	mu                   sync.Mutex // just in case tests run in parallel (is this overkill?)
}

func NewStubbed[SO any](wrapper interface{}, objToStub *SO) *Stubbed[SO] {
	return &Stubbed[SO]{
		wrapper:              wrapper,
		stubbedObj:           objToStub,
		callbackByMethodName: make(map[string]interface{}),
		callArgsByMethodName: make(map[string][][]interface{}),
	}
}

//	NOTE: Thoroughly assert on the typing of the passed fn.
//	This frees us from having to do any type assertions
//	at method invocation time.
func (s *Stubbed[SO]) StubMethod(methodName string, fn interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ValidateStubSignature(s.stubbedObj, methodName, fn)
	s.callbackByMethodName[methodName] = fn
}

func (s *Stubbed[SO]) IsStubbed(methodName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.callbackByMethodName[methodName]
	return ok
}

func (s *Stubbed[SO]) RestoreBehavior(methodName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.callbackByMethodName, methodName)
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

//func (s *Stubbed[SO]) RecordMethodCall(methodName string, args ...interface{}) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	if s.callArgsByMethodName[methodName] == nil {
//		s.callArgsByMethodName[methodName] = make([][]interface{}, 0)
//	}
//	s.callArgsByMethodName[methodName] = append(s.callArgsByMethodName[methodName], args)
//}

func (s *Stubbed[SO]) Call(methodName string, args ...interface{}) []interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return make([]interface{}, 0)
}

func ValidateStubSignature(stubbedObject interface{}, methodName string, fn interface{}) {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		panic("fn must be a function")
	}
	soVal := reflect.ValueOf(stubbedObject)
	soType := soVal.Type()
	soValMethod, foundMethod := soType.MethodByName(methodName)
	if !foundMethod {
		panic(fmt.Sprintf("methodName (%s) must be a method of stubObject", methodName))
	}

	// assert i/o count matches
	fnType := fnVal.Type()
	soFuncType := soValMethod.Func.Type()
	fnNumIn := fnType.NumIn()
	soFuncNumIn := soFuncType.NumIn()
	if soFuncNumIn != fnNumIn {
		//	NOTE: this compares fn to the "under the hood" GENERATED function based upon the method signature
		//	this adds the receiver as the first argument
		panic(fmt.Sprintf("fn must have the same arg count as %s.%s's func signature\nDid you forget to include the receiver arg?", soType.Name(), methodName))
	}
	fnNumOut := fnType.NumOut()
	soFuncNumOut := soFuncType.NumOut()
	if soFuncNumOut != fnNumOut {
		panic(fmt.Sprintf("fn must have the same return count as %s.%s's func signature", soType.Name(), methodName))
	}

	// assert each i/o type match
	for i := 0; i < fnNumIn; i++ {
		if soFuncType.In(i) != fnType.In(i) {
			panic(fmt.Sprintf("fn return #%d must have the same type as %s.%s's func signature", i, soType.Name(), methodName))
		}
	}
	for i := 0; i < fnNumOut; i++ {
		if soFuncType.Out(i) != fnType.Out(i) {
			panic(fmt.Sprintf("fn return #%d must have the same type as %s.%s's func signature", i, soType.Name(), methodName))
		}
	}
}
