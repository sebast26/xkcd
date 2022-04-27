modules:
	cd src && go mod download
lint:
	cd src && golangci-lint --timeout 2m --config ../.golangci.yml run --exclude-use-default=false ./...
test:
	cd src && go test -tags test -parallel=7 ./...
test-coverage:
	cd src && go test -tags test -parallel=7 ./... -coverprofile ../cover.out -covermode=atomic
	go install github.com/boumenot/gocover-cobertura@latest
	cd src && ${GOPATH}/bin/gocover-cobertura < ../cover.out > ../coverage.xml
	cd src && go tool cover -func ../cover.out | grep total