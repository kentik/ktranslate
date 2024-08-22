package kaping

import "net/netip"

type Config struct {
	BindAddr4 netip.Addr
	BindAddr6 netip.Addr
	RawSocket bool
}

func DefaultConfig() Config {
	return Config{
		BindAddr4: netip.MustParseAddr("0.0.0.0"),
		BindAddr6: netip.MustParseAddr("::"),
		RawSocket: false,
	}
}
