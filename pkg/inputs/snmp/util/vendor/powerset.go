package vendor

import (
	"strings"
)

func HandlePowersetStatus(bv []byte) (int64, string) {
	actives := make([]*UpsCode, len(upsCodes))

	for i, v := range bv {
		switch string(v) {
		case "0":
		case "1":
			actives[i] = &upsCodes[i]
		}
	}

	// Track what is set in the system.
	severity := map[int][]string{}

	for _, active := range actives {
		if active == nil {
			continue
		}
		if active.Dependancy > 0 {
			if actives[active.Dependancy] != nil {
				// Ignore because this dependancy is set.
			} else {
				if severity[active.Severity] == nil {
					severity[active.Severity] = []string{active.Status}
				} else {
					severity[active.Severity] = append(severity[active.Severity], active.Status)
				}
			}
		} else { // No dependency so set right away.
			if severity[active.Severity] == nil {
				severity[active.Severity] = []string{active.Status}
			} else {
				severity[active.Severity] = append(severity[active.Severity], active.Status)
			}
		}
	}

	for i := len(severities); i > 0; i-- {
		if len(severity[i]) > 0 { // Return the highest severity present which has flags set.
			return int64(i), severities[i-1] + ": " + strings.Join(severity[i], ",")
		}
	}

	// Generally assume you don't get here.
	return 0, ""
}

type UpsCode struct {
	Severity   int
	Index      int
	Dependancy int
	Status     string
}

var (
	severities = []string{
		"Informational",
		"Warning",
		"High",
		"Disaster",
	}

	upsCodes = []UpsCode{
		UpsCode{Severity: 2, Index: 0, Dependancy: 0, Status: "Abnormal Condition Present"},
		UpsCode{Severity: 3, Index: 1, Dependancy: 28, Status: "On Battery"},
		UpsCode{Severity: 3, Index: 2, Dependancy: 0, Status: "Low Battery"},
		UpsCode{Severity: 1, Index: 3, Dependancy: 0, Status: "On Line"},
		UpsCode{Severity: 3, Index: 4, Dependancy: 0, Status: "Replace Battery"},
		UpsCode{Severity: 1, Index: 5, Dependancy: 0, Status: "Serial Communication Established"},
		UpsCode{Severity: 1, Index: 6, Dependancy: 0, Status: "AVR Boost Active"},
		UpsCode{Severity: 1, Index: 7, Dependancy: 0, Status: "AVR Trim Active"},
		UpsCode{Severity: 3, Index: 8, Dependancy: 28, Status: "Overload"},
		UpsCode{Severity: 1, Index: 9, Dependancy: 0, Status: "Runtime Calibration "},
		UpsCode{Severity: 3, Index: 10, Dependancy: 28, Status: "Batteries Discharged"},
		UpsCode{Severity: 2, Index: 11, Dependancy: 0, Status: "Manual Bypass"},
		UpsCode{Severity: 2, Index: 12, Dependancy: 0, Status: "Software Bypass"},
		UpsCode{Severity: 3, Index: 13, Dependancy: 0, Status: "In Bypass due to Internal Fault"},
		UpsCode{Severity: 3, Index: 14, Dependancy: 0, Status: "In Bypass due to Supply Failure"},
		UpsCode{Severity: 3, Index: 15, Dependancy: 0, Status: "In Bypass due to Fan Failure"},
		UpsCode{Severity: 1, Index: 16, Dependancy: 0, Status: "Sleeping on a Timer"},
		UpsCode{Severity: 2, Index: 17, Dependancy: 0, Status: "Sleeping until Utility Power Returns"},
		UpsCode{Severity: 1, Index: 18, Dependancy: 0, Status: "Powered On"},
		UpsCode{Severity: 2, Index: 19, Dependancy: 0, Status: "Rebooting"},
		UpsCode{Severity: 2, Index: 20, Dependancy: 0, Status: "Battery Communication Lost"},
		UpsCode{Severity: 1, Index: 21, Dependancy: 0, Status: "Graceful Shutdown Initiated"},
		UpsCode{Severity: 1, Index: 22, Dependancy: 0, Status: "Smart Boost or Smart Trim Fault"},
		UpsCode{Severity: 3, Index: 23, Dependancy: 28, Status: "Bad Output Voltage"},
		UpsCode{Severity: 3, Index: 24, Dependancy: 0, Status: "Battery Charger Failure"},
		UpsCode{Severity: 2, Index: 25, Dependancy: 28, Status: "High Battery Temperature"},
		UpsCode{Severity: 3, Index: 26, Dependancy: 0, Status: "Warning Battery Temperature"},
		UpsCode{Severity: 4, Index: 27, Dependancy: 0, Status: "Critical Battery Temperature"},
		UpsCode{Severity: 2, Index: 28, Dependancy: 0, Status: "Self Test In Progress"},
		UpsCode{Severity: 3, Index: 29, Dependancy: 28, Status: "Low Battery / On Battery"},
		UpsCode{Severity: 1, Index: 30, Dependancy: 0, Status: "Graceful Shutdown Issued by Upstream Device"},
		UpsCode{Severity: 1, Index: 31, Dependancy: 0, Status: "Graceful Shutdown Issued by Downstream Device"},
		UpsCode{Severity: 4, Index: 32, Dependancy: 0, Status: "No Batteries Attached"},
		UpsCode{Severity: 1, Index: 33, Dependancy: 0, Status: "Synchronized Command is in Progress"},
		UpsCode{Severity: 1, Index: 34, Dependancy: 0, Status: "Synchronized Sleeping Command is in Progress"},
		UpsCode{Severity: 1, Index: 35, Dependancy: 0, Status: "Synchronized Rebooting Command is in Progress "},
		UpsCode{Severity: 2, Index: 36, Dependancy: 0, Status: "Inverter DC Imbalance "},
		UpsCode{Severity: 3, Index: 37, Dependancy: 0, Status: "Transfer Relay Failure"},
		UpsCode{Severity: 3, Index: 38, Dependancy: 0, Status: "Shutdown or Unable to Transfer"},
		UpsCode{Severity: 3, Index: 39, Dependancy: 0, Status: "Low Battery Shutdown"},
		UpsCode{Severity: 3, Index: 40, Dependancy: 0, Status: "Electronic Unit Fan Failure"},
		UpsCode{Severity: 3, Index: 41, Dependancy: 0, Status: "Main Relay Failure"},
		UpsCode{Severity: 3, Index: 42, Dependancy: 0, Status: "Bypass Relay Failure"},
		UpsCode{Severity: 2, Index: 43, Dependancy: 0, Status: "Temporary Bypass"},
		UpsCode{Severity: 3, Index: 44, Dependancy: 0, Status: "High Internal Temperature"},
		UpsCode{Severity: 3, Index: 45, Dependancy: 0, Status: "Battery Temperature Sensor Fault"},
		UpsCode{Severity: 2, Index: 46, Dependancy: 0, Status: "Input Out of Range for Bypass"},
		UpsCode{Severity: 3, Index: 47, Dependancy: 0, Status: "DC Bus Overvoltage"},
		UpsCode{Severity: 2, Index: 48, Dependancy: 0, Status: "PFC Failure"},
		UpsCode{Severity: 3, Index: 49, Dependancy: 0, Status: "Critical Hardware Fault"},
		UpsCode{Severity: 1, Index: 50, Dependancy: 0, Status: "Green Mode/ECO Mode"},
		UpsCode{Severity: 1, Index: 51, Dependancy: 0, Status: "Hot Standby"},
		UpsCode{Severity: 3, Index: 52, Dependancy: 0, Status: "Emergency Power Off (EPO) Activated"},
		UpsCode{Severity: 2, Index: 53, Dependancy: 0, Status: "Load Alarm Violation"},
		UpsCode{Severity: 3, Index: 54, Dependancy: 0, Status: "Bypass Phase Fault"},
		UpsCode{Severity: 2, Index: 55, Dependancy: 0, Status: "UPS Internal Communication Failure"},
		UpsCode{Severity: 1, Index: 56, Dependancy: 0, Status: "Efficiency Booster Mode"},
		UpsCode{Severity: 2, Index: 57, Dependancy: 0, Status: "Off"},
		UpsCode{Severity: 2, Index: 58, Dependancy: 0, Status: "Standby"},
		UpsCode{Severity: 1, Index: 59, Dependancy: 0, Status: "Not Used"},
		UpsCode{Severity: 1, Index: 60, Dependancy: 0, Status: "Not Used"},
		UpsCode{Severity: 1, Index: 61, Dependancy: 0, Status: "Not Used"},
		UpsCode{Severity: 1, Index: 62, Dependancy: 0, Status: "Not Used"},
		UpsCode{Severity: 1, Index: 63, Dependancy: 0, Status: "Not Used"},
	}
)
