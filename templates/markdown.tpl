# {{ .Title }}
{{ if .HasAbout }}
### About

{{ range .About }}{{ . }}{{ end}}
{{ end }}

{{ if .HasConstants }}
### Constants
{{ range .Constants }}
* `{{ .Name}} = {{ .Value }}` {{ .UnitedDesc }} (_{{ .TypeName 0 }}_){{ end }}
{{ end }}

{{ if .HasVariables }}
### Global Variables
{{ range .Variables }}
* `{{ .Name}} = {{ .Value }}` {{ .UnitedDesc }} (_{{ .TypeName 0 }}_){{ end }}
{{ end }}

{{ if .HasMethods }}
### Methods
{{ range .Methods }}
`{{ .Name }}` - {{ .UnitedDesc }}
{{ range .Arguments }}* {{ .Index }}: {{ .Desc }} {{ if not .IsUnknown }}(_{{ .TypeName 0 }}_){{ end }}{{ if .IsOptional }} [_Optional_]{{ end }}
{{ end }}{{ end }}{{ end }}
