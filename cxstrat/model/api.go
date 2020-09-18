//+build cxstrat, cxstratclient

package model

import (
	"fmt"
	"net/http"
	"regexp"
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
	http.HandleFunc("/action/strat", ActionStrat)
	http.HandleFunc("/action/like", ActionLike)
	http.HandleFunc("/action/restrat", ActionRestrat)
	http.HandleFunc("/action/reply", ActionReply)
	http.HandleFunc("/action/rename", ActionRename)
	http.HandleFunc("/writerename", WriteRename)
	http.HandleFunc("/writereply", WriteReply)
	http.HandleFunc("/action/follow", ActionFollow)
	http.HandleFunc("/accountfollow", AccountFollow)
	//http.HandleFunc("/")
	http.HandleFunc("/writestrat", WriteStrat)
	http.HandleFunc("/all", GetAllStrats)

	go func() {
		err := http.ListenAndServe(":6969", nil)
		panic(err)
	}()
}

func launchApi(prgm *CXProgram) {
	/* write your API here! */
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	cache.myaccount = GetMyAccount()

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

	myaccount       string
	potentialfollow string

	err byte
}

var cache tcache

var stratform string = `<div class="strat">
<div class="boxedName"><h><b><a href="http://127.0.0.1:6969/account?id=%s">%s</a></b></h></div>
<div class="boxedStrat"><p>%s</p></div>
<p>likes: %d     restrats: %d<br>
ID: <a href="http://127.0.0.1:6969/strat?id=%s">%s</a><br>
<a href="http://127.0.0.1:6969/action/like?content=%s">like this strat</a><br>
<a href="http://127.0.0.1:6969/action/restrat?content=%s">restrat this strat</a><br>
<a href="http://127.0.0.1:6969/writereply?id=%s">Reply!</a><br>
</div>`
var accountplus string = `
<p>Followers: %d<br>Following: %d<br></p>
<a href="http://127.0.0.1:6969/account?id=%s">Account</a><br>
<a href="http://127.0.0.1:6969/action/follow?content=@%s">Follow!</a><br>
<a href="http://127.0.0.1:6969/accountlikes?id=%s">Likes</a><br>
<a href="http://127.0.0.1:6969/accountrestrats?id=%s">Restrats</a><br>
<a href="http://127.0.0.1:6969/accounttags?id=%s">Tags</a><br>
<a href="http://127.0.0.1:6969/accountfollow?id=%s">Followed/Following Accounts</a><br>
<a href="http://127.0.0.1:6969/writerename">Change Name</a><br>
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
<a href="http://127.0.0.1:6969/writestrat">Write your own CX Strat!</a> <br>
<img src="https://mcargobe.files.wordpress.com/2010/03/welcome-to-my-site.gif?w=486&h=339"><br>
%s
</body>
</html>
`

/**
 * 1. turns all @account to @Hbj180...
 * 2. checks hashtags
 * 3. checks length (too long fills cache.err with ERRORCODE_STRATTOOLONG)
 * 4. if all is good, then do it!
 */
func actionStrat(t string) {
	var b []byte
	regx := regexp.MustCompile("[0-9A-Za-z_]")
	regt := regexp.MustCompile("[0-9A-Za-z_]")

	for i := 0; i < len(t); i++ {
		if t[i] == '@' {
			b = append(b, byte('@'))
			i++
			var tstring []byte
			for i < len(t) && regx.Match([]byte{byte(t[i])}) {
				tstring = append(tstring, byte(t[i]))
				i++
			}

			//fmt.Println(string(tstring))

			addr_ := getWallet(string(tstring))

			if addr_.Null() {
				cache.err = ERRORCODE_STRATINVALID
				return
			}

			b = append(b, []byte(addr_.String())...)
		} else if t[i] == '#' {
			b = append(b, byte('#'))
			i++
			var tt []byte
			for i < len(t) && regt.Match([]byte{byte(t[i])}) {
				tt = append(tt, byte(t[i]))
				i++
			}
			b = append(b, tt...)
		}
		if i < len(t) {
			b = append(b, byte(t[i]))
		}
		if i >= 512 {
			cache.err = ERRORCODE_STRATTOOLONG
			return
		}
	}
	if len(b) > 512 {
		b = b[:512]
	}
	stratok := makeStrat(string(b))
	if stratok != 0 {
		cache.err = ERRORCODE_STRATINVALID
	}
}

func actionReply(t string) {
	var b []byte
	regx := regexp.MustCompile("[0-9A-Za-z_]")
	regt := regexp.MustCompile("[0-9A-Za-z_]")

	for i := 0; i < len(t); i++ {
		if t[i] == '@' {
			b = append(b, byte('@'))
			i++
			var tstring []byte
			for i < len(t) && regx.Match([]byte{byte(t[i])}) {
				tstring = append(tstring, byte(t[i]))
				i++
			}

			addr_ := getWallet(string(tstring))

			if addr_.Null() {
				cache.err = ERRORCODE_STRATINVALID
				return
			}

			b = append(b, []byte(addr_.String())...)
		} else if t[i] == '#' {
			b = append(b, byte('#'))
			i++
			var tt []byte
			for i < len(t) && regt.Match([]byte{byte(t[i])}) {
				tt = append(tt, byte(t[i]))
				i++
			}
			b = append(b, tt...)
		}
		if i < len(t) {
			b = append(b, byte(t[i]))
		}
		if i >= 412 {
			cache.err = ERRORCODE_STRATTOOLONG
			return
		}
	}
	id, err := cipher.SHA256FromHex(cache.potentialfollow)
	if err != nil || id.Null() {
		cache.err = ERRORCODE_STRATINVALID
	} else {
		if len(b) > 412 {
			b = b[:412]
		}
		stratok := makeReply(id, string(b))
		if stratok != 0 {
			cache.err = ERRORCODE_STRATINVALID
		}
	}
}

func actionLike(t string) {
	bt, err := cipher.SHA256FromHex(t)
	if err != nil {
		cache.err = ERRORCODE_STRATINVALID
		return
	}
	likeok := makeLike(bt)
	if likeok != 0 {
		cache.err = ERRORCODE_STRATINVALID
	}
}

func actionRestrat(t string) {
	bt, err := cipher.SHA256FromHex(t)
	if err != nil {
		cache.err = ERRORCODE_STRATINVALID
		return
	}
	likeok := makeRestrat(bt)
	if likeok != 0 {
		cache.err = ERRORCODE_STRATINVALID
	}
}

func buildAccountPlus(accnt Account) string {
	return fmt.Sprintf(accountplus, accnt.followcnt, accnt.fllwngcnt, accnt.cname, accnt.cname, accnt.cname, accnt.cname, accnt.cname, accnt.cname)
}

func buildStrat(cstrat Strat) string {
	strattxt := ""
	prvidx := 0
	for i := 0; i < len(cstrat.hashidx); i++ {
		strattxt += cstrat.txt[prvidx : cstrat.hashidx[i]-1]
		//fmt.Println(cstrat.txt[prvidx : cstrat.hashidx[i]-1])
		prvidx = cstrat.hashidx[i] + len(cstrat.hashs[i])
		hashtg := "#" + cstrat.hashs[i]
		strattxt += fmt.Sprintf("<a href=\"http://127.0.0.1:6969/hashtag?id=%s\">%s</a>", cstrat.hashs[i], hashtg)
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
			"<a href=\"http://127.0.0.1:6969/account?id="+str_+"\">@"+tagnm+"</a>", -1)
		alreadytagged = append(alreadytagged, str_)
	}
	alreadytagged = []string{}
	//the way we do this is mahjickal.
	//cstrat.hashs = hashsSort(cstrat.hashs)

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
		cstrat.obj.Hex(),
		cstrat.obj.Hex(),
		cstrat.obj.Hex(),
	)
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

func getStrats(hshs []cipher.SHA256) []Strat {
	rv := []Strat{}
	for i := len(hshs) - 1; i >= 0; i-- {
		rv = append(rv, getStrat(hshs[i]))
	}
	return rv
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

func WriteReply(w http.ResponseWriter, r *http.Request) {
	cache.potentialfollow = r.FormValue("id")
	str := `
	<form method="POST" action="http://127.0.0.1:6969/action/reply">
		<label for="writingstrat">Reply:</label>
		<input type="text" id="writingstrat" name="content"><br>
		<input type="submit" value="Submit">
	</form>
	`
	w.Write([]byte(fmt.Sprintf(docform, str)))
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
		rv += "<a href=\"http://127.0.0.1:6969/account?id=" + v.String() + "\">" + getName(v) + "</a><br>"
	}
	rv += "<h1><b>Following:</b></h1><br>"
	for _, v := range cache.account.following {
		rv += "<a href=\"http://127.0.0.1:6969/account?id=" + v.String() + "\">" + getName(v) + "</a><br>"
	}
	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, getName(bt), buildAccountPlus(cache.account))+"<br><br>"+rv)))
	cache.err = ERRORCODE_OK
}

func ActionFollow(w http.ResponseWriter, r *http.Request) {
	fllw := r.FormValue("content")
	if !strings.Contains(fllw, "@") {
		w.Write([]byte(fmt.Sprintf(docform, "<p>Invalid account to follow!<br><a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">back to my account</p>")))
		return
	}
	chk := fllw[1:]
	_, err := cipher.DecodeBase58Address(chk)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(docform, "<p>Invalid follow request!<br><a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">back to my account</p>")))
		return
	}
	retcode := makeFollow(fllw)
	if retcode == 2 {
		w.Write([]byte(fmt.Sprintf(docform, "<p>Already following!<br><a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">back to my account</p>")))
	} else {
		w.Write([]byte(fmt.Sprintf(docform, "<p>Followed!<br><a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">back to my account</p>")))
	}
}

func ActionReply(w http.ResponseWriter, r *http.Request) {
	actionReply(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Strat reply sent!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATINVALID:
		resp = "<p>Strat reply is invalid!<br><a href=\"http://127.0.0.1:6969/writestrat\">back to write strat</p>"
	case ERRORCODE_STRATTOOLONG:
		resp = "<p>Strat reply is too long!<br><a href=\"http://127.0.0.1:6969/writestrat\">back to write strat</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
	cache.err = ERRORCODE_OK
}

func ActionStrat(w http.ResponseWriter, r *http.Request) {
	actionStrat(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Strat sent!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATINVALID:
		resp = "<p>Strat is invalid!<br><a href=\"http://127.0.0.1:6969/writestrat\">back to write strat</p>"
	case ERRORCODE_STRATTOOLONG:
		resp = "<p>Strat is too long!<br><a href=\"http://127.0.0.1:6969/writestrat\">back to write strat</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
	cache.err = ERRORCODE_OK
}

func ActionLike(w http.ResponseWriter, r *http.Request) {
	actionLike(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Strat liked!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATINVALID:
		resp = "<p>Strat like is invalid!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATTOOLONG:
		resp = "<p>Strat like is too long!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
	cache.err = ERRORCODE_OK
}

func ActionRestrat(w http.ResponseWriter, r *http.Request) {
	actionRestrat(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Strat restrated!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATINVALID:
		resp = "<p>Strat restrat is invalid!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_STRATTOOLONG:
		resp = "<p>Strat restrat is too long!<br><a href=\"http://127.0.0.1:6969/account?id=" + cache.myaccount + "\">back to my account</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
	cache.err = ERRORCODE_OK
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(docform, "<a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">my account</p>")))
}

func WriteStrat(w http.ResponseWriter, r *http.Request) {
	str := `
	<form method="POST" action="http://127.0.0.1:6969/action/strat">
		<label for="writingstrat">Strat:</label>
		<input type="text" id="writingstrat" name="content"><br>
		<input type="submit" value="Submit">
	</form>
	`
	w.Write([]byte(fmt.Sprintf(docform, str)))
}

func WriteRename(w http.ResponseWriter, r *http.Request) {
	str := `
	<form method="POST" action="http://127.0.0.1:6969/action/rename">
		<label for="writingstrat">New Name:</label>
		<input type="text" id="writingstrat" name="content"><br>
		<input type="submit" value="Submit">
	</form>
	`
	w.Write([]byte(fmt.Sprintf(docform, str)))
}

func ActionRename(w http.ResponseWriter, r *http.Request) {
	newname := r.FormValue("content")
	if len(newname) > 32 || len(newname) < 3 || !getWallet(newname).Null() {
		w.Write([]byte(fmt.Sprintf(docform, "<p>Name invalid (either taken or too long/short)!<br><a href=\"http://127.0.0.1:6969/account?id="+cache.myaccount+"\">back to my account</p>")))
		return
	} else {
		makeName(newname)
		w.Write([]byte(fmt.Sprintf(docform, "renamed.")))
	}
}

func GetAllStrats(w http.ResponseWriter, r *http.Request) {
	cache.strats = getStrats(getAllStrats())

	w.Write([]byte(fmt.Sprintf(docform, "<br><h1><b>All Strats<br></b></h1><br>"+buildStrats(cache.strats))))
	cache.err = ERRORCODE_OK
}
