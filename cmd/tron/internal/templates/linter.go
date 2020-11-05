package templates

const GolangCILintTemplate = `# More info on config here: https://github.com/golangci/golangci-lint#config-file
run:
  tests: false
linters:
  enable-all: true
  disable:
   - godot
   - gci
   - exhaustivestruct
linters-settings:
  gocritic:
    enabled-tags:
      - style
      - experimental
      - performance
      - diagnostic
      - opinionated
service:
  golangci-lint-version: 1.32.x
`
