package localconfig

import (
	"io/ioutil"
	"testing"

	"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	"github.com/hyperledger/fabric/protos/utils"
)

func TestGenesisBlock(t *testing.T) {
	profile := GetOrdererProfileTemplate()
	t.Log(profile)
	coder := encoder.New(profile)
	block := coder.GenesisBlock()
	blockBytes := utils.MarshalOrPanic(block)
	err := ioutil.WriteFile("test.block", blockBytes, 0666)
	if err != nil {
		t.Error(err)
	}
}

func TestChannelBlock(t *testing.T) {
	profile := GetChannelProfileTemplate()
	t.Log(profile)
	block, err := encoder.MakeChannelCreationTransaction("bclc", nil, profile)
	if err != nil {
		t.Error(err)
	}
	blockBytes := utils.MarshalOrPanic(block)
	err = ioutil.WriteFile("test.channel", blockBytes, 0666)
	if err != nil {
		t.Error(err)
	}
}
