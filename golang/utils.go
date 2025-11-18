package sdk

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/zeromicro/go-zero/core/logx"
)

func ConvertToAntxAddr(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("addr can't be empty")
	}
	conf := sdk.GetConfig()
	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	case strings.HasPrefix(addrString, conf.GetBech32ValidatorAddrPrefix()):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, conf.GetBech32AccountAddrPrefix()):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'",
			conf.GetBech32AccountAddrPrefix(), addrString)
	}

	return sdk.AccAddress(addr).String(), nil
}

func ConvertToEthAddr(addrString string) (string, error) {
	if addrString == "" {
		return "", fmt.Errorf("addr can't be empty")
	}
	conf := sdk.GetConfig()
	var addr []byte
	switch {
	case common.IsHexAddress(addrString):
		addr = common.HexToAddress(addrString).Bytes()
	case strings.HasPrefix(addrString, conf.GetBech32ValidatorAddrPrefix()):
		addr, _ = sdk.ValAddressFromBech32(addrString)
	case strings.HasPrefix(addrString, conf.GetBech32AccountAddrPrefix()):
		addr, _ = sdk.AccAddressFromBech32(addrString)
	default:
		return "", fmt.Errorf("expected a valid hex or bech32 address (acc prefix %s), got '%s'",
			conf.GetBech32AccountAddrPrefix(), addrString)
	}

	return common.BytesToAddress(addr).Hex(), nil
}

func VerifyEthPersonalSignature(address string, data []byte, sig []byte) bool {
	sigHash, _ := accounts.TextAndHash(data)

	if sig[64] > 1 {
		sig[64] -= 27
	}

	sigPublicKey, err := ethCrypto.SigToPub(sigHash, sig)
	if err != nil {
		logx.Errorf("invalid signature crypto.SigToPub err: %v", err)

		return false
	}

	recoverAddress := ethCrypto.PubkeyToAddress(*sigPublicKey)
	return recoverAddress == common.HexToAddress(address)
}
