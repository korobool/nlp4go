.PHONY: test-prescision check-env

test-prescision:
	ONTONOTES_PATH=$(ONTONOTES_PATH) go test -v -timeout 20m tagger_prescision_test.go


check-env:
ifndef ONTONOTES_PATH
    $(error ONTONOTES_PATH is undefined)
endif

