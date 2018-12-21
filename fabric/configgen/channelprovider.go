package configgen

import (
	"os"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	cfg "github.com/tinywell/utils/fabric/configgen/localconfig"
)

type ChannelBlockProvider struct {
	profile   *localconfig.Profile
	ChannelID string
}

func NewChannelBlockProvider(ChannelID string) *ChannelBlockProvider {
	provider := &ChannelBlockProvider{ChannelID: ChannelID}
	provider.initProvider()
	return provider
}

func (provider *ChannelBlockProvider) initProvider() {
	factory.InitFactories(nil)
	provider.profile = cfg.GetChannelProfileTemplate()
}

func (provider *ChannelBlockProvider) SetApplicationOrgs(orgs []*Organization) error {
	var norgs []*localconfig.Organization
	for _, org := range orgs {
		orgcfg, err := CreateOrgCfgByOrg(org)
		if err != nil {
			return err
		}
		norgs = append(norgs, orgcfg)
	}
	provider.profile.Application.Organizations = norgs
	return nil
}

func (provider *ChannelBlockProvider) Marshal() ([]byte, error) {
	cfg, err := CreateChannelBlock(provider.profile, provider.ChannelID)
	if err != nil {
		return nil, err
	}
	os.RemoveAll(TempPath_MSPDIR)
	return cfg, nil
}
