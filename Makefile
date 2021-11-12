generate-pb:
	./scripts/generate-pb.sh
.PHONY: generate-pb

clean:
	./scripts/clean.sh
.PHONY: clean

lint:
	./scripts/lint.sh
.PHONY: lint

build:
	./scripts/build.sh
.PHONY: build

build-clean: clean build
.PHONY: build-clean

docker-down:
	./scripts/docker-down.sh
.PHONY: docker-down

docker-up:
	./scripts/docker-up.sh
.PHONY: docker-up