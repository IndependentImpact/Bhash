package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashgraph/bhash/internal/fluree"
	bhedera "github.com/hashgraph/bhash/internal/hedera"
)

type hederaNetworkFactoryFunc func(bhedera.Config, bool) (bhedera.Network, func(), error)
type flureeClientFactoryFunc func(fluree.Config) flureeTransactor

type flureeTransactor interface {
	Transact(context.Context, fluree.TransactionRequest) (any, error)
}

var (
	hederaNetworkFactory hederaNetworkFactoryFunc = defaultHederaNetworkFactory
	flureeClientFactory  flureeClientFactoryFunc  = defaultFlureeClientFactory
)

func defaultHederaNetworkFactory(cfg bhedera.Config, simulate bool) (bhedera.Network, func(), error) {
	if simulate {
		return bhedera.NewMockNetwork(cfg.Network), func() {}, nil
	}
	sdk, err := bhedera.NewSDKNetwork(cfg)
	if err != nil {
		return nil, nil, err
	}
	return sdk, sdk.Close, nil
}

func defaultFlureeClientFactory(cfg fluree.Config) flureeTransactor {
	return fluree.NewClient(cfg, nil)
}

func runHedera(args []string) {
	if len(args) == 0 {
		hederaUsage()
		os.Exit(1)
	}
	switch args[0] {
	case "bootstrap":
		runHederaBootstrap(args[1:])
	default:
		hederaUsage()
		os.Exit(1)
	}
}

func hederaUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s hedera <bootstrap> [options]\n", filepath.Base(os.Args[0]))
}

func runHederaBootstrap(args []string) {
	fs := flag.NewFlagSet("hedera bootstrap", flag.ExitOnError)
	specPath := fs.String("spec", "", "Path to bootstrap spec JSON")
	ledger := fs.String("ledger", "", "Fluree ledger identifier (owner/dataset)")
	simulate := fs.Bool("simulate", true, "Use the deterministic mock Hedera network")
	commit := fs.Bool("commit", false, "Submit the generated transaction to Fluree")
	networkOverride := fs.String("network", "", "Hedera network (overrides $HEDERA_NETWORK)")
	operatorID := fs.String("operator-id", "", "Hedera operator account ID")
	operatorKey := fs.String("operator-key", "", "Hedera operator private key")
	mirrorURL := fs.String("mirror-url", "", "Hedera mirror network URL")
	apiToken := fs.String("api-token", "", "Fluree API token (defaults to $FLUREE_API_TOKEN)")
	tenant := fs.String("tenant", "", "Fluree tenant handle (defaults to $FLUREE_HANDLE)")
	baseURL := fs.String("base-url", "", "Fluree API base URL (defaults to $FLUREE_BASE_URL)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *specPath == "" {
		fmt.Fprintln(os.Stderr, "spec is required")
		os.Exit(1)
	}

	spec, err := bhedera.LoadBootstrapSpec(*specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	ledgerID := strings.TrimSpace(*ledger)
	if ledgerID == "" {
		ledgerID = strings.TrimSpace(spec.Ledger)
	}
	if ledgerID == "" {
		fmt.Fprintln(os.Stderr, "ledger is required")
		os.Exit(1)
	}

	cfg, err := bhedera.EnvConfigFromLookup(func(key string) (string, bool) {
		value, ok := os.LookupEnv(key)
		return value, ok
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	cfg = cfg.WithOverrides(*networkOverride, *operatorID, *operatorKey, *mirrorURL)
	if spec.Network != "" {
		cfg.Network = spec.Network
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if !*simulate && !cfg.HasOperator() {
		fmt.Fprintln(os.Stderr, "operator credentials are required when simulate=false")
		os.Exit(1)
	}

	network, closer, err := hederaNetworkFactory(cfg, *simulate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if closer != nil {
		defer closer()
	}

	bootstrapper := bhedera.NewBootstrapper(network, cfg.Network)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	result, err := bootstrapper.Execute(ctx, spec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	transaction, err := result.Transaction(ledgerID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	output := map[string]any{
		"network":     result.Network,
		"ledger":      ledgerID,
		"accounts":    result.Accounts,
		"topics":      result.Topics,
		"tokens":      result.Tokens,
		"transaction": transaction,
	}

	if *commit {
		cfg := mustFlureeConfig(*apiToken, *tenant, *baseURL)
		client := flureeClientFactory(cfg)
		respCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		resp, err := client.Transact(respCtx, transaction)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		output["flureeResponse"] = resp
	}

	printJSON(output)
}
