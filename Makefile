PROTOCCMD = protoc
PROTOGEN_PATH = $(shell which protoc-gen-go) 
PROTOGENGRPC_PATH = $(shell which protoc-gen-go-grpc) 

GO_FILES := $(shell find $(SRC_DIR) -name '*.go')

GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean

LDFLAGS := -s -w

ifeq ($(OS), Windows_NT)
	DEFAULT_BUILD_FILENAME := StealthIMProxy.exe
else
	DEFAULT_BUILD_FILENAME := StealthIMProxy
endif

.PHONY: run
run: build
	./bin/$(DEFAULT_BUILD_FILENAME)

StealthIM.FileAPI/fileapi_grpc.pb.go StealthIM.FileAPI/fileapi.pb.go: proto/fileapi.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/fileapi.proto

StealthIM.GroupUser/groupuser_grpc.pb.go StealthIM.GroupUser/groupuser.pb.go: proto/groupuser.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/groupuser.proto

StealthIM.MSAP/msap_grpc.pb.go StealthIM.MSAP/msap.pb.go: proto/msap.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/msap.proto

StealthIM.Session/session_grpc.pb.go StealthIM.Session/session.pb.go: proto/session.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/session.proto

StealthIM.User/user_grpc.pb.go StealthIM.User/user.pb.go: proto/user.proto
	$(PROTOCCMD) --plugin=protoc-gen-go=$(PROTOGEN_PATH) --plugin=protoc-gen-go-grpc=$(PROTOGENGRPC_PATH) --go-grpc_out=. --go_out=. proto/user.proto

.PHONY: proto
proto: \
	StealthIM.FileAPI/fileapi_grpc.pb.go StealthIM.FileAPI/fileapi.pb.go \
	StealthIM.GroupUser/groupuser_grpc.pb.go StealthIM.GroupUser/groupuser.pb.go \
	StealthIM.MSAP/msap_grpc.pb.go StealthIM.MSAP/msap.pb.go \
	StealthIM.Session/session_grpc.pb.go StealthIM.Session/session.pb.go \
	StealthIM.User/user_grpc.pb.go StealthIM.User/user.pb.go

.PHONY: build
build: ./bin/$(DEFAULT_BUILD_FILENAME)

./bin/StealthIMProxy.exe: $(GO_FILES) proto
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/StealthIMProxy.exe

./bin/StealthIMProxy: $(GO_FILES) proto
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/StealthIMProxy

.PHONY: build_linux build_win
build_win: ./bin/StealthIMProxy.exe
build_linux: ./bin/StealthIMProxy

.PHONY: docker_run
docker_run:
	docker-compose up

./bin/StealthIMProxy.docker.zst: $(GO_FILES) proto
	docker-compose build
	docker save stealthimproxy-app > ./bin/StealthIMProxy.docker
	zstd ./bin/StealthIMProxy.docker -19
	@rm ./bin/StealthIMProxy.docker

.PHONY: build_docker
build_docker: ./bin/StealthIMProxy.docker.zst

.PHONY: release
release: build_win build_linux build_docker

.PHONY: clean
clean:
	@rm -rf ./StealthIM.*
	@rm -rf ./bin
	@rm -rf ./__debug*

.PHONY: dev
dev:
	./run_env.sh

.PHONY: debug_proto
debug_proto:
	cd test && python -m grpc_tools.protoc -I. --python_out=. --mypy_out=.  --grpclib_python_out=. --proto_path=../proto user.proto
	cd test && python -m grpc_tools.protoc -I. --python_out=. --mypy_out=.  --grpclib_python_out=. --proto_path=../proto proxy.proto
