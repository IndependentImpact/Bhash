ROBOT ?= robot
PYTHON ?= python3
VENV_DIR := build/venv
PYTHON_BIN := $(VENV_DIR)/bin/python

.PHONY: all reason-core report-core template-example shacl sparql fluree-smoke python-venv clean

all: reason-core report-core shacl sparql

build:
	mkdir -p build build/reports build/templates

reason-core: build
	$(ROBOT) reason --catalog ontology/src/catalog-v001.xml --reasoner ELK --input ontology/src/core.ttl --output build/core-reasoned.ttl

report-core: build
	mkdir -p build/reports
	$(ROBOT) report --catalog ontology/src/catalog-v001.xml --input ontology/src/core.ttl --output build/reports/core-report.tsv

template-example: build
	mkdir -p build/templates
	$(ROBOT) template --template templates/example.csv --output build/templates/example.ttl

python-venv:
	@[ -d $(VENV_DIR) ] || ($(PYTHON) -m venv $(VENV_DIR) && $(PYTHON_BIN) -m pip install --upgrade pip)
	$(PYTHON_BIN) -m pip install -r requirements.txt

shacl: python-venv
	$(PYTHON_BIN) scripts/run_shacl.py

sparql: python-venv
	$(PYTHON_BIN) scripts/run_sparql.py

fluree-smoke:
	go test ./internal/fluree ./scripts/flureeclient

clean:
	rm -rf build
