package serviceManifest_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServiceManifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceManifest Suite")
}
