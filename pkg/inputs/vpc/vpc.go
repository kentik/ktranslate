package vpc

import (
	"context"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"

	"github.com/kentik/ktranslate/pkg/api"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/inputs/vpc/aws"
	"github.com/kentik/ktranslate/pkg/inputs/vpc/gcp"
	"github.com/kentik/ktranslate/pkg/kt"
)

type VpcImpl interface {
	Close()
	HttpInfo() map[string]float64
}

type CloudSource string

const (
	Aws   CloudSource = "aws"
	Gcp               = "gcp"
	Azure             = "azure"
)

func NewVpc(ctx context.Context, cloud CloudSource, log logger.Underlying, registry go_metrics.Registry, jchfChan chan []*kt.JCHF,
	apic *api.KentikApi, maxBatchSize int, lambdaHandler func([]*kt.JCHF, func(error))) (VpcImpl, error) {
	switch cloud {
	case Aws:
		return aws.NewVpc(ctx, log, registry, jchfChan, apic, lambdaHandler)
	case Gcp:
		return gcp.NewVpc(ctx, log, registry, jchfChan, apic, maxBatchSize)
	case Azure:
		return nil, fmt.Errorf("Unimplemented vpc %v", cloud)
	}
	return nil, fmt.Errorf("Unknown vpc %v", cloud)
}
