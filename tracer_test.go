package tracer_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urjitbhatia/tracer"
)

var _ = Describe("Traces http request timing", func() {
	var ts *httptest.Server
	BeforeSuite(func() {
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, client")
		}))

	})

	AfterSuite(func() {
		ts.Close()
	})

	It("Handles requests", func() {
		tests := []struct {
			method string
			url    string
			body   io.Reader
		}{
			{http.MethodGet, ts.URL, nil},
			{http.MethodGet, "https://golang.org", nil},
			{http.MethodPost, ts.URL, bufio.NewReader(bytes.NewBufferString("foo"))},
			{http.MethodHead, ts.URL, nil},
		}

		for _, t := range tests {
			req, err := http.NewRequest(t.method, t.url, t.body)
			Expect(err).To(BeNil())
			req, tr := tracer.AsTraceableReq(req)

			client := http.Client{}
			resp, err := client.Do(req)
			Expect(err).To(BeNil())
			resp.Body.Close()
			tr.Finish()

			Expect(tr.DNSLookup().Seconds()).ToNot(BeZero())
			Expect(tr.TCPDialed().Seconds()).ToNot(BeZero())
			Expect(tr.ConnSetup().Seconds()).ToNot(BeZero())
			Expect(tr.PreTransfer().Seconds()).ToNot(BeZero())
			Expect(tr.TimeToFirstByte().Seconds()).ToNot(BeZero())
			Expect(tr.ServerProcessing().Seconds()).ToNot(BeZero())
			Expect(tr.ContentTransfer().Seconds()).ToNot(BeZero())
			Expect(tr.Total().Seconds()).ToNot(BeZero())
		}
	})

	It("Handles timing request even if request was not executed", func() {
		req, err := http.NewRequest(http.MethodGet, "localhost", nil)
		Expect(err).To(BeNil())
		_, tr := tracer.AsTraceableReq(req)

		Expect(tr.DNSLookup().Seconds()).To(BeZero())
		Expect(tr.TCPDialed().Seconds()).To(BeZero())
		Expect(tr.ConnSetup().Seconds()).To(BeZero())
		Expect(tr.PreTransfer().Seconds()).To(BeZero())
		Expect(tr.ServerProcessing().Seconds()).To(BeZero())
		Expect(tr.ContentTransfer().Seconds()).To(BeZero())
		Expect(tr.Total().Seconds()).To(BeZero())
	})

	It("Handles nil requests", func() {
		_, tr := tracer.AsTraceableReq(nil)

		Expect(tr.DNSLookup().Seconds()).To(BeZero())
		Expect(tr.TCPDialed().Seconds()).To(BeZero())
		Expect(tr.ConnSetup().Seconds()).To(BeZero())
		Expect(tr.PreTransfer().Seconds()).To(BeZero())
		Expect(tr.ServerProcessing().Seconds()).To(BeZero())
		Expect(tr.ContentTransfer().Seconds()).To(BeZero())
		Expect(tr.Total().Seconds()).To(BeZero())
	})
})
