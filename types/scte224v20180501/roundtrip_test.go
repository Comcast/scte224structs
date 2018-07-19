package scte224v20180501

import (
	"encoding/xml"
	"strings"
	"testing"
)

const CALI_XML = `<Media xmlns="http://www.scte.org/schemas/224" id="superflaco.com/media/CALIFORNIA/BROADCAST" description="CALIFORNIA" lastUpdated="2018-05-29T00:44:57Z" source="CURRENT_CHANNEL">
  <MediaPoint xmlns="http://www.scte.org/schemas/224" id="superflaco.com/media/AIRING_ID/program/PROGRAM_IDENTIFIER/start" description="I am a MediaPoint" lastUpdated="2018-05-08T00:49:37Z" effective="2018-05-29T00:00:00Z" expires="2018-05-29T10:00:00Z" source="CURRENT_CHANNEL">
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
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;SIGNAL_UPID&#39;]</Assert>
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;SIGNAL_UPID&#39;]</Assert>
    </MatchSignal>
  </MediaPoint>
</Media>`

func TestRoundtrip(t *testing.T) {

	decoder := xml.NewDecoder(strings.NewReader(CALI_XML))
	var caliScte Media
	decodeErr := decoder.Decode(&caliScte)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}
	if caliScte.Source != "CURRENT_CHANNEL" {
		t.Log(caliScte.Source, "should have been \"CURRENT_CHANNEL\"")
		t.Fail()
	}
	if 1 != len(caliScte.MediaPoints) {
		t.Log("Expected just one MediaPoint")
		t.Fail()
	}

	pretty, marshalErr := xml.MarshalIndent(caliScte, "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	roundtrip := string(pretty)
	if CALI_XML != roundtrip {
		t.Log(roundtrip)
		t.Log("did not match")
		t.Log(CALI_XML)
		t.Fail()
	}
}
