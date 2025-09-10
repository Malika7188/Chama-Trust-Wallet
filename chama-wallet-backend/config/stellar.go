package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

type StellarConfig struct {
	Network            string
	HorizonURL         string
	SorobanRPCURL      string
	NetworkPassphrase  string
	ContractID         string
	IsMainnet          bool
	USDCAssetCode      string
	USDCAssetIssuer    string
}

var Config *StellarConfig

func InitStellarConfig() {
	stellarNetwork := strings.ToLower(os.Getenv("STELLAR_NETWORK"))
	if stellarNetwork == "" {
		stellarNetwork = "testnet" // Default to testnet for safety
	}

	isMainnet := stellarNetwork == "mainnet"

	config := &StellarConfig{
		Network:           stellarNetwork,
		IsMainnet:         isMainnet,
		ContractID:        os.Getenv("SOROBAN_CONTRACT_ID"),
		USDCAssetCode:     os.Getenv("USDC_ASSET_CODE"),
		USDCAssetIssuer:   os.Getenv("USDC_ASSET_ISSUER"),
	}

	if isMainnet {
		// Mainnet configuration
		config.HorizonURL = getEnvOrDefault("STELLAR_HORIZON_URL", "https://horizon.stellar.org")
		config.SorobanRPCURL = getEnvOrDefault("STELLAR_SOROBAN_RPC_URL", "https://soroban-rpc.mainnet.stellar.org:443")
		config.NetworkPassphrase = getEnvOrDefault("STELLAR_NETWORK_PASSPHRASE", "Public Global Stellar Network ; September 2015")
		
		fmt.Println("🌐 Stellar Mainnet Configuration Loaded")
		fmt.Printf("   Horizon: %s\n", config.HorizonURL)
		fmt.Printf("   Soroban RPC: %s\n", config.SorobanRPCURL)
		fmt.Printf("   Contract ID: %s\n", config.ContractID)
	} else {
		// Testnet configuration
		config.HorizonURL = getEnvOrDefault("STELLAR_HORIZON_URL", "https://horizon-testnet.stellar.org")
		config.SorobanRPCURL = getEnvOrDefault("STELLAR_SOROBAN_RPC_URL", "https://soroban-testnet.stellar.org:443")
		config.NetworkPassphrase = getEnvOrDefault("STELLAR_NETWORK_PASSPHRASE", "Test SDF Network ; September 2015")
		
		fmt.Println("🧪 Stellar Testnet Configuration Loaded")
		fmt.Printf("   Horizon: %s\n", config.HorizonURL)
		fmt.Printf("   Soroban RPC: %s\n", config.SorobanRPCURL)
		fmt.Printf("   Contract ID: %s\n", config.ContractID)
	}

	Config = config
}

func GetHorizonClient() *horizonclient.Client {
	if Config.IsMainnet {
		return &horizonclient.Client{
			HorizonURL: Config.HorizonURL,
		}
	}
	return horizonclient.DefaultTestNetClient
}

func GetNetworkPassphrase() string {
	if Config.IsMainnet {
		return network.PublicNetworkPassphrase
	}
	return network.TestNetworkPassphrase
}
