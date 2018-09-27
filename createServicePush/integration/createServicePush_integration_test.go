package createServicePush_integration_test

import (
	"os/exec"

	"code.cloudfoundry.org/cli/plugin/models"

	"code.cloudfoundry.org/cli/cf/util/testhelpers/rpcserver"
	"code.cloudfoundry.org/cli/cf/util/testhelpers/rpcserver/rpcserverfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// This is an integration test of all components of CreateServicePush.  It utilizes a fakeRPC from
// CFCli examples to stub in the CliCommand portion of the plugin.
const validPluginPath = "../../create-services-plugin.exe"

var _ = Describe("CreateServicePush", func() {
	var (
		rpcHandlers *rpcserverfakes.FakeHandlers
		ts          *rpcserver.TestServer
		err         error
	)

	BeforeEach(func() {
		rpcHandlers = new(rpcserverfakes.FakeHandlers)
		ts, err = rpcserver.NewTestRPCServer(rpcHandlers)
		Expect(err).NotTo(HaveOccurred())

		// Fake IsMinCliVersion to ensure Minimum CLI version is met
		rpcHandlers.IsMinCliVersionStub = func(_ string, result *bool) error {
			*result = true
			return nil
		}

		// Fake GetService so to ensure we succeed in our service creation and Progress Reporter doesnt get stuck
		rpcHandlers.GetServiceStub = func(serviceInstance string, retVal *plugin_models.GetService_Model) error {
			*retVal = plugin_models.GetService_Model{
				Guid: "{1234-6789}",
				Name: serviceInstance,
				LastOperation: plugin_models.GetService_LastOperation{
					State: "succeeded",
				},
			}
			return nil
		}

		//set rpc.CallCoreCommand to a successful call
		//rpc.CallCoreCommand is used in both cliConnection.CliCommand() and
		//cliConnection.CliWithoutTerminalOutput()
		rpcHandlers.CallCoreCommandStub = func(_ []string, retVal *bool) error {
			*retVal = true
			return nil
		}

		//set rpc.GetOutputAndReset to return empty string; this is used by CliCommand()/CliWithoutTerminalOutput()
		rpcHandlers.GetOutputAndResetStub = func(_ bool, retVal *[]string) error {
			*retVal = []string{"{}"}
			return nil
		}
	})

	JustBeforeEach(func() {
		err = ts.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ts.Stop()
	})

	It("create service should run and create services from a valid local services-manifest", func() {
		args := []string{ts.Port(), "create-service-push"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("Found Service Manifest File: services-manifest.yml"))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("my-configserver - will now be created as a brokered service."))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("Credentials-UPS - will now be created as a user provided credential service."))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("Route-UPS - will now be created as a user provided route service."))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("LogDrain-UPS - will now be created as a user provided log drain service."))

	})

	It("create service should fail on an invalid file", func() {
		args := []string{ts.Port(), "create-service-push", "--service-manifest", "nonexistent-file.yml"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(1))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("ERROR: The file nonexistent-file.yml was not found"))

	})

	It("create service should fail on an invalid service-manifest argument input", func() {
		args := []string{ts.Port(), "create-service-push", "--service-manifest"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(1))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("ERROR: --service-manifest is missing a manifest filename argument"))
	})

	It("create service should fail on an conflicting service-manifest argument input", func() {
		args := []string{ts.Port(), "create-service-push", "--service-manifest", "myfile", "--no-service-manifest"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(1))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("ERROR: --no-service-manifest cannot be used in conjunction with --service-manifest"))
	})

	It("create service should be able to defer unrecognized inputs to cf push", func() {
		args := []string{ts.Port(), "create-service-push", "-b", "hwc_buildpack", "-p", "push/some/path"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("Performing a CF Push with arguments [ -b hwc_buildpack -p push/some/path ]"))
	})

	It("create service should be able to defer unrecognized inputs to cf push", func() {
		args := []string{ts.Port(), "create-service-push", "-b", "hwc_buildpack", "-p", "push/some/path"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("Performing a CF Push with arguments [ -b hwc_buildpack -p push/some/path ]"))
	})

	It("create service should run and create services from a valid local services-manifest with no push", func() {
		args := []string{ts.Port(), "create-service-push", "--no-push"}
		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("--no-push applied"))
	})

	//!!NOTE: For this test to run and pass, you will need to set an enviroment variable, i.e., export CSPAPP_ENV=blah
	It("create service should run and create services from a valid local services-manifest that expects a local environment variable", func() {
		args := []string{ts.Port(), "create-service-push", "--service-manifest", "services-manifest-env-variable.yml", "--no-push", "--use-env-vars-prefixed-with", "CSPAPP"}

		session, _ := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
		session.Wait()

		println(string(session.Buffer().Contents()[:]))
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("sandbox-configserver - will now be created as a brokered service."))
		Expect(string(session.Buffer().Contents()[:])).Should(ContainSubstring("--no-push applied"))
	})

})
