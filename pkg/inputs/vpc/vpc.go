package vpc

import (
	"context"
	"fmt"

	go_metrics "github.com/kentik/go-metrics"
	"github.com/kentik/ktranslate"

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
	apic *api.KentikApi, maxBatchSize int, lambdaHandler func([]*kt.JCHF, func(error)), cfg *ktranslate.Config) (VpcImpl, error) {
	switch cloud {
	case Aws:
		return aws.NewVpc(ctx, log, registry, jchfChan, apic, lambdaHandler, cfg.AWSVPCInput)
	case Gcp:
		return gcp.NewVpc(ctx, log, registry, jchfChan, apic, maxBatchSize, cfg.GCPVPCInput)
	case Azure:
		return nil, fmt.Errorf("Azure is not yet supported as a VPC: %v.", cloud)
	}
	return nil, fmt.Errorf("You used an unsupported VPC: %v.", cloud)
}
