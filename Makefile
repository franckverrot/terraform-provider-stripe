default: build

.PHONY=build
build:
	go build -o terraform-provider-stripe

test: build
	terraform init
	terraform fmt
	terraform plan -out terraform.tfplan
	terraform apply terraform.tfplan

.PHONY: authors
authors:
	./scripts/generate_authors