{{- define "name" -}}
gardener-extension-backup-s3
{{- end -}}

{{- define "labels.app.key" -}}
app.kubernetes.io/name
{{- end -}}
{{- define "labels.app.value" -}}
{{ include "name" . }}
{{- end -}}

{{- define "labels" -}}
{{ include "labels.app.key" . }}: {{ include "labels.app.value" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{-  define "image" -}}
  {{- if hasPrefix "sha256:" .Values.image.tag }}
  {{- printf "%s@%s" .Values.image.repository .Values.image.tag }}
  {{- else }}
  {{- printf "%s:%s" .Values.image.repository .Values.image.tag }}
  {{- end }}
{{- end }}

{{- define "topologyAwareRouting.enabled" -}}
{{- if and .Values.gardener.seed .Values.gardener.seed.spec.settings.topologyAwareRouting.enabled }}
true
{{- end -}}
{{- end -}}

{{- define "seed.provider" -}}
  {{- if .Values.gardener.seed }}
{{- .Values.gardener.seed.provider }}
  {{- else -}}
""
  {{- end }}
{{- end -}}

{{- define "runtimeCluster.enabled" -}}
{{- if and .Values.gardener.runtimeCluster .Values.gardener.runtimeCluster.enabled }}
true
{{- end }}
{{- end -}}
