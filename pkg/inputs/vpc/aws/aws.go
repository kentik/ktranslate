package aws

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/ic"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kentik/patricia"
)

/**
version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status
2 672740005165 eni-ba88b581 54.191.55.41 10.0.0.53 80 41536 6 27368 39574687 1533851315 1533851375 ACCEPT OK
2 672740005165 eni-ba88b581 10.0.0.53 54.218.137.160 52772 80 6 62 4451 1533851315 1533851375 ACCEPT OK
2 672740005165 eni-ba88b581 159.89.148.43 10.0.0.53 56911 5984 6 1 40 1533851315 1533851375 REJECT OK
2 672740005165 eni-ba88b581 159.89.148.43 10.0.0.53 35789 8545 6 1 40 1533851315 1533851375 REJECT OK
2 672740005165 eni-ba88b581 185.248.100.159 10.0.0.53 40879 22 6 1 40 1533851315 1533851375 ACCEPT OK
2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK

Or, Kinethis:

2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK
{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"1571081770","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"6","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"1571081799"}}

Or, v3
version vpc-id subnet-id instance-id interface-id account-id type srcaddr dstaddr srcport dstport pkt-srcaddr pkt-dstaddr protocol bytes packets start end action tcp-flags log-status
3 vpc-2e0d3a4a subnet-0a279552 - eni-22859f22 047615872826 IPv4 10.232.50.235 10.232.51.166 62527 27016 10.232.50.235 18.162.196.19 6 120 2 1580321301 1580321356 ACCEPT 2 OK

*/

const (
	MIN_AWS_FIELD = 7
	ADDR_LEN      = 17

	AWS_LOG_PREFIX = "AWSLogs"

	AWS_ACTION       = "action"
	AWS_STATUS       = "log-status"
	AWS_VERSION      = "version"
	AWS_VPC_ID       = "vpc-id"
	AWS_SUBNET_ID    = "subnet-id"
	AWS_INSTANCE_ID  = "instance-id"
	AWS_INTERFACE_ID = "interface-id"
	AWS_ACCOUNT_ID   = "account-id"
	AWS_TYPE         = "type"
	AWS_SRC_ADDR     = "srcaddr"
	AWS_DST_ADDR     = "dstaddr"
	AWS_SRC_PORT     = "srcport"
	AWS_DST_PORT     = "dstport"
	AWS_PKT_SRC_ADDR = "pkt-srcaddr"
	AWS_PKT_DST_ADDR = "pkt-dstaddr"
	AWS_PROTOCOL     = "protocol"
	AWS_BYTES        = "bytes"
	AWS_PACKETS      = "packets"
	AWS_START        = "start"
	AWS_END          = "end"
	AWS_LINE_ACTION  = "action"
	AWS_TCP_FLAGS    = "tcp-flags"
	AWS_LOG_STATUS   = "log-status"

	// v4 fields
	AWS_REGION           = "region"
	AWS_AZ_ID            = "az-id"
	AWS_SUBLOCATION_TYPE = "sublocation-type"
	AWS_SUBLOCATION_ID   = "sublocation-id"

	// v5 fields
	AWS_PKT_SRC_AWS_SERVICE = "pkt-src-aws-service"
	AWS_PKT_DST_AWS_SERVICE = "pkt-dst-aws-service"
	AWS_FLOW_DIRECTION      = "flow-direction"
	AWS_TRAFFIC_PATH        = "traffic-path"

	AWS_VPC_TYPE = "AWS_VPC"
)

var (
	baseSplitAWS = `{"messageType":`

	AWS_FLOW_FIELDS = []string{
		AWS_VERSION,
		AWS_INTERFACE_ID,
		AWS_ACCOUNT_ID,
		AWS_SRC_ADDR,
		AWS_DST_ADDR,
		AWS_SRC_PORT,
		AWS_DST_PORT,
		AWS_PROTOCOL,
		AWS_BYTES,
		AWS_PACKETS,
		AWS_START,
		AWS_END,
		AWS_ACTION,
		AWS_STATUS,
		AWS_TYPE,
		AWS_TCP_FLAGS,
		AWS_VPC_ID,
		AWS_SUBNET_ID,
		AWS_INSTANCE_ID,
		AWS_PKT_SRC_ADDR,
		AWS_PKT_DST_ADDR,
		AWS_REGION,
		AWS_AZ_ID,
		AWS_SUBLOCATION_TYPE,
		AWS_SUBLOCATION_ID,
		AWS_PKT_SRC_AWS_SERVICE,
		AWS_PKT_DST_AWS_SERVICE,
		AWS_FLOW_DIRECTION,
		AWS_TRAFFIC_PATH,
	}
)

type AWSLogLine struct {
	Version         int
	AccountID       string
	InterfaceID     string
	SrcAddr         net.IP
	DstAddr         net.IP
	SrcPktAddr      net.IP
	DstPktAddr      net.IP
	TcpFlags        uint32
	SrcPort         uint32
	DstPort         uint32
	Protocol        uint32
	Packets         uint64
	Bytes           uint64
	StartTime       time.Time
	EndTime         time.Time
	Action          string
	Status          string
	Sample          uint32
	VPCID           string
	SubnetID        string
	InstanceID      string
	Region          string
	AzID            string
	SublocationType string
	SublocationID   string
	SrcPktService   string
	DstPktService   string
	FlowDirection   string
	TrafficPath     string
}

type AwsLineMap map[string]int

type KinesisLogWrapper struct {
	MessageType string       `json:"messageType"`
	Owner       string       `json:"owner"`
	LogGroup    string       `json:"logGroup"`
	LogEvents   []KinesisLog `json:"logEvents"`
}

type KinesisLog struct {
	Id              string         `json:"id"`
	Message         string         `json:"message"`
	ExtractedFields ExtractedField `json:"extractedFields"`
}

type ExtractedField struct {
	SrcAddr     string `json:"srcaddr"`
	DstPort     string `json:"dstport"`
	StartTime   string `json:"start"`
	DstAddr     string `json:"dstaddr"`
	Version     string `json:"version"`
	Packets     string `json:"packets"`
	Protocol    string `json:"protocol"`
	AccountId   string `json:"account_id"`
	InterfaceId string `json:"interface_id"`
	Status      string `json:"log_status"`
	Bytes       string `json:"bytes"`
	SrcPort     string `json:"srcport"`
	Action      string `json:"action"`
	EndTime     string `json:"end"`
}

func NewAwsFromKinesis(lineMap AwsLineMap, raw *string, log logger.ContextL) ([]*AWSLogLine, AwsLineMap, error) {

	// Need to split on multiple lines stuck together here.
	res := []*AWSLogLine{}
	basePts := strings.Split(*raw, baseSplitAWS)

	// At this point, roughCount is approx number of flows in this chunk.
	// Now, go through and pick flows to sample.
	for _, b := range basePts {
		if len(b) < 10 {
			continue
		}

		lines := KinesisLogWrapper{}
		err := json.Unmarshal([]byte(baseSplitAWS+b), &lines)
		if err != nil {
			return nil, lineMap, err
		}

		thisRes := make([]*AWSLogLine, len(lines.LogEvents))
		for i, l := range lines.LogEvents {
			line := AWSLogLine{}
			ver, err := strconv.Atoi(l.ExtractedFields.Version)
			if err != nil || ver == 0 {
				if len(l.Message) > 10 && !strings.Contains(*raw, baseSplitAWS) {
					// Recurse because not pulled out for us.
					fLines, _, err := NewAws(lineMap, &l.Message, log)
					if err != nil {
						log.Errorf("Cannot parse basic lines: %v", err)
					} else {
						res = append(res, fLines...)
					}
				}
				//log.Debugf("No parsed lines found: %s, %v", l.Message, l.ExtractedFields)
				continue
			} else {
				line.Version = ver
			}

			line.AccountID = l.ExtractedFields.AccountId
			line.InterfaceID = l.ExtractedFields.InterfaceId
			line.SrcAddr = net.ParseIP(l.ExtractedFields.SrcAddr)
			line.DstAddr = net.ParseIP(l.ExtractedFields.DstAddr)
			base, _ := strconv.Atoi(l.ExtractedFields.SrcPort)
			line.SrcPort = uint32(base)
			base, _ = strconv.Atoi(l.ExtractedFields.DstPort)
			line.DstPort = uint32(base)
			base, _ = strconv.Atoi(l.ExtractedFields.Protocol)
			line.Protocol = uint32(base)
			base, _ = strconv.Atoi(l.ExtractedFields.Packets)
			line.Packets = uint64(base)
			base, _ = strconv.Atoi(l.ExtractedFields.Bytes)
			line.Bytes = uint64(base)
			base, _ = strconv.Atoi(l.ExtractedFields.StartTime)
			line.StartTime = time.Unix(int64(base), 0)
			base, _ = strconv.Atoi(l.ExtractedFields.EndTime)
			line.EndTime = time.Unix(int64(base), 0)
			line.Action = l.ExtractedFields.Action
			line.Status = l.ExtractedFields.Status

			if line.StartTime.Before(time.Now().Add(-7 * 24 * time.Hour)) {
				log.Debugf("Bad Kinesis log line: %v", *raw)
			} else {
				thisRes[i] = &line
			}
		}

		for _, r := range thisRes {
			if r != nil {
				res = append(res, r)
			}
		}
	}

	log.Debugf("Returning %d kinesis lines", len(res))
	return res, lineMap, nil
}

//2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK
func NewAwsFromV2(lineMap AwsLineMap, pts []string, log logger.ContextL) ([]*AWSLogLine, AwsLineMap, error) {
	line := AWSLogLine{}
	line.Version = 2
	line.AccountID = getString(lineMap, pts, AWS_ACCOUNT_ID)
	line.InterfaceID = getString(lineMap, pts, AWS_INTERFACE_ID)
	line.SrcAddr = getIP(lineMap, pts, AWS_SRC_ADDR)
	line.DstAddr = getIP(lineMap, pts, AWS_DST_ADDR)
	line.SrcPort = getUint32(lineMap, pts, AWS_SRC_PORT)
	line.DstPort = getUint32(lineMap, pts, AWS_DST_PORT)
	line.Protocol = getUint32(lineMap, pts, AWS_PROTOCOL)
	line.Packets = getUint64(lineMap, pts, AWS_PACKETS)
	line.Bytes = getUint64(lineMap, pts, AWS_BYTES)
	line.StartTime = getTime(lineMap, pts, AWS_START)
	line.EndTime = getTime(lineMap, pts, AWS_END)
	line.Action = getString(lineMap, pts, AWS_ACTION)
	line.Status = getString(lineMap, pts, AWS_LOG_STATUS)

	done := []*AWSLogLine{}
	if line.StartTime.Before(time.Now().Add(-7 * 24 * time.Hour)) {
		log.Debugf("Bad v2 log line: %v | %v-> %s", lineMap, line.StartTime, strings.Join(pts, "|"))
	} else {
		done = append(done, &line)
	}

	return done, lineMap, nil
}

func getString(lineMap AwsLineMap, pts []string, key string) string {
	if i, ok := lineMap[key]; ok {
		return pts[i]
	}
	return ""
}

func getIP(lineMap AwsLineMap, pts []string, key string) net.IP {
	return net.ParseIP(getString(lineMap, pts, key))
}

func getUint32(lineMap AwsLineMap, pts []string, key string) uint32 {
	v, _ := strconv.Atoi(getString(lineMap, pts, key))
	return uint32(v)
}

func getUint64(lineMap AwsLineMap, pts []string, key string) uint64 {
	v, _ := strconv.Atoi(getString(lineMap, pts, key))
	return uint64(v)
}

func getTime(lineMap AwsLineMap, pts []string, key string) time.Time {
	v, _ := strconv.Atoi(getString(lineMap, pts, key))
	return time.Unix(int64(v), 0)
}

//version vpc-id subnet-id instance-id interface-id account-id type srcaddr dstaddr srcport dstport pkt-srcaddr pkt-dstaddr protocol bytes packets start end action tcp-flags log-status
func NewAwsFromV345(version int, lineMap AwsLineMap, pts []string, log logger.ContextL) ([]*AWSLogLine, AwsLineMap, error) {
	line := AWSLogLine{}
	line.Version = version
	line.VPCID = getString(lineMap, pts, AWS_VPC_ID)
	line.SubnetID = getString(lineMap, pts, AWS_SUBNET_ID)
	line.InstanceID = getString(lineMap, pts, AWS_INSTANCE_ID)
	line.InterfaceID = getString(lineMap, pts, AWS_INTERFACE_ID)
	line.AccountID = getString(lineMap, pts, AWS_ACCOUNT_ID)
	line.SrcAddr = getIP(lineMap, pts, AWS_SRC_ADDR)
	line.DstAddr = getIP(lineMap, pts, AWS_DST_ADDR)
	line.SrcPort = getUint32(lineMap, pts, AWS_SRC_PORT)
	line.DstPort = getUint32(lineMap, pts, AWS_DST_PORT)
	line.SrcPktAddr = getIP(lineMap, pts, AWS_PKT_SRC_ADDR)
	line.DstPktAddr = getIP(lineMap, pts, AWS_PKT_DST_ADDR)
	line.Protocol = getUint32(lineMap, pts, AWS_PROTOCOL)
	line.Bytes = getUint64(lineMap, pts, AWS_BYTES)
	line.Packets = getUint64(lineMap, pts, AWS_PACKETS)
	line.StartTime = getTime(lineMap, pts, AWS_START)
	line.EndTime = getTime(lineMap, pts, AWS_END)
	line.Action = getString(lineMap, pts, AWS_LINE_ACTION)
	line.TcpFlags = getUint32(lineMap, pts, AWS_TCP_FLAGS)
	line.Status = getString(lineMap, pts, AWS_LOG_STATUS)
	line.Region = getString(lineMap, pts, AWS_REGION)
	line.AzID = getString(lineMap, pts, AWS_AZ_ID)
	line.SublocationType = getString(lineMap, pts, AWS_SUBLOCATION_TYPE)
	line.SublocationID = getString(lineMap, pts, AWS_SUBLOCATION_ID)
	line.SrcPktService = getString(lineMap, pts, AWS_PKT_SRC_AWS_SERVICE)
	line.DstPktService = getString(lineMap, pts, AWS_PKT_DST_AWS_SERVICE)
	line.FlowDirection = getString(lineMap, pts, AWS_FLOW_DIRECTION)
	line.TrafficPath = getString(lineMap, pts, AWS_TRAFFIC_PATH)

	done := []*AWSLogLine{}
	if line.StartTime.Before(time.Now().Add(-7 * 24 * time.Hour)) {
		//log.Debugf("Bad v3/4/5 log line: %v | %v-> %s", lineMap, line.StartTime, strings.Join(pts, "|"))
		// Set start time to now.
		line.StartTime = time.Now()
		done = append(done, &line)
	} else {
		done = append(done, &line)
	}

	return done, lineMap, nil
}

func NewAwsHeader(pts []string) ([]*AWSLogLine, AwsLineMap, error) {
	lineMap := AwsLineMap{}
	colMap := map[string]int{}
	for i, col := range pts {
		colMap[col] = i
	}

	for _, field := range AWS_FLOW_FIELDS {
		if i, ok := colMap[field]; ok {
			lineMap[field] = i
		}
	}

	return nil, lineMap, nil
}

func NewAws(lineMap AwsLineMap, raw *string, log logger.ContextL) ([]*AWSLogLine, AwsLineMap, error) {
	pts := strings.Fields(*raw)
	if len(pts) < MIN_AWS_FIELD {
		return nil, lineMap, fmt.Errorf("Invalid line: %s", *raw)
	}

	// This is the header for CSV. Store?
	if pts[0] == "version" {
		return NewAwsHeader(pts)
	} else if len(lineMap) == 0 {
		if !strings.Contains(*raw, baseSplitAWS) {
			return NewAwsHeader(pts)
		}
	}

	ver := getUint32(lineMap, pts, AWS_VERSION)
	switch ver {
	case 2:
		return NewAwsFromV2(lineMap, pts, log)
	case 3:
		return NewAwsFromV345(3, lineMap, pts, log)
	case 4:
		return NewAwsFromV345(4, lineMap, pts, log)
	case 5:
		return NewAwsFromV345(5, lineMap, pts, log)
	default:
		// Try parsing as a Kinesis stream
		if strings.Contains(*raw, baseSplitAWS) {
			return NewAwsFromKinesis(lineMap, raw, log)
		} else { // Error here because we don't know version.
			if len(lineMap) == 0 {
				return nil, lineMap, fmt.Errorf("Bad log version: %d -> %v", ver, *raw)
			} else {
				// Go ahead and try to parse with v5 since we have a line map.
				return NewAwsFromV345(5, lineMap, pts, log)
			}
		}
	}
}

func (m *AWSLogLine) lookupTopo(ip net.IP, topo *AWSTopology) (subnet *ec2.Subnet, vpc *ec2.Vpc, az *ec2.AvailabilityZone) {
	if topo != nil {
		if ip.To4() == nil {
			address := patricia.NewIPv6Address(ip.To16(), 128)
			found, val := topo.Hierarchy.SubnetTrieV6.FindDeepestTag(address)
			if found {
				if s, ok := topo.Entities.Subnets[val]; ok {
					subnet = &s
					if v, ok := topo.Entities.Vpcs[*s.VpcId]; ok {
						vpc = &v
					}
					if a, ok := topo.Entities.AvailabilityZones[*s.AvailabilityZoneId]; ok {
						az = &a
					}
					return
				}
			}
		} else {
			address := patricia.NewIPv4AddressFromBytes(ip.To4(), 32)
			found, val := topo.Hierarchy.SubnetTrieV4.FindDeepestTag(address)
			if found {
				if s, ok := topo.Entities.Subnets[val]; ok {
					subnet = &s
					if v, ok := topo.Entities.Vpcs[*s.VpcId]; ok {
						vpc = &v
					}
					if a, ok := topo.Entities.AvailabilityZones[*s.AvailabilityZoneId]; ok {
						az = &a
					}
					return
				}
			}
		}
	}

	return
}

func (m *AWSLogLine) ToFlow(log logger.ContextL, topo *AWSTopology) (in *kt.JCHF) {

	in = kt.NewJCHF()
	in.CustomStr = make(map[string]string)
	in.CustomInt = make(map[string]int32)
	in.CustomBigInt = make(map[string]int64)
	in.EventType = kt.KENTIK_EVENT_TYPE
	in.Provider = kt.ProviderVPC
	in.Timestamp = m.StartTime.Unix()
	in.InBytes = m.Bytes
	in.InPkts = m.Packets
	in.L4DstPort = m.DstPort
	in.L4SrcPort = m.SrcPort
	in.Protocol = ic.PROTO_NAMES[uint16(m.Protocol)]
	in.SampleRate = 1
	in.TcpFlags = m.TcpFlags
	in.CustomStr["action"] = m.Action
	in.CustomStr["status"] = m.Status
	in.CustomStr["kt.from"] = kt.FromLambda
	in.CustomStr["type"] = AWS_VPC_TYPE
	in.DeviceName = m.VPCID

	if m.Sample > 0 { // Set sample rate here if we are switching.
		in.SampleRate = m.Sample
	}

	in.SrcAddr = getAddr(m.SrcAddr)
	in.DstAddr = getAddr(m.DstAddr)
	in.CustomStr["source_pkt_addr"] = getAddr(m.SrcPktAddr)
	in.CustomStr["dest_pkt_addr"] = getAddr(m.DstPktAddr)
	in.CustomStr["sublocation_type"] = m.SublocationType
	in.CustomStr["sublocation_id"] = m.SublocationID
	in.CustomStr["source_pkt_service"] = m.SrcPktService
	in.CustomStr["dest_pkt_service"] = m.DstPktService
	in.CustomStr["flow_direction"] = m.FlowDirection
	in.CustomStr["traffic_path"] = m.TrafficPath
	in.CustomStr["account_id"] = m.AccountID
	in.CustomBigInt["start_time"] = m.StartTime.Unix()
	in.CustomBigInt["end_time"] = m.EndTime.Unix()

	// Do we know anything more about this conversation?
	srcSubnet, srcVpc, srcAz := m.lookupTopo(m.SrcAddr, topo)
	dstSubnet, dstVpc, dstAz := m.lookupTopo(m.DstAddr, topo)

	// The rest is set up into src and dst parts.
	if m.FlowDirection == "egress" {
		in.OutBytes = in.InBytes // Move these to out since its egress.
		in.OutPkts = in.InPkts
		in.InBytes = 0 // And 0 these out.
		in.InPkts = 0

		in.CustomStr["source_vpc"] = m.VPCID
		in.CustomStr["source_subnet"] = m.SubnetID
		in.CustomStr["source_instance"] = m.InstanceID
		in.CustomStr["source_interface"] = m.InterfaceID
		in.CustomStr["source_az"] = m.AzID
		in.CustomStr["source_region"] = m.Region

		if dstSubnet != nil {
			in.CustomStr["dest_vpc"] = *dstSubnet.VpcId
			in.CustomStr["dest_az"] = *dstAz.ZoneName
			in.CustomStr["dest_region"] = *dstAz.RegionName
		}
		if srcVpc != nil {
			for _, tag := range srcVpc.Tags {
				in.CustomStr[*tag.Key] = *tag.Value
			}
		}
		if srcSubnet != nil {
			for _, tag := range srcSubnet.Tags {
				in.CustomStr[*tag.Key] = *tag.Value
			}
		}
	} else {
		in.CustomStr["dest_vpc"] = m.VPCID
		in.CustomStr["dest_subnet"] = m.SubnetID
		in.CustomStr["dest_instance"] = m.InstanceID
		in.CustomStr["dest_interface"] = m.InterfaceID
		in.CustomStr["dest_az"] = m.AzID
		in.CustomStr["dest_region"] = m.Region

		if srcSubnet != nil {
			in.CustomStr["source_vpc"] = *srcSubnet.VpcId
			in.CustomStr["source_az"] = *srcAz.ZoneName
			in.CustomStr["source_region"] = *srcAz.RegionName
		}
		if dstVpc != nil {
			for _, tag := range dstVpc.Tags {
				in.CustomStr[*tag.Key] = *tag.Value
			}
		}
		if dstSubnet != nil {
			for _, tag := range dstSubnet.Tags {
				in.CustomStr[*tag.Key] = *tag.Value
			}
		}
	}

	// Now add some combo fields.
	in.CustomStr["src_endpoint"] = in.SrcAddr + ":" + strconv.Itoa(int(in.L4SrcPort))
	in.CustomStr["dst_endpoint"] = in.DstAddr + ":" + strconv.Itoa(int(in.L4DstPort))

	return in
}

type FlowSet struct {
	Bucket string
	Key    string
	Lines  []*AWSLogLine `json:"lines"`
}

func (fs *FlowSet) ProcessKey(bucket string, key string) error {
	fs.Bucket = bucket
	fs.Key = key
	return nil
}

// What is the kentik name of this device?
func (fs *FlowSet) GetDeviceKey() (string, error) {
	return getDeviceKey(fs.Bucket, fs.Key)
}

// bucket_ARN/optional_folder/AWSLogs/aws_account_id/vpcflowlogs/region/year/month/day/log_file_name.log.gz
// if optional_folder set, use this, else fall back on bucket name
func getDeviceKey(bucket, path string) (string, error) {
	// fast path: no folder
	if strings.HasPrefix(path, AWS_LOG_PREFIX) {
		return bucket, nil
	}

	var keybuilder strings.Builder
	for i, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, AWS_LOG_PREFIX) {
			return keybuilder.String(), nil
		} else if part != "" {
			if i > 0 {
				keybuilder.WriteRune('/')
			}
			keybuilder.WriteString(part)
		}
	}

	// Just use bucket, flow is right there.
	pts := strings.Split(path, "-")
	if len(pts) >= 7 {
		return bucket, nil
	}
	return "", fmt.Errorf("incorrectly formatted path [%s]", path)
}

func getAddr(addr net.IP) string {
	if addr != nil {
		return addr.String()
	}
	return ""
}
