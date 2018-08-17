package convert

import (
	"encoding/xml"
	"strings"
	"testing"

	scte224 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

func TestDowngradeAudience(t *testing.T) {

	decoder := xml.NewDecoder(strings.NewReader(thoseGuys2018))
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
	if thoseGuys2015 != downgraded {
		t.Log(downgraded)
		t.Log("did not match")
		t.Log(thoseGuys2015)
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
