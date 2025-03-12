package snmp

/**
Use netbox API to pull down list of devices to target. Get all devices matching a given tag.
*/
import (
	"context"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	"github.com/kentik/ktranslate/pkg/kt"

	"github.com/netbox-community/go-netbox/v4"
)

func getDevicesFromNetbox(ctx context.Context, conf *kt.SnmpConfig, log logger.ContextL) error {
	c := netbox.NewAPIClientFor(conf.Disco.NetboxAPIHost, conf.Disco.NetboxAPIToken)

	res, _, err := c.DcimAPI.
		DcimDevicesList(ctx).
		Status([]string{"active"}).
		Limit(10).
		Execute()

	if err != nil {
		return err
	}

	log.Infof("%v", res.Results)

	return nil

}
