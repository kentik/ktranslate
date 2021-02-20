package patricia

// Base class for trees

import (
	"github.com/kentik/golog/logger"
)

// Trees is a pair of IPv4/IPv6 patricia trees
type Trees struct {
	log    *logger.Logger
	Length int
}

// Close cloess up anything that needs to be closed
func (p *Trees) Close() {
	// no-op
}
