run:
	go run -race main.go

build:
	go build -race main.go

test:
	go test -v ./...

lint:
	golangci-lint run --fix

generate:
	go generate -x ./...

#rmcert:
#	rm ./cert/*.pem

cert:
	cd cert; ./gen.sh; cd ..