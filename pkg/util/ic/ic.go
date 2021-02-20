// Contains constants and maps related to interface classification types, i.e.
// network boundary types and connectivity types.
// ic = interface classification
// nb = network boundary
// ct = connectivity type
package ic

import (
	"fmt"
	"strings"
)

// TrafficProfileNumbers holds a list of combinations of traffic profile number pairings
// for a given profile name.
type TrafficProfileNumbers struct {
	OriginationNumbers []uint32
	DestinationNumbers []uint32
}

// populated on init
var trafficProfileNumbersByName map[string][]TrafficProfileNumbers

func NameFromCTInt(ic int) string {
	if v, ok := CONNECTIVITY_TYPE_INT_TO_NAME[ic]; ok {
		return v
	}
	return NONE_NAME
}

func IntFromCTName(ic string) int {
	if v, ok := CONNECTIVITY_TYPE_NAME_TO_INT[ic]; ok {
		return v
	}
	return NONE_INT
}

func NameFromNBInt(ic int) string {
	if v, ok := NETWORK_BOUNDARY_INT_TO_NAME[ic]; ok {
		return v
	}
	return NONE_NAME
}

func IntFromNBName(ic string) int {
	if v, ok := NETWORK_BOUNDARY_NAME_TO_INT[ic]; ok {
		return v
	}
	return NONE_INT
}

// initialize the map of profile names to number combinations
// - each item in the resulting slice represents a way to achieve the match. Each item
//   in OriginationNumbers or DestinationNumbers should be OR'ed together
func initializeTrafficProfileNumbersByName() {
	trafficProfileNumbersByName = make(map[string][]TrafficProfileNumbers)

	// internal
	trafficProfileNumbersByName[NETWORK_SRC_INTERNAL_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_INTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_INTERNAL},
		},
	}

	// through
	trafficProfileNumbersByName[NETWORK_SRC_THROUGH_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
		},
	}

	// from outside, terminated inside
	trafficProfileNumbersByName[NETWORK_SRC_TERMINATED_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_INTERNAL},
		},
	}

	// originated inside, to outside:
	trafficProfileNumbersByName[NETWORK_SRC_ORIGINATED_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_INTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
		},
	}

	// cloud internal:
	trafficProfileNumbersByName[NETWORK_SRC_CLOUD_INTERNAL_NAME] = []TrafficProfileNumbers{
		// AWS -> AWS
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS},
		},
		// Azure -> Azure
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE},
		},
		// GCP -> GCP
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_GCP},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_GCP},
		},
		// IBM -> IBM
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_IBM},
		},
	}

	// multi-cloud:
	trafficProfileNumbersByName[NETWORK_SRC_MULTICLOUD_NAME] = []TrafficProfileNumbers{
		// AWS -> Azure or GCP
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
		},
		// Azure -> AWS or GCP
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
		},
		// GCP -> AWS or Azure
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_GCP},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_IBM},
		},
		// IBM -> AWS or Azure
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP},
		},
		// Azure or GCP -> AWS
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS},
		},
		// AWS or GCP -> Azure
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AZURE},
		},
		// AWS or Azure -> GCP
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_GCP},
		},
		// AWS or Azure -> IBM
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_IBM},
		},
	}

	// from outside to cloud:
	trafficProfileNumbersByName[NETWORK_SRC_FROM_OUTSIDE_TO_CLOUD_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
		},
	}

	// from cloud to outside:
	trafficProfileNumbersByName[NETWORK_SRC_FROM_CLOUD_TO_OUTSIDE_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_EXTERNAL},
		},
	}

	// from cloud to inside:
	trafficProfileNumbersByName[NETWORK_SRC_FROM_CLOUD_TO_INSIDE_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
			DestinationNumbers: []uint32{NETWORK_SRC_INTERNAL},
		},
	}

	// from inside to cloud
	trafficProfileNumbersByName[NETWORK_SRC_FROM_INSIDE_TO_CLOUD_NAME] = []TrafficProfileNumbers{
		TrafficProfileNumbers{
			OriginationNumbers: []uint32{NETWORK_SRC_INTERNAL},
			DestinationNumbers: []uint32{NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_IBM},
		},
	}

	// empty string (default, no match)
	trafficProfileNumbersByName[""] = []TrafficProfileNumbers{}
}

// TrafficProfileNumbersFromName returns the traffic profile numbers matching the input name.
// Because "cloud" represents 3 different clouds, we need multiple values for each.
// - each item in the resulting slice represents a way to achieve the match. Each item
//   in OriginationNumbers or DestinationNumbers should be OR'ed together
func TrafficProfileNumbersFromName(profile string) []TrafficProfileNumbers {
	ret, found := trafficProfileNumbersByName[strings.ToLower(profile)]
	if !found {
		// default case will match nothing
		return trafficProfileNumbersByName[""]
	}
	return ret
}

func TrafficNameFromNumbers(ov uint32, dv uint32) string {
	// helpers
	fromCloud := ov == NETWORK_SRC_CLOUD_AWS || ov == NETWORK_SRC_CLOUD_AZURE || ov == NETWORK_SRC_CLOUD_GCP || ov == NETWORK_SRC_CLOUD_IBM
	fromOutside := ov == NETWORK_SRC_EXTERNAL
	fromInside := ov == NETWORK_SRC_INTERNAL
	toCloud := dv == NETWORK_SRC_CLOUD_AWS || dv == NETWORK_SRC_CLOUD_AZURE || dv == NETWORK_SRC_CLOUD_GCP || dv == NETWORK_SRC_CLOUD_IBM
	toOutside := dv == NETWORK_SRC_EXTERNAL
	toInside := dv == NETWORK_SRC_INTERNAL

	if fromInside && toOutside {
		return NETWORK_SRC_ORIGINATED_NAME
	}

	if fromOutside && toInside {
		return NETWORK_SRC_TERMINATED_NAME
	}

	if fromOutside && toOutside {
		return NETWORK_SRC_THROUGH_NAME
	}

	if fromInside && toInside {
		return NETWORK_SRC_INTERNAL_NAME
	}

	// cloud internal
	if fromCloud && ov == dv {
		return NETWORK_SRC_CLOUD_INTERNAL_NAME
	}

	// multi-cloud
	if fromCloud && toCloud && ov != dv {
		return NETWORK_SRC_MULTICLOUD_NAME
	}

	// from outside to cloud
	if fromOutside && toCloud {
		return NETWORK_SRC_FROM_OUTSIDE_TO_CLOUD_NAME
	}

	// from cloud to outside
	if fromCloud && toOutside {
		return NETWORK_SRC_FROM_CLOUD_TO_OUTSIDE_NAME
	}

	// from cloud to inside
	if fromCloud && toInside {
		return NETWORK_SRC_FROM_CLOUD_TO_INSIDE_NAME
	}

	// from inside to cloud
	if fromInside && toCloud {
		return NETWORK_SRC_FROM_INSIDE_TO_CLOUD_NAME
	}

	return fmt.Sprintf("")
}

const (
	NB_EXTERNAL = 10
	NB_INTERNAL = 20

	NB_EXTERNAL_NAME = "external"
	NB_INTERNAL_NAME = "internal"

	CT_FREE_PNI          = 5
	CT_CUSTOMER          = 15
	CT_HOST              = 25
	CT_BACKBONE          = 35
	CT_PAID_PNI          = 45
	CT_OTHER             = 55
	CT_TRANSIT           = 65
	CT_IX                = 75
	CT_RESERVED          = 85
	CT_AVAILABLE         = 95
	CT_DC_IC             = 105
	CT_AGG_IC            = 115
	CT_EMBEDDED_CACHE_IC = 125
	CT_CLOUD_IX          = 135
	CT_DC_FABRIC         = 145

	CT_FREE_PNI_NAME       = "free_pni"
	CT_CUSTOMER_NAME       = "customer"
	CT_HOST_NAME           = "host"
	CT_BACKBONE_NAME       = "backbone"
	CT_PAID_PNI_NAME       = "paid_pni"
	CT_OTHER_NAME          = "other"
	CT_TRANSIT_NAME        = "transit"
	CT_IX_NAME             = "ix"
	CT_RESERVED_NAME       = "reserved"
	CT_AVAILABLE_NAME      = "available"
	CT_DC_IC_NAME          = "datacenter_interconnect"
	CT_AGG_IC_NAME         = "aggregation_interconnect"
	CT_EMBEDDED_CACHE_NAME = "embedded_cache"
	CT_CLOUD_IX_NAME       = "cloud_interconnect"
	CT_DC_FABRIC_NAME      = "datacenter_fabric"

	NETWORK_CL_INTERNAL_NAME               = "inside"
	NETWORK_CL_EXTERNAL_NAME               = "outside"
	NETWORK_CL_CLOUD_AWS                   = "aws"
	NETWORK_CL_CLOUD_IBM                   = "ibm"
	NETWORK_CL_CLOUD_AZURE                 = "azure"
	NETWORK_CL_CLOUD_GCP                   = "gcp"
	NETWORK_SRC_INTERNAL                   = uint32(10)
	NETWORK_SRC_EXTERNAL                   = uint32(20)
	NETWORK_SRC_CLOUD_AWS                  = uint32(80)
	NETWORK_SRC_CLOUD_AZURE                = uint32(90)
	NETWORK_SRC_CLOUD_GCP                  = uint32(100)
	NETWORK_SRC_CLOUD_IBM                  = uint32(110)
	NETWORK_SRC_INTERNAL_NAME              = "internal"
	NETWORK_SRC_THROUGH_NAME               = "through"
	NETWORK_SRC_TERMINATED_NAME            = "from outside, terminated inside"
	NETWORK_SRC_ORIGINATED_NAME            = "originated inside, to outside"
	NETWORK_SRC_CLOUD_INTERNAL_NAME        = "cloud internal"
	NETWORK_SRC_FROM_OUTSIDE_TO_CLOUD_NAME = "from outside to cloud"
	NETWORK_SRC_FROM_CLOUD_TO_OUTSIDE_NAME = "from cloud to outside"
	NETWORK_SRC_FROM_CLOUD_TO_INSIDE_NAME  = "from cloud to inside"
	NETWORK_SRC_FROM_INSIDE_TO_CLOUD_NAME  = "from inside to cloud"
	NETWORK_SRC_MULTICLOUD_NAME            = "multi-cloud"
	NETWORK_SRC_INBOUND_NAME               = "inbound"
	NETWORK_SRC_OUTBOUND_NAME              = "outbound"
	NETWORK_SRC_OTHER_NAME                 = "other"
	NETWORK_DIR_IN_NAME                    = "in"
	NETWORK_DIR_OUT_NAME                   = "out"
	NETWORK_DIR_NOT_HOST_NAME              = "not_a_host"
	NETWORK_DIR_ERR_NAME                   = "error"
	NETWORK_DIR_IN                         = uint32(40)
	NETWORK_DIR_OUT                        = uint32(50)
	NETWORK_DIR_NOT_HOST                   = uint32(60)
	NETWORK_DIR_ERR                        = uint32(70)

	RPKI_NOT_FOUND               = uint32(1)
	RPKI_EXPLICIT_INVALID        = uint32(2)
	RPKI_INVALID_PREFIX          = uint32(3)
	RPKI_INVALID                 = uint32(4)
	RPKI_VALID                   = uint32(5)
	RPKI_INVALID_COVERING        = uint32(6)
	RPKI_COVERING_NOT_FOUND      = uint32(7)
	RPKI_ERROR                   = uint32(8)
	RPKI_INTERNAL                = uint32(9)
	RPKI_MAX_NUM                 = uint32(32)
	RPKI_NOT_FOUND_NAME          = "RPKI Unknown"
	RPKI_EXPLICIT_INVALID_NAME   = "RPKI Invalid: explicit ASN 0"
	RPKI_INVALID_PREFIX_NAME     = "RPKI Invalid: prefix length out of bounds"
	RPKI_INVALID_NAME            = "RPKI Invalid: incorrect Origin ASN (should be AS%d)"
	RPKI_INVALID_COVERING_NAME   = "RPKI Invalid: valid covering prefix"
	RPKI_COVERING_NOT_FOUND_NAME = "RPKI Invalid: unknown covering prefix"
	RPKI_VALID_NAME              = "RPKI Valid"
	RPKI_ERROR_NAME              = "RPKI Error"
	RPKI_INTERNAL_NAME           = "" // Leave ths off, just note that this is an internal route.

	// For min detail results
	RPKI_UNKNOWN_MIN_NAME         = "RPKI Unknown"
	RPKI_VALID_MIN_NAME           = "RPKI Valid"
	RPKI_INVALID_COVERED_MIN_NAME = "RPKI Invalid - covering Valid/Unknown"
	RPKI_INVALID_MIN_NAME         = "RPKI Invalid - Will be dropped"

	NONE_INT = 0

	NONE_NAME = "none"
)

var (
	NETWORK_BOUNDARY_INT_TO_NAME = map[int]string{
		NB_EXTERNAL: NB_EXTERNAL_NAME,
		NB_INTERNAL: NB_INTERNAL_NAME,
	}

	NETWORK_CLASS_INT_TO_NAME = map[uint32]string{
		NETWORK_SRC_EXTERNAL:    NETWORK_CL_EXTERNAL_NAME,
		NETWORK_SRC_INTERNAL:    NETWORK_CL_INTERNAL_NAME,
		NETWORK_SRC_CLOUD_AWS:   NETWORK_CL_CLOUD_AWS,
		NETWORK_SRC_CLOUD_AZURE: NETWORK_CL_CLOUD_AZURE,
		NETWORK_SRC_CLOUD_GCP:   NETWORK_CL_CLOUD_GCP,
		NETWORK_SRC_CLOUD_IBM:   NETWORK_CL_CLOUD_IBM,
		NETWORK_DIR_IN:          NETWORK_DIR_IN_NAME,
		NETWORK_DIR_OUT:         NETWORK_DIR_OUT_NAME,
		NETWORK_DIR_NOT_HOST:    NETWORK_DIR_NOT_HOST_NAME,
		NETWORK_DIR_ERR:         NETWORK_DIR_ERR_NAME,
	}

	CONNECTIVITY_TYPE_INT_TO_NAME = map[int]string{
		CT_FREE_PNI:          CT_FREE_PNI_NAME,
		CT_CUSTOMER:          CT_CUSTOMER_NAME,
		CT_HOST:              CT_HOST_NAME,
		CT_BACKBONE:          CT_BACKBONE_NAME,
		CT_PAID_PNI:          CT_PAID_PNI_NAME,
		CT_OTHER:             CT_OTHER_NAME,
		CT_TRANSIT:           CT_TRANSIT_NAME,
		CT_IX:                CT_IX_NAME,
		CT_RESERVED:          CT_RESERVED_NAME,
		CT_AVAILABLE:         CT_AVAILABLE_NAME,
		CT_DC_IC:             CT_DC_IC_NAME,
		CT_AGG_IC:            CT_AGG_IC_NAME,
		CT_EMBEDDED_CACHE_IC: CT_EMBEDDED_CACHE_NAME,
		CT_CLOUD_IX:          CT_CLOUD_IX_NAME,
		CT_DC_FABRIC:         CT_DC_FABRIC_NAME,
	}

	RPKI_INT_TO_NAME = map[uint32]string{
		RPKI_NOT_FOUND:          RPKI_NOT_FOUND_NAME,
		RPKI_EXPLICIT_INVALID:   RPKI_EXPLICIT_INVALID_NAME,
		RPKI_INVALID_PREFIX:     RPKI_INVALID_PREFIX_NAME,
		RPKI_INVALID:            RPKI_INVALID_NAME,
		RPKI_VALID:              RPKI_VALID_NAME,
		RPKI_INVALID_COVERING:   RPKI_INVALID_COVERING_NAME,
		RPKI_COVERING_NOT_FOUND: RPKI_COVERING_NOT_FOUND_NAME,
		RPKI_ERROR:              RPKI_ERROR_NAME,
		RPKI_INTERNAL:           RPKI_INTERNAL_NAME,
	}

	// For min detail results
	RPKI_INT_TO_MIN_NAME = map[uint32]string{
		RPKI_NOT_FOUND:          RPKI_UNKNOWN_MIN_NAME,
		RPKI_EXPLICIT_INVALID:   RPKI_INVALID_MIN_NAME,
		RPKI_INVALID_PREFIX:     RPKI_INVALID_MIN_NAME,
		RPKI_INVALID:            RPKI_INVALID_MIN_NAME,
		RPKI_VALID:              RPKI_VALID_MIN_NAME,
		RPKI_INVALID_COVERING:   RPKI_INVALID_COVERED_MIN_NAME,
		RPKI_COVERING_NOT_FOUND: RPKI_INVALID_COVERED_MIN_NAME,
		RPKI_ERROR:              RPKI_ERROR_NAME,
		RPKI_INTERNAL:           RPKI_INTERNAL_NAME,
	}

	NETWORK_BOUNDARY_NAME_TO_INT = func() map[string]int {
		m := make(map[string]int)
		for k, v := range NETWORK_BOUNDARY_INT_TO_NAME {
			m[v] = k
		}
		return m
	}()

	CONNECTIVITY_TYPE_NAME_TO_INT = func() map[string]int {
		m := make(map[string]int)
		for k, v := range CONNECTIVITY_TYPE_INT_TO_NAME {
			m[v] = k
		}
		return m
	}()

	NETWORK_CLASS_NAME_TO_INT = func() map[string]uint32 {
		m := make(map[string]uint32)
		for k, v := range NETWORK_CLASS_INT_TO_NAME {
			m[v] = k
		}
		return m
	}()

	RPKI_NAME_TO_INT = func() map[string]uint32 {
		m := make(map[string]uint32)
		for k, v := range RPKI_INT_TO_NAME {
			m[v] = k
		}
		return m
	}()

	RPKI_MIN_NAME_TO_INT = func() map[string][]uint32 {
		m := make(map[string][]uint32)
		for k, v := range RPKI_INT_TO_MIN_NAME {
			if _, ok := m[v]; !ok {
				m[v] = []uint32{k}
			} else {
				m[v] = append(m[v], k)
			}
		}
		return m
	}()
)

func OverrideConnectivityType(newBase string) int {
	newC := map[string]string{}
	found := 0

	for _, pt := range strings.Split(newBase, ",") {
		pts := strings.Split(pt, ":")
		if len(pts) == 2 {
			newC[pts[0]] = pts[1]
		}
	}

	for old, new := range newC {
		if v, ok := CONNECTIVITY_TYPE_NAME_TO_INT[old]; ok {
			CONNECTIVITY_TYPE_INT_TO_NAME[v] = new
			CONNECTIVITY_TYPE_NAME_TO_INT[new] = v
			found++
		}
	}

	return found
}
