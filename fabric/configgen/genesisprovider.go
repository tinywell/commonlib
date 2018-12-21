package configgen

import (
	"os"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/pkg/errors"

	cfg "github.com/tinywell/utils/fabric/configgen/localconfig"
)

type GenesisBlockProvider struct {
	syschainid string
	profile    *localconfig.Profile
}

func NewGenesisBlockProvider(syschainid string) *GenesisBlockProvider {
	provider := &GenesisBlockProvider{syschainid: syschainid}
	provider.initProvider()
	return provider
}

func (provider *GenesisBlockProvider) initProvider() {
	factory.InitFactories(nil)
	provider.profile = cfg.GetOrdererProfileTemplate()
	ChannelNameSYSTEM = provider.syschainid
}

func (provider *GenesisBlockProvider) SetApplicationOrgs(orgs []*Organization) error {
	var norgs []*localconfig.Organization
	for _, org := range orgs {
		orgcfg, err := CreateOrgCfgByOrg(org)
		if err != nil {
			return err
		}
		norgs = append(norgs, orgcfg)
	}
	_, ok := provider.profile.Consortiums[Consortium_DEFAULT]
	if ok {
		provider.profile.Consortiums[Consortium_DEFAULT].Organizations = norgs
	} else {
		return errors.Errorf("[CreateGenesisBlock] No %s in config", Consortium_DEFAULT)
	}
	return nil
}

func (provider *GenesisBlockProvider) SetOrdererOrgs(orgs []*Organization) error {
	var orderers []*localconfig.Organization
	var addresses []string
	for _, org := range orgs {
		orgcfg, err := CreateOrgCfgByOrg(org)
		if err != nil {
			return err
		}
		orgcfg.AnchorPeers = nil
		orderers = append(orderers, orgcfg)
		addresses = append(addresses, org.Addresses...)
	}
	provider.profile.Orderer.Organizations = orderers
	provider.SetOrdererAddress(addresses)
	return nil
}

func (provider *GenesisBlockProvider) SetOrdererType(ordererType string) {
	provider.profile.Orderer.OrdererType = ordererType
}

func (provider *GenesisBlockProvider) SetKafkas(kafkas []string) {
	if len(kafkas) > 0 {
		provider.profile.Orderer.Kafka.Brokers = kafkas
	}
}

func (provider *GenesisBlockProvider) SetOrdererAddress(addresses []string) {
	if len(addresses) > 0 {
		provider.profile.Orderer.Addresses = addresses
	}
}

func (provider *GenesisBlockProvider) Marshal() ([]byte, error) {
	cfg, err := CreateGenesisBlock(provider.profile)
	if err != nil {
		return nil, err
	}
	os.RemoveAll(TempPath_MSPDIR)
	return cfg, nil
}
