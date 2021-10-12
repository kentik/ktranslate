package aws

import (
	"bufio"
	"os"

	"github.com/kentik/ktranslate/pkg/kt"
)

func (vpc *AwsVpc) handleLocal(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	res := make([]*AWSLogLine, 0)
	size := int64(0)
	scanBuf := make([]byte, 1024*1024)

	scanner := bufio.NewScanner(f)
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
		vpc.Warnf("Could not scan %v", err)
		return err
	}

	if len(res) > 0 {
		record := FlowSet{
			Lines: res,
		}

		// Pull out all the info we can from the key path.
		err := record.ProcessKey("", file)
		if err != nil {
			vpc.Warnf("Could not process %s: %v.", file, err)
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
		vpc.Warnf("No flow data devices found for %s.", file)
	}

	// Get the record back, turn it into flow.
	rec := <-vpc.recs

	dst := make([]*kt.JCHF, len(rec.Lines))
	vpc.Debugf("Found %d logs to send", len(rec.Lines))
	for i, l := range rec.Lines {
		dst[i] = l.ToFlow(vpc, vpc.topo)
	}

	vpc.Infof("Ready to send %d lines", len(dst))
	if len(dst) > 0 {
		vpc.Infof("Row: %v", dst[0])
	}

	return nil
}
