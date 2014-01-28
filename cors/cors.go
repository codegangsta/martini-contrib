package cors

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerMaxAge           = "Access-Control-Max-Age"

	headerRequestMethod  = "Access-Control-Request-Method"
	headerRequestHeaders = "Access-Control-Request-Headers"
)

// Represents Access Control options.
type Opts struct {
	// If set, all origins are allowed.
	AllowAllOrigins bool
	// A list of allowed domain patterns.
	AllowOrigins []string
	// If set, allows to share auth credentials such as cookies.
	AllowCredentials bool
	// A list of allowed HTTP methods.
	AllowMethods []string
	// A list of allowed HTTP headers.
	AllowHeaders []string
	// Max age of the CORS headers.
	MaxAge time.Duration
}

// Converts options into a map of HTTP headers.
func (o *Opts) Header(origin string) map[string]string {
	headers := make(map[string]string)

	// add allow origin
	if o.AllowAllOrigins {
		headers[headerAllowOrigin] = "*"
	} else {
		headers[headerAllowOrigin] = origin
	}

	// add allow credentials
	headers[headerAllowCredentials] = strconv.FormatBool(o.AllowCredentials)

	// add allow methods
	if len(o.AllowMethods) > 0 {
		headers[headerAllowMethods] = strings.Join(o.AllowMethods, ",")
	}

	// add allow headers
	if len(o.AllowHeaders) > 0 {
		// TODO: Add default headers
		headers[headerAllowHeaders] = strings.Join(o.AllowHeaders, ",")
	}
	// add a max age header
	if o.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(o.MaxAge/time.Second), 10)
	}
	return headers
}

// Looks up if origin matches one of the patterns
// provided in Opts.AllowOrigins patterns.
func (o *Opts) IsOriginAllowed(origin string) (allowed bool) {
	for _, pattern := range o.AllowOrigins {
		allowed, _ = regexp.MatchString(pattern, origin)
		if allowed {
			return
		}
	}
	return
}

func Allow(opts *Opts) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" &&
			(req.Header.Get(headerRequestMethod) != "" || req.Header.Get(headerRequestHeaders) != "") {
			// TODO: if preflight, respond with exact headers if allowed
			return
		}

		origin := req.Header.Get("Origin")
		if !opts.AllowAllOrigins && !opts.IsOriginAllowed(origin) {
			return
		}

		for key, value := range opts.Header(origin) {
			res.Header().Set(key, value)
		}
	}
}
