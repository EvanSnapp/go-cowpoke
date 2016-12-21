package configuration

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var allEnvVars = append(RequiredEnvVars, OptionalEnvVars...)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}

var _ = Describe("GetSlackChannels", func() {

	BeforeEach(func() {
		os.Setenv("SLACK_CHANNELS", "")
	})

	It("should handle empty string", func() {
		os.Setenv("SLACK_CHANNELS", ",")
		Expect(GetSlackChannels()).Should(BeNil())
	})

	It("should handle trailing comma string", func() {
		os.Setenv("SLACK_CHANNELS", "abc,,")
		Expect(GetSlackChannels()).To(Equal([]string{"abc"}))
	})

	It("should handle leading comma string", func() {
		os.Setenv("SLACK_CHANNELS", ",abc")
		Expect(GetSlackChannels()).To(Equal([]string{"abc"}))
	})

	It("should handle multiple elements string", func() {
		os.Setenv("SLACK_CHANNELS", "abc,efg")
		Expect(GetSlackChannels()).To(Equal([]string{"abc", "efg"}))
	})

	It("should handle single element string", func() {
		os.Setenv("SLACK_CHANNELS", "abc")
		Expect(GetSlackChannels()).To(Equal([]string{"abc"}))
	})

	It("should handle an all whitespace entry", func() {
		os.Setenv("SLACK_CHANNELS", "abc, ,efg")
		Expect(GetSlackChannels()).To(Equal([]string{"abc", "efg"}))
	})

	It("should trim whitespace", func() {
		os.Setenv("SLACK_CHANNELS", "abc , , efg, hij ")
		Expect(GetSlackChannels()).To(Equal([]string{"abc", "efg", "hij"}))
	})

})
