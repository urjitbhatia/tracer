# tracer
Annotate a golang http request to trace time taken at various steps.

Just as any other instrumentation tool, don't keep this running in production forever without understanding the cost. If you are doing a lot of requests, this might have some (tiny) impact.

# Usage:
```golang
import "github.com/urjitbhatia/tracer"

req, _ := http.NewRequest(http.MethodGet, "http://localhost/foo", body)
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
log.Printlf("Time taken to get connection for request: %v", tr.ConnSetup())
```

##### If you see bugs, please feel free to post issues.
