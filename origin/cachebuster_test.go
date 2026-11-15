package origin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestAddParam(t *testing.T) {
	for _, tc := range []struct{ target, want string }{
		{"https://example.com/", "https://example.com/?_cb=v"},
		{"https://example.com/a?x=1", "https://example.com/a?x=1&_cb=v"},
		{"https://example.com/a?x=1;y=2", "https://example.com/a?x=1;y=2&_cb=v"},
		{"https://example.com/a?", "https://example.com/a?_cb=v"},
		{"https://example.com/a?x=1&", "https://example.com/a?x=1&_cb=v"},
		{"https://example.com/a#top", "https://example.com/a?_cb=v#top"},
		{"https://example.com/a%2520b?q=%252F", "https://example.com/a%2520b?q=%252F&_cb=v"},
	} {
		if got := addParam(tc.target, "_cb", "v"); got != tc.want {
			t.Errorf("addParam(%q) = %q, want %q", tc.target, got, tc.want)
		}
	}
}

// probe calls the wrapper with the raw query and returns the query the wrapped handler sees.
func probe(t *testing.T, rawQuery string) (int, url.Values) {
	t.Helper()
	var seen url.Values
	h := CacheBuster(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { seen = r.URL.Query() }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/probe?"+rawQuery, nil))
	return rec.Code, seen
}

func TestCacheBuster(t *testing.T) {
	target := url.QueryEscape("https://example.com/a?x=1")

	_, q := probe(t, "module=web&target="+target)
	if got := q.Get("target"); got != "https://example.com/a?x=1" {
		t.Errorf("without cache_buster the target changed: %q", got)
	}

	_, q1 := probe(t, "module=web&cache_buster=_cb&target="+target)
	_, q2 := probe(t, "module=web&cache_buster=_cb&target="+target)
	re := regexp.MustCompile(`^https://example\.com/a\?x=1&_cb=[0-9a-f]{16}$`)
	if !re.MatchString(q1.Get("target")) {
		t.Errorf("unexpected target %q", q1.Get("target"))
	}
	if q1.Get("target") == q2.Get("target") {
		t.Errorf("value repeated: %q", q1.Get("target"))
	}
	if q1.Has(CacheBusterParam) || q1.Get("module") != "web" {
		t.Errorf("unexpected query passed on: %v", q1)
	}

	if code, q := probe(t, "module=web&cache_buster=a%26b&target="+target); code != http.StatusBadRequest || q != nil {
		t.Errorf("invalid name: code %d, handler called %v", code, q != nil)
	}
}
