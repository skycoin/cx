package cxdmsg

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/dmsg"
	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/dmsg/disc"
)

// Default cx dmsg values.
const (
	DefaultDiscAddr = "dmsg.discovery.skywire.skycoin.com"
	DefaultPort     = uint16(9090)
)

type Config struct {
	PK       cipher.PubKey
	SK       cipher.SecKey
	DiscAddr string
	DmsgPort uint16
}

func ServeDmsg(ctx context.Context, log logrus.FieldLogger, conf *Config, api *API) {
	dmsgC := dmsg.NewClient(conf.PK, conf.SK, disc.NewHTTP(conf.DiscAddr), nil)
	go dmsgC.Serve(ctx)

	select {
	case <-ctx.Done():
		return
	case <-dmsgC.Ready():
	}

	lis, err := dmsgC.Listen(conf.DmsgPort)
	if err != nil {
		log.WithError(err).Fatalf("Failed to serve on dmsg port '%d'.", conf.DmsgPort)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.WithError(err).Error("Failed to close dmsg listener.")
		}
	}()

	err = http.Serve(lis, api)
	log.WithError(err).Info("Stopped serving.")
}
