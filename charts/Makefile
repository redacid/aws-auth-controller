-include .env
GIT_STAMP ?= $(shell git describe || echo v0.1.0)

KUBEVAL_SCHEMA_LOCATION ?= https://raw.githubusercontent.com/redacid/k8s-json-scheme/refs/heads/main
KUBECONFORM_SCHEMA ?= --schema-location 'https://raw.githubusercontent.com/redacid/k8s-json-scheme/refs/heads/main/{{.NormalizedKubernetesVersion}}-standalone{{.StrictSuffix}}/{{.ResourceKind}}{{.KindSuffix}}.json'
CHART_PATH ?= ./aws-auth-controller
CREATE_NAMESPACE ?= --create-namespace

ENV_NAME ?= TEST

TEST_VALUES ?= -f ./values.yaml
TEST_NAMESPACE ?= aws-auth-ns
TEST_RELEASE_NAME ?= auth-test
TEST_KUBERNETES_VERSION ?= 1.32.2

# colors
GREEN = $(shell tput -Txterm setaf 2)
YELLOW = $(shell tput -Txterm setaf 3)
WHITE = $(shell tput -Txterm setaf 7)
RESET = $(shell tput -Txterm sgr0)
GRAY = $(shell tput -Txterm setaf 6)
TARGET_MAX_CHAR_NUM = 30

.EXPORT_ALL_VARIABLES:

all: help

## Install helm-docs | Docs
.PHONY: install_helm-docs
install_helm-docs:
	wget https://github.com/norwoodj/helm-docs/releases/download/v1.14.2/helm-docs_1.14.2_Linux_x86_64.tar.gz
	tar -zxvf helm-docs_1.14.2_Linux_x86_64.tar.gz helm-docs
	sudo mv helm-docs ./bin/
	rm helm-docs_1.14.2_Linux_x86_64.tar.gz

## Generate Helm Chart Docs
.PHONY: helm-docs-generate
helm-docs-generate:
	helm-docs -c $(CHART_PATH) -o README.md

## Check helm | Linters
helm-lint: helm-dependency-update
	@$(foreach var,$(ENV_NAME), echo "\n\n======= Check $(var) =======\n\n" \
	&& helm lint --namespace=$($(var)_NAMESPACE) $(CHART_PATH) $($(var)_VALUES) || exit;)

## Check helm-kubeval
helm-kubeval: helm-dependency-update
	@$(foreach var,$(ENV_NAME), echo "\n\n======= Check $(var) on kubernetes-version $($(var)_KUBERNETES_VERSION) =======\n\n" \
	&& helm kubeval --exit-on-error --kubernetes-version $($(var)_KUBERNETES_VERSION) \
	--strict --name-template=$($(var)_RELEASE_NAME) --namespace=$($(var)_NAMESPACE) $($(var)_VALUES) \
	$(CHART_PATH) || exit;)

## Check Kubeconform
helm-kubeconform: helm-dependency-update install-kubeconform-plugin
	@$(foreach var,$(ENV_NAME), echo "\n\n======= Check $(var) on kubernetes-version $($(var)_KUBERNETES_VERSION) =======\n\n" \
	&& helm kubeconform --output pretty --goroutines $(shell nproc) --exit-on-error $(KUBECONFORM_SCHEMA) \
		--kube-version $($(var)_KUBERNETES_VERSION) \
		--strict --summary --name-template=$($(var)_RELEASE_NAME) --namespace=$($(var)_NAMESPACE) $($(var)_VALUES) $(CHART_PATH) || exit;)

## Deploy TEST
deploy-test:  helm-dependency-update helm-package
	@helm upgrade -i --namespace=$(TEST_NAMESPACE) $(TEST_RELEASE_NAME) ./aws-auth-controller-$(GIT_STAMP).tgz $(TEST_VALUES) $(CREATE_NAMESPACE)

uninstall-test:
	@helm uninstall --namespace=$(TEST_NAMESPACE) $(TEST_RELEASE_NAME)

## Deploy TEST
template-test:  helm-dependency-update
	@helm template --namespace=$(TEST_NAMESPACE) $(TEST_RELEASE_NAME) $(CHART_PATH) $(TEST_VALUES)

## Dependency update helm | Tools
helm-dependency-update:
	@helm dependency update $(CHART_PATH)
	@helm dependency list $(CHART_PATH)
	@ls -l $(CHART_PATH)

## Build helm package
helm-package:
	@helm package --app-version=$(GIT_STAMP) --version=$(GIT_STAMP) $(CHART_PATH)

## Install kubeconform helm plugin
install-kubeconform-plugin:
	@helm plugin list | grep -q kubeconform || helm plugin install https://github.com/redacid/helm-kubeconform --version 0.6.8.1


## Shows help. | Help
help:
	@echo ''
	@echo 'Usage:'
	@echo ''
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z0-9\-_]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
		    if (index(lastLine, "|") != 0) { \
				stage = substr(lastLine, index(lastLine, "|") + 1); \
				printf "\n ${GRAY}%s: \n\n", stage;  \
			} \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			if (index(lastLine, "|") != 0) { \
				helpMessage = substr(helpMessage, 0, index(helpMessage, "|")-1); \
			} \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@echo ''
