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
	"strings"
	"time"

	"github.com/MovieStoreGuy/detector-doctor/pkg/types"
	"github.com/gorilla/websocket"
)

const (
	// DefaultRealm is used when no realm is given as this is the default realm used by SignalFx
	DefaultRealm = "us0"

	// DefaultAPIEndpoint is the domain to be used when querying the API with the realm set
	DefaultAPIEndpoint = `https://api.%s.signalfx.com/v2`

	// DefaultStreamEndpoint is the domain to be used when using websockets to get data
	DefaultStreamEndpoint = `https://stream.%s.signalfx.com/v2`
)

func toUnixMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// SignalFx is the wrapper around the developer API
type SignalFx struct {
	api         string
	stream      string
	client      *http.Client
	requestFunc func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
	websocFunc  func(ctx context.Context, url string) (*websocket.Conn, error)
}

// NewSignalFxClient returns a configured client that will interact with the API
// using the access token realm set. If client or realm are not configured the defaults are used
func NewSignalFxClient(realm, accessToken string, client *http.Client) *SignalFx {
	sfx := &SignalFx{
		client:      client,
		requestFunc: NewConfiguredRequestFunc(accessToken),
		websocFunc:  NewConfiguredWebsocketFunc(accessToken),
	}
	if realm == "" {
		realm = DefaultRealm
	}
	sfx.api = fmt.Sprintf(DefaultAPIEndpoint, realm)
	sfx.stream = fmt.Sprintf(DefaultStreamEndpoint, realm)
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

func (sfx *SignalFx) readStreamData(ctx context.Context, programText string, opts map[string]interface{}) ([]types.Message, []*types.MetricDataPoint, error) {
	domain, err := url.Parse(sfx.stream)
	if err != nil {
		return nil, nil, err
	}
	domain.Path = path.Join(domain.Path, "signalflow/connect")
	domain.Scheme = strings.Replace(domain.Scheme, "http", "ws", 1)
	conn, err := sfx.websocFunc(ctx, domain.String())
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()
	jobSpec := map[string]interface{}{
		"type":     "execute",
		"channel":  "detdoc-0",
		"program":  programText,
		"timezone": "UTC",
	}
	for field, value := range opts {
		if _, exist := jobSpec[field]; !exist && value != nil {
			jobSpec[field] = value
		}
	}
	now := time.Now().UTC()
	if _, exist := jobSpec["start"]; !exist {
		jobSpec["start"] = toUnixMilliseconds(now.Add(-1 * 24 * time.Hour))
	}
	if _, exist := jobSpec["stop"]; !exist {
		if _, set := jobSpec["immediate"]; !set {
			jobSpec["stop"] = toUnixMilliseconds(now)
		}
	}
	err = conn.WriteJSON(jobSpec)
	if err != nil {
		return nil, nil, err
	}
	messages := make([]types.Message, 0)
	datapoints := make([]*types.MetricDataPoint, 0)
	for stopped := false; !stopped; {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			return nil, nil, err
		}
		if errRead := types.ReadWebsocketError(message); errRead != nil {
			return nil, nil, errRead
		}
		if types.IsEndofMessages(types.ReadControlMessage(message)) {
			break
		}
		switch msgType {
		case websocket.TextMessage:
			if control := types.ReadControlMessage(message); control != nil {
				if types.IsEndofMessages(control) {
					stopped = true
					break
				}
			} else if meta := types.ReadMetadataMessage(message); meta != nil {
				messages = append(messages, meta)

			} else if gen := types.ReadGeneralMessage(message); gen != nil {
				messages = append(messages, gen)
			}
		case websocket.BinaryMessage:
			dp, err := types.ReadMetricDataPoint(message)
			if err != nil {
				return nil, nil, err
			}
			datapoints = append(datapoints, dp)
		case websocket.CloseMessage:
			return messages, datapoints, nil
		}
	}
	return messages, datapoints, nil
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

// GetMetricTimeSeries returns the messages and data provided by the websocket API.
// The list of allowed parameters are documented here: https://developers.signalfx.com/signalflow_analytics/websocket_request_messages.html#_syntax_2
// Some of the values are predefined for you to avoid causing issues with handling computation
// All time values passed should be in UTC and all time values will be configured to be int64 values
func (sfx *SignalFx) GetMetricTimeSeries(ctx context.Context, programText string, params map[string]interface{}) ([]types.Message, []*types.MetricDataPoint, error) {
	for field, value := range params {
		if t, cast := value.(time.Time); cast {
			params[field] = toUnixMilliseconds(t)
		}
	}
	return sfx.readStreamData(ctx, programText, params)
}
