package gitconfig

import (
	"text/template"
)

const dotGitConfigTemplate = `{{ range $section := .Keys }}{{ $section }}{{ range $name := ($.MustGet $section).Keys }}{{ range $val := ($.MustGet $section).MustGet $name}}
	{{ $name }} = {{ $val.Value }}{{ end }}{{ end }}
{{ end }}`

var gitConfigTemplate = template.Must(template.New("gitConfig").Parse(dotGitConfigTemplate))
