package configgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

const (
	AdminRoleAdminPrincipal = "Role.ADMIN"
	TempPath_MSPDIR         = "./tmp/tempmsp/"
	// OrdererGenesis_DEFAULT         = "OneOrgsOrdererGenesis"
	Consortium_DEFAULT = "SampleConsortium"
	// Channel_DEFAULT           = "OneOrgsChannel"
	Prefix string = "CONFIGTX"

	adminBaseName = "Admin"

	// MessageType_COMMON_BLOCK 配置块
	MessageType_COMMON_BLOCK    = "common.Block"
	MessageType_COMMON_CONFIG   = "common.Config"
	MessageType_COMMON_ENVELOPE = "common.Envelope"
)

var (
	ConfigPath_CONSORTIUMSGROUPS = []string{"channel_group", "groups", "Consortiums", "groups", "SampleConsortium", "groups"}
	ConfigPath_APPLICATIONGROUPS = []string{"channel_group", "groups", "Application", "groups"}
	ConfigDefault_ANCHORPEER     = localconfig.AnchorPeer{Host: "DefaultAnchorPeer", Port: 7051}

	ConfigPath_CONFIGUPDATEGROUPS            = []string{"write_set", "groups", "Consortiums", "groups", "SampleConsortium", "groups"}
	ConfigPath_CONFIGUPDATEAPPLICATIONGROUPS = []string{"write_set", "groups", "Application", "groups"}

	// ChannelNameSYSTEM 系统 chainid，创建 genesisprovider 时会更新
	ChannelNameSYSTEM = "SystemChannel" //TODO: 系统通道名称
)

type Organization struct {
	MSPId          string
	Name           string
	TLSCACert      string
	MSPCACert      string
	AdminCert      string
	AdminPrincipal string
	Addresses      []string // 应用组织此字段为AnchorPeers： peer0.example.com:7051  / orderer组织此字段为Orderer address orderer.example.com:7054
	MSPType        string
	// AnchorPeers    []string //127.0.0.1:7051
}

func CreateOrgCfgByOrg(org *Organization) (*genesisconfig.Organization, error) {
	RootDir := TempPath_MSPDIR
	orgmspDir := filepath.Join(RootDir, org.Name, "msp")
	SaveOrgMSPToTemp(org, orgmspDir)

	orgcfg := new(genesisconfig.Organization)
	orgcfg.ID = org.MSPId
	orgcfg.Name = org.Name
	if org.MSPType == "" {
		orgcfg.MSPType = msp.ProviderTypeToString(msp.FABRIC)
	}
	if org.AdminPrincipal == "" {
		orgcfg.AdminPrincipal = AdminRoleAdminPrincipal
	}
	orgcfg.MSPDir = orgmspDir
	for _, anc := range org.Addresses {
		strs := strings.Split(anc, ":")
		anccfg := &genesisconfig.AnchorPeer{}
		anccfg.Host = strs[0]
		port, err := strconv.Atoi(strs[1])
		if err != nil {
			return nil, errors.Wrapf(err, "[createOrg] convert port (%s) string error", strs[1])
		}
		anccfg.Port = port
		orgcfg.AnchorPeers = append(orgcfg.AnchorPeers, anccfg)
	}
	if len(orgcfg.AnchorPeers) == 0 {
		orgcfg.AnchorPeers = append(orgcfg.AnchorPeers, &ConfigDefault_ANCHORPEER)
	}
	return orgcfg, nil
}

func CreateGenesisBlock(config *genesisconfig.Profile) ([]byte, error) {
	pgen := encoder.New(config)
	if config.Consortiums == nil {
		return nil, errors.New("[CreateBlock] does not contain a consortiums group definition")
	}
	genesisBlock := pgen.GenesisBlockForChannel(ChannelNameSYSTEM)
	return utils.MarshalOrPanic(genesisBlock), nil
}

func CreateChannelBlock(config *genesisconfig.Profile, channelID string) ([]byte, error) {
	configtx, err := encoder.MakeChannelCreationTransaction(channelID, nil, config)
	if err != nil {
		return nil, errors.Wrapf(err, "[CreateChannelBlock] MakeChannelCreationTransaction error")
	}
	return utils.MarshalOrPanic(configtx), nil
}

func SaveOrgMSPToTemp(org *Organization, mspDir string) error {
	// mspDir := filepath.Join(TempPath_MSPDIR, org.Name, "msp")
	createFolderStructure(mspDir)
	caCertFileName := fmt.Sprintf("ca.%s-cert.pem", org.Name)
	tlscaCertFileName := fmt.Sprintf("tlsca.%s-cert.pem", org.Name)
	adminCertFileName := fmt.Sprintf("%s@%s-cert.pem", adminBaseName, org.Name)
	// adminCertFileName := fmt.Sprintf("%s@%s", adminBaseName, org.Name)
	err := saveFile(filepath.Join(mspDir, "cacerts"), caCertFileName, org.MSPCACert)
	if err != nil {
		return errors.Wrapf(err, "[SaveOrgMSPToTemp] saveFile (%s) error", caCertFileName)
	}
	err = saveFile(filepath.Join(mspDir, "tlscacerts"), tlscaCertFileName, org.TLSCACert)
	if err != nil {
		return errors.Wrapf(err, "[SaveOrgMSPToTemp] saveFile (%s) error", tlscaCertFileName)
	}
	err = saveFile(filepath.Join(mspDir, "admincerts"), adminCertFileName, org.AdminCert)
	if err != nil {
		return errors.Wrapf(err, "[SaveOrgMSPToTemp] saveFile (%s) error", adminCertFileName)
	}
	// saveFile(TempPath_MSPDIR, adminCertFileName, org.MSPCACert)
	return nil
}

func createFolderStructure(rootDir string) error {
	var folders []string
	folders = []string{
		filepath.Join(rootDir, "admincerts"),
		filepath.Join(rootDir, "cacerts"),
		filepath.Join(rootDir, "tlscacerts"),
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveFile(rootDir, filename string, cert string) error {
	filePath := filepath.Join(rootDir, filename)
	return ioutil.WriteFile(filePath, []byte(cert), 0755)
}
