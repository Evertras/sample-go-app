package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "handlerMessage")
}

var _ = Describe("handlerMessage", func() {
	Context("with a given message of nonzero length and a simple GET request", func() {
		var handler http.HandlerFunc
		const msg = "this is a message"
		var w *httptest.ResponseRecorder
		var r *http.Request

		BeforeEach(func() {
			handler = handlerMessage(msg)
			w = httptest.NewRecorder()
			r = httptest.NewRequest("GET", "/", nil)
			handler(w, r)
		})

		It("writes the given message to the body", func() {
			resultBody := string(w.Body.Bytes())

			Expect(resultBody).To(ContainSubstring(msg))
		})

		It("does not write anything other than the message to the body", func() {
			resultBody := string(w.Body.Bytes())

			Expect(resultBody).To(ContainSubstring(msg))
		})
	})
})
