package main

var usageTemplate = `
{{- define "arg" -}}
  --{{ .Name }}=<arg> {{ "" -}}
{{- end -}}

{{- define "function" -}}
### {{ .Name | ToPascalCase }}

{{ BackTick }}shell
dagger call {{ .Name | ToSnakeCase }} {{ range .Args }}{{ template "arg" . }}{{ end }}
{{ BackTick }}

{{ end -}}
## Usage

{{ range .Functions -}}
  {{ template "function" . }}
{{- end }}
`
