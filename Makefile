VERSION="latest"

docker: 
	docker build -t qdledcontroller:${VERSION} -f Dockerfile.app .

pb: 
	protoc -I . -I /usr/include -I ../googleapis --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/led.proto

gw:
	protoc -I . -I /usr/include -I ../googleapis --grpc-gateway_out ./gateway --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true proto/led.proto && mv gateway/proto/*.go proto/
