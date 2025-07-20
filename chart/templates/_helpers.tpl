{{/*
Generate a name for the application
*/}}
{{- define "msa-broadcaster.name" -}}
msa-broadcaster
{{- end }}

{{/*
Generate the full name including release name
*/}}
{{- define "msa-broadcaster.fullname" -}}
{{ .Release.Name }}-msa-broadcaster
{{- end }}
