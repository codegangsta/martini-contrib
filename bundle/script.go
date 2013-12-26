package bundle

import (
	"net/http"
)

type ScriptBundle struct {
	files     []string
	wrapScope bool
	minify    bool
}

func NewScriptBundle(wrapScope bool, minify bool) *ScriptBundle {
	return &ScriptBundle{
		files:     make([]string, 0),
		wrapScope: wrapScope,
		minify:    minify,
	}
}

func (s *ScriptBundle) Compile() string {
	content := concatFiles(s.files...)

	if s.wrapScope {
		content = "(function () {" + content + "})();"
	}

	if s.minify {
		content = minifyScript(content)
	}

	return content
}

func (s *ScriptBundle) AddFiles(files ...string) {
	s.files = append(s.files, files...)
}

func (s *ScriptBundle) Handler() martini.Handler {
	content := s.Compile()

	return func(res http.ResponseWriter) string {
		res.Header().Set("Content-Type", "text/javascript")
		return content
	}
}
