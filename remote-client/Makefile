
ARCH=amd64

LDFLAGS_STATIC=-ldflags="-extldflags '-static'"

all:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) \
		go build $(CFLAGS) \
		-a $(LDFLAGS_STATIC)
