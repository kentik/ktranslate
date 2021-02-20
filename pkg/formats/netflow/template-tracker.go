package netflow

import (
	"fmt"
	"strings"

	"github.com/kentik/netflow/ipfix"
	"github.com/kentik/netflow/netflow9"
)

type templateTracker map[uint32]map[uint16]string // maps observationDomainID->templateID->template-string

func newTemplateTracker() templateTracker {
	return make(map[uint32]map[uint16]string)
}

func (ipf *NetflowFormat) addTemplate(srcId uint32, templateId uint16, kind string, templateStr string) {

	if ipf.templateTracker[srcId] == nil {
		ipf.templateTracker[srcId] = make(map[uint16]string)
	}

	if templateStr != ipf.templateTracker[srcId][templateId] {
		if ipf.templateTracker[srcId][templateId] == "" {
			ipf.Debugf("received new %s template for observation domain id %d, template id %d:  %s", kind, srcId, templateId, templateStr)
		} else {
			ipf.Warnf("received modified %s template for observation domain id %d, template id %d:  %s", kind, srcId, templateId, templateStr)
		}
		ipf.templateTracker[srcId][templateId] = templateStr
	}
}

func (ipf *NetflowFormat) addIpfixTemplate(srcId uint32, tr ipfix.TemplateRecord) {
	var b strings.Builder
	for _, field := range tr.Fields {
		fmt.Fprintf(&b, "%d:%d (%d), ", field.EnterpriseNumber, field.InformationElementID, field.Length)
	}

	ipf.addTemplate(srcId, tr.TemplateID, "ipfix", b.String())
}

func (ipf *NetflowFormat) addIpfixOptionsTemplate(srcId uint32, or ipfix.OptionsTemplateRecord) {
	var b strings.Builder
	b.WriteString("options scope: ")
	for _, field := range or.ScopeFields {
		fmt.Fprintf(&b, "%d:%d (%d), ", field.EnterpriseNumber, field.InformationElementID, field.Length)
	}

	b.WriteString("data: ")
	for _, field := range or.Fields {
		fmt.Fprintf(&b, "%d:%d (%d), ", field.EnterpriseNumber, field.InformationElementID, field.Length)
	}

	ipf.addTemplate(srcId, or.TemplateID, "ipfix options", b.String())
}

func (ipf *NetflowFormat) addV9Template(srcId uint32, tr netflow9.TemplateRecord) {
	var b strings.Builder
	for _, field := range tr.Fields {
		fmt.Fprintf(&b, "%d (%d), ", field.Type, field.Length)
	}

	ipf.addTemplate(srcId, tr.TemplateID, "v9", b.String())
}

func (ipf *NetflowFormat) addV9OptionsTemplate(srcId uint32, or netflow9.OptionsTemplateRecord) {
	var b strings.Builder
	b.WriteString("options scope: ")
	for _, field := range or.ScopeFields {
		fmt.Fprintf(&b, "%d (%d), ", field.Type, field.Length)
	}

	b.WriteString("data: ")
	for _, field := range or.Fields {
		fmt.Fprintf(&b, "%d (%d), ", field.Type, field.Length)
	}

	ipf.addTemplate(srcId, or.TemplateID, "v9 options", b.String())
}
