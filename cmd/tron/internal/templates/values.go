package templates

type ValuesData struct {
	List []Env
}

type Env struct {
	Key string
	Val string
}

const ValuesTemplate = `{{ range $value := .List -}}
{{ $value.Key }}: {{ $value.Val }}
{{end }}

`

const ValuesLocalTemplate = `# Place realtime config values for LOCAL development here
LOGGER_LEVEL: DEBUG

`

const ValuesDevTemplate = `# Place develop overrides here
LOGGER_LEVEL: INFO

`

const ValuesStgTemplate = `# Place staging overrides here
LOGGER_LEVEL: WARN

`

const ValuesProdTemplate = `# Place production overrides here
LOGGER_LEVEL: WARN

`
