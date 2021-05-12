package translate

func init() {

	// Versa key
	builtin[Key{42359, 519}] = InformationElementEntry{FieldID: 519, Name: "versaAppID", Type: FieldTypes["unsigned32"]}
	builtin[Key{42359, 522}] = InformationElementEntry{FieldID: 522, Name: "versaTenantID", Type: FieldTypes["unsigned16"]}
	builtin[Key{42359, 540}] = InformationElementEntry{FieldID: 540, Name: "versaEventType", Type: FieldTypes["unsigned16"]}
	builtin[Key{42359, 574}] = InformationElementEntry{FieldID: 574, Name: "versaApplianceID", Type: FieldTypes["unsigned16"]}
}
