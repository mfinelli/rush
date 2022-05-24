SOURCES := $(wildcard cmd/*.go)
SOURCES += $(wildcard db/*.go)
SOURCES += $(wildcard server/*.go)

SOURCES += $(wildcard src/*.tmpl)

WEBPACK := webpack.config.js package.json node_modules
LOGIN_INPUTS := src/login.js src/login.scss
LOGIN_OUTPUTS := dist/login.js dist/login.css

CGO_ENABLED := 1

all: rush

clean:
	rm -rf dist rush node_modules

$(LOGIN_OUTPUTS): $(LOGIN_INPUTS) $(WEBPACK)
	npm run build

node_modules: package.json package-lock.json
	npm ci

rush: main.go $(SOURCES) $(LOGIN_OUTPUTS) package.json
	go build -o $@ -tags=nomsgpack main.go

.PHONY: all clean
