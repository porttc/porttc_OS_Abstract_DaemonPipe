BUILDNAME=porttc

.PHONY: all build buildRun fmtAll clean testing release deploy

all: clean
	go build -o ./build/$(BUILDNAME) ./src/

deployBuild:
	go build -trimpath -ldflags="-s -w" -o ./build/$(BUILDNAME) ./src/
	upx --ultra-brute ./build/$(BUILDNAME)

debugBuild:
	go build -tags debug -o ./build/$(BUILDNAME) ./src/

build:
	go build -o ./build/$(BUILDNAME) ./src/

buildRun: clean build
	chmod +x ./build/$(BUILDNAME)
	./build/$(BUILDNAME)

fmtAll:
	go fmt ./src/*

clean:
	rm -f ./build/*