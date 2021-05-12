package translate

func init() {

	// Palo Alto v9 fields
	builtin[Key{0, 56701}] = InformationElementEntry{FieldID: 56701, Name: "paloAltoApplicationID", Type: FieldTypes["string"]}
	builtin[Key{0, 56702}] = InformationElementEntry{FieldID: 56701, Name: "paloAltoUserID", Type: FieldTypes["string"]}
}
