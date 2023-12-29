package gomega_stub_test

import (
	. "github.com/CameronHonis/gomega-stub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type SomeObject struct {
	FieldA string
}

func NewStubbedSomeObject() *Stubbed {

}

var _ = Describe("Stubbed", func() {
	var base *Stubbed
	var stubbedObject *SomeObject
	BeforeEach(func() {
		base = NewStubbed()
	})
	Describe("StubMethod", func() {
		When("the method name passed does not exist on the struct", func() {
			It("panics", func() {
				Expect(func() {
					StubMethod(&struct{}{}, "nonExistentMethod")
				}).To(Panic())
			})
		})
	})
})
