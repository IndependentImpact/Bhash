package hedera

import (
	"context"
	"testing"
	"time"
)

func TestBootstrapperExecute(t *testing.T) {
	now := time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC)
	network := NewMockNetwork("testnet", WithStartingIDs(5000, 7000, 9000), WithNowFunc(func() time.Time { return now }))
	bootstrapper := NewBootstrapper(network, "testnet")
	spec := BootstrapSpec{
		Accounts: []AccountSpec{{Alias: "treasury", Memo: "Treasury account"}},
		Topics:   []TopicSpec{{Alias: "consensus", Memo: "Consensus topic"}},
		Tokens: []TokenSpec{{
			Alias:         "demo",
			Name:          "Demo Token",
			Symbol:        "DEM",
			TreasuryAlias: "treasury",
			SupplyType:    "INFINITE",
			TokenType:     "FUNGIBLE_COMMON",
		}},
	}
	result, err := bootstrapper.Execute(context.Background(), spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Accounts) != 1 || result.Accounts[0].AccountID != "0.0.5000" {
		t.Fatalf("unexpected accounts: %+v", result.Accounts)
	}
	if len(result.Topics) != 1 || result.Topics[0].TopicID != "0.0.7000" {
		t.Fatalf("unexpected topics: %+v", result.Topics)
	}
	if len(result.Tokens) != 1 || result.Tokens[0].TreasuryAccountID != "0.0.5000" {
		t.Fatalf("unexpected tokens: %+v", result.Tokens)
	}
}

func TestBootstrapTransaction(t *testing.T) {
	now := time.Date(2024, 9, 1, 12, 0, 0, 0, time.UTC)
	result := BootstrapResult{
		Network: "testnet",
		Accounts: []AccountRecord{{
			Alias:     "treasury",
			AccountID: "0.0.1001",
			Memo:      "Treasury",
			CreatedAt: now,
			Tags:      []string{"governance"},
		}},
		Topics: []TopicRecord{{
			Alias:     "consensus",
			TopicID:   "0.0.2001",
			Memo:      "Consensus",
			CreatedAt: now,
		}},
		Tokens: []TokenRecord{{
			Alias:             "token",
			TokenID:           "0.0.3001",
			Name:              "Demo",
			Symbol:            "DEM",
			TreasuryAccountID: "0.0.1001",
			Decimals:          2,
			InitialSupply:     1000,
			CreatedAt:         now,
		}},
	}
	tx, err := result.Transaction("tenant/dataset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Ledger != "tenant/dataset" {
		t.Fatalf("unexpected ledger: %s", tx.Ledger)
	}
	if got := tx.Context["hedera"]; got != "https://hashgraphontology.xyz/core/" {
		t.Fatalf("hedera prefix must be the core ontology namespace, got %v", got)
	}
	if len(tx.Insert) != 3 {
		t.Fatalf("expected 3 inserts, got %d", len(tx.Insert))
	}
	const testnet = "https://hashgraphontology.xyz/core/Testnet"
	const accountIRI = "https://hashgraphontology.xyz/resource/testnet/account/0.0.1001"

	accountNode := tx.Insert[0]
	if accountNode["@id"] != accountIRI {
		t.Fatalf("unexpected account IRI: %+v", accountNode)
	}
	if accountNode["hedera:hasAccountId"] != "0.0.1001" || accountNode["hedera:registeredIn"] != testnet {
		t.Fatalf("unexpected account node: %+v", accountNode)
	}

	topicNode := tx.Insert[1]
	if topicNode["@id"] != "https://hashgraphontology.xyz/resource/testnet/topic/0.0.2001" ||
		topicNode["hedera:hasTopicId"] != "0.0.2001" || topicNode["hedera:registeredIn"] != testnet {
		t.Fatalf("unexpected topic node: %+v", topicNode)
	}

	tokenNode := tx.Insert[2]
	if tokenNode["@id"] != "https://hashgraphontology.xyz/resource/testnet/token/0.0.3001" ||
		tokenNode["hedera:hasTokenId"] != "0.0.3001" || tokenNode["hedera:registeredIn"] != testnet {
		t.Fatalf("unexpected token node: %+v", tokenNode)
	}
	if types := tokenNode["@type"].([]string); types[0] != "hedera:FungibleToken" {
		t.Fatalf("an empty token type should map to hedera:FungibleToken, got %v", types)
	}
	if tokenNode["hedera:hasTreasury"] != accountIRI {
		t.Fatalf("expected treasury link, got %+v", tokenNode)
	}
	if tokenNode["hedera:hasSymbol"] != "DEM" || tokenNode["hedera:hasDecimals"] != uint(2) || tokenNode["hedera:hasInitialSupply"] != uint64(1000) {
		t.Fatalf("unexpected token properties: %+v", tokenNode)
	}
	for _, node := range tx.Insert {
		if _, ok := node["hedera:belongsToNetwork"]; ok {
			t.Fatalf("hedera:belongsToNetwork is not an ontology property: %+v", node)
		}
	}
}

func TestBootstrapTransactionNonFungibleToken(t *testing.T) {
	result := BootstrapResult{
		Network: "Mainnet",
		Tokens: []TokenRecord{{
			TokenID:       "0.0.4001",
			TokenType:     "NON_FUNGIBLE_UNIQUE",
			InitialSupply: 5,
		}},
	}
	tx, err := result.Transaction("tenant/dataset")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := tx.Insert[0]
	if node["@id"] != "https://hashgraphontology.xyz/resource/mainnet/token/0.0.4001" ||
		node["hedera:registeredIn"] != "https://hashgraphontology.xyz/core/Mainnet" {
		t.Fatalf("unexpected node: %+v", node)
	}
	if types := node["@type"].([]string); types[0] != "hedera:NonFungibleToken" {
		t.Fatalf("expected hedera:NonFungibleToken, got %v", types)
	}
	if _, ok := node["hedera:hasInitialSupply"]; ok {
		t.Fatalf("hedera:hasInitialSupply has domain hedera:FungibleToken: %+v", node)
	}
	if _, ok := node["ledger:tokenType"]; ok {
		t.Fatalf("a recognised token type should not be kept as ledger:tokenType: %+v", node)
	}
}

func TestBootstrapTransactionRejectsUnknownNetwork(t *testing.T) {
	result := BootstrapResult{Network: "localnet", Accounts: []AccountRecord{{AccountID: "0.0.2"}}}
	if _, err := result.Transaction("tenant/dataset"); err == nil {
		t.Fatal("expected an error for a network without a named individual")
	}
}

func TestBootstrapTransactionRejectsMalformedEntityID(t *testing.T) {
	result := BootstrapResult{Network: "testnet", Topics: []TopicRecord{{TopicID: "topic-1"}}}
	if _, err := result.Transaction("tenant/dataset"); err == nil {
		t.Fatal("expected an error for an ID that is not shard.realm.num")
	}
}
