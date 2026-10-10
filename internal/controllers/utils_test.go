package controllers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/internal/controllers"
)

type Src struct {
	name *string
	age  int
}

type Dst struct {
	Name string
	Age  uint
}

var _ = Describe("Utils", func() {
	It("can copy a struct", func() {
		s := &Src{new("Toto"), 4}
		d := &Dst{}
		CopyPrivateToPublic(s, d)
		Expect(d.Name).To(Equal("Toto"))
		Expect(d.Age).To(Equal(uint(0)))
	})
})
