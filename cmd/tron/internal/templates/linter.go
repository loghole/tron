package templates

const GolangCILintTemplate = `# More info on config here: https://golangci-lint.run/usage/configuration/#config-file
run:
  tests: false
  timeout: 5m

linters:
  enable-all: true
  disable:
    - exhaustivestruct
  fast: false

linters-settings:
  gocritic:
    enabled-tags:
      - style
      - experimental
      - performance
      - diagnostic
      - opinionated

  govet:
    enable-all: true
    disable:
      - shadow
      - fieldalignment

  gci:
    local-prefixes: {{ .Module }}

issues:
  exclude:
    - use MixedCaps in package name
`
