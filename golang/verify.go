package sdk

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	txsigning "cosmossdk.io/x/tx/signing"
	"google.golang.org/protobuf/types/known/anypb"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// verifySignatureComplete performs complete signature verification
// Verifies that the signature actually corresponds to this transaction
func verifySignatureComplete(sig signing.SignatureV2, signer []byte, chainID string, accountNumber uint64,
	signModeHandler *txsigning.HandlerMap, tx sdk.Tx) error {

	// Check if transaction implements V2AdaptableTx interface
	adaptableTx, ok := tx.(authsigning.V2AdaptableTx)
	if !ok {
		return fmt.Errorf("expected tx to implement V2AdaptableTx, got %T", tx)
	}

	// Create signer data
	anyPk, err := codectypes.NewAnyWithValue(sig.PubKey)
	if err != nil {
		return fmt.Errorf("failed to pack public key: %v", err)
	}

	// Get transaction data
	txData := adaptableTx.GetSigningTxData()

	signerData := txsigning.SignerData{
		Address:       sdk.AccAddress(signer).String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sig.Sequence,
		PubKey: &anypb.Any{
			TypeUrl: anyPk.TypeUrl,
			Value:   anyPk.Value,
		},
	}

	err = authsigning.VerifySignature(context.Background(), sig.PubKey, signerData, sig.Data, signModeHandler, txData)
	if err == nil {
		// Verification successful, return nil
		return nil
	}

	// If all common accountNumbers fail, return the last error
	return fmt.Errorf("signature verification failed with accountNumbers: %v, err: %v", accountNumber, err)
}

// VerifyTransactionSignature performs complete signature verification
// Verifies signature format, public key, and whether the signature actually corresponds to this transaction
func VerifyTransactionSignature(tx sdk.Tx, chainID string, accountNumber uint64, signModeHandler *txsigning.HandlerMap) error {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, "transaction does not implement SigVerifiableTx interface")
	}

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrTxDecode, fmt.Sprintf("failed to get signatures: %v", err))
	}

	for i, sig := range sigs {
		signerAddr := sig.PubKey.Address().Bytes()

		// Complete signature verification
		if err := verifySignatureComplete(sig, signerAddr, chainID, accountNumber, signModeHandler, tx); err != nil {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized,
				fmt.Sprintf("signature verification failed for signer %d: %v", i, err))
		}
	}

	return nil
}
