package util

import (
	"regexp"
	"testing"

	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/stretchr/testify/assert"
)

func TestSetAttrDrop(t *testing.T) {
	tests := []struct {
		attr    map[string]interface{}
		in      *kt.JCHF
		metrics map[string]kt.MetricInfo
		lm      kt.LastMetadata
		drop    interface{}
	}{
		{
			map[string]interface{}{},
			kt.NewJCHF(),
			map[string]kt.MetricInfo{},
			kt.LastMetadata{},
			nil,
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
			true,
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
	}

	for i, test := range tests {
		SetAttr(test.attr, test.in, test.metrics, &test.lm)
		assert.Equal(t, test.drop, test.attr[kt.DropMetric], "Test %d", i)
	}
}
