package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "server")
}

var _ = Describe("handlerDelayed", func() {
	Context("with a delay of 100 ms", func() {
		var handler http.HandlerFunc
		const timeout = time.Millisecond * 100
		var w *httptest.ResponseRecorder
		var r *http.Request

		BeforeEach(func() {
			handler = handlerDelayed(timeout)
			w = httptest.NewRecorder()
			r = httptest.NewRequest("GET", "/delay", nil)
		})

		It("completes by the specified timeout", func(done Done) {
			handler(w, r)
			close(done)
		}, timeout.Seconds() + 0.1)

		It("waits for the specified amount of time within 5 milliseconds", func(done Done) {
			now := time.Now()
			expectedDone := now.Add(timeout)
			handler(w, r)
			finishTime := time.Now()
			Expect(finishTime).To(BeTemporally("~", expectedDone, time.Millisecond*5))
			close(done)
		}, timeout.Seconds() + 0.5)
	})
})

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
