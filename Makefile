.PHONY: test fmt showcover lint run bundle

bundle:
	#rm lambda-bundles/$(subdir)/$(func).zip
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lambda-bundles/$(subdir)/$(func) cmd/lambda/$(subdir)/$(func).go
	cd lambda-bundles/$(subdir) && \
		zip -X $(func).zip  $(func) && \
		rm $(func)

run:
	go run cmd/$(subdir)/$(exe).go

test:
	@make lint && go test -v -coverprofile cp.out ./...

showcover:
	go tool cover -html=cp.out
	
fmt:
	gofmt -s -w .

lint:
	golangci-lint run
