package translate

// gigamon values -- it would be nice if clients of the netflow library could push their own keys into translate,
// rather than having to do it globally
func init() {
	builtin[Key{26866, 1}] = InformationElementEntry{FieldID: 1, Name: "gigamonHttpReqUrl", Type: FieldTypes["string"]}
	builtin[Key{26866, 2}] = InformationElementEntry{FieldID: 2, Name: "gigamonHttpRspStatus", Type: FieldTypes["unsigned16"]}
	builtin[Key{26866, 202}] = InformationElementEntry{FieldID: 202, Name: "gigamonDnsOpCode", Type: FieldTypes["unsigned8"]}
	builtin[Key{26866, 212}] = InformationElementEntry{FieldID: 212, Name: "dnsNsCount", Type: FieldTypes["unsigned16"]}
	builtin[Key{26866, 220}] = InformationElementEntry{FieldID: 220, Name: "dnsAuthorityName", Type: FieldTypes["string"]}
	builtin[Key{26866, 221}] = InformationElementEntry{FieldID: 221, Name: "dnsAuthorityType", Type: FieldTypes["unsigned16"]}
}
