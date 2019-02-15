default: build

.PHONY=build
build:
	go build -o terraform-provider-stripe

test: build
	terraform init
	terraform plan
	terraform apply

.PHONY: authors
authors:
	./scripts/generate_authors