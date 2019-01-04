MOD := $(shell head -n1 go.mod | cut -d' ' -f2)
spec:
	go mod vendor
	docker run --rm -it \
		-v $(PWD):/go/src/$(MOD) \
		-w /go/src/$(MOD) \
		quay.io/goswagger/swagger \
		generate spec -o ./swagger.json
