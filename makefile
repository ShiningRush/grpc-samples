SRV_NAME = go-test-grpc-srv
PROTO_NAME = hello
MAIN_PATH = ./cmd/server/
SERVER_BIN_DIST_DIR = ./dist/server/
REMOTE_REPO = git@gitlab.followme.com:deploy/go-crmreport-grpc-srv.git
CURRENT_VERSION := `head -1 releaselog `

BRANCH_NAME = ${PROTO_NAME}_$(CURRENT_VERSION) #it only for push proto

clean:
	rm -rf dist/server
	rm -rf dist/swaggerui
	rm -rf proto/*.go
	rm -rf proto/*.json

gen-swag:
	cd proto && protoc --swagger_out=logtostderr=true,grpc_api_configuration=${PROTO_NAME}.yaml:. ./${PROTO_NAME}.proto
	cd proto && protoc --grpc-gateway_out=logtostderr=true,grpc_api_configuration=${PROTO_NAME}.yaml:std ./${PROTO_NAME}.proto
	go-bindata --nocompress -pkg swaggerui -o dist/swaggerui/datafile.go  thirdparty/swaggerui/...

gen:
	cd proto && protoc --go_out=plugins=grpc:go ./${PROTO_NAME}.proto

build:clean gen
	# GOOS=linux GOARCH=amd64 go build -o $(SERVER_BIN_DIST_DIR)/go-crmreport-grpc-srv ./src/main.go 
	go build -o $(SERVER_BIN_DIST_DIR)/$(SRV_NAME) $(MAIN_PATH)/main.go
	# cd $(SERVER_BIN_DIST_DIR) && chmod +x $(SRV_NAME)
	cp -r $(MAIN_PATH)/configs $(SERVER_BIN_DIST_DIR)

push-proto:
	rm -rf ../Services.Protocol
	cd .. && git clone git@gitlab.followme.com:Followme/Services.Protocol.git
	cd ../Services.Protocol && git checkout -b $(BRANCH_NAME) && cp ../CRM-REPORT/proto/report.proto ./crmreport/report.proto && \
		git add . && git commit -m 'auto publish' && git push origin $(BRANCH_NAME)

push-git:build
	cp ./releaselog ./dist/server/releaselog
	cd dist/server && git init && git remote add origin $(REMOTE_REPO) && git add . && git commit -m 'auto publish' && git tag $(CURRENT_VERSION) && git push --tags

# 仅编译镜像，不推
build-image:
	docker build --no-cache --build-arg build=build --rm -t $(SRV_NAME):$(CURRENT_VERSION) .

# 发布项目
publish:
	docker build --no-cache --rm -t $(SRV_NAME):$(CURRENT_VERSION) .
	docker tag $(SRV_NAME):$(CURRENT_VERSION) $(DH_URL)$(SRV_NAME):$(CURRENT_VERSION)
	docker login -u $(DH_USER) -p $(DH_PASS) $(DH_URL) && docker push $(DH_URL)$(SRV_NAME):$(CURRENT_VERSION)

test:	
	echo $(DH_USER)
	echo $(DH_PASS)
	echo $(DH_URL)

.PHONY: gen build build-image publish clean test