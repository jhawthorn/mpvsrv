WEBPACK=./node_modules/.bin/webpack

all:
	$(WEBPACK)
	go-bindata -debug static/...
	go build

release:
	$(WEBPACK) --optimize-minimize --optimize-dedupe
	go-bindata static/...
	go build

run: all
	./mpvsrv $(DIR)

watch:
	$(WEBPACK) --watch
