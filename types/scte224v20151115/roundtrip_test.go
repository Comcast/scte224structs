package scte224v20151115

import (
	"encoding/xml"
	"strings"
	"testing"
)

const CALI_XML = `<Media xmlns:audience="urn:scte:224:audience" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:action="urn:scte:224:action" xmlns="http://www.scte.org/schemas/224/2015" id="superflaco.com/media/CALIFORNIA/BROADCAST" description="CALIFORNIA" lastUpdated="2018-05-29T00:44:57.000Z" source="CURRENT_CHANNEL">
<MediaPoint id="superflaco.com/media/CALIFORNIA/program/AIRING_ID/BROADCAST/start" description="I am a MediaPoint" lastUpdated="2018-05-08T00:49:37.000Z" effective="2018-05-29T00:00:00.000Z" expires="2018-05-29T10:00:00.000Z" source="CURRENT_CHANNEL">
<AltID>CID:e0be5016-5ac2-4c36-8cef-fc00ce40372e</AltID>
<Metadata>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Sport" provider="superflaco" type="string">Documentary</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Lookback" provider="superflaco" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_enabled" provider="superflaco" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Startover" provider="superflaco" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledStart" provider="superflaco" type="string">2018-05-29T01:00:00Z</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledAiringID" provider="superflaco" type="string">AIRING_ID</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Embargo" provider="superflaco" type="string">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_devices" provider="superflaco" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="SeriesEmbargo" provider="superflaco" type="string">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledEnd" provider="superflaco" type="string">2018-05-29T01:30:00Z</MetadataDetail>
</Metadata>
<Apply>
<Policy xlink:href="superflaco.com/policy/SWITCH_to_CA"/>
</Apply>
<MatchSignal match="ANY">
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()='AIRING_ID']</Assert>
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()='AIRING_ID']</Assert>
</MatchSignal>
</MediaPoint></Media>`

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
	t.Log(string(pretty))
}
