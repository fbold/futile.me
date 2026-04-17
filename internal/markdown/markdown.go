package markdown

import (
	"regexp"
	"strings"
)

func Parse(input string) string {
	output := input

	replacements := []struct {
		pattern *regexp.Regexp
		replacement string
	}{
		{regexp.MustCompile(`(?s)# (.+?)$`), "<h3>$1</h3>"},
		{regexp.MustCompile(`\*\*((?:[^*]|\*[^*])+)\*\*`), "<strong>$1</strong>"},
		{regexp.MustCompile(`\*((?:[^*]|\*[^*])+)\*`), "<em>$1</em>"},
		{regexp.MustCompile(`__((?:[^_]|_[^_])+)__`), "<u>$1</u>"},
		{regexp.MustCompile(`~~((?:[^~]|~[^~])+)~~`), "<del>$1</del>"},
	}

	for _, r := range replacements {
		output = r.pattern.ReplaceAllString(output, r.replacement)
	}

	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result = append(result, "<br>")
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "<br>")
}