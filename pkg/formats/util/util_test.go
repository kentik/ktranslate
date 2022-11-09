package util

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestDropOnFilter(t *testing.T) {
	tests := []struct {
		attr    map[string]interface{}
		in      *kt.JCHF
		metrics map[string]kt.MetricInfo
		lm      kt.LastMetadata
		drop    bool
	}{
		{
			map[string]interface{}{},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{},
			false,
		},
		{
			map[string]interface{}{
				"foo": "bar",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo": regexp.MustCompile("bar"),
				},
			},
			false,
		},
		{
			map[string]interface{}{
				"foo": "ba11",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo": regexp.MustCompile("bar"),
				},
			},
			true,
		},
		{
			map[string]interface{}{
				"foo": "ba11",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo": regexp.MustCompile("^ba"),
				},
			},
			false,
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "down",
				"fooXX":        "bar",
			},
			kt.NewJCHF().SetIFPorts(10),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"fooXX":        regexp.MustCompile("bar"),
					kt.AdminStatus: regexp.MustCompile("up"),
				},
				InterfaceInfo: map[kt.IfaceID]map[string]interface{}{
					10: map[string]interface{}{
						"Description": "myIfDesc",
					},
				},
			},
			true,
		},
		{ // 5
			map[string]interface{}{
				kt.AdminStatus: "up",
				"foo":          "bar",
				"aaa":          "aaa",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo":          regexp.MustCompile("abar"),
					"aaa":          regexp.MustCompile("aaa"),
					kt.AdminStatus: regexp.MustCompile("up"),
				},
			},
			false,
		},
		{ // 6
			map[string]interface{}{
				kt.AdminStatus: "up",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"fooAAA":       regexp.MustCompile("abar"),
					kt.AdminStatus: regexp.MustCompile("up"),
				},
			},
			false, // Let through because status is up and fooAAA doesn't exist in the attribute list.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					kt.AdminStatus: regexp.MustCompile("up"),
				},
			},
			false,
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
				"foo":          "bar",
				"aaa":          "aaa",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo":          regexp.MustCompile("no"),
					"aaa":          regexp.MustCompile("no"),
					kt.AdminStatus: regexp.MustCompile("up"),
				},
			},
			true, // Drop because neither foo or aaa match even though admin is up.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
				"foo":          "bar",
				"aaa":          "aaa",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"foo":          regexp.MustCompile("no"),
					"aaa":          regexp.MustCompile("aa"),
					kt.AdminStatus: regexp.MustCompile("up"),
				},
			},
			false, // Keep because aaa matches and admin is up.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
			},
			kt.NewJCHF().SetIFPorts(20),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"if_Description": regexp.MustCompile("igb3"),
					kt.AdminStatus:   regexp.MustCompile("up"),
					"device_name":    regexp.MustCompile("bart"),
				},
				InterfaceInfo: map[kt.IfaceID]map[string]interface{}{
					20: map[string]interface{}{
						"Description": "igb2",
					},
				},
			},
			true, // Drop because desc doesn't match.
		},
		{
			map[string]interface{}{
				"mib-name": "UDP-MIB",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"if_Description": regexp.MustCompile("igb3"),
					kt.AdminStatus:   regexp.MustCompile("up"),
					"mib-name":       regexp.MustCompile("UDP"),
				},
			},
			false, // keep because mib-name matches and no admin status.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
				"if_Alias":     "foo",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!if_Description": regexp.MustCompile("igb3"),
					kt.AdminStatus:    regexp.MustCompile("up"),
				},
			},
			true, // drop because missing desciption.
		},
		{
			map[string]interface{}{
				kt.AdminStatus:   "up",
				"if_Alias":       "foo",
				"if_Description": "igb3",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!if_Description": regexp.MustCompile("igb3"),
					kt.AdminStatus:    regexp.MustCompile("up"),
				},
			},
			false, // keep because matching desciption.
		},
		{
			map[string]interface{}{
				kt.AdminStatus:   "up",
				"if_Description": "igb4",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!if_Alias":      regexp.MustCompile("igb3"),
					"if_Description": regexp.MustCompile("igb4"),
					kt.AdminStatus:   regexp.MustCompile("up"),
				},
			},
			true, // drop because alias is missing.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
				"if_Alias":     "igb4",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!if_Alias":       regexp.MustCompile("igb4"),
					"!if_Description": regexp.MustCompile("igb4"),
					kt.AdminStatus:    regexp.MustCompile("up"),
				},
			},
			true, // drop because description is missing.
		},
		{
			map[string]interface{}{
				kt.AdminStatus:   "up",
				"if_Description": "igb4",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!(if_Description||if_Alias)": regexp.MustCompile("igb4"),
					kt.AdminStatus:                regexp.MustCompile("up"),
				},
			},
			false, // keep because alias or description match.
		},
		{
			map[string]interface{}{
				kt.AdminStatus:   "up",
				"if_Description": "igb5",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"!(if_Alias||if_Description)": regexp.MustCompile("igb5"),
					kt.AdminStatus:                regexp.MustCompile("up"),
				},
			},
			false, // keep because alias or description match. Test the other way around.
		},
		{
			map[string]interface{}{
				kt.AdminStatus:      "up",
				"if_interface_name": "tu22",
				"if_Alias":          "some other alias",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"DOES_NOT_MATCH":    regexp.MustCompile("true"),
					"if_interface_name": regexp.MustCompile("^(lo|po|vl|tu)"),
					"if_Alias":          regexp.MustCompile("unused|reserved|desktop"),
					kt.AdminStatus:      regexp.MustCompile("up"),
				},
			},
			true, // drop because we are negating alias or description match.
		},
		{
			map[string]interface{}{
				kt.AdminStatus: "up",
			},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{
				MatchAttr: map[string]*regexp.Regexp{
					"DOES_NOT_MATCH":    regexp.MustCompile("true"),
					"if_interface_name": regexp.MustCompile("^(lo|po|vl|tu)"),
					"if_Alias":          regexp.MustCompile("unused|reserved|desktop"),
					kt.AdminStatus:      regexp.MustCompile("up"),
				},
			},
			true, // drop because we are negating alias or description match and alias doesn't exist.
		},
	}

	for i, test := range tests {
		SetAttr(test.attr, test.in, test.metrics, &test.lm, false)
		isIf := false
		for k, _ := range test.attr {
			if strings.HasPrefix(k, "if_") {
				isIf = true
			}
		}
		drop := DropOnFilter(test.attr, &test.lm, isIf)
		assert.Equal(t, test.drop, drop, "Test %d %v", i, test.lm.MatchAttr)
	}
}

func TestCopyAttrforSNMP(t *testing.T) {
	assert := assert.New(t)

	input := map[string]interface{}{}
	for i := 0; i < 10; i++ {
		input[fmt.Sprintf("XXX%d", i)] = i
	}
	name := kt.MetricInfo{Oid: "oid", Mib: "mib"}

	res := CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal(len(input)+3, len(res)) // adds in three keys
	assert.Equal("oid", res["objectIdentifier"])

	for i := 0; i < MAX_ATTR_FOR_SNMP+10; i++ {
		input[fmt.Sprintf("XXX%d", i)] = i
	}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal(MAX_ATTR_FOR_SNMP, len(res)) // truncated at MAX_ATTR_FOR_SNMP
	assert.Equal("oid", res["objectIdentifier"])

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "noMatch"}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal(nil, res["foo"], res)

	input = map[string]interface{}{kt.StringPrefix + "foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "foo"}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{"foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "foo"}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal("one", res["foo"], res)

	input = map[string]interface{}{"foo": "one"}
	name = kt.MetricInfo{Oid: "oid", Mib: "mib", Table: "somethingElse"}
	res = CopyAttrForSnmp(input, "test", name, nil, true)
	assert.Equal(nil, res["foo"], res)
}
