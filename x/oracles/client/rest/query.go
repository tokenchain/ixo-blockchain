package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tokenchain/dp-hub/client/utils"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/tokenchain/dp-hub/x/oracles/internal/keeper"
	"github.com/tokenchain/dp-hub/x/oracles/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/oracles", queryOraclesRequestHandler(cliCtx)).Methods("GET")
}

func queryOraclesRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bz, _, err := utils.QueryWithData(cliCtx, "custom/%s/%s", types.QuerierRoute,
			keeper.QueryOracles)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Couldn't get query data %s", err.Error())))
			return
		}

		var oracles types.Oracles
		if err := cliCtx.Codec.UnmarshalJSON(bz, &oracles); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Couldn't Unmarshal data %s", err.Error())))
			return
		}

		rest.PostProcessResponse(w, cliCtx, oracles)
	}
}
