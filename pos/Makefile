.PHONY: test-accuracy check-env

test-accuracy:
	ONTONOTES_PATH=$(ONTONOTES_PATH) go test -v -timeout 20m tagger_accuracy_test.go


check-env:
ifndef ONTONOTES_PATH
    $(error ONTONOTES_PATH is undefined)
endif

