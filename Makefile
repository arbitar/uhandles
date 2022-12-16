all:
	docker build -t uhandles .
	docker run --rm -v "${PWD}":/host uhandles
	docker rmi uhandles

docker-prep:
	docker build -t uhandles .

docker-build:
	docker run --rm -v "${PWD}":/host uhandles

docker-clean:
	docker rmi uhandles

build:
	go build -o uhandles -tags netgo -installsuffix netgo --ldflags '-w -extldflags "-static"' uhandles.go	

clean:
	rm uhandles
