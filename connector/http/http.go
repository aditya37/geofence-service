package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aditya37/geofence-service/util"
	"golang.org/x/net/context/ctxhttp"
)

type (
	HttpConnector struct {
		client *http.Client
	}
	HttpRequestParam struct {
		BaseURL     string
		Path        string
		Method      string
		AccessToken string
		Body        []byte
		Header      map[string]string
	}
)

var (
	ErrServiceUnhealty = errors.New("Service unhealty")
)

func NewHttpConnector() (*HttpConnector, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &HttpConnector{
		client: client,
	}, nil
}

// httpRequester
func (hc *HttpConnector) HttpRequester(ctx context.Context, param HttpRequestParam) (*http.Response, error) {
	url := fmt.Sprintf(param.BaseURL+"%s", param.Path)
	util.Logger().Info(fmt.Sprintf("Request to: %s", url))

	// HTTP Request instance
	req, err := http.NewRequest(param.Method, url, bytes.NewReader(param.Body))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := hc.addDefinedHeader(req, param); err != nil {
		return nil, err
	}

	// client
	client, err := ctxhttp.Do(ctx, hc.client, req)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// healthChacker...
func (hc *HttpConnector) HealthChecker(param HttpRequestParam) error {
	resp, err := hc.HttpRequester(context.Background(), param)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ErrServiceUnhealty
}
func (hc *HttpConnector) addDefinedHeader(req *http.Request, param HttpRequestParam) error {
	if param.AccessToken != "" && param.Header["Authorization"] != "" {
		req.Header.Add("Authorization", param.AccessToken)
	}
	for header, value := range param.Header {
		req.Header.Add(header, value)
	}
	return nil
}
