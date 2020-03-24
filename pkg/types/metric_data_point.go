package types

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	// Endian is being cached within this package as the default Endian type to use when processing binary requests
	Endian = binary.BigEndian
)

// MetricDataPoint contains all the values needed for a collection of datapoints
type MetricDataPoint struct {
	Version            int8    `json:"-"`
	Mtype              int8    `json:"-"`
	Flags              int8    `json:"-"`
	Channel            string  `json:"channel"`
	Kind               string  `json:"type"`
	LogicalTimestampMs int64   `json:"logicalTimestampMs"`
	MaxDelayMs         int64   `json:"maxDelayMs"`
	Data               []*Data `json:"data"`
}

// Data represents a singular value
type Data struct {
	Version      uint8  `json:"version,omitempty"`
	TimeSeriesID string `json:"tsId"`
	// Value represents either a float32, int32 or nil
	Value interface{} `json:"value"`
}

func min(l, r int8) int {
	if l < r {
		return int(l)
	}
	return int(r)
}

// ReadMetricDataPoint reads a binary message from a websocket response and then
// will convert into a strong golang type
func ReadMetricDataPoint(bin []byte) (*MetricDataPoint, error) {
	// Take a reasonable amount of the logic to decode the binary message
	// by following this: https://github.com/signalfx/signalfx-python/blob/master/signalfx/signalflow/ws.py#L167
	if len(bin) < 20 {
		return nil, errors.New("not enough values for header")
	}
	dp := &MetricDataPoint{
		Kind: "data",
	}
	values := []interface{}{
		&dp.Version,
		&dp.Mtype,
		&dp.Flags,
	}
	header, data := bin[:20], bin[20:]
	for index := 0; index < len(values); index++ {
		if err := binary.Read(bytes.NewReader(header[0:1]), Endian, values[index]); err != nil {
			return nil, err
		}
		header = header[1:]
	}
	if dp.Version > 3 && dp.Version != 0 {
		return nil, fmt.Errorf("unsupported binary version: %v", dp.Version)
	}
	channel := make([]byte, len(header))
	if err := binary.Read(bytes.NewBuffer(header), Endian, channel); err != nil {
		return nil, err
	}

	dp.Channel = string(channel)
	isCompressed := (bool)((dp.Flags & (1 << 0)) != 0)
	isJSON := (bool)((dp.Flags & (1 << 1)) != 0)

	if isCompressed {
		r, err := zlib.NewReader(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		data = buf
	}
	if isJSON {
		if err := json.Unmarshal(data, &dp.Data); err != nil {
			return nil, err
		}
		return dp, nil
	}
	values = []interface{}{
		&dp.LogicalTimestampMs,
		&dp.MaxDelayMs,
	}

	for index := 0; index < min(dp.Version, int8(len(values))); index++ {
		if err := binary.Read(bytes.NewBuffer(data[:8]), Endian, values[index]); err != nil {
			return nil, err
		}
		data = data[8:]
	}
	// Discarding the first four values as per their code base suggests
	// Not entirely sure why but go with it...
	// data = data[4:]
	var size uint32 = 0
	if err := binary.Read(bytes.NewBuffer(data[:4]), Endian, &size); err != nil {
		return nil, err
	}
	data = data[4:]
	for i := 0; i < int(size); i++ {
		point := &Data{}
		if err := binary.Read(bytes.NewBuffer(data[:1]), Endian, &point.Version); err != nil {
			return nil, err
		}
		tsid := make([]byte, 7)
		if err := binary.Read(bytes.NewBuffer(data[1:8]), Endian, tsid); err != nil {
			return nil, err
		}
		point.TimeSeriesID = strings.Replace(base64.URLEncoding.EncodeToString(tsid), "=", "", -1)

		// For reference, refer to this: https://developers.signalfx.com/signalflow_analytics/rest_api_messages/stream_messages_specification.html#table-3-data-message-payload-fields
		switch point.Version {
		case 2: // Is a double c type
			var value float64 = 0.0
			if err := binary.Read(bytes.NewBuffer(data[9:17]), Endian, &value); err != nil {
				return nil, err
			}
			point.Value = value
		case 1, 3: // Is a int c type
			var value int64
			if err := binary.Read(bytes.NewBuffer(data[9:17]), Endian, &value); err != nil {
				return nil, err
			}
			point.Value = value
		}
		dp.Data = append(dp.Data, point)
		data = data[17:]
	}
	return dp, nil
}
