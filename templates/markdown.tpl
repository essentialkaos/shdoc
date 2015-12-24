# {{ .Title }}
{{ if .HasAbout }}
### About

{{ range .About }}{{ . }}{{ end}}
{{ end }}

{{ if .HasConstants }}
### Constants
{{ range .Constants }}
* `{{ .Name}} = {{ .Value }}` {{ .UnitedDesc }} (_{{ .TypeName }}_){{ end }}
{{ end }}

{{ if .HasVariables }}
### Global Variables
{{ range .Variables }}
* `{{ .Name}} = {{ .Value }}` {{ .UnitedDesc }} (_{{ .TypeName }}_){{ end }}
{{ end }}

{{ if .HasMethods }}
### Methods
{{ range .Methods }}
`{{ .Name }}` - {{ .UnitedDesc }}
{{ range .Arguments }}* {{ .Index }}: {{ .Desc }} {{ if not .IsUnknown }}(_{{ .TypeName }}_){{ end }}{{ if .IsOptional }} [_Optional_]{{ end }}
{{ end }}{{ end }}{{ end }}
