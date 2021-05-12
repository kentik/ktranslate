// Tracked Sessions
package trkdsess

import (
	"github.com/kentik/ktranslate/pkg/util/netflow/session"
)

var sessions = []session.Session{}

func New() session.Session {
	s := session.New()
	sessions = append(sessions, s)
	return s
}

func Get() []session.Session {
	return sessions
}
