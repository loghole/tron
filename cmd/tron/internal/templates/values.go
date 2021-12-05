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
# Place default config values here
`

const ValuesLocalTemplate = `# Place config values for LOCAL development here
LOGGER_LEVEL: debug
`

const ValuesDevTemplate = `# Place develop overrides here
LOGGER_LEVEL: debug
`

const ValuesStgTemplate = `# Place staging overrides here
LOGGER_LEVEL: info
`

const ValuesProdTemplate = `# Place production overrides here
LOGGER_LEVEL: warn
`
