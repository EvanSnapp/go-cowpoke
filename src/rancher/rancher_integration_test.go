package rancher

/*
this is just a test harness to flesh out the implementation of the
Rancher integration points. This either needs to be removed or refactored
to somehow run reliably

TODO: will need to create a test suite to mock out calls to test various failures
*/
import (
	"testing"

	"fmt"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRancherIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rancher Integration Suite")
}

var _ = BeforeSuite(func() {
	godotenv.Load("../../.env")
})

var _ = Describe("Catalog", func() {

	It("should get catalog url", func() {
		url, err := GetTemplateURL("lk fork", "cowpoke", "config-improvements-77")

		Expect(err).To(BeNil())
		Expect(url).ToNot(BeNil())
		fmt.Println(url.String())
	})

	It("should get template version", func() {

	})
})
