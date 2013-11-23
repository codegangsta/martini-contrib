package acceptlang

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	ACCEPT_LANGUAGE_HEADER = "Accept-Language"
)

// A single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
	Language string
	Quality  float32
}

// A collection of sortable AcceptLanguage instances.
type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Len() int           { return len(al) }
func (al AcceptLanguages) Swap(i, j int)      { al[i], al[j] = al[j], al[i] }
func (al AcceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }

func (al AcceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range al {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(al)-1 {
			output.WriteString(", ")
		}
	}

	if output.Len() == 0 {
		output.WriteString("[]")
	}

	return output.String()
}

// Creates a new handler that parses the Accept-Language HTTP header.
//
// The parsed structure is a slice of Accept-Language values stored in an
// AcceptLanguages instance, sorted based on the language qualifier.
//
// See the HTTP header fields specification for more details
// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4).
func Languages() martini.Handler {
	return func(context martini.Context, request *http.Request) {
		header := request.Header.Get(ACCEPT_LANGUAGE_HEADER)
		if header != "" {
			acceptLanguageHeaderValues := strings.Split(header, ",")
			acceptLanguages := make(AcceptLanguages, len(acceptLanguageHeaderValues))

			for i, languageRange := range acceptLanguageHeaderValues {
				// Check if a given range is qualified or not
				if qualifiedRange := strings.Split(languageRange, ";q="); len(qualifiedRange) == 2 {
					quality, error := strconv.ParseFloat(qualifiedRange[1], 32)
					if error != nil {
						// When the quality is unparseable, assume it's 1
						acceptLanguages[i] = AcceptLanguage{trimLanguage(qualifiedRange[0]), 1}
					} else {
						acceptLanguages[i] = AcceptLanguage{trimLanguage(qualifiedRange[0]), float32(quality)}
					}
				} else {
					acceptLanguages[i] = AcceptLanguage{trimLanguage(languageRange), 1}
				}
			}

			sort.Sort(acceptLanguages)
			context.Map(acceptLanguages)
		} else {
			// If we have no Accept-Language header just map an empty slice
			context.Map(make(AcceptLanguages, 0))
		}
	}
}

// Trim any excessive characters from the language value.
func trimLanguage(language string) string {
	return strings.Trim(language, " ")
}
