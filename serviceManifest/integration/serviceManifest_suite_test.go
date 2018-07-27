package serviceManifest_integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServiceManifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceManifest Integration Suite")
}
