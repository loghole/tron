module example

go 1.15

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.2
	github.com/loghole/tron v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/utrack/clay/v2 v2.4.9
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a
	google.golang.org/grpc v1.27.1
)

replace github.com/loghole/tron => ../
