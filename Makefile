SOURCES := $(wildcard cmd/*.go)
SOURCES += $(wildcard server/*.go)

all: rush

clean:
	rm -rf rush

rush: main.go $(SOURCES)
	go build -o $@ -tags=nomsgpack main.go

.PHONY: all clean
