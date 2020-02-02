default: compile

.PHONY=compile
compile:
	go get ./...
	go build -o terraform-provider-stripe

test: compile
	terraform init
	terraform fmt
	terraform plan -out terraform.tfplan
	terraform apply terraform.tfplan

.PHONY=install
install: compile
	mkdir -p ~/.terraform.d/plugins
	cp ./terraform-provider-stripe ~/.terraform.d/plugins/

.PHONY: authors
authors:
	./scripts/generate_authors