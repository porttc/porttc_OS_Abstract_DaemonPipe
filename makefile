BUILDNAME=porttcOSABSTR

.PHONY: all build buildRun fmtAll clean testing release deploy

all: clean
	go build -o ./build/$(BUILDNAME) .

deployBuild:
	go build -trimpath -ldflags="-s -w" -o ./build/$(BUILDNAME) .
	upx --ultra-brute ./build/$(BUILDNAME)

debugBuild:
	go build -tags debug -o ./build/$(BUILDNAME) .

build:
	go build -o ./build/$(BUILDNAME) .

buildRun: clean build
	chmod +x ./build/$(BUILDNAME)
	./build/$(BUILDNAME)

test:
	go test ./testing/ -v

fmtAll:
	go fmt ./*

clean:
	rm -f ./build/*