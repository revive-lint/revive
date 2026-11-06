package fixtures

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

// Reported cases.

func discarded() {
	http.Get("http://example.com/") // MATCH /response body must be closed/
}

func blankResponse() error {
	_, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	return err
}

func onlyStatus() int {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

func readNotClosed() ([]byte, error) {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func decodedNotClosed(v any) error {
	resp, err := http.DefaultClient.Do(&http.Request{}) // MATCH /response body must be closed/
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func printed() {
	resp, _ := http.Get("http://example.com/") // MATCH /response body must be closed/
	fmt.Print(resp)
	fmt.Print(resp.Body)
}

func reassigned() {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return
	}
	fmt.Print(resp.Status)

	resp, err = http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return
	}
	fmt.Print(resp.Status)

	resp, _ = http.Get("http://example.com/")
	defer resp.Body.Close()
}

func fetch() (*http.Response, error) {
	return http.Get("http://example.com/")
}

func wrapperIgnored() {
	_, _ = fetch() // MATCH /response body must be closed/
}

func varDeclNotClosed() {
	var resp, _ = fetch() // MATCH /response body must be closed/
	if resp != nil {
		fmt.Println(resp.Header)
	}
}

func blankBodyAlias() {
	resp, _ := http.Get("http://example.com/") // MATCH /response body must be closed/
	_ = resp.Body
}

func inLoop(urls []string) {
	for _, u := range urls {
		resp, err := http.Get(u) // MATCH /response body must be closed/
		if err != nil {
			continue
		}
		fmt.Println(resp.Status)
	}
}

func readResponse(br *bufio.Reader) int {
	resp, err := http.ReadResponse(br, nil) // MATCH /response body must be closed/
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

func ifInitNotClosed() {
	if resp, err := http.Get("http://example.com/"); err == nil { // MATCH /response body must be closed/
		fmt.Println(resp.Status)
	}
}

// Not reported cases.

func deferClose() error {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func plainClose() {
	resp, _ := http.Get("http://example.com/")
	resp.Body.Close()
}

func chainedClose() {
	must(http.Get("http://example.com/")).Body.Close()
}

func bodyAlias() {
	resp, _ := http.Get("http://example.com/")
	body := resp.Body
	defer body.Close()
}

func bodyAliasVar() {
	resp, _ := http.Get("http://example.com/")
	var body io.ReadCloser = resp.Body
	body.Close()
}

func closedInDeferredClosure() {
	resp, _ := http.Get("http://example.com/")
	defer func() {
		_ = resp.Body.Close()
	}()
}

func closedInClosureDeclaredBefore() {
	var resp *http.Response
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	resp, _ = http.Get("http://example.com/")
}

func closedInsideClosure() {
	_ = func() {
		resp, _ := http.Get("http://example.com/")
		resp.Body.Close()
	}
}

func handle(resp *http.Response) {
	defer resp.Body.Close()
}

func passedToHelper() {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	handle(resp)
}

func deferredHelper() {
	resp, _ := http.Get("http://example.com/")
	defer handle(resp)
}

func closeQuietly(c io.Closer) {
	_ = c.Close()
}

func bodyPassedToCloser() {
	resp, _ := http.Get("http://example.com/")
	defer closeQuietly(resp.Body)
}

func returnedResponse() (*http.Response, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func returnedBody() io.ReadCloser {
	resp, _ := http.Get("http://example.com/")
	return resp.Body
}

type holder struct {
	resp *http.Response
	body io.ReadCloser
}

func storedInField(h *holder) {
	h.resp, _ = http.Get("http://example.com/")
	resp, _ := http.Get("http://example.com/")
	h.body = resp.Body
}

func storedInStruct() holder {
	resp, _ := http.Get("http://example.com/")
	return holder{resp: resp}
}

func sentOnChannel(ch chan<- *http.Response) {
	resp, _ := http.Get("http://example.com/")
	ch <- resp
}

var globalResp *http.Response

func storedInGlobal() {
	globalResp, _ = http.Get("http://example.com/")
}

func argumentOfCall() {
	handle(must(http.Get("http://example.com/")))
}

func must(resp *http.Response, err error) *http.Response {
	if err != nil {
		panic(err)
	}
	return resp
}

func recorder() {
	w := httptest.NewRecorder()
	resp := w.Result()
	_, _ = io.ReadAll(resp.Body)
}

func assignedInBranches(useDefault bool, req *http.Request) error {
	var resp *http.Response
	var err error
	if useDefault {
		resp, err = http.DefaultTransport.RoundTrip(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func assignedInSwitch(kind int) {
	var resp *http.Response
	switch kind {
	case 0:
		resp, _ = http.Get("http://example.com/a")
	default:
		resp, _ = http.Get("http://example.com/b")
	}
	defer resp.Body.Close()
}

// bodyclose issue #1: a wrapper returns the response and the caller closes it.
func getWrapped() *http.Response {
	resp, _ := http.Get("http://example.com/")
	return resp
}

func closeWrapped() {
	resp := getWrapped()
	resp.Body.Close()
}

// bodyclose issue #3: the close result is used as an argument.
func closeResultPrinted() {
	resp, _ := http.DefaultClient.Do(&http.Request{})
	defer func() { fmt.Println(resp.Body.Close()) }()
}

func closeResultPassed() {
	resp, _ := http.Get("http://example.com/")
	defer func() { handleClose("test", resp.Body.Close()) }()
}

func handleClose(_ string, _ error) {}

// bodyclose issue #5: the body is drained, then closed.
func drainAndClose() {
	res, _ := http.Get("http://example.com/")
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
}

// bodyclose issue #27: the body is returned as an io.Closer.
func returnedBodyAsCloser(url string) (io.Closer, error) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	return r.Body, nil
}

// bodyclose issue #42: the response is stored in a struct field.
type myResponse struct {
	original *http.Response
}

func (r *myResponse) response() *http.Response {
	return r.original
}

func storedInFieldClosedViaMethod() {
	r := &myResponse{}
	r.original, _ = http.Get("http://example.com/")
	_ = r.response().Body.Close()
}

// bodyclose issue #52: http.ResponseController is not an *http.Response.
func responseController(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})
	_, _ = io.Copy(w, r.Body)
}

// bodyclose issue #61: the response is assigned in a branch and closed after it.
func assignedInIf(cond bool) {
	var resp *http.Response
	if cond {
		resp, _ = http.Get("http://example.com/")
	}
	defer resp.Body.Close()
}

func ifInitClosed() {
	if resp, err := http.Get("http://example.com/"); err == nil {
		defer resp.Body.Close()
	}
}

func closedInGoroutine() {
	resp, _ := http.Get("http://example.com/")
	go func() {
		resp.Body.Close()
	}()
}

func closeMethodValue() {
	resp, _ := http.Get("http://example.com/")
	closeFn := resp.Body.Close
	defer closeFn()
}

func closedInNamedClosure() {
	res, _ := http.Get("http://example.com/")
	closer := func() {
		res.Body.Close()
	}
	closer()
}

func passedToClosure() {
	res, _ := http.Get("http://example.com/")
	closer := func(r *http.Response) {
		r.Body.Close()
	}
	closer(res)
}

// Differences with bodyclose, by design: a body passed to a parameter with a Close method
// is assumed to be closed by the callee, even if the callee does not close it.
func passedToNoopCloser() {
	resp, _ := http.Get("http://example.com/") // bodyclose issue #4 reports this
	noop(resp.Body)
}

func noop(_ io.ReadCloser) {}

func wrappedAndDiscarded() {
	resp, _ := http.Get("http://example.com/") // bodyclose issue #4 reports this
	_ = http.MaxBytesReader(nil, resp.Body, 1024)
}

func passedToDeferredLiteral() {
	resp, _ := http.Get("http://example.com/") // bodyclose defer_crash reports this
	defer func(_ io.ReadCloser) {}(resp.Body)
}

func newResponse() *http.Response {
	r := new(http.Response)
	return r
}

// io.NopCloser cases.

func closedAfterNopCloser() ([]byte, error) {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	defer resp.Body.Close() // closes the no-op wrapper, the original body leaks
	return b, err
}

func deferredLiteralBeforeNopCloser() ([]byte, error) {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return nil, err
	}
	defer func() {
		resp.Body.Close() // runs after the replacement, closes the no-op wrapper
	}()
	b, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return b, err
}

func returnedAfterNopCloserNotClosed() (*http.Response, error) {
	resp, err := http.Get("http://example.com/") // MATCH /response body must be closed/
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, err
}

func replacedWithNoBody() int {
	resp, _ := http.Get("http://example.com/") // MATCH /response body must be closed/
	resp.Body = http.NoBody
	defer resp.Body.Close()
	return resp.StatusCode
}

func bufferedResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, nil
}

func callerOfBufferedResponse() int {
	resp, err := bufferedResponse("http://example.com/") // the body is a no-op closer
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

func deferredCloseBeforeNopCloser() (*http.Response, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // resp.Body is evaluated before the replacement
	b, err := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	return resp, err
}

func aliasClosedAfterNopCloser() ([]byte, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return nil, err
	}
	body := resp.Body
	b, err := io.ReadAll(body)
	resp.Body = io.NopCloser(bytes.NewReader(b))
	defer body.Close()
	return b, err
}

func deferredLiteralAfterNopCloserOfOtherResponse() {
	resp, _ := http.Get("http://example.com/")
	defer func() {
		resp.Body.Close()
	}()
	other := new(http.Response)
	other.Body = io.NopCloser(bytes.NewReader(nil))
}

func delegatingWrapper() error {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		return err
	}
	resp.Body = struct {
		io.Reader
		io.Closer
	}{bufio.NewReader(resp.Body), resp.Body}
	defer resp.Body.Close()
	return nil
}
