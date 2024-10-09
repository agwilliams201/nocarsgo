build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap ./src/scraper.go

deploy: build
	serverless deploy --stage prod


clean:
	rm -rf ./bootstrap .serverless ./vendor Gopkg.lock
