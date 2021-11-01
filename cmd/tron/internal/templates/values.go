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
`

const ValuesDevTemplate = `# Place develop overrides here
`

const ValuesStgTemplate = `# Place staging overrides here
`

const ValuesProdTemplate = `# Place production overrides here
`
