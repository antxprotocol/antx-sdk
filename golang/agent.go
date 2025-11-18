package sdk

import (
	"fmt"
	"strings"
	"time"

	agenttypes "github.com/antxprotocol/antx-proto/gen/go/antx/chain/agent"
	"github.com/antxprotocol/antx-sdk-golang/constants"

	"github.com/ethereum/go-ethereum/crypto"
)

func (c *AntxClient) BindAgent(ethPrivatekeyHex string, chainId string, expireTime uint64) (string, error) {
	ethPrivatekeyHex = strings.TrimPrefix(ethPrivatekeyHex, "0x")
	ethPrivateKey, err := crypto.HexToECDSA(ethPrivatekeyHex)
	if err != nil {
		return "", err
	}
	ethAddress := crypto.PubkeyToAddress(ethPrivateKey.PublicKey).Hex()
	agentAddress := c.agentAddress.String()
	createTime := uint64(time.Now().UnixMilli())
	expireTime = uint64(time.Now().Add(time.Duration(expireTime) * time.Second).UnixMilli())

	message := fmt.Sprintf("Action:BindAgent\nAgentAddress:%s\nCreateTime:%d\nExpireTime:%d\nChainId:%s",
		agentAddress, createTime, expireTime, chainId)

	// Sign message using personal_sign method
	signDoc := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	data := crypto.Keccak256([]byte(signDoc))
	signature, err := crypto.Sign(data, ethPrivateKey)
	if err != nil {
		return "", err
	}
	// Convert to hex string with 0x prefix
	ethSignature := fmt.Sprintf("0x%x", signature)

	msg := agenttypes.MsgBindAgent{
		AgentAddress:   agentAddress,
		ChainType:      agenttypes.ChainType_CHAIN_TYPE_EVM,
		ChainAddress:   ethAddress,
		CreateTime:     createTime,
		ExpireTime:     expireTime,
		ChainSignature: ethSignature,
	}

	txHash, err := c.signAndSendTx(constants.MsgBindAgentTypeURL, &msg, false)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
