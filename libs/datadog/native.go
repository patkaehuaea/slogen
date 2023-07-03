package datadog

import (
	"fmt"

	"github.com/OpenSLO/slogen/libs/datadog/datadogtf"
	"github.com/OpenSLO/slogen/libs/specs"
)

const (
	MonitorKindBurnRate    = "BurnRate"
	MonitorKindSLI         = "SLI"
	ComplianceTypeCalendar = "Calendar"
	ComplianceTypeRolling  = "Rolling"
)

const (
	AnnotationMonitorFolderID = "datadog/monitor-folder-id"
	AnnotationSLOFolderID     = "datadog/slo-folder-id"
	AnnotationTFResourceName  = "datadog/tf-resource-name"
	AnnotationSignalType      = "datadog/signal-type"
)

const (
	// Used to check if type is supported
	SourceTypeMetrics = "datadog-metrics"
)

type SLO struct {
	*datadogtf.SLO
}

func (s SLO) TFResourceName() string {

	if s.ResourceName != "" {
		return s.ResourceName
	}

	return fmt.Sprintf("datadogslo_%s", s.Name)
}

func ConvertToDatadogSLO(slo specs.OpenSLOSpec) (*SLO, error) {

	indicator := slo.Spec.Indicator

	resourceName := slo.Metadata.Annotations[AnnotationTFResourceName]

	tagsMap := make(map[string]string)

	for k := range slo.ObjectHeader.Metadata.Labels {
		if len(slo.ObjectHeader.Metadata.Labels[k]) > 0 {
			tagsMap[k] = slo.ObjectHeader.Metadata.Labels[k][0]

		}
	}

	// TODO: Does not work beyond first element of a list.
	tagsStr := ""
	c := 0
	for k := range tagsMap {
		tagsStr += "\"" + k + ":" + tagsMap[k] + "\""
		if c < len(tagsMap)-1 {
			tagsStr += ", "
		}
		c += 1
	}

	datadogSLO := &SLO{
		&datadogtf.SLO{
			ResourceName: resourceName,
			Name:         slo.SLO.Metadata.Name,
			Type:         "",
			Description:  slo.Spec.Description,

			Query: datadogtf.SLOQuery{
				Numerator:   indicator.Spec.RatioMetric.Good.MetricSource.MetricSourceSpec["query"],
				Denominator: indicator.Spec.RatioMetric.Total.MetricSource.MetricSourceSpec["query"],
			},
			Thresholds: datadogtf.SLOThreshold{
				Timeframe: slo.SLO.Spec.TimeWindow[0].Duration,
				Target:    slo.SLO.Spec.Objectives[0].Target,
			},
			Tags: tagsStr,
		},
	}

	return datadogSLO, nil
}
