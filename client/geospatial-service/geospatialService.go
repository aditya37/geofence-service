package geospatial

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	geospatialSrv "github.com/aditya37/api-contract/geospatial-service/service"
	"github.com/aditya37/geofence-service/client"
	http_conn "github.com/aditya37/geofence-service/connector/http"
	"github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
)

type (
	geospatialHttpclient struct {
		httpClient *http_conn.HttpConnector
		baseurl    string
	}
)

func NewGeospatialServiceClient(httpClient *http_conn.HttpConnector) (client.GeospatialServiceClient, error) {
	//hit health check
	baseURL := getenv.GetString("GEOSPATIAL_SERVICE_BASE_URL", "geospatial-service:7777")
	if err := httpClient.HealthChecker(http_conn.HttpRequestParam{
		BaseURL: baseURL,
		Path:    client.PathServiceHealthCheck.ToString(),
	}); err != nil {
		util.Logger().Error(err)
		return nil, err
	}
	return &geospatialHttpclient{
		httpClient: httpClient,
		baseurl:    baseURL,
	}, nil
}

func (gs *geospatialHttpclient) GetLocationById(ctx context.Context, locationid int64) (*client.ResponseGetLocationById, error) {
	resp, err := gs.httpClient.HttpRequester(
		ctx,
		http_conn.HttpRequestParam{
			BaseURL: gs.baseurl,
			Path:    fmt.Sprintf(client.PathGetLocationById.ToString(), locationid),
			Method:  http.MethodGet,
		},
	)
	if err != nil {
		util.Logger().Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var response client.ResponseGetLocationById
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

//
func (gs *geospatialHttpclient) GetNearbyLocation(ctx context.Context, request geospatialSrv.GetNearbyLocationRequest) (*geospatialSrv.GetNearbyLocationResponse, error) {
	byteBody, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	resp, err := gs.httpClient.HttpRequester(
		ctx,
		http_conn.HttpRequestParam{
			BaseURL: gs.baseurl,
			Path:    client.PathGetNearbyLocation.ToString(),
			Body:    byteBody,
			Method:  http.MethodPost,
		},
	)
	if err != nil {
		util.Logger().Error(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger().Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var response geospatialSrv.GetNearbyLocationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		util.Logger().Error(err)
		return nil, err
	}
	return &response, nil
}
