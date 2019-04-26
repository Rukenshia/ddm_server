.PHONY: proto
proto:
	protoc -I proto/ proto/ddm.proto --go_out=plugins=grpc:proto --dart_out=grpc:dart

.PHONY: build
	go build -ldflags -H=windowsgui