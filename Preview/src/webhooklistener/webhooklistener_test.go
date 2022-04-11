package webhooklistener

import (
	"net/http"
	"net/http/httptest"
	"testing"

	hs "github.com/clarkezone/hookserve/hookserve"
)

func TestGiteaParse(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upper?word=abc", nil)
	//req.Header.Set()
	w := httptest.NewRecorder()

	hs.ServeHttp(req, w)

}
