package scte224v20180501

import (
	"encoding/xml"
	"strings"
	"testing"
)

const CALI_XML = `<Media xmlns:audience="urn:scte:224:audience" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:action="urn:scte:224:action" xmlns="http://www.scte.org/schemas/224" id="nbcuni.com/media/NBCS_CALIFORNIA/BROADCAST" description="NBCS_CALIFORNIA" lastUpdated="2018-05-29T00:44:57.000Z" source="CSN_WEST">
<MediaPoint id="nbcuni.com/media/NBCS_CALIFORNIA/program/00024MA000000008271T0529180100/BROADCAST/start" description="Behind the Mask" lastUpdated="2018-05-08T00:49:37.000Z" effective="2018-05-29T00:00:00.000Z" expires="2018-05-29T10:00:00.000Z" source="CSN_WEST">
<AltID>SH019968500000</AltID>
<AltID>CID:e0be5016-5ac2-4c36-8cef-fc00ce40372e</AltID>
<Metadata>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Sport" provider="nbc" type="string">Documentary</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Lookback" provider="nbc" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_enabled" provider="nbc" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Startover" provider="nbc" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledStart" provider="nbc" type="string">2018-05-29T01:00:00Z</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledAiringID" provider="nbc" type="string">00024MA000000008271T0529180100</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Embargo" provider="nbc" type="string">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="Prg_dai_devices" provider="nbc" type="bool">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="SeriesEmbargo" provider="nbc" type="string">false</MetadataDetail>
<MetadataDetail xmlns="http://ctsrmm.com/ctsesni" name="ScheduledEnd" provider="nbc" type="string">2018-05-29T01:30:00Z</MetadataDetail>
</Metadata>
<Apply>
<Policy xlink:href="nbcuni.com/policy/210_CA_Main_OOM_Cable_to_CA"/>
</Apply>
<MatchSignal match="ANY">
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()='00024MA000000008271T0529180100']</Assert>
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()='30303032344D413030303030303030383237315430353239313830313030']</Assert>
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()='00024MA000000008271T0529180100']</Assert>
<Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()='30303032344D413030303030303030383237315430353239313830313030']</Assert>
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
	if caliScte.Source != "CSN_WEST" {
		t.Log(caliScte.Source, "should have been \"CSN_WEST\"")
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
