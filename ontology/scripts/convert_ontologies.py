#!/usr/bin/env python3
"""Convert ontology source files into deployment artefacts.

The script scans a source directory (default: ``ontology/src``) for files
matching the requested basis extension (default: ``ttl``). Each file is
parsed with ``rdflib`` and serialised to JSON-LD, RDF/XML (``.owl``), and
Turtle in the deployment directory while also emitting a human-readable
HTML catalogue using a Jinja2 template.
"""
from __future__ import annotations

import argparse
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, List, Optional, Sequence, Set, Tuple

import rdflib
from jinja2 import Environment, FileSystemLoader, select_autoescape
from rdflib import Graph, Literal, URIRef
from rdflib.namespace import DCTERMS, OWL, RDF, RDFS, SKOS

# Map filename extensions to rdflib parse formats.
FORMAT_BY_EXT: Dict[str, str] = {
    ".ttl": "turtle",
    ".jsonld": "json-ld",
    ".json": "json-ld",
    ".owl": "xml",
    ".rdf": "xml",
    ".xml": "xml",
}


@dataclass
class OntologyHeader:
    iri: Optional[str]
    title: Optional[str]
    description: Optional[str]


@dataclass
class ClassInfo:
    iri: str
    qname: str
    label: str
    definitions: List[Literal]
    comments: List[Literal]
    examples: List[Literal]
    subClassOf: List[str]


@dataclass
class PropertyInfo:
    iri: str
    qname: str
    label: str
    kind: str
    comments: List[Literal]
    domain: List[str]
    range: List[str]
    subPropertyOf: List[str]
    inverses: List[str]


def parse_args(argv: Optional[Sequence[str]] = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--source-dir",
        default=Path(__file__).resolve().parents[1] / "src",
        type=Path,
        help="Directory containing ontology source files (default: ontology/src)",
    )
    parser.add_argument(
        "--deployment-dir",
        default=Path(__file__).resolve().parents[1] / "deployment",
        type=Path,
        help="Directory to write generated artefacts (default: ontology/deployment)",
    )
    parser.add_argument(
        "--basis",
        default="ttl",
        help="File extension (without dot) that acts as the basis for conversion",
    )
    parser.add_argument(
        "--template",
        type=Path,
        default=Path(__file__).resolve().parent / "templates" / "ontology.html.j2",
        help="HTML Jinja2 template to render human-readable output",
    )
    return parser.parse_args(argv)


def ensure_default_prefix(g: Graph, ontology_iri: Optional[str]) -> None:
    if not ontology_iri:
        return
    base_ns = ontology_iri
    if not base_ns.endswith(("#", "/")):
        base_ns = base_ns + "#"
    try:
        g.namespace_manager.bind("", rdflib.term.URIRef(base_ns), replace=False)
    except Exception:
        pass


def ensure_common_prefixes(graph: Graph) -> None:
    common = {
        "rdf": RDF,
        "rdfs": RDFS,
        "owl": OWL,
        "skos": SKOS,
        "dcterms": DCTERMS,
    }
    for pref, ns in common.items():
        try:
            graph.namespace_manager.bind(pref, ns, replace=False)
        except Exception:
            pass


def compute_used_prefixes(g: Graph) -> List[Dict[str, str]]:
    nm = g.namespace_manager
    used_prefixes: Set[str] = set()
    for s, p, o in g:
        for term in (s, p, o):
            if isinstance(term, rdflib.term.URIRef):
                try:
                    prefix, _ns, _name = nm.compute_qname(term)
                    if prefix is not None:
                        used_prefixes.add(prefix)
                except Exception:
                    continue
    prefixes: List[Dict[str, str]] = []
    for prefix, ns in nm.namespaces():
        if prefix in used_prefixes:
            prefixes.append({"prefix": prefix or ":", "ns": str(ns)})
    prefixes.sort(key=lambda x: x["prefix"])
    return prefixes


def literal_by_lang(values: List[Literal], preferred: Optional[List[str]] = None) -> Tuple[Optional[Literal], List[Literal]]:
    preferred = preferred or ["en"]
    by_lang = {v.language: v for v in values if isinstance(v, Literal)}
    for lang in preferred:
        if lang in by_lang:
            return by_lang[lang], values
    return (values[0] if values else None), values


def qname(graph: Graph, term: URIRef) -> str:
    try:
        return graph.namespace_manager.normalizeUri(term)
    except Exception:
        return str(term)


def get_literals(graph: Graph, s: URIRef, p: URIRef) -> List[Literal]:
    return [o for o in graph.objects(s, p) if isinstance(o, Literal)]


def collect_ontology_info(g: Graph) -> OntologyHeader:
    ontologies = list(g.subjects(RDF.type, OWL.Ontology))
    iri: Optional[str] = None
    title: Optional[str] = None
    description: Optional[str] = None
    if ontologies:
        ont = ontologies[0]
        iri = str(ont)
        titles = get_literals(g, ont, DCTERMS.title)
        title_literal, _ = literal_by_lang(titles)
        if title_literal:
            title = str(title_literal)
        descriptions = get_literals(g, ont, DCTERMS.description)
        desc_literal, _ = literal_by_lang(descriptions)
        if desc_literal:
            description = str(desc_literal)
    return OntologyHeader(iri=iri, title=title, description=description)


def collect_classes(g: Graph) -> List[ClassInfo]:
    classes: Set[URIRef] = set(s for s in g.subjects(RDF.type, OWL.Class))
    items: List[ClassInfo] = []
    for s in classes:
        labels = get_literals(g, s, RDFS.label)
        label_literal, _ = literal_by_lang(labels)
        label = str(label_literal) if label_literal else qname(g, s)
        items.append(
            ClassInfo(
                iri=str(s),
                qname=qname(g, s),
                label=label,
                definitions=get_literals(g, s, SKOS.definition),
                comments=get_literals(g, s, RDFS.comment),
                examples=get_literals(g, s, SKOS.example),
                subClassOf=[qname(g, o) for o in g.objects(s, RDFS.subClassOf) if isinstance(o, URIRef)],
            )
        )
    items.sort(key=lambda x: (x.label.lower(), x.qname))
    return items


def property_kind(g: Graph, term: URIRef) -> str:
    if (term, RDF.type, OWL.ObjectProperty) in g:
        return "ObjectProperty"
    if (term, RDF.type, OWL.DatatypeProperty) in g:
        return "DatatypeProperty"
    return "Property"


def collect_properties(g: Graph) -> List[PropertyInfo]:
    props: Set[URIRef] = set(s for s in g.subjects(RDF.type, OWL.ObjectProperty))
    props |= set(s for s in g.subjects(RDF.type, OWL.DatatypeProperty))
    props |= set(s for s in g.subjects(RDF.type, RDF.Property))

    items: List[PropertyInfo] = []
    for s in props:
        labels = get_literals(g, s, RDFS.label)
        label_literal, _ = literal_by_lang(labels)
        label = str(label_literal) if label_literal else qname(g, s)
        domain = [qname(g, o) for o in g.objects(s, RDFS.domain) if isinstance(o, URIRef)]
        rng = [qname(g, o) for o in g.objects(s, RDFS.range) if isinstance(o, URIRef)]
        sub_props = [qname(g, o) for o in g.objects(s, RDFS.subPropertyOf) if isinstance(o, URIRef)]
        inverses_set: Set[URIRef] = set(o for o in g.objects(s, OWL.inverseOf) if isinstance(o, URIRef))
        inverses_set |= set(x for x in g.subjects(OWL.inverseOf, s) if isinstance(x, URIRef))
        inverses = [qname(g, u) for u in sorted(inverses_set, key=lambda u: qname(g, u))]

        items.append(
            PropertyInfo(
                iri=str(s),
                qname=qname(g, s),
                label=label,
                kind=property_kind(g, s),
                comments=get_literals(g, s, RDFS.comment),
                domain=domain,
                range=rng,
                subPropertyOf=sub_props,
                inverses=inverses,
            )
        )
    items.sort(key=lambda x: (x.label.lower(), x.qname))
    return items


def load_graph(path: Path) -> Graph:
    fmt = FORMAT_BY_EXT.get(path.suffix.lower())
    if fmt is None:
        raise ValueError(f"Unsupported input extension for {path}")
    g = Graph()
    g.parse(path, format=fmt)
    ensure_common_prefixes(g)
    header = collect_ontology_info(g)
    ensure_default_prefix(g, header.iri)
    return g


def render_html(graph: Graph, base_name: str, template_path: Path, output_path: Path, source_path: Path) -> None:
    env = Environment(
        loader=FileSystemLoader(template_path.parent),
        autoescape=select_autoescape(["html", "xml"]),
    )
    tmpl = env.get_template(template_path.name)

    ontology = collect_ontology_info(graph)
    classes = collect_classes(graph)
    properties = collect_properties(graph)
    prefixes = compute_used_prefixes(graph)

    html = tmpl.render(
        ontology=ontology,
        classes=classes,
        properties=properties,
        prefixes=prefixes,
        base_name=base_name,
        source_path=source_path,
    )
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(html, encoding="utf-8")


def serialise_graph(graph: Graph, dest_dir: Path, base_name: str) -> Dict[str, Path]:
    dest_dir.mkdir(parents=True, exist_ok=True)
    outputs = {
        "ttl": dest_dir / f"{base_name}.ttl",
        "jsonld": dest_dir / f"{base_name}.jsonld",
        "owl": dest_dir / f"{base_name}.owl",
    }
    graph.serialize(destination=outputs["ttl"], format="turtle")
    graph.serialize(destination=outputs["jsonld"], format="json-ld")
    graph.serialize(destination=outputs["owl"], format="xml")
    return outputs


def convert_file(path: Path, source_dir: Path, deployment_dir: Path, template_path: Path) -> Dict[str, Path]:
    graph = load_graph(path)
    relative = path.relative_to(source_dir)
    base_name = relative.stem
    dest_dir = deployment_dir / relative.parent

    serialised = serialise_graph(graph, dest_dir, base_name)
    html_path = dest_dir / f"{base_name}.html"
    render_html(graph, base_name, template_path, html_path, relative)
    serialised["html"] = html_path
    return serialised


def run(argv: Optional[Sequence[str]] = None) -> int:
    args = parse_args(argv)
    source_dir = args.source_dir.resolve()
    deployment_dir = args.deployment_dir.resolve()
    basis = args.basis.lstrip(".").lower()
    template_path = args.template.resolve()

    if not template_path.exists():
        raise FileNotFoundError(f"Template not found: {template_path}")

    basis_pattern = f"*.{basis}"
    source_files = sorted(source_dir.rglob(basis_pattern))
    if not source_files:
        print(f"No source files matching {basis_pattern} under {source_dir}")
        return 1

    print(f"Converting {len(source_files)} ontologies from {source_dir} using basis '.{basis}'...")
    deployment_dir.mkdir(parents=True, exist_ok=True)

    for path in source_files:
        print(f"- {path.relative_to(source_dir)}")
        outputs = convert_file(path, source_dir, deployment_dir, template_path)
        for label, out_path in outputs.items():
            print(f"  • {label}: {out_path.relative_to(deployment_dir)}")

    print(f"✅ Finished writing artefacts to {deployment_dir}")
    return 0


if __name__ == "__main__":
    sys.exit(run())
