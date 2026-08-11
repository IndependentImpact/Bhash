package hedera

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// SDKNetwork provides a Network implementation backed by the Hedera Go SDK.
type SDKNetwork struct {
	client      *sdk.Client
	networkName string
}

// NewSDKNetwork initialises a Hedera SDK client using the provided configuration.
func NewSDKNetwork(cfg Config) (*SDKNetwork, error) {
	client, name, err := buildClient(cfg)
	if err != nil {
		return nil, err
	}
	return &SDKNetwork{client: client, networkName: name}, nil
}

func buildClient(cfg Config) (*sdk.Client, string, error) {
	var (
		client *sdk.Client
		name   = strings.ToLower(strings.TrimSpace(cfg.Network))
	)
	switch name {
	case "", "testnet":
		client = sdk.ClientForTestnet()
		name = "testnet"
	case "mainnet":
		client = sdk.ClientForMainnet()
	case "previewnet":
		client = sdk.ClientForPreviewnet()
	default:
		return nil, "", fmt.Errorf("unsupported hedera network %q", cfg.Network)
	}
	if cfg.MirrorNetworkURL != "" {
		client.SetMirrorNetwork([]string{cfg.MirrorNetworkURL})
	}
	if cfg.OperatorAccountID != "" {
		accountID, err := sdk.AccountIDFromString(cfg.OperatorAccountID)
		if err != nil {
			return nil, "", fmt.Errorf("parse operator account id: %w", err)
		}
		key, err := sdk.PrivateKeyFromString(cfg.OperatorPrivateKey)
		if err != nil {
			return nil, "", fmt.Errorf("parse operator private key: %w", err)
		}
		client.SetOperator(accountID, key)
	}
	return client, name, nil
}

func (s *SDKNetwork) CreateAccount(ctx context.Context, spec AccountSpec) (AccountRecord, error) {
	tx := sdk.NewAccountCreateTransaction()
	if spec.PublicKey != "" {
		key, err := sdk.PublicKeyFromString(spec.PublicKey)
		if err != nil {
			return AccountRecord{}, fmt.Errorf("parse public key: %w", err)
		}
		tx.SetKey(key)
	}
	if spec.InitialBalanceTinybar > 0 {
		tx.SetInitialBalance(sdk.HbarFromTinybar(spec.InitialBalanceTinybar))
	}
	if spec.Memo != "" {
		tx.SetAccountMemo(spec.Memo)
	}
	resp, err := tx.Execute(s.client)
	if err != nil {
		return AccountRecord{}, fmt.Errorf("execute account create: %w", err)
	}
	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return AccountRecord{}, fmt.Errorf("fetch account receipt: %w", err)
	}
	record, err := resp.GetRecord(s.client)
	if err != nil {
		return AccountRecord{}, fmt.Errorf("fetch account record: %w", err)
	}
	return AccountRecord{
		Alias:     spec.Alias,
		AccountID: receipt.AccountID.String(),
		PublicKey: spec.PublicKey,
		Memo:      spec.Memo,
		Tags:      append([]string(nil), spec.Tags...),
		CreatedAt: record.ConsensusTimestamp.UTC(),
	}, nil
}

func (s *SDKNetwork) CreateTopic(ctx context.Context, spec TopicSpec) (TopicRecord, error) {
	tx := sdk.NewTopicCreateTransaction()
	if spec.Memo != "" {
		tx.SetTopicMemo(spec.Memo)
	}
	if spec.AdminKey != "" {
		key, err := sdk.PublicKeyFromString(spec.AdminKey)
		if err != nil {
			return TopicRecord{}, fmt.Errorf("parse admin key: %w", err)
		}
		tx.SetAdminKey(key)
	}
	if spec.SubmitKey != "" {
		key, err := sdk.PublicKeyFromString(spec.SubmitKey)
		if err != nil {
			return TopicRecord{}, fmt.Errorf("parse submit key: %w", err)
		}
		tx.SetSubmitKey(key)
	}
	resp, err := tx.Execute(s.client)
	if err != nil {
		return TopicRecord{}, fmt.Errorf("execute topic create: %w", err)
	}
	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return TopicRecord{}, fmt.Errorf("fetch topic receipt: %w", err)
	}
	record, err := resp.GetRecord(s.client)
	if err != nil {
		return TopicRecord{}, fmt.Errorf("fetch topic record: %w", err)
	}
	return TopicRecord{
		Alias:     spec.Alias,
		TopicID:   receipt.TopicID.String(),
		Memo:      spec.Memo,
		Tags:      append([]string(nil), spec.Tags...),
		CreatedAt: record.ConsensusTimestamp.UTC(),
	}, nil
}

func (s *SDKNetwork) CreateToken(ctx context.Context, spec TokenSpec) (TokenRecord, error) {
	if strings.TrimSpace(spec.TreasuryAccountID) == "" {
		return TokenRecord{}, fmt.Errorf("treasury account id is required")
	}
	treasury, err := sdk.AccountIDFromString(spec.TreasuryAccountID)
	if err != nil {
		return TokenRecord{}, fmt.Errorf("parse treasury account id: %w", err)
	}
	tx := sdk.NewTokenCreateTransaction().
		SetTokenName(spec.Name).
		SetTokenSymbol(spec.Symbol).
		SetTokenMemo(spec.Memo).
		SetTreasuryAccountID(treasury).
		SetDecimals(spec.Decimals).
		SetInitialSupply(spec.InitialSupply)
	if spec.MaxSupply != 0 {
		tx.SetMaxSupply(spec.MaxSupply)
	}
	if spec.SupplyType != "" {
		supplyType, err := parseSupplyType(spec.SupplyType)
		if err != nil {
			return TokenRecord{}, err
		}
		tx.SetSupplyType(supplyType)
	}
	if spec.TokenType != "" {
		tokenType, err := parseTokenType(spec.TokenType)
		if err != nil {
			return TokenRecord{}, err
		}
		tx.SetTokenType(tokenType)
	}
	if spec.AdminKey != "" {
		key, err := sdk.PublicKeyFromString(spec.AdminKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse admin key: %w", err)
		}
		tx.SetAdminKey(key)
	}
	if spec.SupplyKey != "" {
		key, err := sdk.PublicKeyFromString(spec.SupplyKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse supply key: %w", err)
		}
		tx.SetSupplyKey(key)
	}
	if spec.KYCKey != "" {
		key, err := sdk.PublicKeyFromString(spec.KYCKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse kyc key: %w", err)
		}
		tx.SetKycKey(key)
	}
	if spec.FreezeKey != "" {
		key, err := sdk.PublicKeyFromString(spec.FreezeKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse freeze key: %w", err)
		}
		tx.SetFreezeKey(key)
	}
	if spec.WipeKey != "" {
		key, err := sdk.PublicKeyFromString(spec.WipeKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse wipe key: %w", err)
		}
		tx.SetWipeKey(key)
	}
	if spec.PauseKey != "" {
		key, err := sdk.PublicKeyFromString(spec.PauseKey)
		if err != nil {
			return TokenRecord{}, fmt.Errorf("parse pause key: %w", err)
		}
		tx.SetPauseKey(key)
	}
	if spec.FreezeDefault != nil {
		tx.SetFreezeDefault(*spec.FreezeDefault)
	}
	resp, err := tx.Execute(s.client)
	if err != nil {
		return TokenRecord{}, fmt.Errorf("execute token create: %w", err)
	}
	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return TokenRecord{}, fmt.Errorf("fetch token receipt: %w", err)
	}
	record, err := resp.GetRecord(s.client)
	if err != nil {
		return TokenRecord{}, fmt.Errorf("fetch token record: %w", err)
	}
	return TokenRecord{
		Alias:             spec.Alias,
		TokenID:           receipt.TokenID.String(),
		Name:              spec.Name,
		Symbol:            spec.Symbol,
		Memo:              spec.Memo,
		TreasuryAccountID: spec.TreasuryAccountID,
		Decimals:          spec.Decimals,
		InitialSupply:     spec.InitialSupply,
		MaxSupply:         spec.MaxSupply,
		SupplyType:        spec.SupplyType,
		TokenType:         spec.TokenType,
		Tags:              append([]string(nil), spec.Tags...),
		CreatedAt:         record.ConsensusTimestamp.UTC(),
	}, nil
}

func parseTokenType(value string) (sdk.TokenType, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "", "FUNGIBLE_COMMON", "TOKEN_TYPE_FUNGIBLE_COMMON":
		return sdk.TokenTypeFungibleCommon, nil
	case "NON_FUNGIBLE_UNIQUE", "TOKEN_TYPE_NON_FUNGIBLE_UNIQUE":
		return sdk.TokenTypeNonFungibleUnique, nil
	default:
		return 0, fmt.Errorf("unsupported token type %q", value)
	}
}

func parseSupplyType(value string) (sdk.TokenSupplyType, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "", "INFINITE", "TOKEN_SUPPLY_TYPE_INFINITE":
		return sdk.TokenSupplyTypeInfinite, nil
	case "FINITE", "TOKEN_SUPPLY_TYPE_FINITE":
		return sdk.TokenSupplyTypeFinite, nil
	default:
		return 0, fmt.Errorf("unsupported token supply type %q", value)
	}
}

// Close releases network resources held by the SDK client.
func (s *SDKNetwork) Close() {
	if s.client != nil {
		_ = s.client.Close()
	}
}

// Client exposes the underlying SDK client for integration scenarios.
func (s *SDKNetwork) Client() *sdk.Client {
	return s.client
}

// NetworkName returns the configured Hedera network name.
func (s *SDKNetwork) NetworkName() string {
	return s.networkName
}
