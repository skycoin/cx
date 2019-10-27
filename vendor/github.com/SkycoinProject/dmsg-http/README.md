# dmsg-http

Http library for dmsg.
It supports all the features from standard golang http library.
Only difference is that instead of tcp as transport protocol it uses DMSG.

In order to run the tests for this project you should run following line

```bash
go get golang.org/x/net
```

## Examples

In order to instantiate the server you can use following code

```golang
// define port where server will listen
serverPort := uint16(8080)

// prepare the server
sPK, sSK := cipher.GenerateKeyPair()
httpS := dmsghttp.Server{
    PubKey: sPK,
    SecKey: sSK,
    Port: testPort,
    Discovery: dmsgD
}

// prepare server route handling
mux := http.NewServeMux()
mux.HandleFunc("/some-route", func(w http.ResponseWriter, _ *http.Request) {
    _, err := w.Write([]byte("Route response goes here"))
    if err != nil {
        panic(err)
    }
})

// run the server
sErr := make(chan error, 1)
go func() {
    sErr <- httpS.Serve(mux)
    close(sErr)
}()
```

If you would like to talk to this server following code will suffice

```golang
// prepare the client
cPK, cSK := cipher.GenerateKeyPair()
c := dmsghttp.DMSGClient(dmsgD, cPK, cSK)

// make request
req, err := http.NewRequest("GET", fmt.Sprintf("dmsg://%v:%d/some-route", sPK.Hex(), testPort), nil)
resp, err := c.Do(req)
respBody, err := ioutil.ReadAll(resp.Body)
fmt.Println(string(respBody))
```
