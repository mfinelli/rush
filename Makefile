SOURCES := $(wildcard cmd/*.go)
SOURCES += $(wildcard db/*.go)
SOURCES += $(wildcard server/*.go)

CGO_ENABLED := 1

all: rush

clean:
	rm -rf rush

rush: main.go $(SOURCES) package.json
	go build -o $@ -tags=nomsgpack main.go

.PHONY: all clean
