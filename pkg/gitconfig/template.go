package gitconfig

import (
	"text/template"
)

const dotGitConfigTemplate = `{{ range $section, $data := .}}{{ $section }}{{ range $key, $vals := $data }}{{ range $vals }}
	{{ $key }} = {{ .Value }}{{ end }}{{ end }}
{{ end }}`

var gitConfigTemplate = template.Must(template.New("gitConfig").Parse(dotGitConfigTemplate))
