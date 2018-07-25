package serviceCreator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServiceCreator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceCreator Suite")
}
