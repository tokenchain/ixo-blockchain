# Key generator
```
package exported

import (
        "bytes"
        "crypto/sha256"
        "encoding/hex"
        "fmt"
        "github.com/btcsuite/btcutil/base58"
        "github.com/cosmos/cosmos-sdk/crypto/keys"
        sdk "github.com/cosmos/cosmos-sdk/types"
        "github.com/cosmos/go-bip39"
        ed25519tm "github.com/tendermint/tendermint/crypto/ed25519"
        edgen "github.com/tokenchain/dp-hub/x/did/ed25519"
        naclBox "golang.org/x/crypto/nacl/box"
)

func NewDidGeneratorBuilder() KeyGenerator {
        return KeyGenerator{
                name: "cosmos",
                mem:  "",
        }
}
func (s KeyGenerator) GetMnemonicString() string {
        return s.mem
}
func (s KeyGenerator) GetSeedString() string {
        return hex.EncodeToString(s.seed[0:32])
}
func (s KeyGenerator) WithName(n string) KeyGenerator {
        s.name = n
        return s
}

func (s KeyGenerator) WithPubKey(n []byte) KeyGenerator {
        s.pubkey = n
        return s
}

func (s KeyGenerator) WithPrivKey(n []byte) KeyGenerator {
        s.privkey = n
        return s
}

func (s KeyGenerator) WithMem(n string) KeyGenerator {
        s.mem = n
        return s
}
func (s KeyGenerator) WithSeed(seed32 [32]byte) KeyGenerator {
        s.seed = seed32
        return s
}
func (s KeyGenerator) generateSeed() KeyGenerator {
        seed := sha256.New()
        seed.Write([]byte(s.mem))
        var seed32 [32]byte
        copy(seed32[:], seed.Sum(nil)[:32])
        s.seed = seed32
        return s
}
func (s KeyGenerator) generateMnemonic() KeyGenerator {
        entropy, _ := bip39.NewEntropy(12)
        mnemonicWords, _ := bip39.NewMnemonic(entropy)
        s.mem = mnemonicWords
        return s
}
func (s KeyGenerator) generateFinal() IxoDid {
        publicKeyBytes, privateKeyBytes, err := edgen.GenerateKey(bytes.NewReader(s.seed[0:32]))
        if err != nil {
                panic(err)
        }
        //head part
        signKey := base58.Encode(privateKeyBytes[:32])
        //keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKey[:]))
        keyPairPublicKey, keyPairPrivateKey, err := naclBox.GenerateKey(bytes.NewReader(privateKeyBytes[:]))

        var pubKey ed25519tm.PubKeyEd25519
        copy(pubKey[:], publicKeyBytes[:])

        sovDid := IxoDid{
                Did:                 dxpDidAddress(base58.Encode(publicKeyBytes[:16])),
                VerifyKey:           base58.Encode(publicKeyBytes[:]),
                EncryptionPublicKey: base58.Encode(keyPairPublicKey[:]),

                Secret: Secret{
                        Seed:                 hex.EncodeToString(s.seed[0:32]),
                        SignKey:              signKey,
                        EncryptionPrivateKey: base58.Encode(keyPairPrivateKey[:]),
                },

                Dpinfo: DpInfo{
                        DpAddress: sdk.AccAddress(pubKey.Address()).String(),
                        PubKey:    sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey),
                        Name:      s.name,
                        Algo:      keys.Ed25519,
                },
        }
        return sovDid
}

func (s KeyGenerator) Build() IxoDid {
        fmt.Println(s.mem)
        if s.mem == "" {
                return s.generateMnemonic().generateSeed().generateFinal()
        } else {
                return s.generateSeed().generateFinal()
        }
}

func (s KeyGenerator) BuildWithCustomSeed(seed32 [32]byte) IxoDid {
        return s.WithSeed(seed32).generateFinal()
}

func (s KeyGenerator) Recover(mem string) IxoDid {
        return s.WithMem(mem).generateSeed().generateFinal()
}

```


#### API Doc
Please find the api document to be located at `*:1317/swagger-ui/` at the LCD.
