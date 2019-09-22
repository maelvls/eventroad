//go:generate go get github.com/golang/protobuf/protoc-gen-go
//go:generate protoc eventpb.proto --go_out=plugins=grpc:.

package eventpb
