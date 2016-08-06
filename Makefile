WEBPACK=./node_modules/.bin/webpack

all:
	$(WEBPACK)
	go build

run: all
	./mpvsrv $(DIR)

watch:
	$(WEBPACK) --watch
