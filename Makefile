build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/bootstrap ./src/scraper.go

deploy: build
	serverless deploy --stage prod

clean:
	rm -rf ./bin ./vendor Gopkg.lock ./serverless
