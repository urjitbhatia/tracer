# tracer
Annotate a golang http request to trace time taken at various steps

# Usage:
```golang
import "github.com/urjitbhatia/tracer"

req, _ := http.NewReques(http.MethodGet, "http://localhost/foo", body)
req, tr = tracer.AsTraceableReq(req)
// ...
resp, err := client.Do(req)
if err != nil {
    log.Fatalf("failed to read response: %v", err)
}
body := readResponseBody(req, resp)
resp.Body.Close()
tr.Finish()

// ..
log.Printlf("Time taken to get connection for request: %v", tr.ConnSetup())
```