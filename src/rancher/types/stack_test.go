package types

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStackType(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stack Type Suite")
}

var _ = Describe("IsUpgradeableTo Tests", func() {
	It("should be upgradeable", func() {
		s := Stack{
			ExternalID: "catalog://catalog:template:1",
		}

		Expect(s.IsUpgradableTo("catalog:template:2")).To(Equal(true))
	})

	It("should not be upgradeable if stack is not created from a catalog template", func() {
		s := Stack{
			ExternalID: "a-weird-id",
		}

		s2 := Stack{}

		Expect(s.IsUpgradableTo("catalog:template:2")).To(Equal(false))
		Expect(s2.IsUpgradableTo("catalog:template:2")).To(Equal(false))
	})

	It("should not be upgradeable if stack is not created from the the same catalog template", func() {

		ids := map[string]string{
			"catalog://catlog:template:1":   "catalog1:template:2",
			"catalog://catalog:template1:1": "catalog:template2:2",
		}

		for stackExID, templateID := range ids {
			s := Stack{
				ExternalID: stackExID,
			}

			Expect(s.IsUpgradableTo(templateID)).To(Equal(false))
		}

	})

	It("should not be upgradeable if stack is created from a higher template version", func() {
		s := Stack{
			ExternalID: "catalog://catalog:template:2",
		}

		Expect(s.IsUpgradableTo("catalog:template:1")).To(Equal(false))
	})

	It("should not be upgradeable if catalog template version id is invalid", func() {
		s := Stack{
			ExternalID: "catalog://catalog:template:1",
		}

		ids := []string{
			"",
			"some-weird-thing",
			"partially:bad-6",
		}

		for _, id := range ids {
			Expect(s.IsUpgradableTo(id)).To(Equal(false))
		}

	})

	It("should not be upgradeable if catalog template version is not a number", func() {
		s := Stack{
			ExternalID: "catalog://catalog:template:1",
		}

		Expect(s.IsUpgradableTo("catalog:template:taco")).To(Equal(false))
	})

	It("should not be upgradeable if stack external id version part is not a number", func() {
		s := Stack{
			ExternalID: "catalog://catalog:template:taco",
		}

		Expect(s.IsUpgradableTo("catalog:template:2")).To(Equal(false))
	})
})
