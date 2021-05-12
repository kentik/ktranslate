package translate

// CERT Assigned (PEN 6871), see https://tools.netsa.cert.org/silk/faq.html#ipfix-fields
func init() {
	builtin[Key{6871, 14}] = InformationElementEntry{FieldID: 14, Name: "initialTCPFlags", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 15}] = InformationElementEntry{FieldID: 15, Name: "unionTCPFlags", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 21}] = InformationElementEntry{FieldID: 21, Name: "reverseFlowDeltaMilliseconds", Type: FieldTypes["dateTimeMilliseconds"]}
	builtin[Key{6871, 22}] = InformationElementEntry{FieldID: 22, Name: "silkTCPState", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 30}] = InformationElementEntry{FieldID: 30, Name: "silkFlowType", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 31}] = InformationElementEntry{FieldID: 31, Name: "silkFlowSensor", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 33}] = InformationElementEntry{FieldID: 33, Name: "silkAppLabel", Type: FieldTypes["unsigned8"]}
	builtin[Key{6871, 40}] = InformationElementEntry{FieldID: 40, Name: "flowAttributes", Type: FieldTypes["unsigned16"]}
}
