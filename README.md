# tracer
Annotate a golang http request to trace time taken at various steps.

Heavily inspired by https://github.com/davecheney/httpstat :)

Note: Just as any other instrumentation tool, don't keep this running in production forever without understanding the cost. If you are doing a lot of requests, this might have some (tiny) impact.

# Usage:
```golang
import "github.com/urjitbhatia/tracer"

req, _ := http.NewRequest(http.MethodGet, "http://golang.org", body)
req, tr = tracer.AsTraceableReq(req)
// ...
resp, err := client.Do(req)
if err != nil {
    log.Fatalf("failed to read response: %v", err)
}
body := readResponseBody(req, resp)
resp.Body.Close()
// Important to call this to calculate total time taken, otherwise it will be negative
tr.Finish()

// ..
log.Printlf("Tracer stats: %v", tr)
```
```
2018/08/10 06:23:48 Tracer stats: {"total":1.917301,"dnsLookup":0.526224, "tcpDialed":0.170330, "connSetup":0.684807, "preTransfer":1.211031, "ttfb":1.917138, "serverProcessing":0.706107, "contentTransfer":0.000162}
```

##### If you see bugs, please feel free to post issues.
