// Package origin holds the additions of the nehoupat fork, see ORIGIN.md.
package origin

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"regexp"
	"strings"
)

// CacheBusterParam is the /probe query parameter with the name of the parameter to add to the target.
const CacheBusterParam = "cache_buster"

var paramName = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,32}$`)

// CacheBuster wraps the /probe handler. With cache_buster=<name> in the request it adds <name>=<random hex> to the
// target URL, a new value for every probe, so the probe misses caches keyed by the full URL. The rest of the target,
// including its query, stays as it is. Without the parameter the request passes unchanged.
func CacheBuster(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		name := params.Get(CacheBusterParam)
		if name == "" {
			next.ServeHTTP(w, r)
			return
		}
		if !paramName.MatchString(name) {
			http.Error(w, fmt.Sprintf("Invalid %s parameter %q", CacheBusterParam, name), http.StatusBadRequest)
			return
		}
		target := params.Get("target")
		if target == "" {
			next.ServeHTTP(w, r)
			return
		}
		params.Set("target", addParam(target, name, fmt.Sprintf("%016x", rand.Uint64())))
		params.Del(CacheBusterParam)
		r2 := r.Clone(r.Context())
		r2.URL.RawQuery = params.Encode()
		next.ServeHTTP(w, r2)
	})
}

// addParam appends name=value to the query of target, before a fragment if there is one. The value and name are plain
// characters, so the extra unescape of the target in prober.Handler leaves them as they are.
func addParam(target, name, value string) string {
	fragment := ""
	if i := strings.IndexByte(target, '#'); i >= 0 {
		target, fragment = target[:i], target[i:]
	}
	sep := "?"
	if strings.Contains(target, "?") {
		sep = "&"
		if strings.HasSuffix(target, "?") || strings.HasSuffix(target, "&") {
			sep = ""
		}
	}
	return target + sep + name + "=" + value + fragment
}
