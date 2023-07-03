package datadogtf

type SLO struct {
	ResourceName string       // terraform resource name to override the one calculated from the SLO name+service
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name"`
	Type         string       `json:"type"` // string^(metric)
	Description  string       `json:"description"`
	Query        SLOQuery     `json:"query"`
	Thresholds   SLOThreshold `json:"thresholds"`
	Tags         string       `json:"tags"`
}

type SLOQuery struct {
	Numerator   string `json:"numerator"`   // Example - "sum:my.custom.count.metric{type:good_events}.as_count()"
	Denominator string `json:"denominator"` // Example - "sum:my.custom.count.metric{*}.as_count()"
}

type SLOThreshold struct {
	Timeframe string  `json:"timeframe"` // string^(7d|30d)|90d|custom)$
	Target    float64 `json:"target"`    // [0..100]
}
