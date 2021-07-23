package metadata

import (
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap"

	"github.com/kentik/gosnmp"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"
	"github.com/kentik/ktranslate/pkg/kt"
)

const (
	SNMP_ifDescr            = "ifDescr"
	SNMP_ifAlias            = "ifAlias"
	SNMP_ifSpeed            = "ifSpeed"
	SNMP_ipAdEntIfIndex     = "ipAdEntIfIndex"
	SNMP_ipAdEntNetMask     = "ipAdEntNetMask"
	SNMP_ipv6AddrPrefix     = "ipv6AddrPrefix"
	SNMP_ifConnectorPresent = "ConnectorPresent"
	SNMP_ifType             = "Type"
	SNMP_ifMtu              = "Mtu"
	SNMP_ifPhysAddress      = "PhysAddress"
	SNMP_ifLastChange       = "LastChange"
	SNMP_ifStackStatus      = "StackStatus"

	SNMP_lldpRemPortIdSubtype = "lldpRemPortIdSubtype"
	SNMP_lldpRemPortId        = "lldpRemPortId"
	SNMP_lldpRemPortDesc      = "lldpRemPortDesc"
	SNMP_lldpRemSysName       = "lldpRemSysName"

	SNMP_Nokia_vRtrIfGlobalIndex = "vRtrIfGlobalIndex"
	SNMP_Nokia_vRtrIfName        = "vRtrIfName"
	SNMP_Nokia_vRtrIfDescription = "vRtrIfDescription"

	// VRFs McNuts
	SNMP_vrfDescr = "mplsL3VpnVrfDescription"
	SNMP_vrfIfIdx = "mplsL3VpnIfVpnClassification"
	SNMP_vrfRD    = "mplsL3VpnVrfRD"
	SNMP_vrfRT    = "mplsL3VpnVrfRT"

	SNMP_HUAWEI_ID_MAP_OID     = "1.3.6.1.4.1.2011.5.25.110.1.2.1.2"
	MIN_WORKING_HUAWEI_VERSION = 5
)

var (
	SNMP_Interface_OIDS = func() *orderedmap.OrderedMap {
		// These values are stored (and retrieved) in the order given.  Don't
		// change the order unless you know what you're doing.  Add new ones on
		// the end.
		m := orderedmap.NewOrderedMap()
		m.Set("1.3.6.1.2.1.2.2.1.2", SNMP_ifDescr)           // Index ifIndex
		m.Set("1.3.6.1.2.1.31.1.1.1.18", SNMP_ifAlias)       // Index ifIndex
		m.Set("1.3.6.1.2.1.31.1.1.1.15", SNMP_ifSpeed)       // Index ifIndex
		m.Set("1.3.6.1.2.1.4.20.1.2", SNMP_ipAdEntIfIndex)   // Index ipAddr
		m.Set("1.3.6.1.2.1.4.20.1.3", SNMP_ipAdEntNetMask)   // Index ipAddr
		m.Set("1.3.6.1.2.1.55.1.8.1.2", SNMP_ipv6AddrPrefix) // Index ipv6IfIndex, ipv6Addr

		// Fetch mplsL3VpnVrfDescription first
		// for subsequent indexing in other vrf polling.
		// VRFs McNuts
		// mplsL3VpnVrfName bellow may be a system generated ID.
		m.Set("1.3.6.1.2.1.10.166.11.1.2.2.1.3", SNMP_vrfDescr) // mplsL3VpnVrfDescription Index: mplsL3VpnVrfName
		m.Set("1.3.6.1.2.1.10.166.11.1.2.2.1.4", SNMP_vrfRD)    // mplsL3VpnVrfRD RD Index: mplsL3VpnVrfName
		m.Set("1.3.6.1.2.1.10.166.11.1.2.3.1.4", SNMP_vrfRT)    // mplsL3VpnVrfRT RT Index: mplsL3VpnVrfName, mplsL3VpnVrfRTIndex, mplsL3VpnVrfRTType
		// Direct vrf->if mapping or vice-versa isn't
		// available in standard mplsl3vpn MIB.
		// Instead use IfVpnClassification OID to
		// identify related interfaces.
		m.Set("1.3.6.1.2.1.10.166.11.1.2.1.1.2", SNMP_vrfIfIdx)      // mplsL3VpnIfVpnClassification Index: mplsL3VpnVrfName
		m.Set("1.3.6.1.2.1.31.1.1.1.17", SNMP_ifConnectorPresent)    // ifConnectorPresent, 1 if physical, 2 otherwise.
		m.Set("1.3.6.1.2.1.2.2.1.3", SNMP_ifType)                    // (ifType)
		m.Set("1.3.6.1.2.1.2.2.1.4", SNMP_ifMtu)                     // (ifMtu)
		m.Set("1.3.6.1.2.1.2.2.1.6", SNMP_ifPhysAddress)             // (ifPhysAddress)
		m.Set("1.3.6.1.2.1.2.2.1.9", SNMP_ifLastChange)              // (ifLastChange)
		m.Set("1.3.6.1.2.1.31.1.2.1.3", SNMP_ifStackStatus)          // (ifStackStatus, indexed by ifStackHigherLayer, ifStackLowerLayer)
		m.Set("1.0.8802.1.1.2.1.4.1.1.6", SNMP_lldpRemPortIdSubtype) // lldpRemPortIdSubtype
		m.Set("1.0.8802.1.1.2.1.4.1.1.7", SNMP_lldpRemPortId)        // lldpRemPortId
		m.Set("1.0.8802.1.1.2.1.4.1.1.8", SNMP_lldpRemPortDesc)      // lldpRemPortDesc
		m.Set("1.0.8802.1.1.2.1.4.1.1.9", SNMP_lldpRemSysName)       // lldpRemSysName
		return m
	}()

	// aka TiMetra oids
	SNMP_Nokia_oids = func() *orderedmap.OrderedMap {
		m := orderedmap.NewOrderedMap()
		m.Set("1.3.6.1.4.1.6527.3.1.2.3.4.1.63", SNMP_Nokia_vRtrIfGlobalIndex)
		m.Set("1.3.6.1.4.1.6527.3.1.2.3.4.1.4", SNMP_Nokia_vRtrIfName)
		m.Set("1.3.6.1.4.1.6527.3.1.2.3.4.1.34", SNMP_Nokia_vRtrIfDescription)
		return m
	}()

	getHuaweiVersionRE = regexp.MustCompile(`Version (\d)\.(\d+)`)
)

type vrf struct {
	Name  string
	Descr string
	RD    string // Route Distinguisher
	RT    string // BGP Route Target
	ExtRD uint64 // Route Distinguisher encoded as ext community RFC 4360
}

type id2vrf map[string]*vrf

// TODO: global?  Really?
var vrfs = make(id2vrf)

type InterfaceMetadata struct {
	log logger.ContextL
}

func NewInterfaceMetadata(interfaceMetadataMibs map[string]*kt.Mib, log logger.ContextL) *InterfaceMetadata {
	if len(interfaceMetadataMibs) > 0 {
		oids := getFromCustomMap(interfaceMetadataMibs)
		for _, oid := range oids {
			_, ok := SNMP_Interface_OIDS.Get(oid)
			if !ok {
				mib := interfaceMetadataMibs[oid]
				name := mib.Name
				if strings.HasPrefix(name, "if") {
					name = name[2:]
				}
				log.Infof("Adding custom interface metadata oid: %s -> %s", oid, name)
				SNMP_Interface_OIDS.Set(oid, name)
			}
		}
	}

	return &InterfaceMetadata{
		log: log,
	}
}

func (im *InterfaceMetadata) Poll(server *gosnmp.GoSNMP) (map[string]*kt.InterfaceData, string, error) {

	// map from ifIndex (as a string) to interface definitions
	intLine := make(map[string]*kt.InterfaceData)

	// map from IP address to ifIndex
	ifcIdx := make(map[string]string)

	// map from ifDescr to interface definition (same structs as in intLine)
	interfacesByDescription := make(map[string]*kt.InterfaceData)

	for el := SNMP_Interface_OIDS.Front(); el != nil; el = el.Next() {
		oidVal := el.Key.(string)
		oidName := el.Value.(string)

		results, err := snmp_util.WalkOID(oidVal, server, im.log, "Interface")
		if err != nil {
			return nil, "", err
		}

		for _, variable := range results {
			parts := strings.Split(variable.Name, oidVal)
			im.log.Debugf("SNMP %v:%+v", variable.Name, variable.Value)
			if len(parts) != 2 || len(parts[1]) == 0 {
				continue
			}

			// variable.Name looks like this: .<oidVal>.<intVal>, e.g.
			// .1.3.6.1.2.1.31.1.1.1.10.219, where .1.3.6.1.2.1.31.1.1.1.10 is
			// the oid and 219 is the intVal.  So splitting on oidVal gives us
			// .intVal.
			oidIdx := parts[1][1:]

			if ip := net.ParseIP(oidIdx); ip != nil &&
				(oidName == SNMP_ipAdEntIfIndex || oidName == SNMP_ipAdEntNetMask) {
				addr := oidIdx
				switch oidName {
				case SNMP_ipAdEntIfIndex:
					switch variable.Type {
					case gosnmp.Integer:
						idx := fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value).Uint64())
						if data, ok := intLine[idx]; ok {
							if data.Address == "" {
								data.Address = addr
							} else {
								data.AliasAddr = append(data.AliasAddr, kt.IPAddr{Address: addr})
							}
							ifcIdx[addr] = idx
						}
					}
				case SNMP_ipAdEntNetMask:
					switch variable.Type {
					case gosnmp.IPAddress:
						if idx, ok := ifcIdx[addr]; ok {
							if data, ok := intLine[idx]; ok {
								if data.Address == addr {
									data.Netmask = variable.Value.(string)
								} else {
									for i := range data.AliasAddr {
										if data.AliasAddr[i].Address == addr {
											data.AliasAddr[i].Netmask = variable.Value.(string)
											break
										}
									}
								}
							}
						}
					}
				}
			} else if oidName == SNMP_ifStackStatus {
				// index is two integer elements -- first is ifIndex of the "higher" layer, second is ifIndex of the "lower" layer
				index := strings.Split(oidIdx, ".")
				if len(index) != 2 {
					im.log.Warnf("incorrectly-indexed ifStackStatus instance: %s -> %v", variable.Name, variable.Value)
					continue
				}
				val := gosnmp.ToBigInt(variable.Value).Uint64()
				if val != 1 && val != 2 {
					im.log.Warnf("bad ifStackStatus value: %s -> %v", variable.Name, variable.Value)
					continue
				}

			} else if oidName == SNMP_lldpRemPortIdSubtype || oidName == SNMP_lldpRemPortId || oidName == SNMP_lldpRemPortDesc || oidName == SNMP_lldpRemSysName {
				// these may contain details about what an interface is connected to:
				//     PortId ~= ifIndex on the remote side, PortDesc ~= ifDescr, and SysName ~= sysName
				//
				// index is three integer elements -- first is a TimeFilter (rfc2021), second is the lldpPortNumber (which is probaby an ifIndex),
				// third is just an incrementing sequence of changes to the lldpRemoteTable.
				// So we'll split the index, walk it from oldest TimeFilter to newestTimeFilter, setting
				// the InterfaceData for any lldpPortNumbers that correspond to ifIndex values to this varbind's value.
				// At the end, every ifIndex that appeared in the list at all should have the value from its most
				// recent TimeFilter value.
				index := strings.Split(oidIdx, ".")
				if len(index) != 3 {
					im.log.Warnf("incorrectly-indexed lldpRemTable instance: %s -> %v", variable.Name, variable.Value)
					continue
				}

				data, ok := intLine[index[1]]
				if ok {
					switch oidName {
					case SNMP_lldpRemPortIdSubtype:
						if variable.Type != gosnmp.Integer {
							im.log.Warnf("unexpected type on lldpRemPortIdSubtype %v %v %v", variable.Name, variable.Type, variable.Value)
							break
						}
						data.ExtraInfo[oidName] = strconv.Itoa(int(gosnmp.ToBigInt(variable.Value).Uint64()))

					case SNMP_lldpRemPortId:
						if variable.Type != gosnmp.OctetString {
							im.log.Warnf("unexpected type on lldpRemPortId %v %v %v", variable.Name, variable.Type, variable.Value)
							break
						}
						value := variable.Value.([]byte)

						// the contents of this variable depend on the
						// subtype value we should have gotten on the previous walk;
						// for values, see the lldpPortIdSubtype definition in
						// http://www.ieee802.org/1/files/public/MIBs/LLDP-MIB-200505060000Z.txt
						switch data.ExtraInfo[SNMP_lldpRemPortIdSubtype] {
						case "":
							// we got lldpRemPortId for this interface, but no lldpRemPortIdSubtype --
							// log and skip
							im.log.Warnf("lldpRemPortId present, but not lldpRemPortIdSubtype: %v %v %v", variable.Name, variable.Type, variable.Value)
						case "1":
							// interfaceAlias -- save the value as-is
							data.ExtraInfo[oidName] = string(value)
						case "3":
							// macAddress
							data.ExtraInfo[oidName] = "0x" + hex.EncodeToString(value)
						case "5":
							// interfaceName -- save the value as-is
							data.ExtraInfo[oidName] = string(value)
						case "7":
							// local -- at least on some Juniper devices, this seems to be
							// the ifIndex of the interface we're connected to.  Just save it.
							data.ExtraInfo[oidName] = string(value)
						default:
							im.log.Warnf("unexpected subtype for lldpRemPortId: %v %v %v %v", variable.Name, variable.Type, variable.Value, data.ExtraInfo[SNMP_lldpRemPortIdSubtype])
						}

					case SNMP_lldpRemPortDesc, SNMP_lldpRemSysName:
						if variable.Type != gosnmp.OctetString {
							im.log.Warnf("unexpected type on %s %v %v %v", oidName, variable.Name, variable.Type, variable.Value)
							break
						}
						data.ExtraInfo[oidName] = string(variable.Value.([]byte))
					}
				}
			} else {
				// intLine is indexed with ifIndex.
				// Ignore OIDs with secondary indices, ie. SNMP_ipv6AddrPrefix.
				data, ok := intLine[oidIdx]
				if !ok {
					data = &kt.InterfaceData{IPAddr: kt.IPAddr{}, Index: oidIdx, ExtraInfo: map[string]string{}}
					switch oidName {
					case SNMP_ipv6AddrPrefix:
					case SNMP_vrfDescr:
					case SNMP_vrfIfIdx:
					case SNMP_vrfRD:
					case SNMP_vrfRT:
					default:
						intLine[oidIdx] = data
					}
				}

				// empty strings are null terminated in OctetString below
				switch oidName {
				case SNMP_ifDescr:
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.TRUNCATE); ok {
						data.Description = value
						interfacesByDescription[value] = data
					}
				case SNMP_ifAlias:
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.TRUNCATE); ok {
						data.Alias = value
					}
				case SNMP_ifSpeed:
					value := gosnmp.ToBigInt(variable.Value).Uint64()
					data.Speed = value
				case SNMP_ipv6AddrPrefix:
					switch variable.Type {
					case gosnmp.Integer:
						s := strings.Split(oidIdx, ".")
						if len(s) != 17 {
							break
						}
						oidIdx = s[0]
						s = s[1:]
						s2 := []string{}
						for i := range s {
							if d, err := strconv.Atoi(s[i]); err == nil {
								s[i] = fmt.Sprintf("%02x", d)
							}
							if i%2 > 0 {
								s2 = append(s2, s[i-1]+s[i])
							}
						}
						cidr := strings.Join(s2, ":") + fmt.Sprintf("/%d", gosnmp.ToBigInt(variable.Value).Uint64())
						if addr6, net6, err := net.ParseCIDR(cidr); err == nil {
							if d, ok := intLine[oidIdx]; !ok {
								// Skip any ipv6Prefixes with unknown ifIndex.
								break
							} else {
								data = d
							}
							mask6 := net.IP(net6.Mask)
							if data.Address == "" {
								data.Address = addr6.String()
								data.Netmask = mask6.String()
							} else {
								data.AliasAddr = append(data.AliasAddr, kt.IPAddr{Address: addr6.String(), Netmask: mask6.String()})
							}
						}
					}
				case SNMP_vrfDescr:
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.NO_TRUNCATE); ok {
						vrfs[oidIdx] = &vrf{Name: value}
					}
				case SNMP_vrfRD:
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.NO_TRUNCATE); ok {
						if v, ok := vrfs[oidIdx]; ok {
							v.RD = value
							v.ExtRD = GetVRFExtRD(value)
							if v.ExtRD == 0 {
								im.log.Warnf("SNMP vrf: Invalid RD:%s", v.RD)
							}
						}
					}
				case SNMP_vrfRT:
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.NO_TRUNCATE); ok {
						// Strip off mplsL3VpnVrfRTIndex, mplsL3VpnVrfRTType
						s := strings.Split(oidIdx, ".")
						if len(s) > 3 {
							vrfId := strings.Join(s[:len(s)-2], ".")
							if v, ok := vrfs[vrfId]; ok {
								v.RT = value
							}
						}
					}
				case SNMP_vrfIfIdx:
					switch variable.Type {
					case gosnmp.Integer:
						s := strings.Split(oidIdx, ".")
						if len(s) > 2 {
							ifId := s[len(s)-1]
							vrfId := strings.Join(s[:len(s)-1], ".")
							im.log.Debugf("SNMP vrf:%s:%s", vrfId, ifId)
							if vrf, ok := vrfs[vrfId]; ok {
								if ifc, ok := intLine[ifId]; ok {
									im.log.Debugf("SNMP vrf->if vrf:%+v if:%+v", vrf, ifc)
									ifc.VrfName = vrf.Name
									ifc.VrfRD = vrf.RD
									ifc.VrfExtRD = vrf.ExtRD
									ifc.VrfRT = vrf.RT
								}
							}
						}
					}
				case SNMP_ifConnectorPresent:
					switch gosnmp.ToBigInt(variable.Value).Uint64() {
					case 1:
						data.ExtraInfo[SNMP_ifConnectorPresent] = kt.PHYSICAL_INTERFACE
					case 2:
						data.ExtraInfo[SNMP_ifConnectorPresent] = kt.LOGICAL_INTERFACE
					default:
						im.log.Warnf("SNMP_ifConnectorPresent: Unexpected value %d", gosnmp.ToBigInt(variable.Value).Uint64())
					}
				case SNMP_ifType:
					data.Type = gosnmp.ToBigInt(variable.Value).Uint64()
				case SNMP_ifMtu:
					val := gosnmp.ToBigInt(variable.Value).Uint64()
					data.ExtraInfo[SNMP_ifMtu] = strconv.Itoa(int(val))
				case SNMP_ifPhysAddress:
					switch variable.Type {
					case gosnmp.OctetString:
						b := variable.Value.([]byte)
						if len(b) > 0 {
							value := "0x" + hex.EncodeToString(b)
							data.ExtraInfo[SNMP_ifPhysAddress] = value
						}
					}
				case SNMP_ifLastChange:
					val := gosnmp.ToBigInt(variable.Value).Uint64()
					data.ExtraInfo[SNMP_ifLastChange] = strconv.Itoa(int(val))
				default:
					switch variable.Type {
					case gosnmp.OctetString:
						data.ExtraInfo[oidName] = string(variable.Value.([]byte))
					case gosnmp.Integer:
						val := gosnmp.ToBigInt(variable.Value).Uint64()
						data.ExtraInfo[oidName] = strconv.Itoa(int(val))
					}
				}
			}
		}
	}

	deviceManufacturer := snmp_util.GetDeviceManufacturer(server, im.log)
	lowerManufacturer := strings.ToLower(deviceManufacturer)
	isNokia := snmp_util.ContainsAny(lowerManufacturer, "timos", "nokia")

	if isNokia || lowerManufacturer == "" {
		if err := im.pollNokiaInterfaceOids(server, intLine, interfacesByDescription); err != nil {
			return nil, "", err
		}
	}

	im.copyInterfaces(intLine, interfacesByDescription, lowerManufacturer)

	return intLine, deviceManufacturer, nil
}

func (im *InterfaceMetadata) copyInterfaces(intLine, interfacesByDescription map[string]*kt.InterfaceData, lowerManufacturer string) {

	isJuniper := snmp_util.ContainsAny(lowerManufacturer, "juniper", "junos")
	isQFX := isJuniper && strings.Contains(lowerManufacturer, "qfx")

	// for devices without an alias or index, see if there's a related interface we can copy one from
	for _, deviceInterface := range intLine {

		if deviceInterface.Alias == "" {
			// See if we're a logical interface, and our parent physical (or LAG) address has an ifAlias
			parentDescr := parentInterfaceDescriptionFromDescription(deviceInterface.Description)
			if parentDescr != "" {
				if parent, ok := interfacesByDescription[parentDescr]; ok {
					if parent.Alias != "" {
						deviceInterface.Alias = fmt.Sprintf("%s (inherited from physical)", parent.Alias)
					}
				}
			}
		}

		getUniqueLogical := func(descr string) (*kt.InterfaceData, bool) {
			// Some customers, including Akamai, seem to describe their logical interfaces, but send sflow
			// that specifies only physical addresses.  To label that flow correctly, need to copy the
			// descriptions and addresses in the opposite direction from above, from logical addresses
			// to physical ones.  But we don't want to do that if there are multiple logical
			// interfaces on a single physical one.
			child0 := interfacesByDescription[descr+".0"]
			if child0 == nil {
				return nil, false
			} else {
				for candidateDescr, candidate := range interfacesByDescription {
					if strings.HasPrefix(candidateDescr, descr+".") && (candidate != child0) {
						return nil, false
					}
				}
				return child0, true
			}
		}

		if isQFX {
			if donor, ok := getUniqueLogical(deviceInterface.Description); ok {
				if deviceInterface.Alias == "" && donor.Alias != "" {
					deviceInterface.Alias = fmt.Sprintf("%s (inherited from logical)", donor.Alias)
					im.log.Infof("QFX copy: setting interface %s(%s) alias from %s(%s); manufacturer %s",
						deviceInterface.Description, deviceInterface.Index,
						donor.Description, donor.Index,
						lowerManufacturer)
				}
				if deviceInterface.Address == "" && donor.Address != "" {
					deviceInterface.Address = donor.Address
					deviceInterface.Netmask = donor.Netmask
					deviceInterface.AliasAddr = donor.AliasAddr
					im.log.Infof("QFX copy: setting interface %s(%s) address from %s(%s); manufacturer %s",
						deviceInterface.Description, deviceInterface.Index,
						donor.Description, donor.Index,
						lowerManufacturer)
				}
			}
		}
	}

}

// determine a parent interface alias from the input alias - intended for Juniper devices that
// are named '<parent alias>.alias', returning "" if it doesn't seem to have a parent
func parentInterfaceDescriptionFromDescription(description string) string {
	parts := strings.SplitN(description, ".", 2)
	if len(parts) > 1 {
		return parts[0]
	}
	// no parent
	return ""
}

func (im *InterfaceMetadata) pollNokiaInterfaceOids(server *gosnmp.GoSNMP,
	intLine, interfacesByDescription map[string]*kt.InterfaceData) error {

	im.log.Debugf("Querying Nokia OIDs")

	// key is localIndex
	toGlobalIndex := map[string]string{}

	for el := SNMP_Nokia_oids.Front(); el != nil; el = el.Next() {
		oidVal := el.Key.(string)
		oidName := el.Value.(string)

		// I'm not sure "oid not found" counts as an error or not.  If it
		// does, this needs to change, since walkOID will try the OID three
		// different ways, and if it just doesn't exist on the device, it'll
		// never work.
		results, err := snmp_util.WalkOID(oidVal, server, im.log, "Nokia-Interface")
		if err != nil {
			im.log.Debugf("Error '%v' on %s/%s", err, oidVal, oidName)
			return err
		}

		if len(results) == 0 {
			im.log.Debugf("No results on %s/%s", oidVal, oidName)
			break
		}

		for _, variable := range results {
			parts := strings.Split(variable.Name, oidVal)
			im.log.Debugf("SNMP %v:%+v", variable.Name, variable.Value)
			if len(parts) != 2 || len(parts[1]) == 0 {
				im.log.Debugf("Could not parse %v", variable.Name)
				continue
			}

			// localIndex is <vRtrId>.<vRtrIfIndex>
			localIndex := parts[1][1:]
			im.log.Debugf("%v localIndex is %s", variable.Name, localIndex)

			switch oidName {
			case SNMP_Nokia_vRtrIfGlobalIndex:
				globalIndex := fmt.Sprintf("%v", variable.Value)
				intLine[globalIndex] = &kt.InterfaceData{
					IPAddr:    kt.IPAddr{},
					Index:     globalIndex,
					ExtraInfo: map[string]string{},
				}
				toGlobalIndex[localIndex] = globalIndex
				im.log.Debugf("%v -> globalIndex %v", localIndex, globalIndex)
			case SNMP_Nokia_vRtrIfName: // same as SNMP_ifDescr above
				if data, ok := intLine[toGlobalIndex[localIndex]]; ok {
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.TRUNCATE); ok {
						im.log.Debugf("globalIndex %v ifName: %s", toGlobalIndex[localIndex], value)
						data.Description = value
						interfacesByDescription[value] = data
					}
				}
			case SNMP_Nokia_vRtrIfDescription: // same as SNMP_ifAlias above
				if data, ok := intLine[toGlobalIndex[localIndex]]; ok {
					if value, ok := snmp_util.ReadOctetString(variable, snmp_util.TRUNCATE); ok {
						im.log.Debugf("globalIndex %v ifDescription: %s", toGlobalIndex[localIndex], value)
						data.Alias = value
					}
				}
			}
		}
	}

	return nil
}

func isBrokenHuawei(manufacturer string) bool {
	if strings.Contains(manufacturer, "Huawei") {
		versionStr := getHuaweiVersionRE.FindStringSubmatch(manufacturer)
		if len(versionStr) > 1 {
			version, err := strconv.Atoi(versionStr[1])
			if err == nil && version < MIN_WORKING_HUAWEI_VERSION {
				return true
			}
		}
	}

	return false
}

func (im *InterfaceMetadata) UpdateForHuawei(server *gosnmp.GoSNMP, d *kt.DeviceData) error {
	results, err := server.WalkAll(SNMP_HUAWEI_ID_MAP_OID)
	if err != nil {
		return err
	}

	id2snmp := map[string]string{}
	for _, variable := range results {
		prts := strings.Split(variable.Name, SNMP_HUAWEI_ID_MAP_OID)
		im.log.Debugf("[API] ", "SNMP %+v:%+v", variable.Name, variable.Value)
		if len(prts) == 2 {
			// variable.Name looks like this: .<oid>.<intVal> (.1.3.6.1.2.1.31.1.1.1.10.219)
			oidIdx := prts[1][1:]
			value := gosnmp.ToBigInt(variable.Value)
			if value != nil {
				id2snmp[oidIdx] = strconv.Itoa(int(value.Uint64()))
			}
		}
	}

	idNew := map[string]*kt.InterfaceData{}
	idNewNotFound := map[string]*kt.InterfaceData{}
	idNewFinal := map[string]*kt.InterfaceData{}
	for id, intD := range d.InterfaceData {
		if _, ok := id2snmp[id]; ok {
			intD.Index = id2snmp[id]
			idNew[intD.Index] = intD
		} else {
			idNewNotFound[id] = intD
		}
	}

	// Loop again, adding any id which are not present in the new set
	for id, intD := range idNewNotFound {
		if _, ok := idNew[id]; !ok {
			idNewFinal[id] = intD
		}
	}

	// Lastly, put the two parts back together again
	for id, intD := range idNew {
		idNewFinal[id] = intD
	}
	d.InterfaceData = idNewFinal

	return nil
}

//
//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  | 0x00          |   Format      |    Global Administrator       |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//  |                     Local Administrator                       |
//  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Format 0x00 - 2B Global ASN, 4B Local ASN
// Format 0x01 - 4B IPv addr, 2B Local ASN
// Format 0x02 - 4B Global ASN, 2B Local ASN
//
func GetVRFExtRD(rd string) uint64 {
	s := strings.Split(rd, ":")
	if len(s) != 2 {
		return 0
	}

	extRD := uint64(0)
	if ip := net.ParseIP(s[0]); ip == nil {
		asn, err := strconv.ParseUint(s[0], 10, 32)
		if err != nil {
			return 0
		}

		if asn <= 0xffff {
			extRD = uint64(0x0)<<48 | uint64(asn)<<32
		} else {
			extRD = uint64(0x2)<<48 | uint64(asn)<<16
		}
	} else {
		if ip4 := ip.To4(); ip4 != nil {
			extRD = uint64(0x1)<<48 | uint64(ip4[0])<<40 | uint64(ip4[1])<<32 | uint64(ip4[2])<<24 | uint64(ip4[3])<<16
		} else {
			return 0
		}
	}
	ext, err := strconv.ParseUint(s[1], 10, 32)
	if err != nil {
		return 0
	} else if ext > 0xffff && extRD>>48 != 0x0 {
		return 0
	}
	extRD = extRD | ext
	return extRD
}
