package localconfig

import (
	tc "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
)

// GetOrdererProfileTemplate 获取创建 Orderer 创始块的基本配置模版
func GetOrdererProfileTemplate() *tc.Profile {
	defaultOrdererProfile := &DefaultProfileOrderer
	return copyProfile(defaultOrdererProfile)

}

// GetChannelProfileTemplate 获取创建 channel 基础交易块的模版
func GetChannelProfileTemplate() *tc.Profile {
	defaultChannelProfile := &DefaultProfileChannel
	return copyProfile(defaultChannelProfile)
}

func GetConsortiumTemplate() *tc.Consortium {
	defaultConsortium := &DefaultConsortium
	return copyConsortium(defaultConsortium)
}

func GetApplicationTemplate() *tc.Application {
	defaultApplication := &DefaultApplication
	return copyApplication(defaultApplication)
}

func copyCapabilities(cap map[string]bool) map[string]bool {
	if cap == nil {
		return nil
	}
	newCap := make(map[string]bool)
	for k, v := range cap {
		newCap[k] = v
	}
	return newCap
}

func copyPolicy(policy *tc.Policy) *tc.Policy {
	if policy == nil {
		return nil
	}
	newPol := &tc.Policy{}
	newPol.Type = policy.Type
	newPol.Rule = policy.Rule
	return newPol
}

func copyPolicies(policies map[string]*tc.Policy) map[string]*tc.Policy {
	if policies == nil {
		return nil
	}
	newPols := make(map[string]*tc.Policy)
	for k, v := range policies {
		pol := copyPolicy(v)
		newPols[k] = pol
	}
	return newPols
}

func copyACLs(acls map[string]string) map[string]string {
	if acls == nil {
		return nil
	}
	newACLs := make(map[string]string)
	for k, v := range acls {
		newACLs[k] = v
	}
	return newACLs
}

func copyResources(resource *tc.Resources) *tc.Resources {
	if resource == nil {
		return nil
	}
	return &tc.Resources{
		DefaultModPolicy: resource.DefaultModPolicy,
	}
}

func copyAncPeer(peer *tc.AnchorPeer) *tc.AnchorPeer {
	if peer == nil {
		return nil
	}
	newPeer := &tc.AnchorPeer{
		Host: peer.Host,
		Port: peer.Port,
	}
	return newPeer
}

func copyAncPeers(ancPeers []*tc.AnchorPeer) []*tc.AnchorPeer {
	var newAncPeers []*tc.AnchorPeer
	for _, anc := range ancPeers {
		newAnc := copyAncPeer(anc)
		newAncPeers = append(newAncPeers, newAnc)
	}
	return newAncPeers
}

func copyOrganization(org *tc.Organization) *tc.Organization {
	if org == nil {
		return nil
	}
	return &tc.Organization{
		Name:           org.Name,
		ID:             org.ID,
		MSPDir:         org.MSPDir,
		MSPType:        org.MSPType,
		Policies:       copyPolicies(org.Policies),
		AnchorPeers:    copyAncPeers(org.AnchorPeers),
		AdminPrincipal: org.AdminPrincipal,
	}
}

func copyOrganizations(orgs []*tc.Organization) []*tc.Organization {
	newOrgs := []*tc.Organization{}
	for _, org := range orgs {
		newOrgs = append(newOrgs, copyOrganization(org))
	}
	return newOrgs
}

func copyOrderer(orderer *tc.Orderer) *tc.Orderer {
	if orderer == nil {
		return nil
	}
	newOrderer := &tc.Orderer{
		OrdererType:  orderer.OrdererType,
		Addresses:    nil,
		BatchTimeout: orderer.BatchTimeout,
		BatchSize: tc.BatchSize{
			MaxMessageCount:   orderer.BatchSize.MaxMessageCount,
			AbsoluteMaxBytes:  orderer.BatchSize.AbsoluteMaxBytes,
			PreferredMaxBytes: orderer.BatchSize.PreferredMaxBytes,
		},
		Kafka: tc.Kafka{
			Brokers: nil,
		},
		EtcdRaft:      nil, //TODO:
		Organizations: copyOrganizations(orderer.Organizations),
		MaxChannels:   orderer.MaxChannels,
		Capabilities:  copyCapabilities(orderer.Capabilities),
		Policies:      copyPolicies(orderer.Policies),
	}
	addr := []string{}
	copy(addr, orderer.Addresses)
	newOrderer.Addresses = addr
	brokers := []string{}
	copy(brokers, orderer.Kafka.Brokers)
	newOrderer.Kafka.Brokers = brokers
	return newOrderer
}

func copyConsortium(consortium *tc.Consortium) *tc.Consortium {
	if consortium == nil {
		return nil
	}
	newCons := &tc.Consortium{
		Organizations: []*tc.Organization{},
	}
	for _, org := range consortium.Organizations {
		newCons.Organizations = append(newCons.Organizations, copyOrganization(org))
	}
	return newCons
}

func copyConsortiums(consortiums map[string]*tc.Consortium) map[string]*tc.Consortium {
	if consortiums == nil {
		return nil
	}
	newConss := make(map[string]*tc.Consortium)
	for k, v := range consortiums {
		newConss[k] = copyConsortium(v)
	}
	return newConss
}

func copyApplication(app *tc.Application) *tc.Application {
	if app == nil {
		return nil
	}
	newApp := &tc.Application{
		Organizations: nil,
		Capabilities:  copyCapabilities(app.Capabilities),
		Resources:     copyResources(app.Resources),
		Policies:      copyPolicies(app.Policies),
		ACLs:          copyACLs(app.ACLs),
	}
	newOrgs := []*tc.Organization{}
	for _, org := range app.Organizations {
		newOrgs = append(newOrgs, copyOrganization(org))
	}
	newApp.Organizations = newOrgs
	return newApp
}

func copyProfile(profile *tc.Profile) *tc.Profile {
	if profile == nil {
		return nil
	}
	newPro := &tc.Profile{
		Consortium:   profile.Consortium,
		Application:  copyApplication(profile.Application),
		Orderer:      copyOrderer(profile.Orderer),
		Consortiums:  copyConsortiums(profile.Consortiums),
		Capabilities: copyCapabilities(profile.Capabilities),
		Policies:     copyPolicies(profile.Policies),
	}
	return newPro
}
