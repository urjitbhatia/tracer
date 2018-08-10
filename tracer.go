package tracer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Tracer records timings for an http request
type Tracer struct {
	dnsStart     time.Time
	dnsDone      time.Time
	tcpConnStart time.Time
	tcpConnDone  time.Time
	gotConn      time.Time
	firstByte    time.Time
	bodyReadDone time.Time
}

// Finish should be called as soon as a request body is read
// This method is important to calculate the total time taken
func (t *Tracer) Finish() {
	t.bodyReadDone = time.Now()
}

// DNSLookup returns dns lookup time taken
func (t *Tracer) DNSLookup() time.Duration {
	return t.dnsDone.Sub(t.dnsStart)
}

// TCPDialed returns time taken to dial the connection over tcp
func (t *Tracer) TCPDialed() time.Duration {
	return t.tcpConnDone.Sub(t.tcpConnStart)
}

// ConnSetup returns time taken to fully create a connection that
// is ready to perform an http request
func (t *Tracer) ConnSetup() time.Duration {
	return t.gotConn.Sub(t.dnsDone)
}

// PreTransfer returns total time taken to setup everything before the
// server has a connection ready to go.
func (t *Tracer) PreTransfer() time.Duration {
	return t.gotConn.Sub(t.dnsStart)
}

// ServerProcessing returns time taken by the server to receive the request
// and start replying (TimeToFirstByte since connection setup)
func (t *Tracer) ServerProcessing() time.Duration {
	return t.firstByte.Sub(t.gotConn)
}

// TimeToFirstByte returns time since start till FirstResponse Byte
func (t *Tracer) TimeToFirstByte() time.Duration {
	return t.firstByte.Sub(t.dnsStart)
}

// ContentTransfer returns time taken from first byte till body is fully read
func (t *Tracer) ContentTransfer() time.Duration {
	return t.bodyReadDone.Sub(t.firstByte)
}

// Total returns time take from very start to when tracer.Finish() was called
func (t *Tracer) Total() time.Duration {
	return t.bodyReadDone.Sub(t.dnsStart)
}

// String representation
func (t *Tracer) String() string {
	return fmt.Sprintf(`{"total":%f,"dnsLookup":%f, "tcpDialed":%f, "connSetup":%f, "preTransfer":%f, "ttfb":%f, "serverProcessing":%f, "contentTransfer":%f}`,
		t.Total().Seconds(), t.DNSLookup().Seconds(), t.TCPDialed().Seconds(), t.ConnSetup().Seconds(),
		t.PreTransfer().Seconds(), t.TimeToFirstByte().Seconds(), t.ServerProcessing().Seconds(), t.ContentTransfer().Seconds())
}

// AsTraceableReq wraps a given *http.Request with tracing timers and returns the wrapped request
// Use the returned Tracer value to query for various timing values that were recorded
func AsTraceableReq(req *http.Request) (*http.Request, *Tracer) {
	tr := &Tracer{}
	if req == nil {
		// return empty state
		return req, tr
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			tr.dnsStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			tr.dnsDone = time.Now()
			if tr.dnsStart.IsZero() {
				tr.dnsStart = tr.dnsDone
			}
		},
		ConnectStart: func(_, _ string) {
			tr.tcpConnStart = time.Now()
			if tr.dnsDone.IsZero() {
				// dns skipped
				tr.dnsDone = tr.tcpConnStart
				tr.dnsStart = tr.tcpConnStart
			}
		},
		ConnectDone: func(net, addr string, err error) {
			tr.tcpConnDone = time.Now()
			if tr.tcpConnStart.IsZero() {
				tr.dnsDone = tr.tcpConnDone
				tr.dnsStart = tr.tcpConnDone
				tr.tcpConnStart = tr.tcpConnDone
			}
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			tr.gotConn = time.Now()
			if tr.tcpConnStart.IsZero() {
				tr.dnsDone = tr.gotConn
				tr.dnsStart = tr.gotConn
				tr.tcpConnStart = tr.gotConn
				tr.tcpConnDone = tr.gotConn
			}
		},
		GotFirstResponseByte: func() {
			tr.firstByte = time.Now()
		},
	}

	return req.WithContext(httptrace.WithClientTrace(context.Background(), trace)), tr
}
