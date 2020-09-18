// +build cxtweetpublic

package cxtweet

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	. "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/skycoin/src/cipher"
)

type Tweet struct {
	owner cipher.Address
	obj   cipher.SHA256
	txt   string
	likes int
	retwt int
	tags  []cipher.Address
	hashs []string
}

type Account struct {
	name     string
	id       cipher.Address
	tweets   []cipher.SHA256
	likes    []cipher.SHA256
	retweets []cipher.SHA256
	tags     []cipher.SHA256
}

func (t Tweet) print() {
	fmt.Println(t.owner.String())
	fmt.Println(t.obj.Hex())
	fmt.Println(t.txt)
	fmt.Println(t.likes, t.retwt)
	for i := 0; i < len(t.hashs); i++ {
		fmt.Println("\t " + t.hashs[i])
	}
	fmt.Println("")
}

func (a Account) print() {
	fmt.Println(a.name)
	fmt.Println(a.id.String())
	for i := 0; i < len(a.tweets); i++ {
		fmt.Println("\t " + a.tweets[i].Hex())
	}
	fmt.Println("")
}

type Hashtag struct {
	name   string
	tlikes int
	trtwts int
	tweets []cipher.SHA256
}

func makeInt(b []byte) int {
	return (int(b[0]) << 24) | (int(b[1]) << 16) | (int(b[2]) << 8) | int(b[3])
}

func breakInt(i int) []byte {
	b := make([]byte, 4)
	b[0] = byte((i & 0xFF000000) >> 24)
	b[1] = byte((i & 0xFF0000) >> 16)
	b[2] = byte((i & 0xFF00) >> 8)
	b[3] = byte(i & 0xFF)
	return b
}

func deserializeTweet(b []byte) Tweet {
	t := Tweet{}
	t.owner, _ = cipher.AddressFromBytes(b[:25])
	b = b[25:]
	t.obj, _ = cipher.SHA256FromBytes(b[:32])
	b = b[32:]
	slen := makeInt(b[:4])
	b = b[4:]
	t.txt = string(b[:slen])
	b = b[slen:]
	t.likes = makeInt(b[:4])
	b = b[4:]
	t.retwt = makeInt(b[:4])
	b = b[4:]
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.AddressFromBytes(b[:25])
		t.tags = append(t.tags, tmp__)
		b = b[25:]
	}
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tlen := makeInt(b[:4])
		b = b[4:]
		t.hashs = append(t.hashs, string(b[:tlen]))
		b = b[tlen:]
	}

	//t.print()

	return t
}

func deserializeAccount(b []byte) Account {
	a := Account{}

	slen := makeInt(b[:4])
	b = b[4:]
	a.name = string(b[:slen])
	b = b[slen:]
	a.id, _ = cipher.AddressFromBytes(b[:25])
	b = b[25:]
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(b[:32])
		b = b[32:]
		a.tweets = append(a.tweets, tmp__)
	}
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(b[:32])
		b = b[32:]
		a.likes = append(a.likes, tmp__)
	}
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(b[:32])
		b = b[32:]
		a.retweets = append(a.retweets, tmp__)
	}
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(b[:32])
		b = b[32:]
		a.tags = append(a.tags, tmp__)
	}

	//a.print()

	return a
}

func deserializeHashtag(b []byte) Hashtag {
	h := Hashtag{}
	slen := makeInt(b[:4])
	b = b[4:]
	h.name = string(b[:slen])
	b = b[slen:]
	h.tlikes = makeInt(b[:4])
	b = b[4:]
	h.trtwts = makeInt(b[:4])
	b = b[4:]
	slen = makeInt(b[:4])
	b = b[4:]
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(b[:32])
		b = b[32:]
		h.tweets = append(h.tweets, tmp__)
	}

	return h
}

var donereq chan bool
var doneres chan bool
var req chan []byte
var res chan []byte

type tcache struct {
	tweet    Tweet
	account  Account
	hashtag  Hashtag
	tweets   []Tweet
	accounts []Account
	hashtags []Hashtag

	myaccount string

	err byte
}

/**
 * Tweet: tweet (obj), account (owner), accounts (tags), hashtags (hashs)
 * Account: account (id), tweets (one of: tweets, likes, retweets, tags)
 * Hashtag: hashtag (name), tweets (tweets)
 */

var cache tcache

const (
	DEFAULT_BUFFER_SIZE = 65536
	TWEETCODE_ISSHA     = 0
	TWEETCODE_ISWAL     = 1
	TWEETCODE_ISSTR     = 2
	TWEETCODE_ISTWT     = 3
	TWEETCODE_ISLIK     = 4
	TWEETCODE_ISRTT     = 5
	TWEETCODE_ISALL     = 6

	ERRORCODE_TWEETTOOLONG = 0
	ERRORCODE_OK           = 1
	ERRORCODE_TWEETINVALID = 2
)

var tweetform string = `<div class="tweet">
<div class="boxedName"><h><b><a href="http://34.122.174.177:80/account?id=%s">%s</a></b></h></div>
<div class="boxedTweet"><p>%s</p></div>
<p>likes: %d     retweets: %d<br>
ID: <a href="http://34.122.174.177:80/tweet?id=%s">%s</a><br>
</div>`
var accountform string = `<div class="account">
<h1><b>%s</b></h1><br></div>
`
var hashtagform string = `<div class="hashtag"><h1><b>%s</b></h1><br><p>total likes: %d<br>total retweets: %d<br></p></div>`

var docform string = `
<!DOCTYPE html>
<html>
<head>
<title>CX Tweet</title>
<style type="text/css">
.boxedTweet {
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
.tweet {
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

func init() {
	donereq = make(chan bool)
	doneres = make(chan bool)
	req = make(chan []byte, 1)
	res = make(chan []byte, 1)
	cache = tcache{}
}

func waitForAPIReturn() []byte {
	<-donereq
	return <-req
}

func waitForPrgmReturn() []byte {
	<-doneres
	return <-res
}

/**
 * 1. Program calls cxtweet.stall(), freeing control to API.
 * 2. API receives a call to an endpoint (functionXXX).
 * 3. functionXXX fills channel with received object.
 * 4. waitUntilAPIReturn() uses select to block until this.
 * 5. ^^^ returns bytes array to program.
 * 6. Program checks byte array to determine kind of API handler to use.
 * 7. Program serves.
 * 8. Program calls expose() to expose value and fill a channel.
 * 9. waitUntilPrgmReturn() was called by endpoint.
 * A. ^ uses select to block until channel is filled (which it is now).
 * B. Passes channel buffer to functionXXX.
 * C. functionXXX can finally return, deserializes channel buffer to return object.
 * D. GOTO step 1.
 */
func stall(prgm *CXProgram) {
	/* stalls the program until value is returned. */
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	byts := waitForAPIReturn()

	//fmt.Println("And now I am free.")
	//fmt.Println("[][][][][][]   : " + string(byts))

	outputSlicePointer := GetFinalOffset(fp, expr.Outputs[0])
	outputSliceOffset := GetPointerOffset(int32(outputSlicePointer))
	outputSliceOffset = int32(SliceResizeEx(outputSliceOffset, int32(len(byts)), 1))
	copy(GetSliceData(outputSliceOffset, 1), byts)
	WriteI32(outputSlicePointer, outputSliceOffset)
}

func expose(prgm *CXProgram) {
	fp := prgm.GetFramePointer()
	expr := prgm.GetExpr()

	inputSliceOffset := GetSliceOffset(fp, expr.Inputs[0])
	byts := GetSliceData(inputSliceOffset, 1)

	/* pass byts into buffer, then pass bool into chan */
	//fmt.Println("and now I returneth to my dusty abode.")
	res <- byts
	doneres <- true
}

func passToPrgm(byts []byte) {
	//fmt.Println("hello?")
	req <- byts
	//fmt.Println("who is it?")
	donereq <- true
	//fmt.Println("ohk")
}

/**
 * States:
 * /account/{id} (or /account/{id}): account, tweets are full
 * /account/{id}/likes: ""
 * /account/{id}/retweets: ""
 * /account/{id}/tags: ""
 * /tweet/{id}: tweet is full
 * /tweet/{id}/tags: tweet, accounts are full      <-- UNUSED
 * /tweet/{id}/hashtags: tweet, hashtags are full  <-- UNUSED
 * /hashtag/{id}: hashtag, tweets are full
 *
 * /action/tweet: tweet is full (is POST)
 * /action/like: tweet is full (is POST)
 * /action/retweet: tweet is full (is POST)
 * /: home
 */

/* /tweet/{id} */
func tweetById(id_ string) {
	/* convert id to string. */
	id, _ := cipher.SHA256FromHex(id_)
	byts_ := id[:]

	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	byts := []byte{TWEETCODE_ISSHA}
	byts = append(byts, byts_...)

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	/* now we are ready to finish the service. */
	cache.tweet = deserializeTweet(byts)
}

func accountById(id_ string, getTweets bool) {
	/* convert id to string. */
	id, _ := cipher.DecodeBase58Address(id_)
	byts_ := id.Bytes()

	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	byts := []byte{TWEETCODE_ISWAL}
	byts = append(byts, byts_...)

	//fmt.Println("all is good with the world.")

	/* return object to stall. */
	passToPrgm(byts)

	//fmt.Println("and all that is right must die.")

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	//fmt.Println("and all that is dead must live again.")

	/* now we are ready to finish the service. */
	cache.account = deserializeAccount(byts)

	/* next, we need to fetch a shit ton of tweets */
	if getTweets {
		cache.tweets = []Tweet{}
		for i := 0; i < len(cache.account.tweets); i++ {
			//fmt.Println(cache.account.tweets[i].Hex())
			tweetById(cache.account.tweets[i].Hex())
			cache.tweets = append(cache.tweets, cache.tweet)
			cache.tweet = Tweet{}
		}
	}
	/* and now we're done! */
}

func accountLikesById(id_ string) {
	if cache.account.id.String() != id_ {
		accountById(id_, false)
	}
	cache.tweets = []Tweet{}
	for i := 0; i < len(cache.account.likes); i++ {
		tweetById(cache.account.likes[i].Hex())
		cache.tweets = append(cache.tweets, cache.tweet)
	}
}

func accountRetweetsById(id_ string) {
	if cache.account.id.String() != id_ {
		accountById(id_, false)
	}
	cache.tweets = []Tweet{}
	for i := 0; i < len(cache.account.retweets); i++ {
		tweetById(cache.account.retweets[i].Hex())
		cache.tweets = append(cache.tweets, cache.tweet)
	}
}

func accountTagsById(id_ string) {
	if cache.account.id.String() != id_ {
		accountById(id_, false)
	}
	cache.tweets = []Tweet{}
	for i := 0; i < len(cache.account.tags); i++ {
		tweetById(cache.account.tags[i].Hex())
		cache.tweets = append(cache.tweets, cache.tweet)
	}
}

func hashtagById(id_ string) {
	/* convert id to string. */
	//fmt.Println("hashtagById: " + id_)
	byts_ := []byte(id_)

	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	byts := []byte{TWEETCODE_ISSTR}
	byts = append(byts, byts_...)

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	/* now we are ready to finish the service. */
	cache.hashtag = deserializeHashtag(byts)

	/* this always fills tweets */
	cache.tweets = []Tweet{}
	for i := 0; i < len(cache.hashtag.tweets); i++ {
		tweetById(cache.hashtag.tweets[i].Hex())
		cache.tweets = append(cache.tweets, cache.tweet)
	}
}

/**
 * 1. turns all @account to <@account>
 * 2. turns all #hashtag to <#hashtag>
 * 3. checks length (too long fills cache.err with ERRORCODE_TWEETTOOLONG)
 * 4. if all is good, then do that shit!
 */
func actionTweet(t string) {
	var b []byte
	regx := regexp.MustCompile("[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]")
	regt := regexp.MustCompile("[0-9A-Za-z_]")

	for i := 0; i < len(t); i++ {
		if t[i] == '@' {
			b = append(b, byte('@'))
			i++
			var tstring []byte
			for regx.Match([]byte{byte(t[i])}) && i < len(t) {
				tstring = append(tstring, byte(t[i]))
				i++
			}
			if i == len(t) {
				cache.err = ERRORCODE_TWEETINVALID
				return
			}
			tstring_ := string(tstring)
			taddr_, err := cipher.DecodeBase58Address(tstring_)
			if err != nil {
				cache.err = ERRORCODE_TWEETINVALID
				return
			}
			b = append(b, []byte(taddr_.String())...)
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
			cache.err = ERRORCODE_TWEETTOOLONG
			return
		}
	}

	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	byts := []byte{TWEETCODE_ISTWT}
	byts = append(byts, b...)

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	/* now we are ready to finish the service. */
	cache.err = byts[0]
}

func actionLike(t string) {
	bt, err := cipher.SHA256FromHex(t)
	if err != nil {
		cache.err = ERRORCODE_TWEETINVALID
		return
	}
	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	byts := []byte{TWEETCODE_ISLIK}
	byts = append(byts, bt[:]...)

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	/* now we are ready to finish the service. */
	cache.err = byts[0]
}

func actionRetweet(t string) {
	/* prepend the appropriate tweetcode (TWEETCODE_ISSHA) */
	bt, err := cipher.SHA256FromHex(t)
	if err != nil {
		cache.err = ERRORCODE_TWEETINVALID
		return
	}
	byts := []byte{TWEETCODE_ISRTT}
	byts = append(byts, bt[:]...)

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	/* now we are ready to finish the service. */
	cache.err = byts[0]
}

func buildTweet(ctweet Tweet) string {
	twttxt := ctweet.txt
	/* we need to go through the tweet text again and hyperlink tags and hashs */
	var alreadytagged []string = []string{}
	for i := 0; i < len(ctweet.tags); i++ {
		str_ := ctweet.tags[i].String()
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

		twttxt = strings.Replace(twttxt, "@"+str_,
			"<a href=\"http://34.122.174.177:80/account?id="+str_+"\">@"+str_+"</a>", -1)
		alreadytagged = append(alreadytagged, str_)
	}
	alreadytagged = []string{}
	for i := 0; i < len(ctweet.hashs); i++ {
		str_ := ctweet.hashs[i]
		var contains bool
		for j := 0; j < len(alreadytagged); j++ {
			if alreadytagged[j] == "#"+str_ {
				contains = true
				break
			}
		}
		if contains {
			continue
		}
		twttxt = strings.Replace(twttxt, "#"+str_,
			"<a href=\"http://34.122.174.177:80/hashtag?id="+str_+"\">#"+str_+"</a>", -1)
		alreadytagged = append(alreadytagged, str_)
	}

	return fmt.Sprintf(
		tweetform,
		ctweet.owner.String(),
		ctweet.owner.String(),
		twttxt,
		ctweet.likes,
		ctweet.retwt,
		ctweet.obj.Hex(),
		ctweet.obj.Hex(),
	)
}

func buildTweets(ctweets []Tweet) string {
	rv := ""
	for i := 0; i < len(ctweets); i++ {
		rv += buildTweet(ctweets[i])
		if i < len(ctweets)-1 {
			rv += "<br>"
		}
	}
	return rv
}

func TweetById(w http.ResponseWriter, r *http.Request) {
	tweetById(r.FormValue("id"))

	w.Write([]byte(fmt.Sprintf(docform, buildTweet(cache.tweet))))
}

func GetAllTweets(w http.ResponseWriter, r *http.Request) {
	byts := []byte{TWEETCODE_ISALL}

	/* return object to stall. */
	passToPrgm(byts)

	/* now simply wait for program to return */
	byts = waitForPrgmReturn()

	slen := makeInt(byts[:4])
	byts = byts[4:]
	cache.tweets = []Tweet{}
	for i := 0; i < slen; i++ {
		tmp__, _ := cipher.SHA256FromBytes(byts[:32])
		byts = byts[32:]
		tweetById(tmp__.Hex())
		cache.tweets = append(cache.tweets, cache.tweet)
	}

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, "All Tweets")+buildTweets(cache.tweets))))
}

func AccountById(w http.ResponseWriter, r *http.Request) {
	accountById(r.FormValue("id"), true)

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, r.FormValue("id"))+buildTweets(cache.tweets))))
}

func AccountLikesById(w http.ResponseWriter, r *http.Request) {
	accountLikesById(r.FormValue("id"))

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, r.FormValue("id"))+buildTweets(cache.tweets))))
}

func AccountRetweetsById(w http.ResponseWriter, r *http.Request) {
	accountRetweetsById(r.FormValue("id"))

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, r.FormValue("id"))+buildTweets(cache.tweets))))
}

func AccountTagsById(w http.ResponseWriter, r *http.Request) {
	accountTagsById(r.FormValue("id"))

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(accountform, r.FormValue("id"))+buildTweets(cache.tweets))))
}

func HashtagById(w http.ResponseWriter, r *http.Request) {
	hashtagById(r.FormValue("id"))

	w.Write([]byte(fmt.Sprintf(docform, fmt.Sprintf(hashtagform, r.FormValue("id"), cache.hashtag.tlikes, cache.hashtag.trtwts)+buildTweets(cache.tweets))))
}

func MyAccount(w http.ResponseWriter, r *http.Request) {

}

func ActionTweet(w http.ResponseWriter, r *http.Request) {
	actionTweet(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Tweet sent!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_TWEETINVALID:
		resp = "<p>Tweet is invalid!<br><a href=\"http://34.122.174.177:80/writetweet\">back to write tweet</p>"
	case ERRORCODE_TWEETTOOLONG:
		resp = "<p>Tweet is too long!<br><a href=\"http://34.122.174.177:80/writetweet\">back to write tweet</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
}

func ActionLike(w http.ResponseWriter, r *http.Request) {
	actionLike(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Tweet liked!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_TWEETINVALID:
		resp = "<p>Tweet like is invalid!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_TWEETTOOLONG:
		resp = "<p>Tweet like is too long!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
}

func ActionRetweet(w http.ResponseWriter, r *http.Request) {
	actionRetweet(r.FormValue("content"))
	var resp string
	switch cache.err {
	case ERRORCODE_OK:
		resp = "<p>Tweet retweeted!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_TWEETINVALID:
		resp = "<p>Tweet retweet is invalid!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	case ERRORCODE_TWEETTOOLONG:
		resp = "<p>Tweet retweet is too long!<br><a href=\"http://34.122.174.177:80/account?id=" + cache.myaccount + "\">back to my account</p>"
	}

	w.Write([]byte(fmt.Sprintf(docform, resp)))
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(docform, "<a href=\"http://34.122.174.177:80/account?id="+cache.myaccount+"\">my account</p>")))
}

func WriteTweet(w http.ResponseWriter, r *http.Request) {
	str := `
	<form method="POST" action="http://34.122.174.177:80/action/tweet">
		<label for="writingtweet">Tweet:</label>
		<input type="text" id="writingtweet" name="content"><br>
		<input type="submit" value="Submit">
	</form>
	`
	w.Write([]byte(fmt.Sprintf(docform, str)))
}

func handleRequests() {
	http.HandleFunc("/", GetAllTweets)
	http.HandleFunc("/account", AccountById)
	http.HandleFunc("/accountlikes", AccountLikesById)
	http.HandleFunc("/accountretweets", AccountRetweetsById)
	http.HandleFunc("/accounttags", AccountTagsById)
	http.HandleFunc("/tweet", TweetById)
	http.HandleFunc("/hashtag", HashtagById)

	go func() {
		err := http.ListenAndServe(":80", nil)
		panic(err)
	}()
}

func launchApi(prgm *CXProgram) {
	expr := prgm.GetExpr()
	fp := prgm.GetFramePointer()

	handleRequests()

	WriteBool(GetFinalOffset(fp, expr.Outputs[0]), true)
}
