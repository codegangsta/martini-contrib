// Package secure is a middleware for Martini that helps enable some quick security wins.
//
//  package main
//
//  import (
//    "github.com/codegangsta/martini"
//    "github.com/codegangsta/martini-contrib/secure"
//  )
//
//  func main() {
//    m := martini.Classic()
//
//    m.Use(secure.Secure(secure.Options{
//      AllowedHosts: []string{"www.example.com", "sub.example.com"},
//      SSLRedirect: true,
//    }))
//
//    m.Get("/", func() string {
//      return "Hello World"
//    })
//
//    m.Run()
//  }
package secure

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

const (
	STSHeader           = "Strict-Transport-Security"
	STSSubdomainString  = "; includeSubdomains"
	FrameOptionsHeader  = "X-Frame-Options"
	FrameOptionsValue   = "DENY"
	ContentTypeHeader   = "X-Content-Type-Options"
	ContentTypeValue    = "nosniff"
	XSSProtectionHeader = "X-XSS-Protection"
	XSSProtectionValue  = "1; mode=block"
	CSPHeader           = "Content-Security-Policy"
)

// Options is a struct for specifying configuration options for the secure.Secure middleware.
type Options struct {
	// AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
	AllowedHosts []string
	// If SSLRedirect is set to true, then only allow https requests. Default is false.
	SSLRedirect bool
	// SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
	SSLHost string
	// SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
	SSLProxyHeaders map[string]string
	// STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
	STSSeconds int64
	// If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
	STSIncludeSubdomains bool
	// If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
	FrameDeny bool
	// CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option.
	CustomFrameOptionsValue string
	// If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
	ContentTypeNosniff bool
	// If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	BrowserXssFilter bool
	// ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
	ContentSecurityPolicy string
}

func Secure(opt Options) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		// AllowedHost check.
		if len(opt.AllowedHosts) > 0 {
			isGoodHost := false
			for _, allowedHost := range opt.AllowedHosts {
				if strings.EqualFold(allowedHost, req.Host) {
					isGoodHost = true
					break
				}
			}

			if isGoodHost == false {
				http.Error(res, "Bad Host", http.StatusInternalServerError)
			}
		}

		// SSL check.
		if opt.SSLRedirect {
			isSSL := false
			if strings.EqualFold(req.URL.Scheme, "https") {
				isSSL = true
			} else {
				for hKey, hVal := range opt.SSLProxyHeaders {
					if req.Header.Get(hKey) == hVal {
						isSSL = true
						break
					}
				}
			}

			if isSSL == false {
				url := req.URL
				url.Scheme = "https"
				url.Host = req.Host

				if opt.SSLHost != "" {
					url.Host = opt.SSLHost
				}

				http.Redirect(res, req, url.String(), http.StatusMovedPermanently)
			}
		}

		//---------------------------------------------------------------------
		c.Next()
		//---------------------------------------------------------------------

		// Strict Transport Security header.
		if opt.STSSeconds != 0 {
			stsSub := ""
			if opt.STSIncludeSubdomains {
				stsSub = STSSubdomainString
			}

			res.Header().Add(STSHeader, fmt.Sprintf("max-age=%d%s", opt.STSSeconds, stsSub))
		}

		// Frame Options header.
		if opt.CustomFrameOptionsValue != "" {
			res.Header().Add(FrameOptionsHeader, opt.CustomFrameOptionsValue)
		} else if opt.FrameDeny {
			res.Header().Add(FrameOptionsHeader, FrameOptionsValue)
		}

		// Content Type Options header.
		if opt.ContentTypeNosniff {
			res.Header().Add(ContentTypeHeader, ContentTypeValue)
		}

		// XSS Protection header.
		if opt.BrowserXssFilter {
			res.Header().Add(XSSProtectionHeader, XSSProtectionValue)
		}

		// Content Security Policy header.
		if opt.ContentSecurityPolicy != "" {
			res.Header().Add(CSPHeader, opt.ContentSecurityPolicy)
		}
	}
}
