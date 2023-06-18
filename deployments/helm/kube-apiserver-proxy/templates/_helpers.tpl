{{/*
Expand the name of the chart.
*/}}
{{- define "kube-apiserver-proxy.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kube-apiserver-proxy.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kube-apiserver-proxy.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kube-apiserver-proxy.labels" -}}
helm.sh/chart: {{ include "kube-apiserver-proxy.chart" . }}
{{ include "kube-apiserver-proxy.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kube-apiserver-proxy.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kube-apiserver-proxy.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kube-apiserver-proxy.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kube-apiserver-proxy.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Prints a slice of strings.
It correctly outputs values such as ["*"].
*/}}
{{- define "kube-apiserver-proxy.stringSlice" -}}
[
{{- range $i, $el := . -}}
    {{- if ne $i 0 -}}, {{- end -}}
    {{- . | quote -}}
{{- end -}}
]
{{- end }}
