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
EXAMPLE_VALUE: example_value
`

const ValuesLocalTemplate = `# Place config values for LOCAL development here
LOGGER_LEVEL: debug
`

const ValuesDevTemplate = `# Place develop overrides here
EXAMPLE_VALUE: example_value_dev
`

const ValuesStgTemplate = `# Place staging overrides here
EXAMPLE_VALUE: example_value_stg
`

const ValuesProdTemplate = `# Place production overrides here
EXAMPLE_VALUE: example_value_prod
`
