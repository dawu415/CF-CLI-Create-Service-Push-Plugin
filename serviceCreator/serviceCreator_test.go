package serviceCreator_test

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceCreator"
	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceCreator/mock"
)

var _ = Describe("ServiceCreator", func() {
	var mockCFPlugin *MockCliConnection
	var mockServiceManifest *serviceManifest.ServiceManifest
	var serviceCreatorCmd *ServiceCreator
	BeforeEach(func() {
		mockCFPlugin = NewMockCliConnection()
		mockServiceManifest = &serviceManifest.ServiceManifest{}
		serviceCreatorCmd = &ServiceCreator{}
	})

	It("serviceCreator should still work without errors on an empty manifest", func() {
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("serviceCreator should fail with an invalid service type", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "thatservicethingytype",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)

		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())

	})

	It("serviceCreator should be able to create a brokered service with a blank type succesfully", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"create-service", "p-mysql", "standard", "MyService",
				"-t", "\"blah, cool\"", "-c", "{\"git\":\"www.git.com\"}"}))
	})

	It("serviceCreator should be able to create a brokered service, even with updateServices true, with a blank type succesfully", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  true,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"create-service", "p-mysql", "standard", "MyService",
				"-t", "\"blah, cool\"", "-c", "{\"git\":\"www.git.com\"}"}))
	})

	It("serviceCreator should fail on create-service if cf plugin wasn't able to query the services from CloudFoundry", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.SimulateErrorOnGetServices = true
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should fail on create-service if cf plugin wasn't able to query the services from CloudFoundry by a service name", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.SimulateErrorOnGetServiceByName = true
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator should fail on create-service if it had trouble talking to Cloud Foundry on a CliCommand", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.SimulateErrorOnCliCommand = true
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)

		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeTrue())
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator should be able to create a brokered service with a 'brokered' type succesfully", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "brokered",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"create-service", "p-mysql", "standard", "MyService",
				"-t", "\"blah, cool\"", "-c", "{\"git\":\"www.git.com\"}"}))
	})

	It("serviceCreator should not create the brokered service again if it already exists and we don't want to update it", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "brokered",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should be able to update brokered service if it already exists", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "brokered",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  true,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"update-service", "MyService",
				"-t", "\"blah, cool\"", "-c", "{\"git\":\"www.git.com\"}"}))
	})

	It("serviceCreator should not get stuck in progress loop if an error occurred", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:    serviceName,
			Type:           "brokered",
			Broker:         "p-mysql",
			PlanName:       "standard",
			URL:            "www.blah.com",
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
			Tags:           "blah, cool",
		}
		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "failed",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator can create credential user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName: serviceName,
			Type:        "credentials",
			Credentials: map[string]string{
				"host":  "www.david.com",
				"user":  "abcd1234",
				"psswd": "ooosupersecret",
			},
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-p",
				"{\"host\":\"www.david.com\",\"psswd\":\"ooosupersecret\",\"user\":\"abcd1234\"}"}))
	})

	It("serviceCreator can create credential user provided service with UpdateServices: true when service doesn't exist", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName: serviceName,
			Type:        "credentials",
			Credentials: map[string]string{
				"host":  "www.david.com",
				"user":  "abcd1234",
				"psswd": "ooosupersecret",
			},
			UpdateService:  true,
			JSONParameters: "{\"git\":\"www.git.com\"}",
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-p",
				"{\"host\":\"www.david.com\",\"psswd\":\"ooosupersecret\",\"user\":\"abcd1234\"}"}))
	})

	It("serviceCreator should fail on create user provided credential service if cf plugin wasn't able to query the services from CloudFoundry", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName: serviceName,
			Type:        "credentials",
			Credentials: map[string]string{
				"host":  "www.david.com",
				"user":  "abcd1234",
				"psswd": "ooosupersecret",
			},
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
		}
		mockCFPlugin.SimulateErrorOnGetServices = true
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should not create the user provided credential service again if it already exists and we don't want to update it", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName: serviceName,
			Type:        "credentials",
			Credentials: map[string]string{
				"host":  "www.david.com",
				"user":  "abcd1234",
				"psswd": "ooosupersecret",
			},
			UpdateService:  false,
			JSONParameters: "{\"git\":\"www.git.com\"}",
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should be able to update user provided credential service if it already exists", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName: serviceName,
			Type:        "credentials",
			Credentials: map[string]string{
				"host":  "www.david.com",
				"user":  "abcd1234",
				"psswd": "ooosupersecret",
			},
			UpdateService:  true,
			JSONParameters: "{\"git\":\"www.git.com\"}",
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"uups", "MyService",
				"-p",
				"{\"host\":\"www.david.com\",\"psswd\":\"ooosupersecret\",\"user\":\"abcd1234\"}"}))
	})

	It("serviceCreator can create log-drain user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "drain",
			URL:           "drain://www.drainme.com",
			UpdateService: false,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-l",
				"drain://www.drainme.com"}))
	})

	It("serviceCreator can create log-drain user provided service when it doesnt exist and UpdateService is true", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "drain",
			URL:           "drain://www.drainme.com",
			UpdateService: true,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-l",
				"drain://www.drainme.com"}))
	})

	It("serviceCreator should fail on create user provided log-drain service if cf plugin wasn't able to query the services from CloudFoundry", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "drain",
			URL:           "drain://www.drainme.com",
			UpdateService: false,
		}

		mockCFPlugin.SimulateErrorOnGetServices = true
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should not create the user provided log-drain service again if it already exists and we don't want to update it", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "drain",
			URL:           "drain://www.drainme.com",
			UpdateService: false,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should be able to update user provided log-drain service if it already exists", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "drain",
			URL:           "drain://www.drainme.com",
			UpdateService: true,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"uups", "MyService",
				"-l",
				"drain://www.drainme.com"}))
	})

	It("serviceCreator can create route user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "https://www.route.com",
			UpdateService: false,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-r",
				"https://www.route.com"}))
	})

	It("serviceCreator can create route user provided service when it doesnt exist and updateService is true", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "https://www.route.com",
			UpdateService: true,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"cups", "MyService",
				"-r",
				"https://www.route.com"}))
	})

	It("serviceCreator should fail on creation of route user provided service if cf couldn't query the services", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "https://www.route.com",
			UpdateService: false,
		}

		mockCFPlugin.SimulateErrorOnGetServices = true
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should not create the user provided route service again if it already exists and we don't want to update it", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "https://www.route.com",
			UpdateService: false,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should be able to update user provided route service if it already exists", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "https://www.route.com",
			UpdateService: true,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(mockCFPlugin.CommandOutput).Should(Equal(
			[]string{
				"uups", "MyService",
				"-r",
				"https://www.route.com"}))
	})

	It("serviceCreator should not be able to update user provided route service if the URL is http schema", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "http://www.route.com",
			UpdateService: true,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should not be able to update user provided route service if the URL is invalid", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "%gh&%ij",
			UpdateService: true,
		}

		// Add the mock service in and then try to create the serice
		mockCFPlugin.GetServicesModels = append(mockCFPlugin.GetServicesModels,
			plugin_models.GetServices_Model{
				Name: serviceName})

		mockCFPlugin.GetServiceExists = true
		mockCFPlugin.GetServiceModel = plugin_models.GetService_Model{
			Name: serviceName,
			LastOperation: plugin_models.GetService_LastOperation{
				State: "succeeded",
			},
		}
		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
		Expect(mockCFPlugin.CliCommandWasCalled).Should(BeFalse())
	})

	It("serviceCreator should fail on http route user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "http://www.route.com",
			UpdateService: false,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator should fail on no scheme route user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "www.route.com",
			UpdateService: false,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator should fail on an invalid route user provided service", func() {
		serviceName := "MyService"
		brokeredService := serviceManifest.Service{
			ServiceName:   serviceName,
			Type:          "route",
			URL:           "%gh&%ij",
			UpdateService: false,
		}

		(*mockServiceManifest).Services = append((*mockServiceManifest).Services, brokeredService)
		err := serviceCreatorCmd.CreateServices(mockServiceManifest, mockCFPlugin)
		Expect(err).Should(HaveOccurred())
	})

	It("serviceCreator's Progress reporter should function correctly", func() {
		var outputBufferString string
		mockLogFunction := func(format string, a ...interface{}) (n int, err error) {
			outputBufferString += format
			return 0, nil
		}

		progressReporter := NewProgressReporterWithLoggerOut(mockLogFunction)
		progressReporter.Step("")
		Expect(progressReporter).ShouldNot(BeNil())
		Expect(outputBufferString).Should(Equal("|\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("/\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("-\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("\\\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("|\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("/\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("-\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("\\\r"))
		outputBufferString = ""
		progressReporter.Step("")
		Expect(outputBufferString).Should(Equal("|\r"))
		outputBufferString = ""
		progressReporter.Step("Finished!")
		Expect(outputBufferString).Should(Equal("Finished!\n/\r"))
	})
})
