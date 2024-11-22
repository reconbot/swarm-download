BINARY_NAME=swarm
MAIN_SOURCE=.

all: build

build:
	# env CGO_ENABLED=1 go build -trimpath -buildmode=pie -ldflags '-extldflags "-static -s -w -Wl,--allow-multiple-definition"' -o $(BINARY_NAME) $(MAIN_SOURCE)
	go build -trimpath -buildmode=pie -o $(BINARY_NAME) $(MAIN_SOURCE)
run: build
	./${BINARY_NAME}
