package recovery

import (
	"bufio"
	"github.com/codegangsta/martini"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// This is development-only middleware. It was more or less copied from github.com/pilu/traffic
// and github.com/gocraft/web
func Recovery() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, logger *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				const size = 4096
				stack := make([]byte, size)
				stack = stack[:runtime.Stack(stack, false)]

				logger.Printf("PANIC: %s\n%s", err, string(stack))
				renderRecovery(res, req, err, stack)
			}
		}()

		c.Next()
	}
}

func renderRecovery(res http.ResponseWriter, req *http.Request, err interface{}, stack []byte) {
	_, filePath, line, _ := runtime.Caller(5)

	data := map[string]interface{}{
		"Error":    err,
		"Stack":    string(stack),
		"Params":   req.URL.Query(),
		"Method":   req.Method,
		"FilePath": filePath,
		"Line":     line,
		"Lines":    readErrorFileLines(filePath, line),
	}

	rw := res.(martini.ResponseWriter)
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusInternalServerError)

	tpl := template.Must(template.New("ErrorPage").Parse(panicPageTpl))
	tpl.Execute(rw, data)
}

func readErrorFileLines(filePath string, errorLine int) map[int]string {
	lines := make(map[int]string)

	file, err := os.Open(filePath)
	if err != nil {
		return lines
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	currentLine := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil || currentLine > errorLine+5 {
			break
		}

		currentLine++

		if currentLine >= errorLine-5 {
			lines[currentLine] = strings.Replace(line, "\n", "", -1)
		}
	}

	return lines
}

const panicPageTpl string = `
  <html>
    <head>
      <title>Panic</title>
      <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
      <style type="text/css">
        html, body{ font-family: "Helvetica Neue", Arial, Helvetica, sans-serif; padding: 0; margin: 0; }
        header { background: #e74c3c; color: white; border-bottom: 2px solid #c0392b; }
        h1 { padding: 10px 0; margin: 0; }
        h2 { font-size: 20px; }
        p { font-size: 14px; }
        .container { margin: 0 20px; }
        .error { background: #ffd3e0; color: #c0392b; padding: 10px 0; box-shadow: 0 1px 3px rgba(192, 57, 43, .9); }
        .file-info .file-name { font-weight: bold; }
        .stack { height: 300px; overflow-y: scroll; border: 1px solid #e5e5e5; padding: 10px; }

        table.source {
          border-collapse: collapse;
          border: 1px solid #e5e5e5;
          width: 100%;
        }

        table.source td {
          padding: 0;
        }

        table.source .numbers {
          font-size: 14px;
          vertical-align: top;
          width: 1%;
          color: rgba(0,0,0,0.3);
          text-align: right;
        }

        table.source .numbers .number {
          display: block;
          padding: 0 5px;
          border-right: 1px solid #e5e5e5;
        }

        table.source .numbers .number.line-{{ .Line }} {
          border-right: 1px solid #ffcccc;
          font-weight: bold;
        }

        table.source .numbers pre {
          white-space: pre-wrap;
        }

        table.source .code {
          font-size: 14px;
          vertical-align: top;
        }

        table.source .code .line {
          padding-left: 10px;
          display: block;
        }

        table.source .numbers .number,
        table.source .code .line {
          padding-top: 1px;
          padding-bottom: 1px;
        }

        table.source .code .line:hover {
          background-color: #f6f6f6;
        }

        table.source .line-{{ .Line }},
        table.source line-{{ .Line }},
        table.source .code .line.line-{{ .Line }}:hover {
          background: #ffd3e0;
          color: #c0392b;
        }
      </style>
    </head>
  <body>
    <header>
      <div class="container">
        <h1>Error</h1>
      </div>
    </header>

    <div class="error">
      <p class="container">{{ .Error }}</p>
    </div>

    <div class="container">
      <p class="file-info">
        In <span class="file-name">{{ .FilePath }}:{{ .Line }}</span></p>
      </p>

      <table class="source">
        <tr>
          <td class="numbers">
            <pre>{{ range $lineNumber, $line :=  .Lines }}<span class="number line-{{ $lineNumber }}">{{ $lineNumber }}</span>{{ end }}</pre>
          </td>
          <td class="code">
            <pre>{{ range $lineNumber, $line :=  .Lines }}<span class="line line-{{ $lineNumber }}">{{ $line }}<br /></span>{{ end }}</pre>
          </td>
        </tr>
      </table>
      <h2>Stack</h2>
      <pre class="stack">{{ .Stack }}</pre>
      <h2>Request</h2>
      <p><strong>Method:</strong> {{ .Method }}</p>
      <p>Parameters:</p>
      <ul>
        {{ range $key, $value := .Params }}
          <li><strong>{{ $key }}:</strong> {{ $value }}</li>
        {{ end }}
      </ul>
    </div>
  </body>
  </html>
  `
