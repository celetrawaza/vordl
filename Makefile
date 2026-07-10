
.PHONY: autorun rebuild build clean prepare client server data collect run

rebuild: clean build

autorun: rebuild run

build: prepare client server data collect

clean:
	rm -rf out

prepare:
	mkdir -p out/

client:
	cd client && npm run all

server:
	cd server && go build -o vordl .

data: ;

collect:
	# client
	mkdir -p out/static
	cp -r client/dist/* out/static/
	# server
	cp server/vordl out/
	# data
	mkdir -p out/data/
	cp data/* out/data/

run:
	cd out/ && ./vordl
