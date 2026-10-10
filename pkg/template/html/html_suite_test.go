package html_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHtml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Template/Html Suite")
}
