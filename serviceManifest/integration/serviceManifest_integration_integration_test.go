package serviceManifest_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

var _ = Describe("ServiceManifest", func() {
	var realParser *ParseData
	BeforeEach(func() {
		realParser = NewParser()
	})

	It("A parser should be able to open a blank file and Parse it", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-blank.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})

		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest).ShouldNot(BeNil())
	})

	It("A parser should fail on a non-existent file", func() {
		p, err := realParser.CreateParser("./fixtures/somewhere-in-the-universe.yml")

		Expect(err).Should(HaveOccurred())
		Expect(p.Reader).Should(BeNil())
	})

	It("A parser be able to open a yml file but should fail on a invalid file", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-invalid.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		_, err = p.Parse([]string{}, map[string]string{})
		Expect(err).Should(HaveOccurred())
	})

	It("A parser be able to open a yml file with valid output", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(4))
	})

	It("A parser can open a valid yml broker service definition", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-broker.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(BeEmpty())
		Expect(manifest.Services[0].ServiceName).Should(Equal("my-database-service"))
		Expect(manifest.Services[0].Broker).Should(Equal("p-mysql"))
		Expect(manifest.Services[0].PlanName).Should(Equal("1gb"))
		Expect(manifest.Services[0].JSONParameters).Should(Equal("{\"RAM\": 4gb }"))
		Expect(manifest.Services[0].Tags).Should(Equal("test1, test2"))
		Expect(manifest.Services[0].UpdateService).Should(BeTrue())
	})

	It("A parser can open a valid yml user provided credential service definition", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-credential.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("credentials"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("CUPS"))
		Expect(manifest.Services[0].Credentials["host"]).Should(Equal("https://abc.mydatabase.com/abcd"))
		Expect(manifest.Services[0].Credentials["username"]).Should(Equal("david"))
		Expect(manifest.Services[0].Credentials["password"]).Should(Equal("12.23@123password"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())
	})

	It("A parser can open a valid yml user provided log drain service definition", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-logdrain.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("drain"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("LDUPS"))
		Expect(manifest.Services[0].URL).Should(Equal("syslog-tls://server.myapp.com:1020"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())
	})

	It("A parser can open a valid yml user provided route service definition", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-route.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("route"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("RUPS"))
		Expect(manifest.Services[0].URL).Should(Equal("https://www.google.com"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())
	})

	It("A parser can successfully evaluate a --vars variable", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-route-variable.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		vars := map[string]string{
			"environment": "sandbox",
			"endpoint":    "/apps/test",
		}

		manifest, err := p.Parser.Parse([]string{}, vars)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("route"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("sandbox-RUPS"))
		Expect(manifest.Services[0].URL).Should(Equal("https://www.google.com/apps/test"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())

	})

	It("A parser can successfully evaluate --vars variables even if the map had extra elements inside", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-route-variable.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		vars := map[string]string{
			"environment":  "sandbox",
			"endpoint":     "/apps/test",
			"extraelement": "Boo!",
		}

		manifest, err := p.Parser.Parse([]string{}, vars)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("route"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("sandbox-RUPS"))
		Expect(manifest.Services[0].URL).Should(Equal("https://www.google.com/apps/test"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())
	})

	It("A parser will error when trying to evaluate the yaml file that does not have all variables defined", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-route-variable.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		vars := map[string]string{
			"environment": "sandbox",
		}

		_, err = p.Parser.Parse([]string{}, vars)

		Expect(err).Should(HaveOccurred())
	})

	It("A parser can successfully evaluate a --vars-file variable", func() {
		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-route-variable.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		manifest, err := p.Parser.Parse([]string{"./fixtures/service-manifest-test-variables.yml"}, map[string]string{})

		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("route"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("sandbox-RUPS"))
		Expect(manifest.Services[0].URL).Should(Equal("https://www.google.com/apps/test"))
		Expect(manifest.Services[0].UpdateService).Should(BeFalse())
	})

	It("A parser can successfully evaluate both a --vars and --vars-file", func() {

		p, err := realParser.CreateParser("./fixtures/service-manifest-valid-with-variables.yml")

		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Reader).ShouldNot(BeNil())

		vars := map[string]string{
			"username": "david",
			"password": "qwerty9876",
		}

		manifest, err := p.Parser.Parse([]string{"./fixtures/service-manifest-test-variables.yml"}, vars)

		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest.Services).ShouldNot(BeNil())

		Expect(len(manifest.Services)).Should(Equal(1))

		Expect(manifest.Services[0].Type).Should(Equal("credentials"))
		Expect(manifest.Services[0].ServiceName).Should(Equal("Credentials-UPS"))
		Expect(manifest.Services[0].Credentials).Should(HaveKeyWithValue("username", "david"))
		Expect(manifest.Services[0].Credentials).Should(HaveKeyWithValue("password", "qwerty9876"))
		Expect(manifest.Services[0].Credentials).Should(HaveKeyWithValue("host", "https://sandbox.mydatabase.com/apps/test"))

	})
})
