package hedera

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
)

const (
	// coreNS is the namespace of the Bhash core ontology and its service modules.
	coreNS = "https://hashgraphontology.xyz/core/"
	// resourceNS is the root of the canonical instance IRIs (D-0002).
	resourceNS = "https://hashgraphontology.xyz/resource/"
	// ledgerNS holds ledger fields that have no Bhash ontology property yet.
	// It is the namespace scripts/hedera_topic_to_fluree.py already uses for
	// such fields; terms in it are not part of the ontology.
	ledgerNS = "https://hashgraphontology.xyz/ledger#"
)

// networkIRIs maps a Hedera network name, which is also the {network} segment
// of the canonical resource IRIs, to its named individual in the core
// ontology (D-0005).
var networkIRIs = map[string]string{
	"mainnet":    coreNS + "Mainnet",
	"testnet":    coreNS + "Testnet",
	"previewnet": coreNS + "Previewnet",
}

// entityIDPattern matches a native Hedera entity ID (shard.realm.num).
var entityIDPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

var defaultContext = map[string]any{
	"hedera":                  coreNS,
	"ledger":                  ledgerNS,
	"prov":                    "http://www.w3.org/ns/prov#",
	"schema":                  "http://schema.org/",
	"xsd":                     "http://www.w3.org/2001/XMLSchema#",
	"prov:generatedAtTime":    map[string]any{"@type": "xsd:dateTime"},
	"hedera:registeredIn":     map[string]any{"@type": "@id"},
	"hedera:hasAccountId":     map[string]any{"@type": "xsd:string"},
	"hedera:hasTopicId":       map[string]any{"@type": "xsd:string"},
	"hedera:hasTokenId":       map[string]any{"@type": "xsd:string"},
	"hedera:hasSymbol":        map[string]any{"@type": "xsd:string"},
	"hedera:hasTreasury":      map[string]any{"@type": "@id"},
	"hedera:hasDecimals":      map[string]any{"@type": "xsd:integer"},
	"hedera:hasInitialSupply": map[string]any{"@type": "xsd:decimal"},
	"hedera:hasMaxSupply":     map[string]any{"@type": "xsd:decimal"},
	"schema:keywords":         map[string]any{"@container": "@set"},
}

// network identifies the Hedera network a bootstrap result was recorded on.
type network struct {
	slug string // {network} segment of the canonical resource IRIs
	iri  string // named individual in the core ontology
}

func resolveNetwork(name string) (network, error) {
	slug := strings.ToLower(strings.TrimSpace(name))
	iri, ok := networkIRIs[slug]
	if !ok {
		return network{}, fmt.Errorf("network %q has no named individual in the core ontology (expected mainnet, testnet or previewnet)", name)
	}
	return network{slug: slug, iri: iri}, nil
}

// resourceIRI mints the canonical instance IRI
// https://hashgraphontology.xyz/resource/{network}/{kind}/{shard}.{realm}.{num}.
func (n network) resourceIRI(kind, entityID string) (string, error) {
	if !entityIDPattern.MatchString(entityID) {
		return "", fmt.Errorf("%s ID %q is not a shard.realm.num entity ID", kind, entityID)
	}
	return resourceNS + n.slug + "/" + kind + "/" + entityID, nil
}

// Transaction builds a Fluree transaction that inserts JSON-LD nodes for every
// artefact recorded in the result. Nodes use the canonical resource IRIs and
// the core ontology vocabulary, and link to the network's named individual
// with hedera:registeredIn.
func (r BootstrapResult) Transaction(ledger string) (fluree.TransactionRequest, error) {
	net, err := resolveNetwork(r.Network)
	if err != nil {
		return fluree.TransactionRequest{}, err
	}

	ctx := make(map[string]any, len(defaultContext))
	for k, v := range defaultContext {
		ctx[k] = v
	}

	req := fluree.TransactionRequest{Ledger: ledger, Context: ctx}
	for _, account := range r.Accounts {
		node, err := account.asJSONLD(net)
		if err != nil {
			return fluree.TransactionRequest{}, err
		}
		req.Insert = append(req.Insert, node)
	}
	for _, topic := range r.Topics {
		node, err := topic.asJSONLD(net)
		if err != nil {
			return fluree.TransactionRequest{}, err
		}
		req.Insert = append(req.Insert, node)
	}
	for _, token := range r.Tokens {
		node, err := token.asJSONLD(net)
		if err != nil {
			return fluree.TransactionRequest{}, err
		}
		req.Insert = append(req.Insert, node)
	}
	return req, nil
}

func (a AccountRecord) asJSONLD(net network) (map[string]any, error) {
	id, err := net.resourceIRI("account", a.AccountID)
	if err != nil {
		return nil, err
	}
	node := map[string]any{
		"@id":                 id,
		"@type":               []string{"hedera:Account", "prov:Entity"},
		"hedera:hasAccountId": a.AccountID,
		"hedera:registeredIn": net.iri,
	}
	if !a.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(a.CreatedAt)
	}
	if a.Alias != "" {
		node["schema:name"] = a.Alias
	}
	if a.Memo != "" {
		node["schema:description"] = a.Memo
	}
	if a.PublicKey != "" {
		node["ledger:publicKey"] = a.PublicKey
	}
	if len(a.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), a.Tags...)
	}
	return node, nil
}

func (t TopicRecord) asJSONLD(net network) (map[string]any, error) {
	id, err := net.resourceIRI("topic", t.TopicID)
	if err != nil {
		return nil, err
	}
	node := map[string]any{
		"@id":                 id,
		"@type":               []string{"hedera:ConsensusTopic", "prov:Entity"},
		"hedera:hasTopicId":   t.TopicID,
		"hedera:registeredIn": net.iri,
	}
	if !t.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(t.CreatedAt)
	}
	if t.Memo != "" {
		node["schema:description"] = t.Memo
	}
	if len(t.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), t.Tags...)
	}
	if t.Sequence > 0 {
		node["ledger:initialSequence"] = t.Sequence
	}
	return node, nil
}

func (t TokenRecord) asJSONLD(net network) (map[string]any, error) {
	id, err := net.resourceIRI("token", t.TokenID)
	if err != nil {
		return nil, err
	}
	tokenClass, fungible, known := tokenClassFor(t.TokenType)
	node := map[string]any{
		"@id":                 id,
		"@type":               []string{tokenClass, "prov:Entity"},
		"hedera:hasTokenId":   t.TokenID,
		"hedera:registeredIn": net.iri,
	}
	if !known {
		node["ledger:tokenType"] = t.TokenType
	}
	if t.Name != "" {
		node["schema:name"] = t.Name
	}
	if t.Symbol != "" {
		node["hedera:hasSymbol"] = t.Symbol
	}
	if !t.CreatedAt.IsZero() {
		node["prov:generatedAtTime"] = formatTime(t.CreatedAt)
	}
	if t.Memo != "" {
		node["schema:description"] = t.Memo
	}
	if t.TreasuryAccountID != "" {
		treasury, err := net.resourceIRI("account", t.TreasuryAccountID)
		if err != nil {
			return nil, fmt.Errorf("token %s treasury: %w", t.TokenID, err)
		}
		node["hedera:hasTreasury"] = treasury
	}
	if t.Decimals > 0 {
		node["hedera:hasDecimals"] = t.Decimals
	}
	// hedera:hasInitialSupply has domain hedera:FungibleToken.
	if fungible && t.InitialSupply > 0 {
		node["hedera:hasInitialSupply"] = t.InitialSupply
	}
	if t.MaxSupply != 0 {
		node["hedera:hasMaxSupply"] = t.MaxSupply
	}
	if t.SupplyType != "" {
		node["ledger:supplyType"] = t.SupplyType
	}
	if len(t.Tags) > 0 {
		node["schema:keywords"] = append([]string(nil), t.Tags...)
	}
	return node, nil
}

// tokenClassFor maps a Hedera token type, as accepted by parseTokenType, to
// the matching Token Service class. An empty type is fungible, as in the SDK.
// Unrecognised types fall back to hedera:Token with known == false.
func tokenClassFor(tokenType string) (class string, fungible, known bool) {
	switch strings.ToUpper(strings.TrimSpace(tokenType)) {
	case "", "FUNGIBLE_COMMON", "TOKEN_TYPE_FUNGIBLE_COMMON":
		return "hedera:FungibleToken", true, true
	case "NON_FUNGIBLE_UNIQUE", "TOKEN_TYPE_NON_FUNGIBLE_UNIQUE":
		return "hedera:NonFungibleToken", false, true
	default:
		return "hedera:Token", false, false
	}
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
