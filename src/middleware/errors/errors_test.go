package errors

import (
	"net/http"
	"net/http/httptest"
	"rancher/types"
	"testing"

	"io"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var errorsRouter *gin.Engine
var aboutStatusCode = http.StatusBadRequest
var apiErrorStub = types.APIError{Status: 503, Message: "something bad"}
var apiErrorStub2 = types.APIError{Status: 505, Message: "some error"}

func TestPublicErrorMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public Error Middleware Suite")
}

var _ = BeforeSuite(func() {
	gin.SetMode(gin.TestMode)
	errorsRouter = gin.New()
	errorsRouter.Use(HandlePublicError())

	errorsRouter.GET("/api-error", func(c *gin.Context) {
		c.Error(apiErrorStub).SetType(gin.ErrorTypePublic)
	})

	errorsRouter.GET("/system-error", func(c *gin.Context) {
		err := io.ErrClosedPipe
		c.Error(err).SetType(gin.ErrorTypePublic)
	})

	errorsRouter.GET("/non-public-error", func(c *gin.Context) {
		c.Error(apiErrorStub)
	})

	errorsRouter.GET("/multiple-public-errors", func(c *gin.Context) {
		c.Error(apiErrorStub).SetType(gin.ErrorTypePublic)
		c.Error(apiErrorStub2).SetType(gin.ErrorTypePublic)
	})

	errorsRouter.GET("/error-with-abort", func(c *gin.Context) {
		c.Error(apiErrorStub).SetType(gin.ErrorTypePublic)
		c.AbortWithStatus(aboutStatusCode)
	})

})

var _ = Describe("Rancher API Errors", func() {
	It("should set response from api error", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/api-error", nil)
		errorsRouter.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(apiErrorStub.Status))
	})

	It("should set response from system error", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/system-error", nil)
		errorsRouter.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(http.StatusInternalServerError))
	})

	It("should not set repsonse for non public error", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/non-public-error", nil)
		errorsRouter.ServeHTTP(w, r)

		Expect(w.Code).NotTo(Equal(apiErrorStub.Status))
	})

	It("should set response for last public error", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/multiple-public-errors", nil)
		errorsRouter.ServeHTTP(w, r)

		Expect(w.Code).To(Equal(apiErrorStub2.Status))
	})

	It("should not set response if request has been aborted", func() {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/non-public-error", nil)
		errorsRouter.ServeHTTP(w, r)

		Expect(w.Code).NotTo(Equal(aboutStatusCode))
	})

})
