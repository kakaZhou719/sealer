GitTag=$(version)

build: clean ## 构建二进制
	@echo "build sealer and sealutil bin"
	hack/build.sh $(GitTag)

install:
	@echo "install sealer bin"
	hack/install-sealer.sh

test:
	@echo "test sealer bin"
	hack/test-sealer.sh

clean: ## clean
	@rm -rf _output
