package exported

import (
	"bufio"
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/spf13/cobra"
	"io"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/go-bip39"
	"golang.org/x/crypto/ed25519"
	naclBox "golang.org/x/crypto/nacl/box"
)

const (
	flagUserEntropy     = "unsafe-entropy"
	mnemonicEntropySize = 256
)

type SovrinSecret struct {
	Seed                 string `json:"seed" yaml:"seed"`
	SignKey              string `json:"signKey" yaml:"signKey"`
	EncryptionPrivateKey string `json:"encryptionPrivateKey" yaml:"encryptionPrivateKey"`
}

func (ss SovrinSecret) String() string {
	output, err := json.MarshalIndent(ss, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

type SovrinDid struct {
	Did                 string       `json:"did" yaml:"did"`
	VerifyKey           string       `json:"verifyKey" yaml:"verifyKey"`
	EncryptionPublicKey string       `json:"encryptionPublicKey" yaml:"encryptionPublicKey"`
	Secret              SovrinSecret `json:"secret" yaml:"secret"`
}

func (sd SovrinDid) String() string {
	output, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%v", string(output))
}

func GenerateMnemonic() string {
	entropy, _ := bip39.NewEntropy(12)
	mnemonicWords, _ := bip39.NewMnemonic(entropy)
	return mnemonicWords
}

func fromJsonString(jsonSovrinDid string) (IxoDid, error) {
	var did IxoDid
	err := json.Unmarshal([]byte(jsonSovrinDid), &did)
	if err != nil {
		err := fmt.Errorf("Could not unmarshal did into struct. Error: %s", err.Error())
		return IxoDid{}, err
	}

	return did, nil
}

func mnemonicToDid(mnemonic string) IxoDid {
	seed := sha256.New()
	seed.Write([]byte(mnemonic))
	var seed32 [32]byte
	copy(seed32[:], seed.Sum(nil)[:32])
	return fromSeedToDid(seed32)
}

func dxpDidAddress(document string) string {
	return fmt.Sprintf("did:dxp:%s", document)
}

func fromSeedToDid(seed [32]byte) IxoDid {

	publicKeyBytes, privateKeyBytes, err := ed25519.GenerateKey(bytes.NewReader(seed[0:32]))
	if err != nil {
		panic(err)
	}
	publicKey := []byte(publicKeyBytes)
	privateKey := []byte(privateKeyBytes)

	signKey := base58.Encode(privateKey[:32])
	keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))

	sovDid := IxoDid{
		Did:                 dxpDidAddress(base58.Encode(publicKey[:16])),
		VerifyKey:           base58.Encode(publicKey),
		EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

		Secret: Secret{
			Seed:                 hex.EncodeToString(seed[0:32]),
			SignKey:              signKey,
			EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
		},
	}

	return sovDid
}

/*
func Gen() IxoDid {
	var seed [32]byte
	if _, err := io.ReadFull(cryptoRand.Reader, seed[:]); err != nil {
		panic(err)
	}
	did, _ := fromJsonString(seed)
	return did
}
*/
func SignMessage(message []byte, signKey string, verifyKey string) []byte {
	// Force the length to 64
	privateKey := make([]byte, ed25519.PrivateKeySize)
	fullPrivKey := ed25519.PrivateKey(privateKey)
	copy(fullPrivKey[:], getArrayFromKey(signKey))
	copy(fullPrivKey[32:], getArrayFromKey(verifyKey))

	return ed25519.Sign(fullPrivKey, message)
}

func VerifySignedMessage(message []byte, signature []byte, verifyKey string) bool {
	publicKey := ed25519.PublicKey{}
	copy(publicKey[:], getArrayFromKey(verifyKey))
	result := ed25519.Verify(publicKey, message, signature)

	return result
}

func GetNonce() [24]byte {
	var nonce [24]byte
	if _, err := io.ReadFull(cryptoRand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func getArrayFromKey(key string) []byte {
	return base58.Decode(key)
}

func GetKeyPairFromSignKey(signKey string) ([32]byte, [32]byte) {
	publicKey, privateKey, err := naclBox.GenerateKey(bytes.NewReader(getArrayFromKey(signKey)))
	if err != nil {
		panic(err)
	}
	return *publicKey, *privateKey
}

func RunMnemonicCmd(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	userEntropy, _ := flags.GetBool(flagUserEntropy)
	var entropySeed []byte
	if userEntropy {
		// prompt the user to enter some entropy
		buf := bufio.NewReader(cmd.InOrStdin())
		inputEntropy, err := input.GetString("> WARNING: Generate at least 256-bits of entropy and enter the results here:", buf)
		if err != nil {
			return err
		}
		if len(inputEntropy) < 43 {
			return fmt.Errorf("256-bits is 43 characters in Base-64, and 100 in Base-6. You entered %v, and probably want more", len(inputEntropy))
		}
		conf, err := input.GetConfirmation(fmt.Sprintf("> Input length: %d", len(inputEntropy)), buf)
		if err != nil {
			return err
		}
		if !conf {
			return nil
		}

		// hash input entropy to get entropy seed
		hashedEntropy := sha256.Sum256([]byte(inputEntropy))
		entropySeed = hashedEntropy[:]
	} else {
		// read entropy seed straight from crypto.Rand
		var err error
		entropySeed, err = bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			return err
		}
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return err
	}
	//	cmd.Println(mnemonic)
	cmd.Println("======= The passphrase please keep in the secured place:")
	cmd.Println(mnemonic)
	cmd.Println("===========================================================================================")
	did_document := mnemonicToDid(mnemonic)

	cmd.Println("======= DID account address")
	cmd.Println(did_document.DidAddress())

	cmd.Println("======= generated a new DID document with the above passphrase")
	cmd.Println(did_document.String())

	return nil
}
