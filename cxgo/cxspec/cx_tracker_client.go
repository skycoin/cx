package cxspec

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SkycoinProject/cx-chains/src/cipher"
	"github.com/sirupsen/logrus"
)

type CXTrackerClient struct {
	log  logrus.FieldLogger
	c    *http.Client
	addr string
}

func NewCXTrackerClient(log logrus.FieldLogger, c *http.Client, addr string) *CXTrackerClient {
	if log == nil {
		l := logrus.New()
		l.Level = logrus.FatalLevel
		log = l
	}
	if c == nil {
		c = http.DefaultClient
	}
	addr = strings.TrimSuffix(addr, "/")

	return &CXTrackerClient{log: log, c: c, addr: addr}
}

func (c *CXTrackerClient) AllSpecs(ctx context.Context) ([]SignedChainSpec, error) {
	log := c.log.WithField("func", "Spec")

	addr := fmt.Sprintf("%s/api/specs", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeRespBody(log, resp)

	if err := checkRespCode(resp); err != nil {
		return nil, err
	}

	var specs []SignedChainSpec
	if err := json.NewDecoder(resp.Body).Decode(&specs); err != nil {
		return nil, err
	}

	for i, spec := range specs {
		if err := spec.Verify(); err != nil {
			return nil, fmt.Errorf("failed to verify returned spec [%d]%s: %w", i, spec.Spec.ChainPubKey, err)
		}
	}

	return specs, nil
}

func (c *CXTrackerClient) SpecByPK(ctx context.Context, pk cipher.PubKey) (SignedChainSpec, error) {
	log := c.log.WithField("func", "SpecByPK")

	addr := fmt.Sprintf("%s/api/specs/pk:%s", c.addr, pk.Hex())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return SignedChainSpec{}, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return SignedChainSpec{}, err
	}
	defer closeRespBody(log, resp)

	if err := checkRespCode(resp); err != nil {
		return SignedChainSpec{}, err
	}

	var spec SignedChainSpec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return SignedChainSpec{}, err
	}

	if err := spec.Verify(); err != nil {
		return SignedChainSpec{}, fmt.Errorf("failed to verify returned spec: %w", err)
	}

	return spec, nil
}

func (c *CXTrackerClient) SpecByTicker(ctx context.Context, ticker string) (SignedChainSpec, error) {
	log := c.log.WithField("func", "SpecByTicker")

	addr := fmt.Sprintf("%s/api/specs/ticker:%s", c.addr, strings.TrimSpace(strings.ToUpper(ticker)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return SignedChainSpec{}, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return SignedChainSpec{}, err
	}
	defer closeRespBody(log, resp)

	if err := checkRespCode(resp); err != nil {
		return SignedChainSpec{}, err
	}

	var spec SignedChainSpec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return SignedChainSpec{}, err
	}

	if err := spec.Verify(); err != nil {
		return SignedChainSpec{}, fmt.Errorf("failed to verify returned spec: %w", err)
	}

	return spec, nil
}

func (c *CXTrackerClient) PostSpec(ctx context.Context, spec SignedChainSpec) error {
	log := c.log.WithField("func", "PostSpec")

	if err := spec.Verify(); err != nil {
		return err
	}

	r, w := io.Pipe()
	go func() {
		if err := json.NewEncoder(w).Encode(spec); err != nil {
			log.WithError(err).Error("Failed to encode spec to json: %w", err)
		}
		_ = w.Close() //nolint:errcheck
	}()

	addr := fmt.Sprintf("%s/api/specs", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, r)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer closeRespBody(log, resp)

	return checkRespCode(resp)
}

func (c *CXTrackerClient) DelSpec(ctx context.Context, pk cipher.PubKey) error {
	log := c.log.WithField("func", "DelSpec")

	addr := fmt.Sprintf("%s/api/spec/pk:%s", c.addr, pk.Hex())
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, addr, nil)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer closeRespBody(log, resp)

	return checkRespCode(resp)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func closeRespBody(log logrus.FieldLogger, resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.WithError(err).Error("Failed to close HTTP response body.")
	}
}

func httpStatusError(addr string, code int) error {
	return fmt.Errorf("request to '%s' failed with code '%d %s'", addr, code, http.StatusText(code))
}

func checkRespCode(resp *http.Response) error {
	code := resp.StatusCode

	if code == http.StatusOK {
		return nil
	}

	var errMsg string
	if err := json.NewDecoder(resp.Body).Decode(&errMsg); err != nil {
		// unexpected server response
		return fmt.Errorf("failed to decode server response with code '%d %s': %w",
			code, http.StatusText(code), err)
	}

	return fmt.Errorf("server responded with '%d %s': %s",
		code, http.StatusText(code), errMsg)
}