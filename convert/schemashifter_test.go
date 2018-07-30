package convert

import (
	"encoding/xml"
	"strings"
	"testing"

	scte224_2015 "github.comcast.com/jcolwe200/scte224/types/scte224v20151115"
)

const CALI_2015_XML = `<Media xmlns="http://www.scte.org/schemas/224/2015" id="superflaco.com/media/CALIFORNIA/BROADCAST" description="CALIFORNIA" lastUpdated="2018-05-29T00:44:57Z" source="CURRENT_CHANNEL">
  <MediaPoint xmlns="http://www.scte.org/schemas/224/2015" id="superflaco.com/media/CALIFORNIA/program/AIRING_ID/BROADCAST/start" description="I am a MediaPoint" lastUpdated="2018-05-08T00:49:37Z" effective="2018-05-29T00:00:00Z" expires="2018-05-29T10:00:00Z" source="CURRENT_CHANNEL">
    <AltID xmlns="http://www.scte.org/schemas/224/2015">CID:e0be5016-5ac2-4c36-8cef-fc00ce40372e</AltID>
    <Metadata xmlns="http://www.scte.org/schemas/224/2015">
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Sport" type="string" provider="superflaco">Documentary</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Lookback" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_enabled" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Startover" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledStart" type="string" provider="superflaco">2018-05-29T01:00:00Z</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledAiringID" type="string" provider="superflaco">AIRING_ID</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Embargo" type="string" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_devices" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="SeriesEmbargo" type="string" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledEnd" type="string" provider="superflaco">2018-05-29T01:30:00Z</MetadataDetail>
    </Metadata>
    <Apply xmlns="http://www.scte.org/schemas/224/2015">
      <Policy xmlns="http://www.scte.org/schemas/224/2015" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="superflaco.com/policy/SWITCH_to_CA"></Policy>
    </Apply>
    <MatchSignal xmlns="http://www.scte.org/schemas/224/2015" match="ANY">
      <Assert xmlns="http://www.scte.org/schemas/224/2015">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;AIRING_ID&#39;]</Assert>
      <Assert xmlns="http://www.scte.org/schemas/224/2015">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;AIRING_ID&#39;]</Assert>
    </MatchSignal>
  </MediaPoint>
</Media>`

const CALI_2018_XML = `<Media xmlns="http://www.scte.org/schemas/224" id="superflaco.com/media/CALIFORNIA/BROADCAST" description="CALIFORNIA" lastUpdated="2018-05-29T00:44:57Z" source="CURRENT_CHANNEL">
  <MediaPoint xmlns="http://www.scte.org/schemas/224" id="superflaco.com/media/CALIFORNIA/program/AIRING_ID/BROADCAST/start" description="I am a MediaPoint" lastUpdated="2018-05-08T00:49:37Z" effective="2018-05-29T00:00:00Z" expires="2018-05-29T10:00:00Z" source="CURRENT_CHANNEL">
    <AltID xmlns="http://www.scte.org/schemas/224">CID:e0be5016-5ac2-4c36-8cef-fc00ce40372e</AltID>
    <Metadata xmlns="http://www.scte.org/schemas/224">
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Sport" type="string" provider="superflaco">Documentary</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Lookback" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_enabled" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Startover" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledStart" type="string" provider="superflaco">2018-05-29T01:00:00Z</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledAiringID" type="string" provider="superflaco">AIRING_ID</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Embargo" type="string" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_devices" type="bool" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="SeriesEmbargo" type="string" provider="superflaco">false</MetadataDetail>
      <MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledEnd" type="string" provider="superflaco">2018-05-29T01:30:00Z</MetadataDetail>
    </Metadata>
    <Apply xmlns="http://www.scte.org/schemas/224">
      <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="superflaco.com/policy/SWITCH_to_CA"></Policy>
    </Apply>
    <MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY">
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;AIRING_ID&#39;]</Assert>
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;AIRING_ID&#39;]</Assert>
    </MatchSignal>
  </MediaPoint>
</Media>`

const thoseGuys2015 = `<Audience xmlns="http://www.scte.org/schemas/224/2015" id="superflaco.com/audience/ThoseGuys" description="ThoseGuys" lastUpdated="2018-07-17T17:14:32.359Z" match="ANY">
  <Vird xmlns="urn:scte:224:audience">ThoseGuys</Vird>
</Audience>`

const thoseGuys2018 = `<Audience xmlns="http://www.scte.org/schemas/224" id="superflaco.com/audience/ThoseGuys" description="ThoseGuys" lastUpdated="2018-07-17T17:14:32.359Z" match="ANY">
  <Vird xmlns="urn:scte:224:audience">ThoseGuys</Vird>
</Audience>`

func TestUpgradeAudience(t *testing.T) {

	decoder := xml.NewDecoder(strings.NewReader(thoseGuys2015))
	var cmcScte2015 scte224_2015.Audience
	decodeErr := decoder.Decode(&cmcScte2015)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}

	pretty, marshalErr := xml.MarshalIndent(UpgradeAudience(cmcScte2015), "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	upgraded := string(pretty)
	if thoseGuys2018 != upgraded {
		t.Log(upgraded)
		t.Log("did not match")
		t.Log(thoseGuys2018)
		t.Fail()
	}
}

func TestUpgradeMedia(t *testing.T) {

	decoder := xml.NewDecoder(strings.NewReader(CALI_2015_XML))
	var caliScte2015 scte224_2015.Media
	decodeErr := decoder.Decode(&caliScte2015)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}

	pretty, marshalErr := xml.MarshalIndent(UpgradeMedia(caliScte2015), "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	upgraded := string(pretty)
	if CALI_2018_XML != upgraded {
		t.Log(upgraded)
		t.Log("did not match")
		t.Log(CALI_2018_XML)
		t.Fail()
	}
}
