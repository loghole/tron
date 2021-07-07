package templates

const Buf = `# Below is not the lint and breaking configuration we recommend!
# This just just what googleapis passes.
# For lint, we recommend having the single value "DEFAULT" in "use"
# with no values in "except".
# For breaking, we recommend having the single value "FILE" in use.
# See https://docs.buf.build/lint-usage
# See https://docs.buf.build/breaking-usage
version: v1beta1
build:
  roots:
    - vendor.pb
lint:
  use:
    - BASIC
    - FILE_LOWER_SNAKE_CASE
  except:
    - ENUM_NO_ALLOW_ALIAS
    - IMPORT_NO_PUBLIC
    - PACKAGE_AFFINITY
    - PACKAGE_DIRECTORY_MATCH
    - PACKAGE_SAME_DIRECTORY
breaking:
  use:
    - WIRE_JSON
`

const BufGen = `version: v1beta1
plugins:
  - name: go
    path: bin/protoc-gen-go
    out: .
    opt: paths=source_relative
    strategy: directory
  - name: go-grpc
    path: bin/protoc-gen-go-grpc
    out: .
    opt: paths=source_relative
    strategy: directory
  - name: grpc-gateway
    path: bin/protoc-gen-grpc-gateway
    out: .
    opt:
      - generate_unbound_methods=true
      - logtostderr=true
      - paths=source_relative
    strategy: directory
  - name: openapiv2
    path: bin/protoc-gen-openapiv2
    out: .
    opt:
      - generate_unbound_methods=true
      - fqn_for_openapi_name=true
    strategy: directory
  - name: tron
    path: bin/protoc-gen-tron
    out: .
    opt:
      - paths=source_relative
    strategy: directory
`
