package tracer

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Tracer struct {
	dnsStart     time.Time
	dnsDone      time.Time
	tcpConnStart time.Time
	tcpConnDone  time.Time
	gotConn      time.Time
	firstByte    time.Time
	bodyReadDone time.Time
}

func (t *Tracer) Finish() {
	t.bodyReadDone = time.Now()
}

func (t *Tracer) DNSLookup() time.Duration {
	return t.dnsDone.Sub(t.dnsStart)
}

func (t *Tracer) TCPDialed() time.Duration {
	return t.tcpConnDone.Sub(t.tcpConnStart)
}

func (t *Tracer) ConnSetup() time.Duration {
	return t.gotConn.Sub(t.dnsDone)
}

func (t *Tracer) PreTransfer() time.Duration {
	return t.gotConn.Sub(t.dnsStart)
}

func (t *Tracer) ServerProcessing() time.Duration {
	return t.firstByte.Sub(t.gotConn)
}

func (t *Tracer) ContentTransfer() time.Duration {
	return t.bodyReadDone.Sub(t.firstByte)
}

func (t *Tracer) Total() time.Duration {
	return t.bodyReadDone.Sub(t.dnsStart)
}

func AsTraceableReq(req *http.Request) (*http.Request, *Tracer) {
	var tr *Tracer

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { tr.dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { tr.dnsDone = time.Now() },
		ConnectStart: func(_, _ string) {
			tr.tcpConnStart = time.Now()
			if tr.dnsDone.IsZero() {
				tr.dnsDone = tr.tcpConnStart
			}
		},
		ConnectDone:          func(net, addr string, err error) { tr.tcpConnDone = time.Now() },
		GotConn:              func(_ httptrace.GotConnInfo) { tr.gotConn = time.Now() },
		GotFirstResponseByte: func() { tr.firstByte = time.Now() },
	}

	return req.WithContext(httptrace.WithClientTrace(context.Background(), trace)), tr
}
