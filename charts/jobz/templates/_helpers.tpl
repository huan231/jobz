{{/*
Expand the name of the chart.
*/}}
{{- define "jobz.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified server component name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "jobz.server.fullname" -}}
{{- if .Values.server.fullnameOverride }}
{{- .Values.server.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- printf "%s-%s" .Release.Name .Values.server.name | trunc 63 | trimSuffix "-" -}}
{{- else }}
{{- printf "%s-%s-%s" .Release.Name $name .Values.server.name | trunc 63 | trimSuffix "-" -}}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified ui component name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "jobz.ui.fullname" -}}
{{- if .Values.ui.fullnameOverride }}
{{- .Values.ui.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- printf "%s-%s" .Release.Name .Values.ui.name | trunc 63 | trimSuffix "-" -}}
{{- else }}
{{- printf "%s-%s-%s" .Release.Name $name .Values.ui.name | trunc 63 | trimSuffix "-" -}}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "jobz.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "jobz.labels" -}}
helm.sh/chart: {{ include "jobz.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "jobz.selectorLabels" -}}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Selector labels for the server component
*/}}
{{- define "jobz.server.selectorLabels" -}}
app.kubernetes.io/name: {{ printf "%s-%s" (include "jobz.name" .) .Values.server.name | trunc 63 | trimSuffix "-" }}
{{ include "jobz.selectorLabels" . }}
{{- end }}

{{/*
Selector labels for the ui component
*/}}
{{- define "jobz.ui.selectorLabels" -}}
app.kubernetes.io/name: {{ printf "%s-%s" (include "jobz.name" .) .Values.ui.name | trunc 63 | trimSuffix "-" }}
{{ include "jobz.selectorLabels" . }}
{{- end }}

{{/*
Common labels for the server component
*/}}
{{- define "jobz.server.labels" -}}
{{ include "jobz.labels" . }}
{{ include "jobz.server.selectorLabels" . }}
{{- end }}

{{/*
Common labels for the ui component
*/}}
{{- define "jobz.ui.labels" -}}
{{ include "jobz.labels" . }}
{{ include "jobz.ui.selectorLabels" . }}
{{- end }}

{{/*
Create the name of the service account to use for the server component
*/}}
{{- define "jobz.server.serviceAccountName" -}}
{{- default (include "jobz.server.fullname" .) .Values.server.serviceAccount.name }}
{{- end }}
