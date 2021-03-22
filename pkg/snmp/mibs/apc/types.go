package apc

import (
	"encoding/xml"
)

type Oid struct {
	Oid    string   `xml:"oid,attr"`
	IsTree string   `xml:"istree,attr"`
	RuleId string   `xml:"ruleid,attr"`
	Values []string `xml:"valueList>value"`
}

type Rule struct {
	RuleId     string `xml:"ruleid,attr"`
	SuppressId string `xml:"suppressid,attr"`
}

type Val struct {
	Oid string `xml:"getOid"`
}

type Value struct {
	Val
	Id       string `xml:"mapid,attr"`
	IsBinary string `xml:"is-binary"`
	Mults    []Val  `xml:"mult>op"`
}

type Sensor struct {
	RuleId    string `xml:"ruleid,attr"`
	Type      string `xml:"type"`
	Value     Value  `xml:"value"`
	SensorId  string `xml:"sensorId"`
	Label     string `xml:"label"`
	Enum      string `xml:"enum"`
	SensorSet string `xml:"sensorSet"`
}

type ValueMap struct {
	Value
	RuleId    string   `xml:"ruleid,attr"`
	ValuesIn  []string `xml:"valueIn"`
	ValuesOut []string `xml:"valueOut"`
}

type EnumMap struct {
	RuleId string   `xml:"ruleid,attr"`
	Labels []string `xml:"label"`
}

type AlarmFlag struct {
	RuleId   string   `xml:"ruleid,attr"`
	ValueMap ValueMap `xml:"value>mapValue"`
}

type Op struct {
	Value
}

type SetData struct {
	RuleId string `xml:"ruleid,attr"`
	Field  string `xml:"field,attr"`
	ReOps  []Op   `xml:"regex>op"`
	Value
}

type Device struct {
	Id              string      `xml:"deviceid,attr"`
	OidMustExist    []Oid       `xml:"oidMustExist"`
	OidMustMatch    []Oid       `xml:"oidMustMatch"`
	SuppressRule    Rule        `xml:"suppressRule"`
	SetProductData  []SetData   `xml:"setProductData"`
	SetLocationData []SetData   `xml:"setLocationData"`
	StateSensors    []Sensor    `xml:"stateSensor"`
	NumSensors      []Sensor    `xml:"numSensor"`
	ValueMaps       []ValueMap  `xml:"valueMap"`
	EnumMaps        []EnumMap   `xml:"enumMap"`
	AlarmFlags      []AlarmFlag `xml:"alarmFlags"`
}

type APC struct {
	XMLName xml.Name `xml:"APC_DDF"`
	Devices []Device `xml:"device"`
}
