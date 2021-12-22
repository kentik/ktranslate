package gcp

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"
	"github.com/kentik/ktranslate/pkg/util/ic"
)

const (
	RECV_WINDOW  = -1 * 5 * 60 * time.Second
	GCP_VPC_TYPE = "GCP_VPC"
)

type GCELogLine struct {
	InsertID  string    `json:"insertId"`
	Payload   *Payload  `json:"jsonPayload"`
	LogName   string    `json:"logName"`
	RecvTs    time.Time `json:"receiveTimestamp"`
	Resource  *Resource `json:"resource"`
	Timestamp time.Time `json:"timestamp"`
	BytesRaw  int64     `json:"bytesRaw"`
}

type Connection struct {
	DestIP   string `json:"dest_ip"`
	DestPort int    `json:"dest_port"`
	Protocol int    `json:"protocol"`
	SrcIP    string `json:"src_ip"`
	SrcPort  int    `json:"src_port"`
}

type Instance struct {
	ProjectID string `json:"project_id"`
	Region    string `json:"region"`
	VMName    string `json:"vm_name"`
	Zone      string `json:"zone"`
}

type VPC struct {
	ProjectID      string `json:"project_id"`
	SubnetworkName string `json:"subnetwork_name"`
	Name           string `json:"vpc_name"`
}

type Payload struct {
	Bytes        string      `json:"bytes_sent"`
	Connection   *Connection `json:"connection"`
	DestInstance *Instance   `json:"dest_instance"`
	SrcInstance  *Instance   `json:"src_instance"`
	DestVPC      *VPC        `json:"dest_vpc"`
	SrcVPC       *VPC        `json:"src_vpc"`
	EndTime      time.Time   `json:"end_time"`
	Pkts         string      `json:"packets_sent"`
	Reporter     string      `json:"reporter"`
	RTT          string      `json:"rtt_msec"`
	SrcLocation  *Location   `json:"src_location"`
	DstLocation  *Location   `json:"dest_location"`
	StartTime    time.Time   `json:"start_time"`
}

type Location struct {
	City      string `json:"city"`
	Continent string `json:"continent"`
	Country   string `json:"country"`
	Region    string `json:"region"`
}

type Resource struct {
	Labels *Label `json:"labels"`
	Type   string `json:"type"`
}

type Label struct {
	Location       string `json:"location"`
	ProjectID      string `json:"project_id"`
	SubnetworkID   string `json:"subnetwork_id"`
	SubnetworkName string `json:"subnetwork_name"`
}

func (m *GCELogLine) GetVMName() (host string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if j, e := json.Marshal(m); e != nil {
				err = e
			} else {
				err = fmt.Errorf("%v -> %s", r, string(j))
			}
		}
	}()

	if m.IsIn() {
		host = m.Payload.SrcInstance.VMName
	} else {
		host = m.Payload.DestInstance.VMName
	}

	return host, nil
}

func (m *GCELogLine) IsValid() bool {
	if m.Payload != nil {
		return m.Payload.EndTime.After(time.Now().Add(RECV_WINDOW))
	}

	return false
}

func (m *GCELogLine) IsIn() bool {
	return m.Payload.SrcInstance != nil && m.Payload.SrcInstance.VMName != ""
}

func (m *GCELogLine) IsInternal() bool {
	return (m.Payload.SrcInstance != nil && m.Payload.SrcInstance.VMName != "") && (m.Payload.DestInstance != nil && m.Payload.DestInstance.VMName != "")
}

func (m *GCELogLine) ToJson() []byte {
	j, _ := json.Marshal(m)
	return j
}

func (m *GCELogLine) ToFlow(log logger.ContextL) (in *kt.JCHF, err error) {
	defer func() {
		if r := recover(); r != nil {
			if j, e := json.Marshal(m); e != nil {
				err = e
			} else {
				err = fmt.Errorf("%v -> %s", r, string(j))
			}
		}
	}()

	in = kt.NewJCHF()
	in.CustomStr = make(map[string]string)
	in.CustomInt = make(map[string]int32)
	in.CustomBigInt = make(map[string]int64)
	in.EventType = kt.KENTIK_EVENT_TYPE
	in.Provider = kt.ProviderVPC
	in.SampleRate = 1
	vmname, _ := m.GetVMName()
	in.DeviceName = vmname

	in.CustomStr["kt.from"] = kt.FromGCP
	in.CustomStr["type"] = GCP_VPC_TYPE
	in.CustomStr["insert_id"] = m.InsertID
	in.CustomStr["log_name"] = m.LogName
	in.CustomBigInt["rcv_time"] = m.RecvTs.Unix()

	m.Payload.Save(in)
	m.Resource.Save(in)

	return in, err
}

func (l *Label) Save(in *kt.JCHF) {
	if l == nil {
		return
	}

	in.CustomStr["label_location"] = l.Location
	in.CustomStr["label_project_id"] = l.ProjectID
	in.CustomStr["label_subnetwork_id"] = l.SubnetworkID
	in.CustomStr["label_subnetwork_name"] = l.SubnetworkName
}

func (r *Resource) Save(in *kt.JCHF) {
	if r == nil {
		return
	}

	in.CustomStr["type"] = r.Type
	r.Labels.Save(in)
}

func (c *Connection) Save(in *kt.JCHF) {
	if c == nil {
		return
	}

	in.L4DstPort = uint32(c.DestPort)
	in.L4SrcPort = uint32(c.SrcPort)
	in.Protocol = ic.PROTO_NAMES[uint16(uint32(c.Protocol))]
	in.SrcAddr = c.SrcIP
	in.DstAddr = c.DestIP
}

func (i *Instance) Save(in *kt.JCHF, direction string) {
	if i == nil {
		return
	}

	in.CustomStr[direction+"instance_project_id"] = i.ProjectID
	in.CustomStr[direction+"instance_region"] = i.Region
	in.CustomStr[direction+"instance_vm_name"] = i.VMName
	in.CustomStr[direction+"instance_zone"] = i.Zone
}

func (v *VPC) Save(in *kt.JCHF, direction string) {
	if v == nil {
		return
	}

	in.CustomStr[direction+"vpc_project_id"] = v.ProjectID
	in.CustomStr[direction+"vpc_subnetwork_name"] = v.SubnetworkName
	in.CustomStr[direction+"vpc_name"] = v.Name
}

func (l *Location) Save(in *kt.JCHF, direction string) {
	if l == nil {
		return
	}

	in.CustomStr[direction+"city"] = l.City
	in.CustomStr[direction+"continent"] = l.Continent
	in.CustomStr[direction+"country"] = l.Country
	in.CustomStr[direction+"region"] = l.Region
}

func (p *Payload) Save(in *kt.JCHF) {
	if p == nil {
		return
	}

	in.InBytes = getUInt64(p.Bytes)
	in.InPkts = getUInt64(p.Pkts)
	in.CustomBigInt["rtt_msec"] = getInt64(p.RTT)
	in.CustomStr["reporter"] = p.Reporter

	// Do we want all these times?
	in.Timestamp = p.StartTime.Unix()
	in.CustomBigInt["start_time"] = p.StartTime.Unix()
	in.CustomBigInt["end_time"] = p.EndTime.Unix()

	// Now save the rest.
	p.Connection.Save(in)
	p.DestInstance.Save(in, "dst_")
	p.SrcInstance.Save(in, "src_")
	p.DestVPC.Save(in, "dst_")
	p.SrcVPC.Save(in, "src_")
	p.DstLocation.Save(in, "dst_")
	p.SrcLocation.Save(in, "src_")
}

func getUInt32(s string) uint32 {
	n, _ := strconv.Atoi(s)
	return uint32(n)
}

func getMSUInt32(s string) uint32 {
	n, _ := strconv.Atoi(s)
	nms := float64(n) / 1000
	return uint32(math.Floor(nms))
}

func getUInt64(s string) uint64 {
	n, _ := strconv.Atoi(s)
	return uint64(n)
}

func getInt64(s string) int64 {
	n, _ := strconv.Atoi(s)
	return int64(n)
}
