package createServicePush_integration_test

import (
	"testing"
	"time"

	"code.cloudfoundry.org/cli/cf/util/testhelpers/pluginbuilder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCreateServicePush(t *testing.T) {
	RegisterFailHandler(Fail)

	// Build the binary to be tested
	pluginbuilder.BuildTestBinary("../../", "create-services-plugin")

	RunSpecs(t, "CreateServicePush Integration Suite")
}

var _ = BeforeEach(func() {
	SetDefaultEventuallyTimeout(3 * time.Second)
})
