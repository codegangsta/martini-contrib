package bundle

import (
	"fmt"
	"io/ioutil"
)

func concatFiles(files ...string) string {
	content := ""

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			content += fmt.Sprintf("/* ERR: %v */", err)
			continue
		}

		content += string(data)
	}

	return content
}
