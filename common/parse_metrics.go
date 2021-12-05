package common

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"io"
)

func ParseMF(reader io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}
	return mf, nil
}
