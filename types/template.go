package types

import (
	"regexp"
	"strings"
)

type Template struct {
	ID      string `json:"id"`
	UserId  string `json:"-" boltholdIndex:"UserId"`
	Content string `json:"content"`
}

func (t *Template) FillContent(vars map[string]string) string {
	content := t.Content
	for k, v := range vars {
		re := regexp.MustCompile("{{\\s*" + k + "\\s*}}")
		content = re.ReplaceAllString(content, v)
	}
	return content
}

func (t *Template) GuessIsHtml() bool {
	return strings.Contains(t.Content, "<html>")
}
