package bundle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func minifyScript(script string) string {
	const closureApiUrl = "http://closure-compiler.appspot.com/compile"

	formData := url.Values{
		"js_code":           {script},
		"compilation_level": {"SIMPLE_OPTIMIZATIONS"},
		"output_format":     {"text"},
		"output_info":       {"compiled_code"},
	}

	resp, err := http.PostForm(closureApiUrl, formData)
	if err != nil {
		return prependError(err, script)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return prependError(err, script)
	}

	return string(respData)
}

func prependError(err error, script string) string {
	return fmt.Sprintf("/* failed to minify script: %v */ %v", err, script)
}
