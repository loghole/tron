package templates

const GolangCILintTemplate = `# More info on config here: https://github.com/golangci/golangci-lint#config-file
run:
  tests: false
  skip-dirs:
  - scripts
linters:
  enable-all: true
  disable:
   - godot
   - gci
linters-settings:
  gocritic:
    enabled-tags:
      - style
      - experimental
      - performance
      - diagnostic
      - opinionated
service:
  golangci-lint-version: 1.31.x
`
