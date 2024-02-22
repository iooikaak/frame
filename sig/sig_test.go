package sig

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	param := url.Values{}
	param.Set("id1", "1")
	param.Set("id2", "1")
	req, _ := http.NewRequest("GET", "/", strings.NewReader(param.Encode()))
	a, b := GetSign("test", req)
	t.Logf("%v---%v", a, b)
}
