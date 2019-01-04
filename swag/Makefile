MOD := $(shell head -n1 ../go.mod | cut -d' ' -f2)

run:
	go run main.go

spec:
	go mod vendor
	docker run --rm -it \
		-v $(PWD)/..:/go/src/$(MOD) \
		-w /go/src/$(MOD)/swag \
		discobuilder \
		swag init
