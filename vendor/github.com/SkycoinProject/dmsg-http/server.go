package dmsghttp

import (
	"context"
	"log"
	"net/http"

	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/dmsg/disc"
	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/dmsg"
)

// Server holds relevant data for server to run properly
// Data includes Public / Secret key pair that identifies the server.
// There is also port on which server will listen.
// Optional parameter is Discovery, if none is provided default one will be used.
// Default dicovery URL is stored as dmsghttp.DefaultDiscoveryURL
type Server struct {
	PubKey    cipher.PubKey
	SecKey    cipher.SecKey
	Port      uint16
	Discovery disc.APIClient

	hs *http.Server
}

// Serve handles request to dmsg server
// Accepts handler holding routes for the current instance
func (s *Server) Serve(handler http.Handler) error {
	s.hs = &http.Server{Handler: handler}

	client := dmsg.NewClient(s.PubKey, s.SecKey, s.Discovery, dmsg.SetLogger(logging.MustGetLogger("dmsgC_httpS")))
	if err := client.InitiateServerConnections(context.Background(), 1); err != nil {
		log.Fatalf("Error initiating server connections by initiator: %v", err)
	}

	list, err := client.Listen(s.Port)
	if err != nil {
		return err
	}
	return s.hs.Serve(list)
}

// Close closes active Listeners and Connections by invoking http's Close func
func (s *Server) Close() error {
	return s.hs.Close()
}
