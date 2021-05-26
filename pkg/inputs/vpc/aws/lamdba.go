package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/kentik/ktranslate/pkg/kt"
)

func (vpc *AwsVpc) handleLamdba(ctx context.Context, s3Event events.S3Event) error {
	var wg sync.WaitGroup

	outputCB := func(err error) {
		if err != nil {
			vpc.Errorf("Cannot send: %v", err)
		}
		wg.Done()
	}

	for _, record := range s3Event.Records {
		if record.EventName != "ObjectCreated:Put" {
			vpc.Warnf("Skipping non put operation: %s", record.EventName)
			continue
		}
		obj := &s3.Object{
			Key: aws.String(record.S3.Object.Key),
		}
		err := vpc.processObject(record.S3.Bucket.Name, obj)
		if err != nil {
			return err
		}

		// Get the record back, turn it into flow.
		rec := <-vpc.recs

		dst := make([]*kt.JCHF, len(rec.Lines))
		for i, l := range rec.Lines {
			dst[i] = l.ToFlow(vpc, vpc.topo)
		}

		if len(dst) > 0 {
			wg.Add(1)
			vpc.lambdaHandler(dst, outputCB) // Decrement in the callback.
		}
	}

	vpc.Infof("Waiting on %d records to finish sending", len(s3Event.Records))
	wg.Wait()

	return nil
}
