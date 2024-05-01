{{/*
Expand the name of the chart.
*/}}
{{- define "msockperf-chart.name" -}}
{{- default .Chart.Name .Values.msockperf.appName | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "msockperf-chart.fullname" -}}
{{- if .Values.msockperf.fullnameOverride }}
{{- .Values.msockperf.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.msockperf.appName }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

 __    __  ______  ______  ______  __  __   ______  ______  ______  ______  
/\ "-./  \/\  ___\/\  __ \/\  ___\/\ \/ /  /\  == \/\  ___\/\  == \/\  ___\ 
\ \ \-./\ \ \___  \ \ \/\ \ \ \___\ \  _"-.\ \  _-/\ \  __\\ \  __<\ \  __\ 
 \ \_\ \ \_\/\_____\ \_____\ \_____\ \_\ \_\\ \_\   \ \_____\ \_\ \_\ \_\   
  \/_/  \/_/\/_____/\/_____/\/_____/\/_/\/_/ \/_/    \/_____/\/_/ /_/\/_/   
                                                                           
Get ready to sockperf your things with latency with SRE ⚡️ monitoring principals in mind!
AKA MSOCKPERF - Monitoring + Sockperf ⚡️⚡️⚡️

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "msockperf-chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "msockperf-chart.labels" -}}
helm.sh/chart: {{ include "msockperf-chart.chart" . }}
{{ include "msockperf-chart.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "msockperf-chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "msockperf-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Annotations
*/}}
{{- define "msockperf-chart.annotations" -}}
meta.helm.sh/release-name: {{ .Release.Name }}
meta.helm.sh/release-namespace: {{ .Release.Namespace }}
{{- end }}
