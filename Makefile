SUBDIRS := sm2 sm3 sm4
all: dep lint $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

dep:
	go mod verify

lint:
	go vet
	go fmt

.PHONY: all dep lint $(SUBDIRS)
