package templates

//nolint:dupword // skip
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
  service_suffix: API
  use:
    - DEFAULT
  except:
    - RPC_REQUEST_STANDARD_NAME # not configurable and the request suffix is too long.
    - RPC_RESPONSE_STANDARD_NAME # not configurable and the response suffix is too long.
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
      - generate_unbound_methods=true
      - paths=source_relative
    strategy: directory
`
