package guidTranslator

import (
	"bufio"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestGuidTranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guid Translator Suite")
}

var _ = Describe("guidTranslator", func() {

	Describe("guidTranslator", func() {

		Context("Find Instance Name", func() {
			It("Instance name should match", func() {
				file, err := os.Open("../test/service_instances.txt")
				if err != nil {
					fmt.Println(err)
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}
				result := FindInstanceName(nil, "8912303d-3cdf-476e-b864-47f008b5ba5e", output)
				Expect(result).To(Equal("\"demo-blueprint\""))
			})
		})

		Context("Find Service Id", func() {
			It("Service Id should match", func() {
				file, err := os.Open("../test/services.txt")
				if err != nil {
					fmt.Println("File Reading Error!")
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}
				result := FindServiceId(nil, "232fb15a-462b-43f1-b48a-fb0327d65fe6", output)
				Expect(result).To(Equal("\"24731fb8-7b84-4f57-914f-c3d55d793dd4\""))
			})
		})

		Context("Find Service Plan Id", func() {
			It("Service Plan Id should match", func() {
				file, err := os.Open("../test/service_plans.txt")
				if err != nil {
					fmt.Println("File Reading Error!")
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}
				result := FindServicePlanId(nil, "9c67ab74-66f1-4abf-a098-8dce06a02362", output)
				Expect(result).To(Equal("\"bc158c9a-7934-401e-94ab-057082a5073f\""))
			})
		})

		Context("Find Service Guid", func() {
			It("Service Guid should match", func() {
				file, err := os.Open("../test/service_plans.txt")
				if err != nil {
					fmt.Println("File Reading Error!")
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}
				result := FindServiceGUId(nil, "9c67ab74-66f1-4abf-a098-8dce06a02362", output)
				Expect(result).To(Equal("\"232fb15a-462b-43f1-b48a-fb0327d65fe6\","))
			})
		})

		Context("Find Instance Guid", func() {
			It("Instance Guid should match", func() {
				var instanceName string = "demo-blueprint"
				userSpaceGuid := "b0728cce-2eef-4a8b-ac57-b480f2c48461"

				file, err := os.Open("../test/service_instances.txt")
				if err != nil {
					fmt.Println("File Reading Error!")
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}

				result := FindInstanceGuid(nil, instanceName, output, userSpaceGuid)
				Expect(result).To(Equal("8912303d-3cdf-476e-b864-47f008b5ba5e"))
			})
		})

		Context("Find Service Plan Guid", func() {
			It("Service Plan Guid should match", func() {
				var instanceName string = "demo-blueprint"
				userSpaceGuid := "b0728cce-2eef-4a8b-ac57-b480f2c48461"

				file, err := os.Open("../test/service_instances.txt")
				if err != nil {
					fmt.Println("File Reading Error!")
				}

				defer file.Close()

				var output []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					output = append(output, scanner.Text())
				}

				result := FindServicePlanGuid(nil, instanceName, output, userSpaceGuid)
				Expect(result).To(Equal("\"9c67ab74-66f1-4abf-a098-8dce06a02362\","))
			})
		})

	})
})
