package transactionsigner

import (
	"fmt"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	sighash "github.com/bitcoin-sv/go-sdk/transaction/sighash"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type TransactionSigner struct {
	xPriv *bip32.ExtendedKey
}

func New(xPriv *bip32.ExtendedKey) *TransactionSigner {
	return &TransactionSigner{
		xPriv: xPriv,
	}
}

func (ts *TransactionSigner) GetSignedHex(dt *response.DraftTransaction) (string, error) {
	// Create transaction from hex
	tx, err := trx.NewTransactionFromHex(dt.Hex)
	if err != nil {
		return "", fmt.Errorf("failed to parse hex, %w", err)
	}
	// we need to reset the inputs as we are going to add them via tx.AddInputFrom (ts-sdk method) and then sign
	tx.Inputs = make([]*trx.TransactionInput, 0)

	// Enrich inputs
	for _, draftInput := range dt.Configuration.Inputs {
		lockingScript, err := prepareLockingScript(&draftInput.Destination)
		if err != nil {
			return "", fmt.Errorf("failed to prepare locking script, %w", err)
		}

		unlockScript, err := prepareUnlockingScript(ts.xPriv, &draftInput.Destination)
		if err != nil {
			return "", fmt.Errorf("failed to prepare unlocking script, %w", err)
		}

		err = tx.AddInputFrom(draftInput.TransactionID, draftInput.OutputIndex, lockingScript.String(), draftInput.Satoshis, unlockScript)
		if err != nil {
			return "", fmt.Errorf("failed to add inputs to transaction, %w", err)
		}
	}

	err = tx.Sign()
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction, %w", err)
	}

	return tx.String(), nil
}

func prepareLockingScript(dst *response.Destination) (*script.Script, error) {
	lockingScript, err := script.NewFromHex(dst.LockingScript)
	if err != nil {
		return nil, fmt.Errorf("failed to create locking script from hex for destination: %w", err)
	}

	return lockingScript, nil
}

func prepareUnlockingScript(xPriv *bip32.ExtendedKey, dst *response.Destination) (*p2pkh.P2PKH, error) {
	key, err := getDerivedKeyForDestination(xPriv, dst)
	if err != nil {
		return nil, fmt.Errorf("failed to get derived key for destination: %w", err)
	}

	return getUnlockingScript(key)
}

func getDerivedKeyForDestination(xPriv *bip32.ExtendedKey, dst *response.Destination) (*ec.PrivateKey, error) {
	// Derive the child key (m/chain/num)
	derivedKey, err := bip32.GetHDKeyByPath(xPriv, dst.Chain, dst.Num)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key for unlocking input, %w", err)
	}

	// Handle paymail destination derivation if applicable
	if dst.PaymailExternalDerivationNum != nil {
		derivedKey, err = derivedKey.Child(*dst.PaymailExternalDerivationNum)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key for unlocking paymail input, %w", err)
		}
	}

	// Get the private key from the derived key
	priv, err := bip32.GetPrivateKeyFromHDKey(derivedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key for unlocking paymail input, %w", err)
	}

	return priv, nil
}

func getUnlockingScript(privateKey *ec.PrivateKey) (*p2pkh.P2PKH, error) {
	sigHashFlags := sighash.AllForkID
	unlocked, err := p2pkh.Unlock(privateKey, &sigHashFlags)
	if err != nil {
		return nil, fmt.Errorf("failed to create unlocking script, %w", err)
	}

	return unlocked, nil
}
