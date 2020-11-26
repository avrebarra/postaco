package markdown

import (
	"regexp"
)

var regexHeader = regexp.MustCompile(`#+\s.+`)
var regexHeaderTag = regexp.MustCompile(`{\d+}$`)

func ExtractHeaders(b []byte) (m []string) {
	m = []string{}
	rs := regexHeader.FindAll(b, -1)
	if rs == nil {
		return
	}
	for _, r := range rs {
		m = append(m, string(r))
	}
	return
}

func ExtractHeaderTag(str string) (tag string) {
	rs := regexHeaderTag.FindAllString(str, -1)
	if rs == nil {
		return
	}
	tag = rs[len(rs)-1]
	return
}
