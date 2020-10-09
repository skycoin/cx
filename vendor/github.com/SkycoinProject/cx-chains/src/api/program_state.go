package api

import (
	"fmt"
	"net/http"

	wh "github.com/SkycoinProject/cx-chains/src/util/http"
	"github.com/SkycoinProject/cx-chains/src/visor"
)

// Returns the program state of a CX chain.
// Method: GET
// URI: /api/v1/programState
// Args:
//     addrs: Comma separated addresses [optional, returns all transactions if no address provided]
//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
//	   verbose: [bool] include verbose transaction input data
func programStateHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		// Gets 'addrs' parameter value
		addrs, err := parseAddressesFromStr(r.FormValue("addrs"))
		if err != nil {
			wh.Error400(w, fmt.Sprintf("parse parameter: 'addrs' failed: %v", err))
			return
		}

		// Initialize transaction filters
		flts := []visor.TxFilter{
			visor.NewAddrsFilter(addrs),
			visor.NewConfirmedTxFilter(true),
		}

		prgrmState, err := gateway.GetProgramState(flts)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, prgrmState)
	}
}
