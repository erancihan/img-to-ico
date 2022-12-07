GO_BUILD_CMD=go build -o ./build/

build:
	${GO_BUILD_CMD}

build-windows:
	GOOS=windows ${GO_BUILD_CMD}

.PHONY: build build-windows
