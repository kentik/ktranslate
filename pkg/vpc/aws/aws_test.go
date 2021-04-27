package aws

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/stretchr/testify/assert"
)

var (
	AWS_V4_FIELDS = map[string]bool{
		AWS_REGION:           true,
		AWS_AZ_ID:            true,
		AWS_SUBLOCATION_TYPE: true,
		AWS_SUBLOCATION_ID:   true,
	}
)

func TestGetDeviceKeyDefault(t *testing.T) {
	key, err := getDeviceKey("mybucket", AWS_LOG_PREFIX+"/stuff.gz")
	assert.Equal(t, key, "mybucket")
	assert.NoError(t, err)
}
func TestGetDeviceKeyFolder(t *testing.T) {
	key, err := getDeviceKey("mybucket", "cool/"+AWS_LOG_PREFIX+"/stuff.gz")
	assert.Equal(t, key, "cool")
	assert.NoError(t, err)
}

func TestGetDeviceKeySubFolder(t *testing.T) {
	key, err := getDeviceKey("mybucket", "cool/lol/"+AWS_LOG_PREFIX+"/stuff.gz")
	assert.Equal(t, key, "cool/lol")
	assert.NoError(t, err)
}

func TestGetDeviceKeyBadFormat(t *testing.T) {
	_, err := getDeviceKey("mybucket", "stuff.gz")
	assert.Error(t, err, "incorrectly formatted")
}

func TestGetDeviceKeyPan(t *testing.T) {
	_, err := getDeviceKey("mybucket", "VPC-Log-Export-us-east-1/")
	assert.Error(t, err)
}

func TestGetDeviceKeyExodus(t *testing.T) {
	_, err := getDeviceKey("mybucket", "2019-05-02-18-27-30-80DD9ED25F98887C")
	assert.NoError(t, err)
}

func TestGetDeviceKeyKineses(t *testing.T) {
	key, err := getDeviceKey("abd-vpc-flow-logs-nonprod", "2019/10/11/19/VPC-Flow-Logs-NonProd-KinesesFH-1-2019-10-11-19-19-28-3c0b7660-77df-4a82-9257-ed9e1c79f70d")
	assert.NoError(t, err)
	assert.Equal(t, key, "abd-vpc-flow-logs-nonprod")
}

func TestParseFlows(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	now := time.Now().Unix()
	input := []string{
		"version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status",
		fmt.Sprintf("2 672740005165 eni-ba88b581 54.191.55.41 10.0.0.53 80 41536 6 27368 39574687 %d %d ACCEPT OK", now, now+300),
		fmt.Sprintf("2 672740005165 eni-ba88b581 10.0.0.53 54.218.137.160 52772 80 6 62 4451 %d %d ACCEPT OK", now, now+300),
		fmt.Sprintf("2 672740005165 eni-ba88b581 159.89.148.43 10.0.0.53 56911 5984 6 1 40 %d %d REJECT OK", now, now+300),
		fmt.Sprintf("2 672740005165 eni-ba88b581 159.89.148.43 10.0.0.53 35789 8545 6 1 40 %d %d REJECT OK", now, now+300),
		"version vpc-id subnet-id instance-id interface-id account-id type srcaddr dstaddr srcport dstport pkt-srcaddr pkt-dstaddr protocol bytes packets start end action tcp-flags log-status",
		fmt.Sprintf("3 vpc-2e0d3a4a subnet-0a279552 - eni-22859f22 672740005165 IPv4 10.232.50.235 10.232.51.166 62527 27016 10.232.50.235 18.162.196.19 6 120 2 %d %d ACCEPT 2 OK", now, now+300),
		fmt.Sprintf("3 vpc-6081b704 subnet-feb1cc88 - eni-05ed3ab55eb4374f8 672740005165 IPv4 159.203.161.141 10.232.83.71 46017 8088 159.203.161.141 10.232.83.71 6 120 2 %d %d REJECT 2 OK", now, now+300),
		"version vpc-id subnet-id instance-id interface-id account-id type srcaddr dstaddr srcport dstport pkt-srcaddr pkt-dstaddr protocol bytes packets start end action tcp-flags log-status region az-id sublocation-type sublocation-id",
		fmt.Sprintf("4 vpc-6081b704 subnet-feb1cc88 - eni-05ed3ab55eb4374f8 672740005165 IPv4 159.203.161.141 10.232.83.71 46017 8088 159.203.161.141 10.232.83.71 6 120 2 %d %d REJECT 2 OK us-east 2323232 outpost 234234234", now, now+300),
		"version vpc-id subnet-id instance-id interface-id account-id type srcaddr dstaddr srcport dstport pkt-srcaddr pkt-dstaddr protocol bytes packets start pkt-src-aws-service action tcp-flags log-status region az-id sublocation-type sublocation-id",
		fmt.Sprintf("5 vpc-6081b704 subnet-feb1cc88 - eni-05ed3ab55eb4374f8 672740005165 IPv4 159.203.161.141 10.232.83.71 46017 8088 159.203.161.141 10.232.83.71 6 120 2 %d AMAZON REJECT 2 OK us-east 2323232 outpost 234234234", now),
	}

	lineMap := AwsLineMap{}
	for _, i := range input {
		res, lineMapOut, err := NewAws(lineMap, &i, l)
		assert.NoError(t, err)
		lineMap = lineMapOut
		if strings.HasPrefix(i, "version") {
			// Version case here. noop.
			assert.Equal(t, 0, len(res))
			assert.True(t, len(lineMap) >= 14, "%v", lineMap)
			for i, field := range AWS_FLOW_FIELDS {
				_, ok := lineMap[field]
				if len(lineMap) == 14 { // v2 case
					if i <= 13 {
						assert.True(t, ok, "missing v2 line field %s %d", field, i)
					}
				} else {
					if len(lineMap) == 14 { // v2 case
						if i <= 13 {
							assert.True(t, ok, "missing v2 line field %s %d", field, i)
						}
					} else {
						// v3 case
						if lineMap[AWS_VERSION] == 3 && !AWS_V4_FIELDS[field] {
							assert.True(t, ok, "missing v3 line field %s", field)
						} else if lineMap[AWS_VERSION] == 4 {
							assert.True(t, ok, "missing v4 line field %s", field)
						}
					}
				}
			}
		} else {
			assert.Equal(t, 1, len(res))
			assert.Equal(t, "672740005165", res[0].AccountID)
			if res[0].Version == 4 {
				assert.Equal(t, int(2), int(res[0].TcpFlags))
				assert.Equal(t, int(6), int(res[0].Protocol))
				assert.Equal(t, int(120), int(res[0].Bytes))
				assert.Equal(t, int(2), int(res[0].Packets))
				assert.Equal(t, "us-east", res[0].Region)
				assert.Equal(t, "2323232", res[0].AzID)
			} else if res[0].Version == 5 {
				assert.Equal(t, int(2), int(res[0].TcpFlags))
				assert.Equal(t, int(6), int(res[0].Protocol))
				assert.Equal(t, int(120), int(res[0].Bytes))
				assert.Equal(t, int(2), int(res[0].Packets))
				assert.Equal(t, "AMAZON", res[0].SrcPktService)
				assert.Equal(t, "2323232", res[0].AzID)
			} else if res[0].Version == 3 {
				assert.Equal(t, int(2), int(res[0].TcpFlags))
				assert.Equal(t, int(6), int(res[0].Protocol))
				assert.Equal(t, int(120), int(res[0].Bytes))
				assert.Equal(t, int(2), int(res[0].Packets))
			} else {
				assert.Equal(t, int(2), int(res[0].Version))
				if res[0].DstPort == 80 {
					assert.Equal(t, int(6), int(res[0].Protocol))
					assert.Equal(t, int(4451), int(res[0].Bytes))
					assert.Equal(t, int(62), int(res[0].Packets))
				}
			}
			assert.Equal(t, "OK", res[0].Status)
			assert.True(t, res[0].DstPort != 0)
			assert.True(t, res[0].SrcPort != 0)
			assert.True(t, res[0].Packets != 0)
			assert.True(t, res[0].Bytes != 0)
		}
	}
}

func TestParseFlowKinesis(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	now := time.Now().Unix()
	lineMap := AwsLineMap{}
	input := fmt.Sprintf(`{"messageType": "DATA_MESSAGE",
    "owner": "111111111111",
    "logGroup": "CloudTrail",
    "logStream": "111111111111_CloudTrail_us-east-1",
    "subscriptionFilters": [
        "Destination"
    ],
    "logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"6","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300)

	res, lineMap, err := NewAws(lineMap, &input, l)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "391389995465", res[0].AccountID)
	assert.Equal(t, int(0), int(res[0].Sample))
}

func TestParseFlowKinesisWithSampling(t *testing.T) {
	l := lt.NewTestContextL(logger.NilContext, t)
	now := time.Now().Unix()
	lineMap := AwsLineMap{}
	inputs := []string{
		fmt.Sprintf(`{"messageType": "DATA_MESSAGE","owner": "111111111111","logGroup": "CloudTrail","logStream": "111111111111_CloudTrail_us-east-1","subscriptionFilters": ["Destination"],"logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"6","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300),
		fmt.Sprintf(`{"messageType": "DATA_MESSAGE","owner": "111111111111","logGroup": "CloudTrail","logStream": "111111111111_CloudTrail_us-east-1","subscriptionFilters": ["Destination"],"logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"7","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300),
		fmt.Sprintf(`{"messageType": "DATA_MESSAGE","owner": "111111111111","logGroup": "CloudTrail","logStream": "111111111111_CloudTrail_us-east-1","subscriptionFilters": ["Destination"],"logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"8","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300),
		fmt.Sprintf(`{"messageType": "DATA_MESSAGE","owner": "111111111111","logGroup": "CloudTrail","logStream": "111111111111_CloudTrail_us-east-1","subscriptionFilters": ["Destination"],"logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"9","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300),
	}

	together := strings.Join(inputs, "")
	res, lineMap, err := NewAws(lineMap, &together, l)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, int(6), int(res[0].Protocol))
	assert.Equal(t, int(0), int(res[0].Sample))

	inputs = []string{
		fmt.Sprintf(`{"messageType": "DATA_MESSAGE","owner": "111111111111","logGroup": "CloudTrail","logStream": "111111111111_CloudTrail_us-east-1","subscriptionFilters": ["Destination"],"logEvents": [{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"6","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}},{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"7","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}},{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"8","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}},{"id":"35036294238826492806306218206507857848794444175103892012","timestamp":1571081770000,"message":"2 391389995465 eni-0939c7c9e1255db73 10.236.54.140 10.236.57.28 31547 27068 6 2 112 1571081770 1571081799 ACCEPT OK","extractedFields":{"srcaddr":"10.236.54.140","dstport":"27068","start":"%d","dstaddr":"10.236.57.28","version":"2","packets":"2","protocol":"9","account_id":"391389995465","interface_id":"eni-0939c7c9e1255db73","log_status":"OK","bytes":"112","srcport":"31547","action":"ACCEPT","end":"%d"}}]}`, now, now+300, now, now+300, now, now+300, now, now+300),
	}

	together = strings.Join(inputs, "")
	res, lineMap, err = NewAws(lineMap, &together, l)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, int(7), int(res[1].Protocol))
	assert.Equal(t, int(0), int(res[0].Sample))
}
