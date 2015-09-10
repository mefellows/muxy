package symptom

import (
	"net/http"
	"testing"
	"time"
)

func TestCreateCookieHeader(t *testing.T) {
	cookies := []http.Cookie{
		http.Cookie{
			Domain:   "foo.com",
			Expires:  time.Now(),
			HttpOnly: true,
			MaxAge:   1,
			Name:     "Foo",
			Path:     "/",
			Secure:   false,
			Value:    "blah",
		},
	}
	for _, c := range cookies {
		t.Logf("cookie: %s", c.String())

	}
}
