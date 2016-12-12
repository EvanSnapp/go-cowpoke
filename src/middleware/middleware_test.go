package middleware

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
)

func TestMIddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}

var router *gin.Engine

var _ = BeforeSuite(func() {
	gin.SetMode(gin.TestMode)
	router = gin.New()
	router.Use(Authenticate())

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})
})

var _ = Describe("Authentication", func() {

	BeforeEach(func() {
		os.Setenv("API_KEY", "")
	})

	It("should allow anonymous authentication if no api key is set", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(200))
	})

	It("should allow anonymous authentication if a bearer token is sent", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		r.Header.Set("bearer", "suh bruh")
		router.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(200))
	})

	It("should not authenticate user if header is not present", func() {
		os.Setenv("API_KEY", "muh-key")
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		r.Header.Set("bearer", "suh bruh")
		router.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(401))
	})

	It("should not authenticate if bearer token does not match", func() {
		os.Setenv("API_KEY", "muh-key")
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		r.Header.Set("bearer", "a bad key")
		router.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(401))
	})

	It("should authenticate", func() {
		apiKey := "muh-key"
		os.Setenv("API_KEY", apiKey)
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		r.Header.Set("bearer", apiKey)
		router.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(200))
	})
})
