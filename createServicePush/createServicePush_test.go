package createServicePush_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/createServicePush"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/createServicePush/mock"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceCreator/mock"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

var _ = Describe("CreateServicePush", func() {
	var mockCFPlugin *serviceCreatorMock.MockCliConnection
	var mockServiceManifest *serviceManifest.ServiceManifest
	var mockCreateServiceInterfaces *createService_mock.MockCreateService
	var mockCSP *CreateServicePush
	var mockExitHandler *createService_mock.MockExitHandler
	BeforeEach(func() {
		mockCFPlugin = serviceCreatorMock.NewMockCliConnection()
		mockServiceManifest = &serviceManifest.ServiceManifest{}
		mockCreateServiceInterfaces = createService_mock.NewMockCreateService()
		mockExitHandler = createService_mock.NewMockExitHandler()
		// mockCreateService will hold the implementation for all
		// interfaces within CreateServicePush
		mockCSP = &CreateServicePush{
			Parser:         mockCreateServiceInterfaces,
			ArgProcessor:   mockCreateServiceInterfaces,
			ServiceCreator: mockCreateServiceInterfaces,
			Exit:           mockExitHandler,
		}
	})

	It("create service should fail if ArgProcessor Failed", func() {
		mockCreateServiceInterfaces.ArgumentHasError = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeTrue())
	})

	It("create service should fail if CreateParser Failed", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeTrue())
	})

	It("create service should fail if Parse Failed", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeTrue())
	})

	It("create service should fail if CreateService Failed", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeTrue())
	})

	It("create service should fail if CliCommand Failed", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false
		mockCFPlugin.SimulateErrorOnCliCommand = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeTrue())
	})

	It("create service should succeed if there were no problems", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false
		mockCFPlugin.SimulateErrorOnCliCommand = false
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeFalse())
	})

	It("create service should succeed if there were no problems with no services created if DoNotCreateServices was true", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false
		mockCFPlugin.SimulateErrorOnCliCommand = false

		mockCreateServiceInterfaces.DoNotCreateServices = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeFalse())
		Expect(mockCreateServiceInterfaces.ServicesCreated).Should(BeFalse())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeTrue())
	})

	It("create service should succeed if there were no problems with no services created if DoNotCreateServices was true and DoNotPush was true", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false
		mockCFPlugin.SimulateErrorOnCliCommand = false

		mockCreateServiceInterfaces.DoNotCreateServices = true
		mockCreateServiceInterfaces.DoNotPush = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeFalse())
		Expect(mockCreateServiceInterfaces.ServicesCreated).Should(BeFalse())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("create service is uninstalling", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false

		mockCreateServiceInterfaces.PlugIsUninstalling = true

		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeFalse())
		Expect(mockExitHandler.Exit0WasCalled).Should(BeTrue())
	})

	It("create service should succeed if there were no problems with no services created if DoNotCreateServices was false and DoNotPush was true", func() {
		mockCreateServiceInterfaces.ArgumentHasError = false
		mockCreateServiceInterfaces.CreateParserHasError = false
		mockCreateServiceInterfaces.ParseHasError = false
		mockCreateServiceInterfaces.CreateServiceHasError = false
		mockCFPlugin.SimulateErrorOnCliCommand = false

		mockCreateServiceInterfaces.DoNotCreateServices = false
		mockCreateServiceInterfaces.DoNotPush = true
		mockCSP.Run(mockCFPlugin, []string{})
		Expect(mockExitHandler.Exit1WasCalled).Should(BeFalse())
		Expect(mockCreateServiceInterfaces.ServicesCreated).Should(BeTrue())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})
})
