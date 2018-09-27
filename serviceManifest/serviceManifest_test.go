package serviceManifest_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
	. "github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest/mock"
)

var _ = Describe("ServiceManifest", func() {
	var mockParser *ParseData
	var mockFileIO *MockFileIO
	BeforeEach(func() {
		mockFileIO = NewMockFileIO()
		mockParser = &ParseData{
			Reader:  bytes.NewBufferString("Test"),
			Decoder: NewMockDecoder(),
			FileIO:  mockFileIO,
		}
	})

	It("A parser should fail with an non-existent/Invalid file", func() {
		mockFileIO.FileNotExist = true
		mockFileIO.FileCanOpen = true

		_, err := mockParser.CreateParser("blah")
		Expect(err).Should(HaveOccurred())

		// Set the File to exist but make it unopeneable.
		mockFileIO.FileNotExist = false
		mockFileIO.FileCanOpen = false
		_, err = mockParser.CreateParser("blah2")
		Expect(err).Should(HaveOccurred())

	})

	It("A parser should succeed with a working file", func() {
		mockFileIO.FileNotExist = false
		mockParser, err := mockParser.CreateParser("workingfile")

		bytesRead, err := ioutil.ReadAll(mockParser.Reader)
		outputString := string(bytesRead[:])

		Expect(err).ShouldNot(HaveOccurred())
		Expect(outputString).Should(Equal("Opened_workingfile"))
	})

	It("Parse should return a manifest struct file with a service name set to the buffer content", func() {
		manifest, err := mockParser.Parse([]string{}, map[string]string{})
		Expect(err).ShouldNot(HaveOccurred())
		Expect(manifest).ShouldNot(BeNil())
		Expect(manifest.Services[0].ServiceName).Should(Equal("Test"))

	})
})
