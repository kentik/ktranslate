package aws

import (
	"bufio"
	"compress/gzip"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

/**
{"Records":[{"eventVersion":"2.1","eventSource":"aws:s3","awsRegion":"us-west-2","eventTime":"2020-09-17T18:24:24.619Z","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"AWS:AROAV2CJ256E23ZOBRSAF:prod.pdx.dbs.datafeeds.aws.internal"},"requestParameters":{"sourceIPAddress":"172.19.15.211"},"responseElements":{"x-amz-request-id":"AD5D47CDA09091B3","x-amz-id-2":"gMLgqszXsmKN41Ou3/l330BEXa+ARrbIP8UkW9VqW21WdC42ie4Ki1WMP5Zm8M6R1TuAgkjPJFtlFB2HX+Ui3yeNqrJ5oBz0"},"s3":{"s3SchemaVersion":"1.0","configurationId":"Flow","bucket":{"name":"kentik-test-orangeflow","ownerIdentity":{"principalId":"A2L4QHGC7GJYP3"},"arn":"arn:aws:s3:::kentik-test-orangeflow"},"object":{"key":"AWSLogs/451031991406/vpcflowlogs/us-west-2/2020/09/17/451031991406_vpcflowlogs_us-west-2_fl-0ac5de8260cdc0575_20200917T1820Z_c127bbf5.log.gz","size":1797,"eTag":"1638e127fe977dacb06958f144f0b549","sequencer":"005F63A9DAF94F4D9E"}}}]}
*/
type SQSEvent struct {
	Records []SQSRecord `json:"Records"`
}

type SQSRecord struct {
	EventName string `json:"eventName"`
	S3        SQSS3  `json:"s3"`
}

type SQSS3 struct {
	Bucket SQSBucket `json:"bucket"`
	Object SQSObject `json:"object"`
}

type SQSBucket struct {
	Name string `json:"name"`
}

type SQSObject struct {
	Key string `json:"key"`
}

func (vpc *AwsVpc) processObject(bucket string, mdata *s3.Object) error {
	vpc.Debugf("Processing %s %s", bucket, *mdata.Key)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    mdata.Key,
	}
	res := make([]*AWSLogLine, 0)
	size := int64(0)
	scanBuf := make([]byte, 1024*1024)
	result, err := vpc.client.GetObject(input)
	if err != nil {
		vpc.Errorf("Cannot process object: %s %s -> %v", bucket, *mdata.Key, err)
		return err
	} else {
		if result.ContentLength != nil {
			size = *result.ContentLength
		}
		zr, err := gzip.NewReader(result.Body)
		if err != nil {
			if !strings.HasSuffix(*mdata.Key, ".gz") {
				// if possibly not gz encoded, try looping directly.
				scanner := bufio.NewScanner(result.Body)
				if size > bufio.MaxScanTokenSize {
					scanner.Buffer(scanBuf, len(scanBuf)*1024)
				}
				lineMap := AwsLineMap{}
				for scanner.Scan() {
					rawLine := scanner.Text()
					lines, lineMapOut, err := NewAws(lineMap, &rawLine, vpc)
					if err != nil {
						// weird case here so just skip.
						if len(rawLine) > 80 {
							vpc.Errorf("Error A reading line: %v -> %s", err, rawLine[0:80])
						} else {
							vpc.Errorf("Error A reading line: %v -> %s", err, rawLine)
						}
					} else {
						if lineMapOut != nil {
							lineMap = lineMapOut
						}
						if lines != nil {
							res = append(res, lines...)
						}
					}
				}
				if err := scanner.Err(); err != nil {
					vpc.Warnf("cannot scan clear %s %s: %v", bucket, *mdata.Key, err)
					return err
				}
			} else { // Can't un-gzip here.
				vpc.Warnf("cannot gz %s %s: %v", bucket, *mdata.Key, err)
				return err
			}
		} else {
			scanner := bufio.NewScanner(zr)
			if size > bufio.MaxScanTokenSize {
				scanner.Buffer(scanBuf, len(scanBuf)*1024)
			}
			lineMap := AwsLineMap{}
			for scanner.Scan() {
				rawLine := scanner.Text()
				lines, lineMapOut, err := NewAws(lineMap, &rawLine, vpc)
				if err != nil {
					if len(rawLine) > 80 {
						vpc.Errorf("Error reading line: %v -> %s", err, rawLine[0:80])
					} else {
						vpc.Errorf("Error reading line: %v -> %s", err, rawLine)
					}
				} else {
					if lineMapOut != nil {
						lineMap = lineMapOut
					}
					if lines != nil {
						res = append(res, lines...)
					}
				}
			}
			if err := scanner.Err(); err != nil {
				vpc.Warnf("cannot scan %s %s: %v", bucket, *mdata.Key, err)
				return err
			}
			if err := zr.Close(); err != nil {
				vpc.Warnf("cannot close zr %s %s: %v", bucket, *mdata.Key, err)
				return err
			}
		}
	}

	if len(res) > 0 {
		record := FlowSet{
			Lines: res,
		}

		// Pull out all the info we can from the key path.
		err := record.ProcessKey(bucket, *mdata.Key)
		if err != nil {
			vpc.Warnf("cannot process %s %s: %v", bucket, *mdata.Key, err)
			return err
		}

		// Send this record on to be processed.
		select {
		case vpc.recs <- &record:
			vpc.metrics.Flows.Mark(int64(len(record.Lines)))
		default:
			vpc.metrics.DroppedFlows.Mark(int64(len(record.Lines)))
		}
	} else {
		vpc.Warnf("No flows found for %s %s", bucket, *mdata.Key)
	}

	return nil
}

// Get message from sqs
func (vpc *AwsVpc) checkQueue(ctx context.Context) error {
	maxMessages := int64(10)
	visTimout := int64(10)
	waitTime := int64(5)
	res, err := vpc.sqsCli.ReceiveMessage(&sqs.ReceiveMessageInput{
		MaxNumberOfMessages: &maxMessages,
		QueueUrl:            &vpc.awsQUrl,
		VisibilityTimeout:   &visTimout,
		WaitTimeSeconds:     &waitTime,
	})
	if err != nil {
		return err
	}

	for _, m := range res.Messages {
		vpc.Debugf("Got message: %s", *m.Body)

		// Delete for only once processing
		_, err := vpc.sqsCli.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &vpc.awsQUrl,
			ReceiptHandle: m.ReceiptHandle,
		})
		if err != nil {
			vpc.Errorf("Error deleting message: %v", err)
		}

		// Now, process this out.
		var record SQSEvent
		if err := json.Unmarshal([]byte(*m.Body), &record); err != nil {
			vpc.Errorf("Error parsing message: %v", err)
		} else {
			for _, rec := range record.Records {
				if rec.S3.Bucket.Name != "" && rec.S3.Object.Key != "" {
					obj := &s3.Object{
						Key: aws.String(rec.S3.Object.Key),
					}
					go vpc.processObject(rec.S3.Bucket.Name, obj)
				} else {
					vpc.Errorf("Invalid message: %s", *m.Body)
				}
			}
		}
	}

	return nil
}
