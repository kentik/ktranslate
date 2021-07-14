package rule

// Match holds IP/ASN match results in a bitmask
// - can match multiple types
type Match uint32

// MatchNone is a Match that represents no match
var MatchNone Match = 0

// MatchPrivateIP is a Match that represents a match from a static list of private IPs
var MatchPrivateIP Match = 1

// MatchPrivateASN is a Match that represents a match from a static list of private ASNs (non-cloud)
var MatchPrivateASN Match = 7

// NewMatch returns a new Match
func NewMatch(ruleMatch uint32) Match {
	return Match(ruleMatch)
}

// Uint32 returns a uint32 version
func (m Match) Uint32() uint32 {
	return uint32(m)
}
