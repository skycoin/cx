package dmsghttp

import (
	"bufio"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/skycoin/dmsg"
	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/dmsg/disc"
	"github.com/skycoin/skycoin/src/util/logging"
)

// Defaults for dmsg configuration, such as discovery URL
const (
	DefaultDiscoveryURL = "https://messaging.discovery.skywire.skycoin.net"
)

// DMSGTransport holds information about client who is initiating communication.
type DMSGTransport struct {
	Discovery disc.APIClient
	PubKey    cipher.PubKey
	SecKey    cipher.SecKey
}

// RoundTrip implements golang's http package support for alternative transport protocols.
// In this case DMSG is used instead of TCP to initiate the communication with the server.
func (t DMSGTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// init client
	dmsgC := dmsg.NewClient(t.PubKey, t.SecKey, t.Discovery, dmsg.SetLogger(logging.MustGetLogger("dmsgC_httpC")))

	// connect to the DMSG server
	if err := dmsgC.InitiateServerConnections(context.Background(), 1); err != nil {
		log.Fatalf("Error initiating server connections by initiator: %v", err)
	}

	// process remote pub key and port from dmsg-addr request header
	addrSplit := strings.Split(req.Host, ":")
	if len(addrSplit) != 2 {
		return nil, errors.New("Invalid server Pub Key or Port")
	}
	var pk cipher.PubKey
	if err := pk.Set(addrSplit[0]); err != nil {
		return nil, err
	}
	rPort, _ := strconv.Atoi(addrSplit[1])
	port := uint16(rPort)

	conn, err := dmsgC.Dial(context.Background(), pk, port)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := req.Write(conn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(conn), req)
}
