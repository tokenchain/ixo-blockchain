package dap

import (
	"github.com/tokenchain/ixo-blockchain/x/dap/types"
)

const (
	IxoNativeToken = types.NativeToken

)

var (
	// Auth

	//ApproximateFeeForTx     = auth.ApproximateFeeForTx
	//	GenerateOrBroadcastMsgs = auth.GenerateOrBroadcastMsgs
	//CompleteAndBroadcastTxRest       = auth.CompleteAndBroadcastTxRest
	//SignAndBroadcastTxFromStdSignMsg = auth.SignAndBroadcastTxFromStdSignMsg



	//SignAndBroadcastTxRest = auth.SignAndBroadcastTxRest

	// Types
	//IxoDecimals = types.IxoDecimals
)

/*
func SignIxoMessage(signBytes []byte, privKey [64]byte) types.IxoSignature {
	signatureBytes := ed25519.Sign(&privKey, signBytes)
	signature := *signatureBytes
	return types.NewSignature(time.Now(), signature)
}

func VerifySignature(signBytes []byte, publicKey [32]byte, sig types.IxoSignature) bool {
	result := ed25519.Verify(&publicKey, signBytes, &sig.SignatureValue)
	if !result {
		fmt.Println("******* VERIFY_MSG: Failed ******* ")
	}
	return result
}

func LookupEnv(name string, defaultValue string) string {
	val, found := os.LookupEnv(name)
	if found && len(val) > 0 {
		return val
	}
	return defaultValue
}
*/