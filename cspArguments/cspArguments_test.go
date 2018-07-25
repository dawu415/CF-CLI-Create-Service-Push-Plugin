package cspArguments_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/cspArguments"
)

var _ = Describe("CspArguments", func() {
	var cspArgs *CSPArguments
	BeforeEach(func() {
		cspArgs = NewCSPArguments()
	})

	It("Should fail without the create-service-push", func() {
		_, err := cspArgs.Process([]string{"--no-push", "blah"})
		Expect(err).Should(HaveOccurred())
	})

	It("Should pass with the create-service-push and have a normal service-manifest name", func() {
		cspArgs, err := cspArgs.Process([]string{"create-service-push", "--no-push", "blah"})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(cspArgs.ServiceManifestFilename).Should(Equal("service-manifest.yml"))
		Expect(cspArgs.OtherCFArgs).ShouldNot(ContainElement("create-service-push"))
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

})
