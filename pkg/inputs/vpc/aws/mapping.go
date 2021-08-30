package aws

import (
	"context"
	"time"
)

func (vpc *AwsVpc) checkMappings(ctx context.Context) {
	checkTicker := time.NewTicker(MappingCheckDuration)
	defer func() {
		checkTicker.Stop()
	}()

	err := vpc.updateMapping(ctx)
	if err != nil {
		vpc.Errorf("Cannot get mapping %v", err)
	}

	vpc.Infof("checkMapping Online")
	for {
		select {
		case <-checkTicker.C:
			vpc.Infof("Starting to update mapping")
			err := vpc.updateMapping(ctx)
			if err != nil {
				vpc.Errorf("Cannot get mapping %v", err)
			}

		case <-ctx.Done():
			vpc.Infof("checkMapping Done")
			return
		}
	}
}

func (vpc *AwsVpc) updateMapping(ctx context.Context) error {
	start := time.Now()

	topo, allGood := FetchAllEntities(ctx, vpc, *IamRole, vpc.regions)
	if !allGood {
		vpc.Warnf("There was an error when fetching mappings.")
	}

	// And if good, update our mapping set.
	vpc.topo = &topo
	vpc.Infof("Updated mappings -- %d vpcs, %d regions in %v", len(topo.Entities.Vpcs), len(topo.Hierarchy.Regions), time.Now().Sub(start))

	return nil
}
