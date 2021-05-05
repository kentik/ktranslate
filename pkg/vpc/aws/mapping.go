package aws

import (
	"context"
	"time"

	tree "github.com/kentik/patricia/string_tree"
)

type MappingSet struct {
	SubnetTrieV4 *tree.TreeV4
	SubnetTrieV6 *tree.TreeV6
}

func NewMappingSet() *MappingSet {
	return &MappingSet{
		SubnetTrieV4: tree.NewTreeV4(),
		SubnetTrieV6: tree.NewTreeV6(),
	}
}

func (vpc *AwsVpc) checkMappings(ctx context.Context) {
	checkTicker := time.NewTicker(MappingCheckDuration)
	defer func() {
		checkTicker.Stop()
	}()

	err := vpc.updateMapping(ctx)
	if err != nil {
		vpc.Errorf("Cannot get mapping %v", err)
	}

	go func() {
		vpc.Infof("checkMapping Online")
		for {
			select {
			case <-checkTicker.C:
				err := vpc.updateMapping(ctx)
				if err != nil {
					vpc.Errorf("Cannot get mapping %v", err)
				}

			case <-ctx.Done():
				vpc.Infof("checkMapping Done")
				return
			}
		}
	}()
}

func (vpc *AwsVpc) updateMapping(ctx context.Context) error {
	start := time.Now()

	topo, allGood := FetchAllEntities(ctx, vpc, *IamRole, vpc.regions)
	if !allGood {
		vpc.Warnf("Failed to fetch some mappings")
	}

	// And if good, update our mapping set.
	vpc.mux.Lock()
	vpc.topo = &topo
	vpc.mux.Unlock()
	vpc.Infof("Updated mappings -- %d vpcs, %d regions in %v", len(topo.Entities.Vpcs), len(topo.Hierarchy.Regions), time.Now().Sub(start))

	return nil
}
