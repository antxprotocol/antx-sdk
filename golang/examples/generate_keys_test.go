package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestGenerateKeys(t *testing.T) {
	fmt.Println("=== Antx SDK Key Generator ===")

	// Generate ethPrivateKey
	ethPrivKey, err := ethCrypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate eth private key: %v", err)
	}
	ethPrivKeyHex := hex.EncodeToString(ethCrypto.FromECDSA(ethPrivKey))
	ethAddress := ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey).Hex()

	fmt.Println("\n--- ETH Private Key ---")
	fmt.Printf("ETH_PRIVATE_KEY=%s\n", ethPrivKeyHex)
	fmt.Printf("ETH_ADDRESS=%s\n", ethAddress)

	// Generate agentPrivateKey
	agentKeyBytes := make([]byte, 32)
	if _, err := rand.Read(agentKeyBytes); err != nil {
		t.Fatalf("Failed to generate agent private key: %v", err)
	}
	agentPrivKey := &secp256k1.PrivKey{Key: agentKeyBytes}
	agentPrivKeyHex := hex.EncodeToString(agentKeyBytes)
	agentAddress := sdk.AccAddress(agentPrivKey.PubKey().Address()).String()

	fmt.Println("\n--- Agent Private Key ---")
	fmt.Printf("AGENT_PRIVATE_KEY=%s\n", agentPrivKeyHex)
	fmt.Printf("AGENT_ADDRESS=%s\n", agentAddress)

	fmt.Println("\n=== Done ===")
	fmt.Println("Export these as environment variables before running examples:")
	fmt.Printf("  export ETH_PRIVATE_KEY=%s\n", ethPrivKeyHex)
	fmt.Printf("  export ETH_ADDRESS=%s\n", ethAddress)
	fmt.Printf("  export AGENT_PRIVATE_KEY=%s\n", agentPrivKeyHex)
}
