discobuilder:
	docker build -t discobuilder -f Dockerfile.builder .

msgp:
	./build.sh go generate ./...
