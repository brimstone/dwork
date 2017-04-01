pb/dwork.pb.go: pb/dwork.proto
	protoc -I . $< --go_out=plugins=grpc:.
