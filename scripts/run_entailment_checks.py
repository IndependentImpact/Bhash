#!/usr/bin/env python3
"""Check reasoner-backed entailments over the example datasets.

The SPARQL regression harness queries the example graphs without the ontology
sources, so it can only see explicitly asserted triples. This script applies
the OWL 2 RL closure (owlrl) to the ontology sources merged with the example
data and runs every violation query in ``tests/verify/`` — a query returning
any row is a failure. This is what catches a removed or broken axiom (e.g.
owl:inverseOf) that the plain SPARQL checks cannot.
"""
from __future__ import annotations

import pathlib
import sys

import owlrl
from rdflib import Graph

REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent
ONTOLOGIES = [
    REPO_ROOT / "ontology" / "src" / "core.ttl",
    REPO_ROOT / "ontology" / "src" / "consensus.ttl",
]
DATASETS = [
    REPO_ROOT / "ontology" / "examples" / "core-consensus.ttl",
]
VERIFY_DIR = REPO_ROOT / "tests" / "verify"


def main() -> int:
    graph = Graph()
    for path in ONTOLOGIES + DATASETS:
        graph.parse(path)
    owlrl.DeductiveClosure(owlrl.OWLRL_Semantics).expand(graph)

    failures = []
    for query_path in sorted(VERIFY_DIR.glob("*.rq")):
        rows = list(graph.query(query_path.read_text(encoding="utf-8")))
        if rows:
            failures.append(query_path.name)
            print(f"FAIL {query_path.name}:")
            for row in rows:
                print(f"  {', '.join(str(value) for value in row)}")
        else:
            print(f"PASS {query_path.name}")

    if failures:
        print(f"entailment check failures: {', '.join(failures)}")
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
