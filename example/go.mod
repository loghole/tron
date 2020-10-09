module github.com/loghole/example

go 1.15

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.2
	github.com/loghole/tron v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.7.1
	github.com/utrack/clay/v2 v2.4.9
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215
	google.golang.org/grpc v1.33.0
)

replace github.com/loghole/tron => ../
