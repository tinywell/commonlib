package localconfig

import (
	"time"

	tc "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/msp"
)

const (
	// TestChainID is the channel name used for testing purposes when one is
	// not given
	TestChainID = "testchainid"

	DefaultSysChainID = "systemchainid"

	// SampleConsortiumName is the sample consortium from the
	// sample configtx.yaml
	SampleConsortiumName = "SampleConsortium"
	// SampleOrgName is the name of the sample org in the sample profiles
	SampleOrgName = "SampleOrg"

	// AdminRoleAdminPrincipal is set as AdminRole to cause the MSP role of
	// type Admin to be used as the admin principal default
	AdminRoleAdminPrincipal = "Role.ADMIN"
	// MemberRoleAdminPrincipal is set as AdminRole to cause the MSP role of
	// type Member to be used as the admin principal default
	MemberRoleAdminPrincipal = "Role.MEMBER"
)

////////////////////////////////////////////////////////////////////////////////
//
//   CAPABILITIES
//
//   This section defines the capabilities of fabric network. This is a new
//   concept as of v1.1.0 and should not be utilized in mixed networks with
//   v1.0.x peers and orderers.  Capabilities define features which must be
//   present in a fabric binary for that binary to safely participate in the
//   fabric network.  For instance, if a new MSP type is added, newer binaries
//   might recognize and validate the signatures from this type, while older
//   binaries without this support would be unable to validate those
//   transactions.  This could lead to different versions of the fabric binaries
//   having different world states.  Instead, defining a capability for a channel
//   informs those binaries without this capability that they must cease
//   processing transactions until they have been upgraded.  For v1.0.x if any
//   capabilities are defined (including a map with all capabilities turned off)
//   then the v1.0.x peer will deliberately crash.
//
////////////////////////////////////////////////////////////////////////////////
var DefaultCapabilitiesOrderer = map[string]bool{
	"V1_1": true,
}

var DefaultCapabilitiesChannel = map[string]bool{
	"V1.3": true,
}

var DefaultCapabilitiesApplication = map[string]bool{
	"V1_3": true,
	"V1_2": false,
	"V1_1": false,
}

var DefaultACLs = map[string]string{
	"lscc/ChaincodeExists":           "/Channel/Application/Readers",
	"lscc/GetDeploymentSpec":         "/Channel/Application/Readers",
	"lscc/GetChaincodeData":          "/Channel/Application/Readers",
	"lscc/GetInstantiatedChaincodes": "/Channel/Application/Readers",
	"qscc/GetChainInfo":              "/Channel/Application/Readers",
	"qscc/GetBlockByNumber":          "/Channel/Application/Readers",
	"qscc/GetBlockByHash":            "/Channel/Application/Readers",
	"qscc/GetTransactionByID":        "/Channel/Application/Readers",
	"qscc/GetBlockByTxID":            "/Channel/Application/Readers",
	"cscc/GetConfigBlock":            "/Channel/Application/Readers",
	"cscc/GetConfigTree":             "/Channel/Application/Readers",
	"cscc/SimulateConfigTreeUpdate":  "/Channel/Application/Readers",
	"peer/Propose":                   "/Channel/Application/Writers",
	"peer/ChaincodeToChaincode":      "/Channel/Application/Readers",
	"event/Block":                    "/Channel/Application/Readers",
	"event/FilteredBlock":            "/Channel/Application/Readers",
}

var DefaultPoliciesApplication = map[string]*tc.Policy{
	"Readers": &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
	"Writers": &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
	"Admins":  &tc.Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
}
var DefaultPoliciesOrderer = map[string]*tc.Policy{
	"Readers":         &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
	"Writers":         &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
	"Admins":          &tc.Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
	"BlockValidation": &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
}
var DefaultPoliciesChannel = map[string]*tc.Policy{
	"Readers": &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
	"Writers": &tc.Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
	"Admins":  &tc.Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
}
var DefaultPoliciesOrganization = map[string]*tc.Policy{
	"Readers": &tc.Policy{Type: "Signature", Rule: "OR('SampleOrg.member')"},
	"Writers": &tc.Policy{Type: "Signature", Rule: "OR('SampleOrg.member')"},
	"Admins":  &tc.Policy{Type: "Signature", Rule: "OR('SampleOrg.member')"},
}

var DefaultOrganization = tc.Organization{
	Name:     SampleOrgName,
	ID:       SampleOrgName,
	MSPDir:   "msp",
	MSPType:  msp.ProviderTypeToString(msp.FABRIC),
	Policies: DefaultPoliciesOrganization,
	AnchorPeers: []*tc.AnchorPeer{
		{Host: "127.0.0.1", Port: 7051},
	},
	AdminPrincipal: AdminRoleAdminPrincipal,
}

var DefaultApplication = tc.Application{
	Organizations: nil,
	Capabilities:  DefaultCapabilitiesApplication,
	Resources:     nil,
	Policies:      DefaultPoliciesApplication,
	ACLs:          DefaultACLs,
}

var DefaultOrderer = tc.Orderer{
	OrdererType:  "solo",
	Addresses:    []string{"127.0.0.1:7050"},
	BatchTimeout: time.Second * 2,
	BatchSize: tc.BatchSize{
		MaxMessageCount:   10,
		AbsoluteMaxBytes:  98 * 1024 * 1024,
		PreferredMaxBytes: 512 * 1024,
	},
	Kafka: tc.Kafka{
		Brokers: []string{"kafka0:9092"},
	},
	EtcdRaft:      nil,
	Organizations: nil,
	MaxChannels:   0,
	Capabilities:  DefaultCapabilitiesOrderer,
	Policies:      DefaultPoliciesOrderer,
}

var DefaultConsortium = tc.Consortium{
	Organizations: []*tc.Organization{&DefaultOrganization},
}

var DefaultProfileChannel = tc.Profile{
	Consortium:  SampleConsortiumName,
	Application: &DefaultApplication,
}

var DefaultProfileOrderer = tc.Profile{
	Consortium:  SampleConsortiumName,
	Application: nil,
	Orderer:     &DefaultOrderer,
	Consortiums: map[string]*tc.Consortium{
		SampleConsortiumName: &DefaultConsortium,
	},
	Capabilities: DefaultCapabilitiesOrderer,
	Policies:     DefaultPoliciesOrderer,
}
