# Vendored external ontologies

Local copies of the transitive `owl:imports` closure of the Bhash import
profiles, so `robot reason`/`robot report` resolve everything from the
repository via `ontology/src/catalog-v001.xml` without network access.

| File | Import IRI | Retrieved |
|---|---|---|
| `prov-o-20130430.ttl` | `http://www.w3.org/ns/prov-o-20130430`, also mapped from `http://www.w3.org/ns/prov-o#` (imported by DCAT) | 2026-08-10 |
| `dcat3.ttl` | `http://www.w3.org/ns/dcat` (serves DCAT 3; un-versioned upstream, pinned here) | 2026-08-10 |
| `dcterms.ttl` | `http://purl.org/dc/terms/` (imported by DCAT) | 2026-08-10 |
| `skos.ttl` | `http://www.w3.org/2004/02/skos/core` (imported by DCAT) | 2026-08-10 |

Both PROV IRIs map to the single archived document: the live `ns/prov-o` and
the archived `ns/prov-o-20130430` differ only in one revision assertion but
declare the same ontology IRI (`http://www.w3.org/ns/prov#`), and the OWL API
refuses to load two different documents with the same ontology ID.

All fetched with `Accept: text/turtle`. The W3C vocabularies are
redistributed under the [W3C Document License](https://www.w3.org/copyright/document-license/);
DCMI terms under the [DCMI licence](https://www.dublincore.org/about/copyright/).
Do not edit these files; refresh them from the IRIs above instead.
