.PHONY: test fmt showcover lint run bundle

bundle:
	#rm lambda-bundles/$(subdir)/$(func).zip
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lambda-bundles/$(func) cmd/lambda/$(func)/main.go
	cd lambda-bundles && \
		zip -X $(func).zip  $(func) && \
		rm $(func)

run:
	go run cmd/lambda/$(func)/main.go

test:
	@make lint && go test -v -coverprofile cp.out ./...

showcover:
	go tool cover -html=cp.out
	
fmt:
	gofmt -s -w .

lint:
	golangci-lint run

build:
	sam build

deploy:
	sam deploy
