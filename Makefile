# =====================
# Tools

setup:
ifeq (, $(shell which xk6))
	@echo installing xk6
	@go install go.k6.io/xk6/cmd/xk6@latest
	@echo xk6 successfully installed
endif

# =====================
# Targets

build: setup
	@xk6 build --with github.com/golioth/xk6-coap=.
