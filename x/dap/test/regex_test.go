package test

import (
	"github.com/stretchr/testify/require"
	"github.com/tokenchain/ixo-blockchain/x/did/exported"
	"testing"
)

func TestDXP(t *testing.T){
	require.True(t, true, exported.IsValidDid(sample_did_2))
}