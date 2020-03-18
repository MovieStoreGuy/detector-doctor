package types

// Detector defines the response object from the SignalFx API
// for further details, check: https://developers.signalfx.com/detectors_reference.html
type Detector struct {
	// Metadata blocks definitions

	ID            string `json:"id"`
	Created       int64  `json:"created"`
	LastUpdated   int64  `json:"lastUpdated"`
	Creator       string `json:"creator"`
	LastUpdatedBy string `json:"lastUpdatedBy"`

	CustomProperties string           `json:"customProperties"`
	Locked           bool             `json:"locked"`
	LabelResolutions map[string]int64 `json:"labelResolution"`
	MaxDelay         int32            `json:"maxDelay"`
	OverMTSLimit     bool             `json:"overMTSLimit"`
	Timezone         string           `json:"timezone"`

	// User Supplied values

	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	Teams                []string        `json:"teams"`
	Tags                 []string        `json:"tags"`
	ProgramText          string          `json:"programText"`
	Rules                []Rule          `json:"rules"`
	VisualizationOptions []Visualization `json:"visualizationOptions"`
}

// Rule defines the response object from the SignalFx API
type Rule struct {
	Description          string              `json:"description"`
	DetectLabel          string              `json:"detectLabel"`
	Disabled             bool                `json:"disabled"`
	Notifications        []map[string]string `json:"notifications,omitempty"`
	ParameterizedBody    string              `json:"parameterizedBody,omitempty"`
	ParameterizedSubject string              `json:"parameterizedSubject,omitempty"`
	RunbookUrl           string              `json:"runbookUrl,omitempty"`
	Severity             string              `json:"severity"`
	Tip                  string              `json:"tip,omitempty"`
}

type Visualization struct {
	DisableSampling     bool          `json:"disableSampling"`
	PublishLabelOptions []interface{} `json:"publishLabelOptions"`
	ShowDataMarkers     bool          `json:"showDataMarkers"`
	ShowEventLines      bool          `json:"showEventLines"`
	Time                Time          `json:"time"`
}

type Time struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Range int64  `json:"range"`
	Type  string `json:"type"`
}
