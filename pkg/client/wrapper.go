package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
)

const (
	// DefaultRealm is used when no realm is given as this is the default realm used by SignalFx
	DefaultRealm = "us0"

	// DefaultAPIEndpoint is the domain to be used when querying the API with the realm set
	DefaultAPIEndpoint = `https://api.%s.signalfx.com/v2`
)

// SignalFx is the wrapper around the developer API
type SignalFx struct {
	api         string
	client      *http.Client
	requestFunc func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
}

// NewSignalFxClient returns a configured client that will interact with the API
// using the access token realm set. If client or realm are not configured the defaults are used
func NewSignalFxClient(realm, accessToken string, client *http.Client) *SignalFx {
	sfx := &SignalFx{
		client:      client,
		requestFunc: NewConfiguredRequestFunc(accessToken),
	}
	if realm == "" {
		realm = DefaultRealm
	}
	sfx.api = fmt.Sprintf(DefaultAPIEndpoint, realm)
	if sfx.client == nil {
		sfx.client = NewConfiguredClient()
	}
	return sfx
}

// makeRequest abstracts needing to handle the requests to and from SignalFx and returns the buffer read from the body
func (sfx *SignalFx) makeRequest(ctx context.Context, method string, data io.Reader, queryParams map[string]interface{}, pathByParts ...string) ([]byte, error) {
	domain, err := url.Parse(sfx.api)
	if err != nil {
		return nil, err
	}
	domain.Path = path.Join(domain.Path, path.Join(pathByParts...))
	q := domain.Query()
	for name, value := range queryParams {
		q.Set(name, fmt.Sprint(value))
	}
	domain.RawQuery = q.Encode()
	req, err := sfx.requestFunc(ctx, method, domain.String(), data)
	if err != nil {
		return nil, err
	}
	resp, err := sfx.client.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		// Do Nothing
	case http.StatusBadRequest:
		return nil, types.ErrAPIIssue
	case http.StatusNotFound:
		return nil, types.ErrNoDetectorFound
	}
	return ioutil.ReadAll(resp.Body)
}

// GetDetectorByID retrives the provided detector as defined by https://developers.signalfx.com/detectors_reference.html#tag/Retrieve-Detector-ID
func (sfx *SignalFx) GetDetectorByID(ctx context.Context, detectorID string) (*types.Detector, error) {
	buff, err := sfx.makeRequest(ctx, http.MethodGet, nil, nil, "detector", detectorID)
	if err != nil {
		return nil, err
	}
	var det types.Detector
	return &det, json.Unmarshal(buff, &det)
}

// GetIncidentsByDetectorID retrives the provided incidents as defined by https://developers.signalfx.com/detectors_reference.html#tag/Retrieve-Incidents-Single-Detector
func (sfx *SignalFx) GetIncidentsByDetectorID(ctx context.Context, detectorID string, query map[string]interface{}) ([]*types.Incident, error) {
	return nil, types.ErrNotImplemented
}
