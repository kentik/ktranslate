package translate

// nprobe values -- it would be nice if clients of the netflow library could push their own keys into translate,
// rather than having to do it globally
func init() {
	builtin[Key{0, 57581}] = InformationElementEntry{FieldID: 57581, Name: "RETRANSMITTED_IN_PKTS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57582}] = InformationElementEntry{FieldID: 57582, Name: "RETRANSMITTED_OUT_PKTS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57583}] = InformationElementEntry{FieldID: 57583, Name: "OOORDER_IN_PKTS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57584}] = InformationElementEntry{FieldID: 57584, Name: "OOORDER_OUT_PKTS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57552}] = InformationElementEntry{FieldID: 57552, Name: "FRAGMENTS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57595}] = InformationElementEntry{FieldID: 57595, Name: "CLIENT_NW_LATENCY_MS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57596}] = InformationElementEntry{FieldID: 57596, Name: "SERVER_NW_LATENCY_MS", Type: FieldTypes["unsigned32"]}
	builtin[Key{0, 57597}] = InformationElementEntry{FieldID: 57597, Name: "APPL_LATENCY_MS", Type: FieldTypes["unsigned32"]}
}
