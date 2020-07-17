package types

import (
	"encoding/json"
	"fmt"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgAddDid        = "add-did"
	TypeMsgAddCredential = "add-credential"
)

var (
	_ IxoMsg = MsgAddDid{}
	_ IxoMsg = MsgAddCredential{}
)

type MsgAddDid struct {
	DidDoc BaseDidDoc `json:"didDoc" yaml:"didDoc"`
}

func (msg *MsgAddDid) UnmarshalJSON(bytes []byte) error {
	var msg2 struct {
		DidDoc BaseDidDoc `json:"didDoc" yaml:"didDoc"`
	}
	err := json.Unmarshal(bytes, &msg2)
	if err != nil {
		return err
	}

	if msg2.DidDoc.Credentials == nil {
		msg2.DidDoc.Credentials = []exported.DidCredential{}
	}

	*msg = msg2
	return nil
}

func NewMsgAddDid(did string, publicKey string) MsgAddDid {
	return MsgAddDid{
		DidDoc: NewBaseDidDoc(did, publicKey),
	}
}

func (msg MsgAddDid) Type() string            { return TypeMsgAddDid }
func (msg MsgAddDid) Route() string           { return RouterKey }
func (msg MsgAddDid) GetSignerDid() exported.Did { return msg.DidDoc.Did }
func (msg MsgAddDid) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{DidToAddr(msg.GetSignerDid())}
}
func (msg MsgAddDid) ValidateBasic() error {
	// Check that not empty
	if strings.TrimSpace(msg.DidDoc.Did) == "" {
		return exported.Invalid("did should not be empty")
	} else if strings.TrimSpace(msg.DidDoc.PubKey) == "" {
		return exported.InvalidPubKey("pubKey should not be empty")
	}

	// Check DidDoc credentials for empty fields
	for _, cred := range msg.DidDoc.Credentials {
		if strings.TrimSpace(cred.Issuer) == "" {
			return exported.InvalidIssuer( "issuer should not be empty")
		} else if strings.TrimSpace(cred.Claim.Id) == "" {
			return exported.Invalid( "claim id should not be empty")
		}
	}

	// Check that DID valid
	if !exported.IsValidDid(msg.DidDoc.Did) {
		return exported.Invalid( "did is invalid")
	}

	return nil
}

func (msg MsgAddDid) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}

func (msg MsgAddDid) String() string {
	return fmt.Sprintf("MsgAddDid{Did: %v, publicKey: %v}", string(msg.DidDoc.GetDid()), msg.DidDoc.GetPubKey())
}

type MsgAddCredential struct {
	DidCredential exported.DidCredential `json:"credential" yaml:"credential"`
}

func NewMsgAddCredential(did string, credType []string, issuer string, issued string) MsgAddCredential {
	didCredential := exported.DidCredential{
		CredType: credType,
		Issuer:   issuer,
		Issued:   issued,
		Claim: exported.Claim{
			Id:           did,
			KYCValidated: true,
		},
	}

	return MsgAddCredential{
		DidCredential: didCredential,
	}
}
func (msg MsgAddCredential) Type() string            { return TypeMsgAddCredential }
func (msg MsgAddCredential) Route() string           { return RouterKey }
func (msg MsgAddCredential) GetSignerDid() exported.Did { return msg.DidCredential.Issuer }
func (msg MsgAddCredential) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{DidToAddr(msg.GetSignerDid())}
}
func (msg MsgAddCredential) String() string {
	return fmt.Sprintf("MsgAddCredential{Did: %v, Type: %v, Signer: %v}",
		string(msg.DidCredential.Claim.Id), msg.DidCredential.CredType, string(msg.DidCredential.Issuer))
}
func (msg MsgAddCredential) ValidateBasic() error {
	// Check if empty
	if strings.TrimSpace(msg.DidCredential.Claim.Id) == "" {
		return exported.Invalid("claim id should not be empty")
	} else if strings.TrimSpace(msg.DidCredential.Issuer) == "" {
		return exported.InvalidIssuer( "issuer should not be empty")
	}
	// Check that DID valid
	if !exported.IsValidDid(msg.DidCredential.Issuer) {
		return exported.Invalid("issuer id is invalid")
	}
	return nil
}
func (msg MsgAddCredential) GetSignBytes() []byte {
	if bz, err := json.Marshal(msg); err != nil {
		panic(err)
	} else {
		return sdk.MustSortJSON(bz)
	}
}
