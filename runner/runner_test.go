package runner

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusAPI(t *testing.T) {
	r := getRouter()

	r.GET("/", statusResponseHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOk := w.Code == http.StatusOK
		return statusOk
	})
}
