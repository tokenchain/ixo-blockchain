package keeper

import (
	"fmt"
	errors "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

const (
	QueryOrder = "order"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err errors.Error) {
		switch path[0] {
		case QueryOrder:
			return queryOrder(ctx, path[1:], req, keeper)
		default:
			return nil,
				errors.ErrUnknownRequest("unknown token name query endpoint")
		}
	}
}

// queryOrder is a query function to get order by order ID.
func queryOrder(
	ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper,
) ([]byte, errors.Error) {
	if len(path) == 0 {
		return nil,
			errors.ErrUnknownRequest("must specify the order id")
	}
	id, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		return nil,
			errors.ErrUnknownRequest(fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	order, err := keeper.GetOrder(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	return keeper.cdc.MustMarshalJSON(order), nil
}