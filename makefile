.PHONY:  default  refresh  test  test-coverage  test-docker

default: test

bootstrap:
	cd terraform/bootstrap/ && terraform apply

deploy:
	GOOS=linux GOARCH=amd64 go build -o historian main.go
	zip -9 historian.zip historian
	cd terraform/deploy/ && terraform apply -auto-approve

destroy:
	cd terraform/deploy/ && terraform destroy

refresh:
	cookiecutter gh:sjansen/cookiecutter-golang --output-dir .. --config-file .cookiecutter.yaml --no-input --overwrite-if-exists
	git checkout go.mod go.sum

coverage:
	mkdir -p dist
	test -f dist/coverage.txt || go test -coverpkg ./... -coverprofile=dist/coverage.txt ./...
	go tool cover -html=dist/coverage.txt

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage: test coverage

test-docker:
	docker-compose --version
	docker-compose build --pull go
	docker-compose up --abort-on-container-exit --exit-code-from=go --force-recreate
