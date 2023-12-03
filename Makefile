GOPATH := $(shell go env GOPATH)
GOSRC := ${GOPATH}/src
PB_IDL_PATH:= ${GOPATH}/src/github.com/ergedage/web/proto

idl-gen:
    # pb -> go代码
	${GOPATH}/src/code.byted.org/bytelingoi18n/idl/bin/osx/protoc --proto_path=${PB_IDL_PATH} --go_out=${GOPATH}/src/github.com/ergedage/web/  ${PB_IDL_PATH}/message/biz_api/*.proto
#	find . -name "*.pb.go" | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}' && cd ..
#	# 代码生成器: 更新router.go 、xx_handler.go
#	(cd $(mktemp -d); GO111MODULE=on go install code.byted.org/bytelingo/lingo_gen@feat_ai_tutor)
#	lingo_gen  aiTutorApiGen --idl_path ${PB_IDL_PATH}/service/ai_tutor_service.proto --project_path  ${GOPATH}/src/code.byted.org/bytelingoi18n/ai_tutor_api








