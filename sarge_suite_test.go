package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSarge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sarge Suite")
}
