package fs

import (
	"fmt"
	"strconv"

	provModel "github.com/afzalabbasi/testrepo/api_call"
	"github.com/callmylist/callmylist-provapi/model"
	"github.com/callmylist/callmylist-provapi/system/cmlconstants"
	//"github.com/callmylist/callmylist-provapi/system/keys"
	"github.com/callmylist/callmylist-provapi/utils/cmlutils"
)

const TestCampaignName = "TEST"

// private method
// this defines call plan string to be passed to freeswitch
func campaignDialString(carrierToken string, campaign provModel.Call, callId string) string {
	// sip relates
	originateTimeoutValue := "originate_timeout=45"
	//domainValue := "sip_invite_domain=" + domain
	//carrierTokenValue := "sip_h_X-Telnyx-Token=" + carrierToken
	ignoreEarlyMedia := "ignore_early_media=true"
	systemValues := originateTimeoutValue + "," + ignoreEarlyMedia
	callIdKey := "origination_uuid"
	callerIdKey := "origination_caller_id_number" //from _did
	callerIdName := "origination_caller_id_name"  //from_did
	tx_DID := "tx_DID"
	campaignTypeKey := "campaign_type" //
	if isTest {
		campaignId = TestCampaignName
	}
	// sound keys
	human_Audio := "soundfile_id"
	// dndkeys
	dndKey := "dnd_DTMF"
	// vm keys
	vmDetect := "vm_detect"
	vmDrop := "vm_drop"
	vmAudio := "vm_audio"
	// transfer keys
	transferNumberKey := "transfer_DTMF"
	//
	return fmt.Sprintf("{%s,%s=%s,%s=%s,%s=%s,%s=%s,%s=%s,%s=%s,%s=%s,%s=%s,%s=%s}", systemValues,
		callIdKey, callId,
		callerIdKey, campaign.From_DID,
		callerIdName, campaign.From_DID,
		tx_DID, campaign.Tx_DID,
		human_Audio, campaign.Human_audio,
		dndKey, campaign.DND_DTMF,
		vmDetect, campaign.Vm_detect,
		vmDrop, campaign.Vm_drop,
		vmAudio, campaign.Vm_audio,
		transferNumberKey, campaign.Tx_DTMF,
		campaignTypeKey, campaign.Compaign_type,
	)

}
