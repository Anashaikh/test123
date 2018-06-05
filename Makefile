files := $(shell find . -path ./vendor -prune -path ./pb -prune -o -name '*.go' -print)
pkgs := $(govendor list -no-status +local)

.PHONY: all format test vet lint checkformat check

all : check
check : checkformat vet lint test

format :
	@echo "== format"
	@goimports -w $(files)
	@sync

unformatted = $(shell goimports -l $(files))

checkformat :
	@echo "== check formatting"
ifneq "$(unformatted)" ""
	@echo "needs formatting: $(unformatted)"
	@echo "run 'make format'"
	@exit 1
endif

vet :
	@echo "== vet"
	@govendor vet +local

pkgs = $(shell govendor list -no-status +local)

lint :
	@echo "== lint"
	@for pkg in $(pkgs); do \
		golint -set_exit_status $$pkg || exit 1; \
	done;

test :
	@echo "== run tests"
	govendor test -race
