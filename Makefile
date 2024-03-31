GO_BUILD_CMD=go build -o ./build/imgtoico
GO_CMD_FILE=cmd/img-to-ico/main.go

build:
	${GO_BUILD_CMD} ${GO_CMD_FILE}

build-windows:
	GOOS=windows ${GO_BUILD_CMD} ${GO_CMD_FILE}

.PHONY: build build-windows
