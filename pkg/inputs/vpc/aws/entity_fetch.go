package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/kentik/patricia"
	tree "github.com/kentik/patricia/string_tree"
)

// AWSEntities holds all of the entities fetched for a company
type AWSEntities struct {
	Vpcs                      map[string]ec2.Vpc                      `json:"Vpcs"`
	AvailabilityZones         map[string]ec2.AvailabilityZone         `json:"AvailabilityZones"`
	Subnets                   map[string]ec2.Subnet                   `json:"Subnets"`
	InternetGateways          map[string]ec2.InternetGateway          `json:"InternetGateways"`
	NatGateways               map[string]ec2.NatGateway               `json:"NatGateways"`
	TransitGateways           map[string]ec2.TransitGateway           `json:"TransitGateways"`
	TransitGatewayAttachments map[string]ec2.TransitGatewayAttachment `json:"TransitGatewayAttachments"`
	VpnGateways               map[string]ec2.VpnGateway               `json:"VpnGateways"`
	VpcPeeringConnections     map[string]ec2.VpcPeeringConnection     `json:"VpcPeeringConnections"`
	// TODO: Firewalls
	// TODO: LoadBalancers             map[string]ec2.ClassicLoadBalancer      `json:"load_balancers"`
}

type AWSTopology struct {
	Hierarchy AWSHierarchy `json:"Hierarchy"`
	Entities  AWSEntities  `json:"Entities"`
}

func NewAWSTopology() AWSTopology {
	return AWSTopology{
		Hierarchy: NewAWSHierarchy(),
		Entities:  NewAWSEntities(),
	}
}

type AWSHierarchy struct {
	Regions      map[string]RegionSkel `json:"Regions"`
	SubnetTrieV4 *tree.TreeV4
	SubnetTrieV6 *tree.TreeV6
}

func NewAWSHierarchy() AWSHierarchy {
	return AWSHierarchy{
		Regions:      make(map[string]RegionSkel),
		SubnetTrieV4: tree.NewTreeV4(),
		SubnetTrieV6: tree.NewTreeV6(),
	}
}

type VpcPeeringConnectionSkel struct {
	RequesterVpcId string `json:"RequesterVpcId"`
	AccepterVpcId  string `json:"AccepterVpcId"`
}

type SubnetSkel struct {
	SubnetId    string                    `json:"SubnetId"`
	NatGateways map[string]NatGatewaySkel `json:"NatGateways"`
}

func NewSubnetSkel(subnetId string) SubnetSkel {
	return SubnetSkel{
		SubnetId:    subnetId,
		NatGateways: make(map[string]NatGatewaySkel),
	}
}

type InternetGatewaySkel struct {
	InternetGatewayId          string                                   `json:"InternetGatewayId"`
	InternetGatewayAttachments map[string]InternetGatewayAttachmentSkel `json:"InternetGatewayAttachments"`
}

func NewInternetGatewaySkel(id string) InternetGatewaySkel {
	return InternetGatewaySkel{
		InternetGatewayId:          id,
		InternetGatewayAttachments: make(map[string]InternetGatewayAttachmentSkel),
	}
}

type InternetGatewayAttachmentSkel struct {
	InternetGatewayId string `json:"InternetGatewayId"`
	VpcId             string `json:"VpcId"`
	State             string `json:"State"`
}

type NatGatewaySkel struct {
	NatGatewayId string `json:"NatGatewayId"`
}

type TransitGatewaySkel struct {
	TransitGatewayId          string                                  `json:"TransitGatewayId "`
	TransitGatewayAttachments map[string]TransitGatewayAttachmentSkel `json:"TransitGatewayAttachment"`
}

func NewTransitGatewaySkel(id string) TransitGatewaySkel {
	return TransitGatewaySkel{
		TransitGatewayId:          id,
		TransitGatewayAttachments: make(map[string]TransitGatewayAttachmentSkel),
	}
}

type TransitGatewayAttachmentSkel struct {
	TransitGatewayAttachmentId string `json:"TransitGatewayAttachmentId"`
}

type VpnGatewaySkel struct {
	VpnGatewayId string `json:"VpnGatewayId"`
}

type VpcSkel struct {
	VpcId                      string                                   `json:"VpcId"`
	Subnets                    map[string]SubnetSkel                    `json:"Subnets"`
	TransitGatewayAttachments  map[string]TransitGatewayAttachmentSkel  `json:"TransitGatewayAttachments"`  // by GatewayAttachmentId
	InternetGatewayAttachments map[string]InternetGatewayAttachmentSkel `json:"InternetGatewayAttachments"` // by InternetGatewayId
	VpcPeeringConnections      map[string]VpcPeeringConnectionSkel      `json:"VpcPeeringConnections"`      // by VpcPeeringConnectionId
}

func NewVpcSkel(id string) VpcSkel {
	return VpcSkel{
		VpcId:                      id,
		Subnets:                    make(map[string]SubnetSkel),
		TransitGatewayAttachments:  make(map[string]TransitGatewayAttachmentSkel),
		InternetGatewayAttachments: make(map[string]InternetGatewayAttachmentSkel),
	}
}

type AvailabilityZoneSkel struct {
	ZoneId string `json:"ZoneId"`
}

type RegionSkel struct {
	Name              string                          `json:"Name"`
	Vpcs              map[string]VpcSkel              `json:"Vpcs"`
	AvailabilityZones map[string]AvailabilityZoneSkel `json:"AvailabilityZones"`
	InternetGateways  map[string]InternetGatewaySkel  `json:"InternetGateways"`
	TransitGateways   map[string]TransitGatewaySkel   `json:"TransitGateways"`
	VpnGateways       map[string]VpnGatewaySkel       `json:"VpnGateways"`
}

func NewRegionSkel(name string) RegionSkel {
	return RegionSkel{
		Name:              strings.ToLower(name),
		Vpcs:              make(map[string]VpcSkel),
		AvailabilityZones: make(map[string]AvailabilityZoneSkel),
		InternetGateways:  make(map[string]InternetGatewaySkel),
		TransitGateways:   make(map[string]TransitGatewaySkel),
		VpnGateways:       make(map[string]VpnGatewaySkel),
	}
}

func NewAWSEntities() AWSEntities {
	return AWSEntities{
		Vpcs:                      make(map[string]ec2.Vpc),
		AvailabilityZones:         make(map[string]ec2.AvailabilityZone),
		Subnets:                   make(map[string]ec2.Subnet),
		InternetGateways:          make(map[string]ec2.InternetGateway),
		NatGateways:               make(map[string]ec2.NatGateway),
		TransitGateways:           make(map[string]ec2.TransitGateway),
		TransitGatewayAttachments: make(map[string]ec2.TransitGatewayAttachment),
		VpnGateways:               make(map[string]ec2.VpnGateway),
		VpcPeeringConnections:     make(map[string]ec2.VpcPeeringConnection),
	}
}

// FetchAllEntities fetches all the things, and returns whether it was a total success
func FetchAllEntities(ctx context.Context, log logger.ContextL, arnName string, regions []string) (AWSTopology, bool) {
	topology := NewAWSTopology()
	sleepTime := 50 * time.Millisecond
	totalSuccess := true
	start := time.Now()
	vpcToRegion := make(map[string]string)
	for _, regionName := range regions {
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(regionName),
		}))

		if _, found := topology.Hierarchy.Regions[regionName]; !found {
			topology.Hierarchy.Regions[regionName] = NewRegionSkel(regionName)
		}

		ec2cli := ec2.New(session,
			&aws.Config{
				Credentials: stscreds.NewCredentials(session, arnName),
			})

		// Vpcs
		if err := fetchVpcs(ctx, ec2cli, &topology, regionName, vpcToRegion, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// Availability Zones
		if err := fetchAvailabilityZones(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// InternetGateways
		if err := fetchInternetGateways(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// TransitGateways
		if err := fetchTransitGateways(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// VpnGateways
		if err := fetchVPNGateways(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// Subnets
		if err := fetchSubnets(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// NatGateways
		if err := fetchNATGateways(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// TransitGatewayAttachments
		if err := fetchTransitGatewayAttachments(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// VpcPeeringConnections
		if err := fetchVpcPeeringConnections(ctx, ec2cli, &topology, regionName, sleepTime); err != nil {
			totalSuccess = false
			log.Errorf("region: %s, arn: %s - %s", regionName, arnName, err)
		}

		// TODO: Firewalls
		// TODO: LoadBalancers
	}

	// VpcPeeringConnections cross VPCs - we've loaded them all up in the entities collections,
	// but need to tie them to VPCs after crawling all regions, since they can span regions
	for pcId, pc := range topology.Entities.VpcPeeringConnections {
		requesterId := ""
		accepterId := ""
		if pc.RequesterVpcInfo != nil {
			requesterId = safeStr(pc.RequesterVpcInfo.VpcId)
		}
		if pc.AccepterVpcInfo != nil {
			accepterId = safeStr(pc.AccepterVpcInfo.VpcId)
		}
		if requesterId == "" && accepterId == "" {
			continue
		}

		skel := VpcPeeringConnectionSkel{
			RequesterVpcId: requesterId,
			AccepterVpcId:  accepterId,
		}

		attachSkelToVpc := func(vpcId string) {
			if vpcId == "" {
				return
			}
			if regionName, found := vpcToRegion[vpcId]; found {
				if _, found := topology.Hierarchy.Regions[regionName].Vpcs[vpcId]; !found {
					// I don't think this should happen
					topology.Hierarchy.Regions[regionName].Vpcs[vpcId] = NewVpcSkel(vpcId)
				}
				topology.Hierarchy.Regions[regionName].Vpcs[vpcId].VpcPeeringConnections[pcId] = skel
			}
		}

		attachSkelToVpc(requesterId)
		attachSkelToVpc(accepterId)
	}

	log.Infof("Finished fetching entites in %s", time.Since(start))
	return topology, totalSuccess
}

func fetchVpcs(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, vpcToRegion map[string]string, sleepTime time.Duration) error {
	err := ec2cli.DescribeVpcsPagesWithContext(ctx, &ec2.DescribeVpcsInput{}, func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
		for _, vpc := range page.Vpcs {
			if vpc == nil || vpc.VpcId == nil {
				continue
			}
			vpcToRegion[*vpc.VpcId] = regionName
			topology.Entities.Vpcs[*vpc.VpcId] = *vpc
			topology.Hierarchy.Regions[regionName].Vpcs[*vpc.VpcId] = NewVpcSkel(*vpc.VpcId)
		}
		time.Sleep(sleepTime)
		return true
	})
	if err != nil {
		return fmt.Errorf("There was an error when fetching VPCS: %s.", err)
	}
	return nil
}

func fetchAvailabilityZones(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	out, err := ec2cli.DescribeAvailabilityZonesWithContext(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return fmt.Errorf("There was an error when fetching availability zones: %s.", err)
	}
	for _, az := range out.AvailabilityZones {
		if az == nil || az.ZoneId == nil {
			continue
		}
		topology.Entities.AvailabilityZones[*az.ZoneId] = *az
		topology.Hierarchy.Regions[regionName].AvailabilityZones[*az.ZoneId] = AvailabilityZoneSkel{ZoneId: *az.ZoneId}
	}
	time.Sleep(sleepTime)
	return nil
}

func fetchSubnets(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeSubnetsPagesWithContext(ctx,
		&ec2.DescribeSubnetsInput{},
		func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			for _, subnet := range page.Subnets {
				if subnet == nil || subnet.SubnetId == nil {
					continue
				}
				topology.Entities.Subnets[*subnet.SubnetId] = *subnet
				if subnet.VpcId == nil {
					continue
				}
				if _, found := topology.Hierarchy.Regions[regionName].Vpcs[*subnet.VpcId]; !found {
					// shouldn't happen - we crawled VPCs first - do this just in case
					topology.Hierarchy.Regions[regionName].Vpcs[*subnet.VpcId] = NewVpcSkel(*subnet.VpcId)
				}
				topology.Hierarchy.Regions[regionName].Vpcs[*subnet.VpcId].Subnets[*subnet.SubnetId] = NewSubnetSkel(*subnet.SubnetId)

				// Now, add to out lookup.
				ip4, ip6, err := patricia.ParseIPFromString(*subnet.CidrBlock)
				if err != nil {
					continue
				}
				if ip4 != nil {
					topology.Hierarchy.SubnetTrieV4.Set(*ip4, *subnet.SubnetId)
				} else {
					topology.Hierarchy.SubnetTrieV6.Set(*ip6, *subnet.SubnetId)
				}
			}
			time.Sleep(sleepTime)
			return true
		})
	if err != nil {
		return fmt.Errorf("There was an error when fetching subnets: %s.", err)
	}
	return nil
}

func fetchInternetGateways(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeInternetGatewaysPagesWithContext(ctx,
		&ec2.DescribeInternetGatewaysInput{},
		func(page *ec2.DescribeInternetGatewaysOutput, lastPage bool) bool {
			for _, gw := range page.InternetGateways {
				if gw == nil || gw.InternetGatewayId == nil {
					continue
				}
				topology.Entities.InternetGateways[*gw.InternetGatewayId] = *gw
				topology.Hierarchy.Regions[regionName].InternetGateways[*gw.InternetGatewayId] = NewInternetGatewaySkel(*gw.InternetGatewayId)

				// add the attachments to both the internet gateway, and the VPC
				for _, gwa := range gw.Attachments {
					if gwa.VpcId == nil {
						continue
					}
					skel := InternetGatewayAttachmentSkel{
						InternetGatewayId: *gw.InternetGatewayId,
						VpcId:             *gwa.VpcId,
						State:             safeStr(gwa.State),
					}
					topology.Hierarchy.Regions[regionName].InternetGateways[*gw.InternetGatewayId].InternetGatewayAttachments[*gwa.VpcId] = skel

					if _, found := topology.Hierarchy.Regions[regionName].Vpcs[*gwa.VpcId]; !found {
						// shouldn't happen - we crawled VPCs first - do this just in case
						topology.Hierarchy.Regions[regionName].Vpcs[*gwa.VpcId] = NewVpcSkel(*gwa.VpcId)
					}
					topology.Hierarchy.Regions[regionName].Vpcs[*gwa.VpcId].InternetGatewayAttachments[*gw.InternetGatewayId] = skel
				}
			}
			time.Sleep(sleepTime)
			return true
		})
	if err != nil {
		return fmt.Errorf("There was an error when fetching internet gateways: %s.", err)
	}
	return nil
}

func fetchNATGateways(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeNatGatewaysPagesWithContext(ctx,
		&ec2.DescribeNatGatewaysInput{},
		func(page *ec2.DescribeNatGatewaysOutput, lastPage bool) bool {
			for _, gw := range page.NatGateways {
				if gw == nil || gw.NatGatewayId == nil {
					continue
				}
				topology.Entities.NatGateways[*gw.NatGatewayId] = *gw
				if gw.SubnetId == nil || gw.VpcId == nil {
					continue
				}
				if _, found := topology.Hierarchy.Regions[regionName].Vpcs[*gw.VpcId]; !found {
					// shouldn't happen, but just in case
					topology.Hierarchy.Regions[regionName].Vpcs[*gw.VpcId] = NewVpcSkel(*gw.VpcId)
				}
				if _, found := topology.Hierarchy.Regions[regionName].Vpcs[*gw.VpcId].Subnets[*gw.SubnetId]; !found {
					// shouldn't happen, but just in case
					topology.Hierarchy.Regions[regionName].Vpcs[*gw.VpcId].Subnets[*gw.SubnetId] = NewSubnetSkel(*gw.SubnetId)
				}
				topology.Hierarchy.Regions[regionName].Vpcs[*gw.VpcId].Subnets[*gw.SubnetId].NatGateways[*gw.NatGatewayId] = NatGatewaySkel{NatGatewayId: *gw.NatGatewayId}
			}
			time.Sleep(sleepTime)
			return true
		})
	if err != nil {
		return fmt.Errorf("There was an error when fetching NAT gateways: %s.", err)
	}
	return nil
}

func fetchTransitGateways(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeTransitGatewaysPagesWithContext(ctx,
		&ec2.DescribeTransitGatewaysInput{},
		func(page *ec2.DescribeTransitGatewaysOutput, lastPage bool) bool {
			for _, gw := range page.TransitGateways {
				if gw == nil || gw.TransitGatewayId == nil {
					continue
				}
				topology.Entities.TransitGateways[*gw.TransitGatewayId] = *gw
				topology.Hierarchy.Regions[regionName].TransitGateways[*gw.TransitGatewayId] = NewTransitGatewaySkel(*gw.TransitGatewayId)
			}
			time.Sleep(sleepTime)
			return true
		})
	if err != nil {
		return fmt.Errorf("There was an error when fetching transit gateways: %s.", err)
	}
	return nil
}

func fetchTransitGatewayAttachments(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeTransitGatewayAttachmentsPagesWithContext(ctx,
		&ec2.DescribeTransitGatewayAttachmentsInput{},
		func(page *ec2.DescribeTransitGatewayAttachmentsOutput, lastPage bool) bool {
			for _, at := range page.TransitGatewayAttachments {
				if at == nil || at.TransitGatewayAttachmentId == nil {
					continue
				}

				topology.Entities.TransitGatewayAttachments[*at.TransitGatewayAttachmentId] = *at

				// add the attachment to the transit gateway
				if at.TransitGatewayId != nil {
					if _, found := topology.Hierarchy.Regions[regionName].TransitGateways[*at.TransitGatewayId]; !found {
						// shouldn't happen
						topology.Hierarchy.Regions[regionName].TransitGateways[*at.TransitGatewayId] = NewTransitGatewaySkel(*at.TransitGatewayId)
					}
					topology.Hierarchy.Regions[regionName].TransitGateways[*at.TransitGatewayId].TransitGatewayAttachments[*at.TransitGatewayAttachmentId] = TransitGatewayAttachmentSkel{TransitGatewayAttachmentId: *at.TransitGatewayId}
				}

				// also add the attachment to the appropriate resource, if we model it
				// valid values for ResourceType: vpc, vpn, direct-connect-gateway, peering, connect
				if at.ResourceId != nil && at.ResourceType != nil && strings.ToLower(*at.ResourceType) == "vpc" {
					// VPCs for this region should already have been loaded
					if _, found := topology.Hierarchy.Regions[regionName].Vpcs[*at.ResourceId]; !found {
						// shouldn't happen - we crawled VPCs first - do this just in case
						topology.Hierarchy.Regions[regionName].Vpcs[*at.ResourceId] = NewVpcSkel(*at.ResourceId)
					}
					topology.Hierarchy.Regions[regionName].Vpcs[*at.ResourceId].TransitGatewayAttachments[*at.TransitGatewayAttachmentId] = TransitGatewayAttachmentSkel{TransitGatewayAttachmentId: *at.TransitGatewayId}
				}
			}

			time.Sleep(sleepTime)
			return true
		})
	if err != nil {
		return fmt.Errorf("There was an error when fetching transit gateways attachments: %s.", err)
	}
	return nil
}

func fetchVPNGateways(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	output, err := ec2cli.DescribeVpnGatewaysWithContext(ctx, &ec2.DescribeVpnGatewaysInput{})
	if err != nil {
		return fmt.Errorf("There was an error when fetching VPN gateways: %s.", err)
	}

	for _, gw := range output.VpnGateways {
		if gw == nil || gw.VpnGatewayId == nil {
			continue
		}
		topology.Entities.VpnGateways[*gw.VpnGatewayId] = *gw
		topology.Hierarchy.Regions[regionName].VpnGateways[*gw.VpnGatewayId] = VpnGatewaySkel{VpnGatewayId: *gw.VpnGatewayId}
	}
	time.Sleep(sleepTime)
	return nil
}

func safeStr(strPtr *string) string {
	if strPtr == nil {
		return ""
	}
	return *strPtr
}

func fetchVpcPeeringConnections(ctx context.Context, ec2cli *ec2.EC2, topology *AWSTopology, regionName string, sleepTime time.Duration) error {
	err := ec2cli.DescribeVpcPeeringConnectionsPages(&ec2.DescribeVpcPeeringConnectionsInput{}, func(page *ec2.DescribeVpcPeeringConnectionsOutput, lastPage bool) bool {
		for _, pc := range page.VpcPeeringConnections {
			if pc.VpcPeeringConnectionId == nil {
				continue
			}

			topology.Entities.VpcPeeringConnections[*pc.VpcPeeringConnectionId] = *pc
		}

		time.Sleep(sleepTime)
		return true
	})
	if err != nil {
		return fmt.Errorf("There was an error when fetching VPC peering connections gateways: %s.", err)
	}
	return nil
}
