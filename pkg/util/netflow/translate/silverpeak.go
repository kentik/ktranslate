package translate

func init() {

	// Silverpeak key
	builtin[Key{23867, 1}] = InformationElementEntry{FieldID: 1, Name: "silverpeakClientIpv4Address", Type: FieldTypes["ipv4Address"]}
	builtin[Key{23867, 2}] = InformationElementEntry{FieldID: 2, Name: "silverpeakServerIpv4Address", Type: FieldTypes["ipv4Address"]}
	builtin[Key{23867, 3}] = InformationElementEntry{FieldID: 3, Name: "silverpeakConnectionServerOctetDeltaCount", Type: FieldTypes["unsigned64"]}
	builtin[Key{23867, 4}] = InformationElementEntry{FieldID: 4, Name: "silverpeakConnectionServerPacketDeltaCount", Type: FieldTypes["unsigned64"]}
	builtin[Key{23867, 5}] = InformationElementEntry{FieldID: 5, Name: "silverpeakConnectionClientOctetDeltaCount", Type: FieldTypes["unsigned64"]}
	builtin[Key{23867, 6}] = InformationElementEntry{FieldID: 6, Name: "silverpeakConnectionClientPacketDeltaCount", Type: FieldTypes["unsigned64"]}
	builtin[Key{23867, 7}] = InformationElementEntry{FieldID: 7, Name: "silverpeakConnectionInitiator", Type: FieldTypes["ipv4Address"]}
	builtin[Key{23867, 8}] = InformationElementEntry{FieldID: 8, Name: "silverpeakApplicationHttHost", Type: FieldTypes["string"]}
	builtin[Key{23867, 9}] = InformationElementEntry{FieldID: 9, Name: "silverpeakConnectionNumberOfConnections", Type: FieldTypes["unsigned8"]}
	builtin[Key{23867, 10}] = InformationElementEntry{FieldID: 10, Name: "silverpeakConnectionServerResponsesCount", Type: FieldTypes["unsigned8"]}
	builtin[Key{23867, 11}] = InformationElementEntry{FieldID: 11, Name: "silverpeakConnectionServerResponseDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 12}] = InformationElementEntry{FieldID: 12, Name: "silverpeakConnectionNetworkToServerDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 13}] = InformationElementEntry{FieldID: 13, Name: "silverpeakConnectionNetworkToClientDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 14}] = InformationElementEntry{FieldID: 14, Name: "silverpeakConnectionClientPacketRetransmissionCount", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 15}] = InformationElementEntry{FieldID: 15, Name: "silverpeakConnectionClientToServerNetworkDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 16}] = InformationElementEntry{FieldID: 16, Name: "silverpeakConnectionApplicationDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 17}] = InformationElementEntry{FieldID: 17, Name: "silverpeakConnectionClientToServerResponseDelay", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 18}] = InformationElementEntry{FieldID: 18, Name: "silverpeakConnectionTransactionDuration", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 19}] = InformationElementEntry{FieldID: 19, Name: "silverpeakConnectionTransactionDurationMin", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 20}] = InformationElementEntry{FieldID: 20, Name: "silverpeakConnectionTransactionDurationMax", Type: FieldTypes["unsigned32"]}
	builtin[Key{23867, 21}] = InformationElementEntry{FieldID: 21, Name: "silverpeakConnectionTransactionCompleteCount", Type: FieldTypes["unsigned8"]}
	builtin[Key{23867, 22}] = InformationElementEntry{FieldID: 22, Name: "silverpeakFromZone", Type: FieldTypes["string"]}
	builtin[Key{23867, 23}] = InformationElementEntry{FieldID: 23, Name: "silverpeakToZone", Type: FieldTypes["string"]}
	builtin[Key{23867, 24}] = InformationElementEntry{FieldID: 24, Name: "silverpeakTag", Type: FieldTypes["string"]}
	builtin[Key{23867, 25}] = InformationElementEntry{FieldID: 25, Name: "silverpeakOverlay", Type: FieldTypes["string"]}
	builtin[Key{23867, 26}] = InformationElementEntry{FieldID: 26, Name: "silverpeakDirection", Type: FieldTypes["string"]}
	builtin[Key{23867, 27}] = InformationElementEntry{FieldID: 27, Name: "silverpeakApplicationCategory", Type: FieldTypes["string"]}
}
