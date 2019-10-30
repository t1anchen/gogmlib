SUBDIRS := $(wildcard */.)

all: lint $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

lint:
	go vet
	go fmt

.PHONY: all $(SUBDIRS)
