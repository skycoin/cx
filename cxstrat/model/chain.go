//+build cxstrat, cxstratlnx

package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	. "github.com/SkycoinProject/cx/cx"
	//"github.com/SkycoinProject/cx/cxstrat/model"
	"github.com/SkycoinProject/skycoin/src/api"
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/readable"
)

var client_ *api.Client
var wlname_ string
var walid_ string
var addr_ string
var addrh_ cipher.Ripemd160
var wallet_ *api.WalletResponse
var wbal_ uint64
var wbalh_ uint64
var acname_ string

func postMessageBlockchain(prgm *CXProgram, content_ []byte) {
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	csrfToken, err := client_.CSRF()
	if err != nil {
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	url_ := "http://127.0.0.1:6422/api/v1/wallet/transaction"

	var dataMap map[string]interface{}
	dataMap = make(map[string]interface{}, 0)
	dataMap["tweet"] = string(content_[:])
	//fmt.Println(string(content_[:]))
	dataMap["hours_selection"] = map[string]string{"type": "manual"}
	// dataMap["wallet_id"] = map[string]string{"id": options.walletId}
	dataMap["wallet_id"] = walid_
	dataMap["to"] = []interface{}{map[string]string{"address": addr_, "coins": strconv.FormatUint(wbal_/1000000, 10), "hours": strconv.FormatUint(wbalh_-1, 10)}}
	//fmt.Println(wbal_, wbalh_)

	jsonStr, err := json.Marshal(dataMap)
	if err != nil {
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	req, err := http.NewRequest("POST", url_, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	var enctxn api.CreateTransactionResponse
	if err := json.Unmarshal(body, &enctxn); err != nil {
		// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
		fmt.Println(string(body))
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	csrfToken, err = client_.CSRF()
	if err != nil {
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	url_ = "http://127.0.0.1:6422/api/v1/injectTransaction"
	dataMap = make(map[string]interface{}, 0)
	dataMap["rawtx"] = enctxn.EncodedTransaction

	jsonStr, err = json.Marshal(dataMap)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	req, err = http.NewRequest("POST", url_, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), true)

	wbalh_--
}

func makeAccount(prgm *CXProgram) {
	/* makes account from 12 word seed*/
	/* I also need to add opcode for changing names. */

	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	csrfToken, err := client_.CSRF()
	if err != nil {
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	/* first FIRST: first check if wallet exists */
	reqc, err := http.NewRequest("GET", "http://127.0.0.1:6422/api/v1/wallets", nil)
	reqc.Header.Set("X-CSRF-Token", csrfToken)
	reqc.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respc, err := client.Do(reqc)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	defer respc.Body.Close()
	bodyc, err := ioutil.ReadAll(respc.Body)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	var respBodyc []api.WalletResponse = []api.WalletResponse{}
	if err := json.Unmarshal(bodyc, &respBodyc); err != nil {
		// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
		fmt.Println(string(bodyc))
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	if len(respBodyc) > 0 {
		wlname_ = respBodyc[0].Meta.Label
		addr_ = respBodyc[0].Entries[0].Address
		addrht_, _ := cipher.DecodeBase58Address(addr_)
		addrh_ = addrht_.Key
		walid_ = respBodyc[0].Meta.Filename
		wallet_ = &respBodyc[0]
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), true)
	} else {
		csrfToken, err = client_.CSRF()
		if err != nil {
			WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
			return
		}

		mne_ := bip39.MustNewDefaultMnemonic()

		dat_ := url.Values{}
		dat_.Set("label", "defwalname")
		dat_.Set("type", "deterministic")
		dat_.Set("seed", mne_)

		url_ := "http://127.0.0.1:6422/api/v1/wallet/create"
		req, err := http.NewRequest("POST", url_, strings.NewReader(dat_.Encode()))
		req.Header.Set("X-CSRF-Token", csrfToken)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
			return
		}

		var respBody api.WalletResponse
		if err := json.Unmarshal(body, &respBody); err != nil {
			// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
			fmt.Println(string(body))
			WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
			return
		}

		wlname_ = respBody.Meta.Label
		addr_ = respBody.Entries[0].Address
		addrht_, _ := cipher.DecodeBase58Address(addr_)
		addrh_ = addrht_.Key
		walid_ = respBody.Meta.Filename
		wallet_ = &respBody
	}

	csrfToken, err = client_.CSRF()
	if err != nil {
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	reqcc, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:6422/api/v1/wallet/balance?id=%s", walid_), nil)
	reqc.Header.Set("X-CSRF-Token", csrfToken)
	reqc.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	respcc, err := client.Do(reqcc)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	defer respcc.Body.Close()
	bodycc, err := ioutil.ReadAll(respcc.Body)
	if err != nil {
		fmt.Println(err.Error())
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	var respBodycc api.BalanceResponse
	if err := json.Unmarshal(bodycc, &respBodycc); err != nil {
		// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
		fmt.Println(string(bodycc))
		WriteBool(GetFinalOffset(fp, expr.Outputs[0]), false)
		return
	}

	myaccount = addr_

	wbal_ = respBodycc.Confirmed.Coins
	wbalh_ = respBodycc.Confirmed.Hours
	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), true)
}

var myaccount string

func GetMyAccount() string {
	return myaccount
}

func getBlock(prgm *CXProgram) {
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	csrfToken, err := client_.CSRF()
	if err != nil {
		WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
		return
	}

	url_ := fmt.Sprintf("http://127.0.0.1:6422/api/v1/block?seq=%d&verbose=1", int(ReadI32(prgm.GetFramePointer(), prgm.GetExpr().Inputs[0])))
	req, err := http.NewRequest("GET", url_, nil)
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
		return
	}

	var respBody readable.BlockVerbose
	if err := json.Unmarshal(body, &respBody); err != nil {
		// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
		fmt.Println(string(body))
		WriteI32(GetFinalOffset(fp, expr.Outputs[0]), int32(0))
		return
	}

	txid, _ := cipher.SHA256FromHex(respBody.Body.Transactions[0].Hash)
	wal_, _ := cipher.DecodeBase58Address(respBody.Body.Transactions[0].In[0].Address)
	//fmt.Println(wal_.String())
	//fmt.Println(wal_.Bytes())
	//fmt.Println(txid[:])
	txt := respBody.Body.Transactions[0].Tweet
	var fllwwal_ cipher.Address
	var hashs_ []string
	var pname string
	hashs_ = makeHashs(txt)
	if strings.HasPrefix(txt, "<|::[]CXSTRATUS_FOLLOWF[]::|>") {
		fllwwal_ = makeTags(txt)[0]
	}
	if strings.HasPrefix(txt, "<|::[]CXSTRATUS_NAMENAM[]::|>") {
		b := []byte{}
		i := 29
		for txt[i] != 0 {
			b = append(b, byte(txt[i]))
			i++
		}
		pname = string(b)
	}
	taggs := makeTags(txt)
	//for i, _ := range taggs {
	//fmt.Println(taggs[i].Bytes())
	//fmt.Println(taggs[i].String())
	//}
	strat := Strat{owner: wal_, obj: txid, txt: txt, pname: pname, tags: taggs, hashs: hashs_, hashidx: makeHashIdx(txt), follow: fllwwal_}
	val := serializeStrat(strat)
	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(val)), 1))
	copy(GetSliceData(outputSliceOffset, 1), val)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func chainLen(prgm *CXProgram) {
	csrfToken, err := client_.CSRF()
	if err != nil {
		fmt.Println(err.Error())
		WriteI32(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), int32(-1))
		return
	}

	url_ := "http://127.0.0.1:6422/api/v1/blockchain/metadata"
	req, err := http.NewRequest("GET", url_, nil)
	req.Header.Set("X-CSRF-Token", csrfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err)
		WriteI32(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), int32(-1))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err)
		WriteI32(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), int32(-1))
		return
	}

	var respBody readable.BlockchainMetadata
	if err := json.Unmarshal(body, &respBody); err != nil {
		// Printing the body instead of `err`. Body has the error generated in the Skycoin API.
		fmt.Println(string(body))
		WriteI32(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), int32(-1))
		return
	}

	seqno := respBody.Head.BkSeq
	WriteI32(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), int32(seqno))
}

func transact(prgm *CXProgram) {
	/* all POST requests use the wallet file generated on startup. */
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	/* we need to read the tweet */
	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts_ := GetSliceData(inputSliceOffset, 1)

	postMessageBlockchain(prgm, byts_)
}

//1. launches the website
//2. launches the Tweetchain
//3.
func launchCXStrat(prgm *CXProgram) {
	cmd_ := exec.Command("./tweetcoin", "-enable-all-api-sets",
		"-disable-default-peers",
		"-custom-peers-file=localhost-peers.txt",
		"-download-peerlist=false", "-launch-browser=false",
		"-genesis-address=pZ9RPocV6CGoQYbwe36wyChSVtKmXfoQW2",
		"-genesis-signature=b3e3f118307ab7729b5480634cad554a1888c9944125ca676b63102ca250f5e440113c11360f3c8bf5d4f3ff167738e33814a15c7c6a765fb851de1565b8bd0001",
		"-blockchain-public-key=022c67fe688692a73f98386bca1cf1d75e72212692397ae0f3425c46ac66773784",
		"-block-publisher=false",
		"-web-interface-port=6422", "-port=6002",
		"> debug.txt",
	)

	cmd_.Start()

	//cmd2_ := exec.Command("rm", "-f", "cxdb.db")
	//cmd2_.Start()

	client_ = api.NewClient("http://127.0.0.1:6422")

	time.Sleep(time.Second * 5)

	WriteBool(GetFinalOffset(prgm.GetFramePointer(), prgm.GetExpr().Outputs[0]), true)
}
