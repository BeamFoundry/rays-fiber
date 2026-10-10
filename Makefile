.PHONY: all
all: test

.PHONY: fmt
fmt:
	go vet ./...
	go fmt ./...

.PHONY: test-rays-fiber
test-rays-fiber: fmt
	ginkgo run --cover ./middleware ./security ./storage

.PHONY: test-storages
test-storages: fmt
	cd storage/memcache && ginkgo run --cover ./
	cd storage/redis && ginkgo run --cover ./

.PHONY: test
test: test-rays-fiber test-storages
