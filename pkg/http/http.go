package http

import (
	"encoding/json"
	"io"
	"main/pkg/types"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Client struct {
	logger zerolog.Logger
	chain  string
}

func NewClient(logger zerolog.Logger, chain string) *Client {
	return &Client{
		logger: logger.With().
			Str("component", "http").
			Str("chain", chain).
			Logger(),
		chain: chain,
	}
}

func (c *Client) GetInternal(
	url string,
	query string,
) (io.ReadCloser, types.QueryInfo, error) {
	var transport http.RoundTripper

	transportRaw, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		transport = transportRaw.Clone()
	} else {
		transport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	start := time.Now()

	queryInfo := types.QueryInfo{
		Success: false,
		Chain:   c.chain,
		Query:   query,
		URL:     url,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, queryInfo, err
	}

	req.Header.Set("User-Agent", "astronomer")

	c.logger.Debug().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	queryInfo.Duration = time.Since(start)
	if err != nil {
		c.logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return nil, queryInfo, err
	}

	c.logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return res.Body, queryInfo, err
}

func (c *Client) GetPlain(
	url string,
	query string,
) ([]byte, types.QueryInfo, error) {
	body, queryInfo, err := c.GetInternal(url, query)
	if err != nil {
		return nil, queryInfo, err
	}

	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, queryInfo, err
	}

	return bytes, queryInfo, nil
}

func (c *Client) Get(
	url string,
	query string,
	target interface{},
) (types.QueryInfo, error) {
	body, queryInfo, err := c.GetInternal(url, query)
	if err != nil {
		return queryInfo, err
	}

	if jsonErr := json.NewDecoder(body).Decode(target); jsonErr != nil {
		c.logger.Warn().Str("url", url).Err(jsonErr).Msg("Error decoding JSON from response")
		return queryInfo, jsonErr
	}

	return queryInfo, body.Close()
}
