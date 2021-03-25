package convert

import (
	"encoding/xml"
	"strings"
	"testing"

	scte224 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

func TestDowngradeAudience(t *testing.T) {
	testAudienceDowngrade(t, thoseGuys2018, thoseGuys2015)
	testAudienceDowngrade(t, nestedAudience2018, nestedAudience2015)

}

func testAudienceDowngrade(t *testing.T, new, old string) {

	decoder := xml.NewDecoder(strings.NewReader(new))
	var cmcScte2018 scte224.Audience
	decodeErr := decoder.Decode(&cmcScte2018)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}

	pretty, marshalErr := xml.MarshalIndent(DowngradeAudience(cmcScte2018), "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	downgraded := string(pretty)
	if old != downgraded {
		t.Log(downgraded)
		t.Log("did not match")
		t.Log(old)
		t.Fail()
	}
}

func TestDowngradeMedia(t *testing.T) {

	decoder := xml.NewDecoder(strings.NewReader(CALI_2018_XML))
	var caliScte2018 scte224.Media
	decodeErr := decoder.Decode(&caliScte2018)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}

	pretty, marshalErr := xml.MarshalIndent(DowngradeMedia(caliScte2018), "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	downgraded := string(pretty)
	if CALI_2015_XML != downgraded {
		t.Log(downgraded)
		t.Log("did not match")
		t.Log(CALI_2015_XML)
		t.Fail()
	}
}

const expected_viewingpolicy2015 = `<ViewingPolicy xmlns="http://www.scte.org/schemas/224/2015" id="test.com/viewingpolicy/viewingpolicyXYZ" description="XYZ viewingpolicy" lastUpdated="2018-07-17T17:14:32.359Z">
  <Content xmlns="urn:scte:224:action">Slate</Content>
</ViewingPolicy>`

const expected_viewingpolicy2018 = `<ViewingPolicy xmlns="http://www.scte.org/schemas/224" id="test.com/viewingpolicy/viewingpolicyXYZ" description="XYZ viewingpolicy" lastUpdated="2018-07-17T17:14:32.359Z">
  <Content xmlns="urn:scte:224:action">Slate</Content>
</ViewingPolicy>`

func TestDowngradeViewingPolicy(t *testing.T) {
	decoder := xml.NewDecoder(strings.NewReader(expected_viewingpolicy2018))
	var vp2018 scte224.ViewingPolicy
	decodeErr := decoder.Decode(&vp2018)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}

	pretty, marshalErr := xml.MarshalIndent(DowngradeViewingPolicy(vp2018), "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	downgraded := string(pretty)
	if expected_viewingpolicy2015 != downgraded {
		t.Log(downgraded)
		t.Log("did not match")
		t.Log(expected_viewingpolicy2015)
		t.Fail()
	}
}
