package cspArguments_test

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/cspArguments"
)

var _ = Describe("CspArguments", func() {
	var cspArgs *CSPArguments
	var uninstallArg string
	BeforeEach(func() {
		cspArgs = NewCSPArguments()

		// This is set by the CF CLI Command
		uninstallArg = "CLI-MESSAGE-UNINSTALL"
	})

	It("Should fail without the create-service-push", func() {
		_, err := cspArgs.Process([]string{"--no-push", "blah"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should pass without any arguments", func() {
		_, err := cspArgs.Process([]string{"create-service-push"})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("Should pass with just the app name any arguments", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "myapp"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("myapp"))
	})

	It("Should pass with the create-service-push and have a normal service-manifest name", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--no-push", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.ServiceManifestFilename).Should(Equal("services-manifest.yml"))
		Expect(cspArgs.OtherCFArgs).ShouldNot(ContainElement("create-service-push"))
	})

	It("Should pass with the create-service-push uninstalling", func() {
		cspArgs, err := cspArgs.Process([]string{uninstallArg, "create-service-push", "--no-push", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.IsUninstallingPlugin).Should(BeTrue())
	})

	It("Should fail with invalid --service-manifest inputs", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--service-manifest", "--no-push", "blah"})
		Expect(err).Should(HaveOccurred())

		_, err = cspArgs.Process([]string{"create-service-push", "--service-manifest", "myFile", "--no-service-manifest", "blah"})
		Expect(err).Should(HaveOccurred())

		_, err = cspArgs.Process([]string{"create-service-push", "--service-manifest"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should pass with valid --service-manifest inputs", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--service-manifest", "myfile", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("Should pass with valid --no-service-manifest inputs", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--no-service-manifest", "myfile", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.DoNotCreateServices).Should(BeTrue())
	})

	It("Should fail with invalid --no-service-manifest inputs", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--no-service-manifest", "--service-manifest", "myfile", "blah"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should have a doNotPush flag if --no-push is set", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--no-push", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.DoNotPush).Should(BeTrue())
	})

	It("Should have the correct number of remaining arguments, if it is not recognised", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--no-push", "blah", "foo", "bar"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(cspArgs.OtherCFArgs)).Should(Equal(3))
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("blah"))
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("foo"))
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("bar"))
		Expect(cspArgs.OtherCFArgs).ShouldNot(ContainElement("--no-push"))
	})

	It("Should handle multiple inputs of --vars-file commands and should not pass this to CF push", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "myapp", "--vars-file", "someVar.yml", "--vars-file", "params.yml"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(cspArgs.StaticVariablesFilePaths)).Should(Equal(2))
		Expect(cspArgs.StaticVariablesFilePaths).Should(ContainElement("someVar.yml"))
		Expect(cspArgs.StaticVariablesFilePaths).Should(ContainElement("params.yml"))
		Expect(len(cspArgs.OtherCFArgs)).Should(Equal(1))
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("myapp"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).ShouldNot(ContainSubstring("--vars-file someVar.yml"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).ShouldNot(ContainSubstring("--vars-file params.yml"))
	})

	It("Should handle bad inputs of --vars-file commands", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--vars-file", "--vars-file", "params.yml"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle bad inputs of --vars-file commands where no filename is specified", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--vars-file"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle multiple inputs of --var commands and should not pass this to CF push", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--var", "explode=false", "--var", "IsGood=true"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(cspArgs.StaticVariables)).Should(Equal(2))
		Expect(cspArgs.StaticVariables).Should(HaveKeyWithValue("explode", "false"))
		Expect(cspArgs.StaticVariables).Should(HaveKeyWithValue("IsGood", "true"))
		Expect(len(cspArgs.OtherCFArgs)).Should(Equal(0))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).ShouldNot(ContainSubstring("--var explode=false"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).ShouldNot(ContainSubstring("--var IsGood=true"))
	})

	It("Should handle bad input where there's no key value pair ", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "--var"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle bad input where key value pair has a space before the = sign ", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--var", "preboom", "=true"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle bad input where key value pair has a space after the = sign ", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--var", "postboom=", "true"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle --push-as-subprocess ", func() {
		csp, err := cspArgs.Process([]string{"create-service-push", "myapp", "--push-as-subprocess"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(csp.PushAsSubProcess).To(BeTrue())
	})

	It("Should give error when --push-as-subprocess is combined with --no-push", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--push-as-subprocess", "--no-push"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should give error when --no-push is combined with --push-as-subprocess ", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--no-push", "--push-as-subprocess"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should pass --vars-file when --push-as-subprocess is used", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "myapp", "--vars-file", "someVar.yml", "--vars-file", "params.yml", "--push-as-subprocess"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(cspArgs.StaticVariablesFilePaths)).Should(Equal(2))
		Expect(cspArgs.StaticVariablesFilePaths).Should(ContainElement("someVar.yml"))
		Expect(cspArgs.StaticVariablesFilePaths).Should(ContainElement("params.yml"))
		Expect(len(cspArgs.OtherCFArgs)).Should(Equal(5))
		Expect(cspArgs.OtherCFArgs).Should(ContainElement("myapp"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).Should(ContainSubstring("--vars-file someVar.yml"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).Should(ContainSubstring("--vars-file params.yml"))
		Expect(cspArgs.PushAsSubProcess).To(BeTrue())
	})

	It("Should handle multiple inputs of --var commands and should pass this to CF push", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--var", "explode=false", "--var", "IsGood=true", "--push-as-subprocess"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(cspArgs.StaticVariables)).Should(Equal(2))

		Expect(cspArgs.StaticVariables).Should(HaveKeyWithValue("explode", "false"))
		Expect(cspArgs.StaticVariables).Should(HaveKeyWithValue("IsGood", "true"))
		Expect(len(cspArgs.OtherCFArgs)).Should(Equal(4))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).Should(ContainSubstring("--var explode=false"))
		Expect(strings.Join(cspArgs.OtherCFArgs, " ")).Should(ContainSubstring("--var IsGood=true"))
		Expect(cspArgs.PushAsSubProcess).To(BeTrue())
	})

	It("Should return usage text properly", func() {
		usageInstructions := cspArgs.GetUsage()
		println(usageInstructions)
		Expect(usageInstructions).ShouldNot(BeNil())
		Expect(usageInstructions).ShouldNot(BeEmpty())
	})

	It("Should return the flag descriptions", func() {
		argDescriptions := cspArgs.GetArgumentsDescription()
		Expect(argDescriptions).Should(BeAssignableToTypeOf(map[string]string{}))
	})

	It("Should handle invalid input of --use-env-vars-prefixed-with", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--use-env-vars-prefixed-with"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should handle input --use-env-vars-prefixed-with even though there were no environment variables prefixed", func() {
		_, err := cspArgs.Process([]string{"create-service-push", "myapp", "--use-env-vars-prefixed-with", "BLAHBLAHBLAH"})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("Should be able to handle input --use-env-vars-prefixed-with correctly", func() {

		os.Setenv("CSPENV_VariableA", "12345")
		os.Setenv("CSPENV_VariableB", "David")
		csp, err := cspArgs.Process([]string{"create-service-push", "myapp", "--use-env-vars-prefixed-with", "CSPENV"})
		Expect(err).ShouldNot(HaveOccurred())

		Expect(csp.StaticVariables).Should(HaveKeyWithValue("CSPENV_VariableA", "12345"))
		Expect(csp.StaticVariables).Should(HaveKeyWithValue("CSPENV_VariableB", "David"))
	})
})
