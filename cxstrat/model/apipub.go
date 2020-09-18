//+build cxstrat, cxstratpublic

package model

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/skycoin/src/cipher"
	//. "github.com/SkycoinProject/cx/cx/base"
)

const (
	ERRORCODE_OK = iota
	ERRORCODE_STRATINVALID
	ERRORCODE_STRATTOOLONG
)

func handleRequests() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/account", AccountById)
	http.HandleFunc("/accountlikes", AccountLikesById)
	http.HandleFunc("/accountrestrats", AccountRestratsById)
	http.HandleFunc("/accounttags", AccountTagsById)
	http.HandleFunc("/strat", StratById)
	http.HandleFunc("/hashtag", HashtagById)
	http.HandleFunc("/accountfollow", AccountFollow)
	http.HandleFunc("/all", GetAllStrats)

	go func() {
		err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
		panic(err)
	}()
}

func launchApi(prgm *CXProgram) {
	/* write your API here! */
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	//cache.myaccount = GetMyAccount()

	handleRequests()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), true)
}

type tcache struct {
	strat    Strat
	account  Account
	hashtag  Hashtag
	strats   []Strat
	accounts []Account
	hashtags []Hashtag

	err byte
}

var cache tcache

var stratform string = `<div class="strat">
<div class="boxedName"><h><b><a href="https://cx.cxtweetcoin.com/account?id=%s">%s</a></b></h></div>
<div class="boxedStrat"><p>%s</p></div>
<p>likes: %d     restrats: %d<br>
ID: <a href="https://cx.cxtweetcoin.com/strat?id=%s">%s</a><br>
</div>`
var accountplus string = `
<p>Followers: %d<br>Following: %d<br></p>
<a href="https://cx.cxtweetcoin.com/account?id=%s">Account</a><br>
<a href="https://cx.cxtweetcoin.com/accountlikes?id=%s">Likes</a><br>
<a href="https://cx.cxtweetcoin.com/accountrestrats?id=%s">Restrats</a><br>
<a href="https://cx.cxtweetcoin.com/accounttags?id=%s">Tags</a><br>
<a href="https://cx.cxtweetcoin.com/accountfollow?id=%s">Followed/Following Accounts</a><br>
`

var accountform string = `<div class="account">
<h1><b>%s</b></h1></div><br>%s<br><br>
`
var hashtagform string = `<div class="hashtag"><h1><b>%s</b></h1><br><p>total likes: %d<br>total restrats: %d<br></p></div>`

var docform string = `
<!DOCTYPE html>
<html>
<head>
<title>CX Stratus</title>
<style type="text/css">
.boxedStrat {
	border: 4px solid green;
	width: 500px;
	padding: 25px;
	margin: 1px;
}
.boxedName {
	border: 1px solid black;
	width: 500px;
	padding: 0px;
	margin: 1px;
}
.strat {
	border: 10px solid black;
	margin: 10px;
}
.account {
	margin: 25px;
}
.hashtag {
	margin: 25px;
}
</style>
</head>
<body>
<img src="https://mcargobe.files.wordpress.com/2010/03/welcome-to-my-site.gif?w=486&h=339"><br>
%s
</body>
</html>
`

func buildAccountPlus(accnt Account) string {
	return fmt.Sprintf(accountplus, accnt.followcnt, accnt.fllwngcnt, accnt.cname, accnt.cname, accnt.cname, accnt.cname, accnt.cname)
}

func buildStrat(cstrat Strat) string {
	strattxt := ""
	prvidx := 0
	for i := 0; i < len(cstrat.hashidx); i++ {
		strattxt += cstrat.txt[prvidx : cstrat.hashidx[i]-1]
		//fmt.Println(cstrat.txt[prvidx : cstrat.hashidx[i]-1])
		prvidx = cstrat.hashidx[i] + len(cstrat.hashs[i])
		hashtg := "#" + cstrat.hashs[i]
		strattxt += fmt.Sprintf("<a href=\"https://cx.cxtweetcoin.com/hashtag?id=%s\">%s</a>", cstrat.hashs[i], hashtg)
	}
	/* we need to go through the strat text again and hyperlink tags and hashs */
	strattxt += cstrat.txt[prvidx:len(cstrat.txt)]
	var alreadytagged []string = []string{}
	for i := 0; i < len(cstrat.tags); i++ {
		str_ := cstrat.tags[i].String()
		var contains bool
		for j := 0; j < len(alreadytagged); j++ {
			if alreadytagged[j] == "@"+str_ {
				contains = true
				break
			}
		}
		if contains {
			continue
		}
		tagnm := getName(cstrat.tags[i])
		strattxt = strings.Replace(strattxt, "@"+str_,
			"<a href=\"https://cx.cxtweetcoin.com/account?id="+str_+"\">@"+tagnm+"</a>", -1)
		alreadytagged = append(alreadytagged, str_)
	}

	accnm := getName(cstrat.owner)

	return fmt.Sprintf(
		stratform,
		cstrat.owner.String(),
		accnm,
		strattxt,
		cstrat.likes,
		cstrat.restrats,
		cstrat.obj.Hex(),
		cstrat.obj.Hex()[:7]+"...",
	)
}

func getStrats(hshs []cipher.SHA256) []Strat {
	rv := []Strat{}
	for i := len(hshs) - 1; i >= 0; i-- {
		rv = append(rv, getStrat(hshs[i]))
	}
	return rv
}

func buildStrats(cstrats []Strat) string {
	rv := ""
	for i := 0; i < len(cstrats); i++ {
		rv += buildStrat(cstrats[i])
		if i < len(cstrats)-1 {
			rv += "<br>"
		}
	}
	return rv
}

func StratById(w http.ResponseWriter, r *http.Request) {
	bt, _ := cipher.SHA256FromHex(r.FormValue("id"))
	cache.strat = getStrat(bt)
	locache := getReplies(bt)
	cache.strats = []Strat{}
	for i, _ := range locache {
		cache.strats = append(cache.strats, getStrat(locache[i]))
	}
	w.Write([]byte(fmt.Sprintf(docform, buildStrat(cache.strat)+"<br><br><h1><b>Replies:</b></h1><br>"+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func AccountById(w http.ResponseWriter, r *http.Request) {
	bt, err := cipher.DecodeBase58Address(r.FormValue("id"))
	if err != nil {
		w.Write([]byte("<p>Error: this account could not be displayed.</p><br>"))
		return
	}
	cache.account = getAccount(bt)
	cache.strats = getStrats(cache.account.strats)
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func AccountLikesById(w http.ResponseWriter, r *http.Request) {
	bt, err := cipher.DecodeBase58Address(r.FormValue("id"))
	if err != nil {
		w.Write([]byte("<p>Error: this account could not be displayed.</p><br>"))
		return
	}
	cache.account = getAccount(bt)
	cache.strats = getStrats(cache.account.likes)
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func AccountRestratsById(w http.ResponseWriter, r *http.Request) {
	bt, err := cipher.DecodeBase58Address(r.FormValue("id"))
	if err != nil {
		w.Write([]byte("<p>Error: this account could not be displayed.</p><br>"))
		return
	}
	cache.account = getAccount(bt)
	cache.strats = getStrats(cache.account.restrats)
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func AccountTagsById(w http.ResponseWriter, r *http.Request) {
	bt, err := cipher.DecodeBase58Address(r.FormValue("id"))
	if err != nil {
		w.Write([]byte("<p>Error: this account could not be displayed.</p><br>"))
		return
	}
	cache.account = getAccount(bt)
	cache.strats = getStrats(cache.account.tags)
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func HashtagById(w http.ResponseWriter, r *http.Request) {
	cache.hashtag = getHashtag(r.FormValue("id"))
	cache.strats = getStrats(cache.hashtag.strats)
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(hashtagform, r.FormValue("id"), cache.hashtag.tlikes, cache.hashtag.trestrats)+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}

func AccountFollow(w http.ResponseWriter, r *http.Request) {
	bt, err := cipher.DecodeBase58Address(r.FormValue("id"))
	if err != nil {
		w.Write([]byte("<p>Error: this account could not be displayed.</p><br>"))
		return
	}
	cache.account = getAccount(bt)
	rv := "<h1><b>Followers:</b></h1><br>"
	for _, v := range cache.account.followers {
		rv += "<a href=\"https://cx.cxtweetcoin.com/account?id=" + v.String() + "\">" + getName(v) + "</a><br>"
	}
	rv += "<h1><b>Following:</b></h1><br>"
	for _, v := range cache.account.following {
		rv += "<a href=\"https://cx.cxtweetcoin.com/account?id=" + v.String() + "\">" + getName(v) + "</a><br>"
	}
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+"<br><br>"+rv)))
	cache.err = ERRORCODE_OK
}

func Home(w http.ResponseWriter, r *http.Request) {
	GetAllStrats(w, r)
}

func GetAllStrats(w http.ResponseWriter, r *http.Request) {
	cache.strats = getStrats(getAllStrats())

	w.Write([]byte(fmt.Sprintf(docform, "<br><h1><b>All Strats<br></b></h1><br>"+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}
