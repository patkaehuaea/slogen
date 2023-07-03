package datadog

import (
	"bytes"
	"embed"
	"text/template"

	oslo "github.com/OpenSLO/oslo/pkg/manifest/v1"
	"github.com/OpenSLO/slogen/libs/specs"
)

//go:embed templates/*.tf.gotf
var tmplFiles embed.FS

var tfTemplates *template.Template

const (
	SLOTmplName = "slo.tf.gotf"
)

func init() {
	var err error

	tfTemplates, err = template.ParseFS(tmplFiles, "templates/*.tf.gotf")
	if err != nil {
		panic(err)
	}
}

func GiveTerraform(apMap map[string]oslo.AlertPolicy, ntMap map[string]oslo.AlertNotificationTarget,
	slo specs.OpenSLOSpec) (string, error) {
	sloStr, err := GiveSLOTerraform(slo)

	if err != nil {
		return "", err
	}

	return sloStr, err
}

func GiveSLOTerraform(s specs.OpenSLOSpec) (string, error) {

	sumoSLO, err := ConvertToDatadogSLO(s)

	if err != nil {
		return "", err
	}

	tmpl := tfTemplates.Lookup(SLOTmplName)

	buff := &bytes.Buffer{}
	err = tmpl.Execute(buff, sumoSLO)

	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func IsSource(slo specs.OpenSLOSpec) bool {
	indicator := slo.Spec.Indicator

	sourceType := ""

	if indicator.Spec.RatioMetric != nil {
		sourceType = indicator.Spec.RatioMetric.Total.MetricSource.Type
	}

	return sourceType == SourceTypeMetrics
}
