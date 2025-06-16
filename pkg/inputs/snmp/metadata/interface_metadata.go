package metadata

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap"

	"github.com/gosnmp/gosnmp"

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
	SNMP_ifPhysAddress      = "PhysAddress"
	SNMP_ifLastChange       = "LastChange"

	SNMP_HUAWEI_ID_MAP_OID     = "1.3.6.1.4.1.2011.5.25.110.1.2.1.2"
	MIN_WORKING_HUAWEI_VERSION = 5
)

var (
	SNMP_Interface_OIDS = func() *orderedmap.OrderedMap {
		// These values are stored (and retrieved) in the order given.  Don't
		// change the order unless you know what you're doing.  Add new ones on
		// the end.
		m := orderedmap.NewOrderedMap()
		m.Set("1.3.6.1.2.1.31.1.1.1.15", SNMP_ifSpeed)       // Index ifIndex
		m.Set("1.3.6.1.2.1.4.20.1.2", SNMP_ipAdEntIfIndex)   // Index ipAddr
		m.Set("1.3.6.1.2.1.4.20.1.3", SNMP_ipAdEntNetMask)   // Index ipAddr
		m.Set("1.3.6.1.2.1.55.1.8.1.2", SNMP_ipv6AddrPrefix) // Index ipv6IfIndex, ipv6Addr
		m.Set("1.3.6.1.2.1.2.2.1.6", SNMP_ifPhysAddress)     // (ifPhysAddress)
		return m
	}()

	getHuaweiVersionRE = regexp.MustCompile(`Version (\d)\.(\d+)`)
)

type InterfaceMetadata struct {
	log           logger.ContextL
	mibs          map[string]*kt.Mib
	interfaceOids *orderedmap.OrderedMap
	debug         bool
}

func NewInterfaceMetadata(interfaceMetadataMibs map[string]*kt.Mib, log logger.ContextL, debug bool) *InterfaceMetadata {
	mibs := map[string]*kt.Mib{}
	m := orderedmap.NewOrderedMap()

	// Copy the global values into this map which is per device.
	//for el := SNMP_Interface_OIDS.Front(); el != nil; el = el.Next() {
	//	m.Set(el.Key, el.Value)
	//}

	if len(interfaceMetadataMibs) > 0 {
		oids := getFromCustomMap(interfaceMetadataMibs)
		for _, oid := range oids {
			mib := interfaceMetadataMibs[oid]
			name := mib.GetName()
			if strings.HasPrefix(name, "if_") {
				name = name[3:]
			}
			if strings.HasPrefix(name, "if") {
				name = name[2:]
			}
			mibs[name] = mib
			_, ok := m.Get(oid)
			if !ok {
				log.Infof("Adding custom interface metadata oid: %s -> %s %v", oid, name, mib.OtherTables)
				m.Set(oid, name)
			}
		}
	}

	return &InterfaceMetadata{
		log:           log,
		debug:         debug,
		mibs:          mibs,
		interfaceOids: m,
	}
}

func (im *InterfaceMetadata) Poll(ctx context.Context, conf *kt.SnmpDeviceConfig, server *gosnmp.GoSNMP) (map[string]*kt.InterfaceData, string, error) {

	// map from ifIndex (as a string) to interface definitions
	intLine := make(map[string]*kt.InterfaceData)

	// map from IP address to ifIndex
	ifcIdx := make(map[string]string)

	// map from ifDescr to interface definition (same structs as in intLine)
	interfacesByDescription := make(map[string]*kt.InterfaceData)

	for el := im.interfaceOids.Front(); el != nil; el = el.Next() {
		oidVal := el.Key.(string)
		oidName := el.Value.(string)
		mib := im.mibs[oidName]

		results, err := snmp_util.WalkOID(ctx, conf, oidVal, server, im.log, "Interface")
		if err != nil {
			return nil, "", err
		}

		for _, variable := range results {
			parts := strings.Split(variable.Name, oidVal)
			if im.debug {
				im.log.Debugf("SNMP %v:%+v", variable.Name, variable.Value)
			}
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
			} else {
				// intLine is indexed with ifIndex.
				// Ignore OIDs with secondary indices, ie. SNMP_ipv6AddrPrefix.
				data, ok := intLine[oidIdx]
				if !ok {
					data = &kt.InterfaceData{IPAddr: kt.IPAddr{}, Index: oidIdx, ExtraInfo: map[string]string{}}
					switch oidName {
					case SNMP_ipv6AddrPrefix:
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
					val := gosnmp.ToBigInt(variable.Value).Int64()
					if mib != nil && mib.EnumRev != nil {
						if ev, ok := mib.EnumRev[val]; ok {
							data.Type = ev
						} else {
							data.Type = fmt.Sprintf("unknown: %d", val)
						}
					} else {
						data.Type = fmt.Sprintf("%d", val)
					}
				case SNMP_ifPhysAddress:
					switch variable.Type {
					case gosnmp.OctetString:
						_, sval, _ := snmp_util.GetFromConv(variable, "hwaddr", im.log)
						data.ExtraInfo[SNMP_ifPhysAddress] = sval
					}
				case SNMP_ifLastChange:
					val := gosnmp.ToBigInt(variable.Value).Uint64()
					data.ExtraInfo[SNMP_ifLastChange] = strconv.Itoa(int(val))
				default:
					switch variable.Type {
					case gosnmp.OctetString:
						data.ExtraInfo[oidName] = string(variable.Value.([]byte))
					case gosnmp.ObjectIdentifier:
						data.ExtraInfo[oidName] = string(variable.Value.(string))
					case gosnmp.Integer:
						val := gosnmp.ToBigInt(variable.Value).Uint64()
						if mib != nil && mib.EnumRev != nil {
							if ev, ok := mib.EnumRev[int64(val)]; ok {
								data.ExtraInfo[oidName] = ev
							} else {
								data.ExtraInfo[oidName] = strconv.Itoa(int(val))
							}
						} else {
							data.ExtraInfo[oidName] = strconv.Itoa(int(val))
						}
					}
				}
			}
		}
	}

	deviceManufacturer := snmp_util.GetDeviceManufacturer(server, im.log)
	lowerManufacturer := strings.ToLower(deviceManufacturer)

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
		if im.debug {
			im.log.Debugf("[API] ", "SNMP %+v:%+v", variable.Name, variable.Value)
		}
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
