.PHONY: help

compare:
	go test -v -bench='^\QBenchmark_Compare_' -run=none .
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "This mechanism is used to compare various time mechanisms in golang, "
	@echo "using makefile for automation."
	@echo ""
	@echo "Available targets:"
	@echo "  compare     - benchmark test"
	@echo ""