WEBPACK=./node_modules/.bin/webpack

all:
	$(WEBPACK)
	go build

run: all
	./mpvsrv

watch:
	$(WEBPACK) --watch
