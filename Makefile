
.PHONY: all build rebuild clean prepare client server data run

all: rebuild run

rebuild: clean build

build: prepare client server data

clean:
	rm -rf out

prepare:
	mkdir out/
	mkdir out/static/

client:
	cd client && npm run all
	cp client/dist/* out/static/

server:
	cd server && go build
	cp server/vordl out/

data:
	cp data/* out/

run:
	cd out/ && ./vordl
