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
    <Apply xmlns="http://www.scte.org/schemas/224" priority="0">
      <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="superflaco.com/policy/SWITCH_to_CA"></Policy>
    </Apply>
    <Apply xmlns="http://www.scte.org/schemas/224">
      <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="superflaco.com/policy/AnotherPolicy"></Policy>
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

const topLevelNamspaces = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Media source="NBCS_BOSTON" id="nbcuni.com/media/NBCS_BOSTON_ALT/BROADCAST" description="NBCS_BOSTON_ALT" lastUpdated="2020-03-19T08:53:41.951Z" xmlns="http://www.scte.org/schemas/224" xmlns:lrm="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" xmlns:xlink="http://www.w3.org/1999/xlink">
  <MediaPoint source="NBCS_BOSTON" order="1" effective="2020-03-18T13:00:00Z" expires="2020-03-18T23:00:00Z" id="nbcuni.com/media/NBCS_BOSTON_ALT/program/00044MA000000037610T0318201400/BROADCAST/start" description="Zolak and Bertrand" lastUpdated="2020-03-19T08:53:41.951Z">
    <AltID>SH029993960000</AltID>
    <Metadata>
      <lrm:MetadataDetail name="ScheduledStart" type="string" provider="nbc">2020-03-18T14:00:00Z</lrm:MetadataDetail>
      <lrm:MetadataDetail name="ScheduledEnd" type="string" provider="nbc">2020-03-18T18:00:00Z</lrm:MetadataDetail>
      <lrm:MetadataDetail name="ScheduledAiringID" type="string" provider="nbc">00044MA000000037610T0318201400</lrm:MetadataDetail>
      <lrm:MetadataDetail name="Sport" type="string" provider="nbc">Sports talk</lrm:MetadataDetail>
    </Metadata>
    <Apply>
      <Policy xlink:href="nbcuni.com/policy/NBCS_BOSTON_ALT/cleared"/>
    </Apply>
    <MatchSignal match="ANY">
      <Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()='00044MA000000037610T0318201400']</Assert>
      <Assert>/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()='00044MA000000037610T0318201400']</Assert>
    </MatchSignal>
  </MediaPoint>
</Media>`

const inlinedNamespaces = `<Media xmlns="http://www.scte.org/schemas/224" id="nbcuni.com/media/NBCS_BOSTON_ALT/BROADCAST" description="NBCS_BOSTON_ALT" lastUpdated="2020-03-19T08:53:41.951Z" source="NBCS_BOSTON">
  <MediaPoint xmlns="http://www.scte.org/schemas/224" id="nbcuni.com/media/NBCS_BOSTON_ALT/program/00044MA000000037610T0318201400/BROADCAST/start" description="Zolak and Bertrand" lastUpdated="2020-03-19T08:53:41.951Z" effective="2020-03-18T13:00:00Z" expires="2020-03-18T23:00:00Z" source="NBCS_BOSTON" order="1">
    <AltID xmlns="http://www.scte.org/schemas/224">SH029993960000</AltID>
    <Metadata xmlns="http://www.scte.org/schemas/224">
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="ScheduledStart" type="string" provider="nbc">2020-03-18T14:00:00Z</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="ScheduledEnd" type="string" provider="nbc">2020-03-18T18:00:00Z</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="ScheduledAiringID" type="string" provider="nbc">00044MA000000037610T0318201400</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="Sport" type="string" provider="nbc">Sports talk</MetadataDetail>
    </Metadata>
    <Apply xmlns="http://www.scte.org/schemas/224">
      <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="nbcuni.com/policy/NBCS_BOSTON_ALT/cleared"></Policy>
    </Apply>
    <MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY">
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;00044MA000000037610T0318201400&#39;]</Assert>
      <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and text()=&#39;00044MA000000037610T0318201400&#39;]</Assert>
    </MatchSignal>
  </MediaPoint>
</Media>`

func TestMetadataNamespace(t *testing.T) {
	decoder := xml.NewDecoder(strings.NewReader(topLevelNamspaces))
	var metaScte Media
	decodeErr := decoder.Decode(&metaScte)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}
	if metaScte.Source != "NBCS_BOSTON" {
		t.Log(metaScte.Source, "should have been \"NBCS_BOSTON\"")
		t.Fail()
	}
	if 1 != len(metaScte.MediaPoints) {
		t.Log("Expected just one MediaPoint")
		t.Fail()
	}

	pretty, marshalErr := xml.MarshalIndent(metaScte, "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	roundtrip := string(pretty)
	if inlinedNamespaces != roundtrip {
		t.Log(roundtrip)
		t.Log("did not match")
		t.Log(inlinedNamespaces)
		t.Fail()
	}
}

const originalADIMetadata = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Media source="HBOHD" id="hbo.com/media/HBOHD" description="HBO HD" lastUpdated="2020-03-31T18:00:20.5Z" xmlns="http://www.scte.org/schemas/224" xmlns:signaling="http://www.scte.org/schemas/236/2017/signaling" xmlns:ad="http://www.scte.org/schemas/236/2017/ad" xmlns:lrm="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" xmlns:title="http://www.scte.org/schemas/236/2017/title" xmlns:content="http://www.scte.org/schemas/236/2017/content" xmlns:offer="http://www.scte.org/schemas/236/2017/offer" xmlns:core="http://www.scte.org/schemas/236/2017/core" xmlns:terms="http://www.scte.org/schemas/236/2017/terms" xmlns:ns9="http://www.scte.org/schemas/35/2017" xmlns:xlink="http://www.w3.org/1999/xlink">
  <MediaPoint matchTime="2020-03-31T01:00:00Z" source="HBOHD" order="0" effective="2020-03-31T00:55:00Z" expires="2020-03-31T01:05:00Z" id="hbo.com/media/HBOHD/program/f2MhzKZcfiosBNjDM2hHGK0F8_Jv5fUR/start" description="The Plot Against America" lastUpdated="2020-03-31T18:00:20.5Z">
    <AltID>ow7qaht5qULfWFcl5sbvzJXOVJEHsP6a</AltID>
    <Metadata>
      <core:ADI3>
        <core:Asset xsi:type="content:MovieType" providerVersionNum="20" internalVersionNum="20" creationDateTime="2020-03-30T00:00:00" startDateTime="2020-03-30T18:00:00.000-07:00" endDateTime="2020-03-30T19:01:00.000-07:00" lastModifiedDateTime="2020-03-24T23:21:02.000-07:00" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <core:Provider>HBO</core:Provider>
          <content:Language bitStreamMode="2">eng</content:Language>
        </core:Asset>
        <core:Asset xsi:type="title:TitleType" uriId="hbo.com/title/HBO" providerVersionNum="20" internalVersionNum="20" creationDateTime="2020-03-30T00:00:00" startDateTime="2020-03-30T18:00:00.000-07:00" endDateTime="2020-03-30T19:01:00.000-07:00" lastModifiedDateTime="2020-03-24T23:21:02.000-07:00" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
          <core:Provider>HBO</core:Provider>
          <core:Ext>
            <core:App_Data name="Season_Number" value="1"/>
            <core:App_Data name="Season_Episode_Number" value="3"/>
            <core:App_Data name="Series_Episode_Number" value="3"/>
            <core:App_Data name="Category" value="Other"/>
            <core:App_Data name="Genre" value="Drama"/>
          </core:Ext>
          <title:LocalizableTitle>
            <title:TitleBrief>Plot Against America</title:TitleBrief>
            <title:TitleMedium>The Plot Against America</title:TitleMedium>
            <title:TitleLong>The Plot Against America</title:TitleLong>
            <title:SummaryShort>Based on Philip Roth's novel of an alternate American history during World War II.</title:SummaryShort>
          </title:LocalizableTitle>
          <title:Rating ratingSystem="urn:v-chip">TVMA</title:Rating>
          <title:IsClosedCaptioning>true</title:IsClosedCaptioning>
          <title:DisplayRunTime>60</title:DisplayRunTime>
          <title:Year>2020</title:Year>
        </core:Asset>
      </core:ADI3>
      <lrm:MetadataDetail name="ScheduledStart" type="timestamp" provider="HBO">2020-03-31T01:00:00Z</lrm:MetadataDetail>
      <lrm:MetadataDetail name="ScheduledEnd" type="timestamp" provider="HBO">2020-03-31T02:01:00Z</lrm:MetadataDetail>
      <lrm:MetadataDetail name="StartOver" type="Boolean" provider="HBO">false</lrm:MetadataDetail>
      <lrm:MetadataDetail name="StartOverDuration" type="String" provider="HBO"></lrm:MetadataDetail>
      <lrm:MetadataDetail name="Lookback" type="Boolean" provider="HBO">false</lrm:MetadataDetail>
    </Metadata>
  </MediaPoint>
</Media>`

const expectedInlinedADIMetadata = `<Media xmlns="http://www.scte.org/schemas/224" id="hbo.com/media/HBOHD" description="HBO HD" lastUpdated="2020-03-31T18:00:20.5Z" source="HBOHD">
  <MediaPoint xmlns="http://www.scte.org/schemas/224" id="hbo.com/media/HBOHD/program/f2MhzKZcfiosBNjDM2hHGK0F8_Jv5fUR/start" description="The Plot Against America" lastUpdated="2020-03-31T18:00:20.5Z" effective="2020-03-31T00:55:00Z" expires="2020-03-31T01:05:00Z" matchTime="2020-03-31T01:00:00Z" source="HBOHD" order="0">
    <AltID xmlns="http://www.scte.org/schemas/224">ow7qaht5qULfWFcl5sbvzJXOVJEHsP6a</AltID>
    <Metadata xmlns="http://www.scte.org/schemas/224">
      <ADI3 xmlns="http://www.scte.org/schemas/236/2017/core">
        <Asset xmlns="http://www.scte.org/schemas/236/2017/core" xmlns:XMLSchema-instance="http://www.w3.org/2001/XMLSchema-instance" XMLSchema-instance:type="content:MovieType" uriId="" providerVersionNum="20" internalVersionNum="20" creationDateTime="2020-03-30T00:00:00" startDateTime="2020-03-30T18:00:00.000-07:00" endDateTime="2020-03-30T19:01:00.000-07:00" lastModifiedDateTime="2020-03-24T23:21:02.000-07:00">
          <Provider xmlns="http://www.scte.org/schemas/236/2017/core">HBO</Provider>
          <Language xmlns="http://www.scte.org/schemas/236/2017/content" bitStreamMode="2">eng</Language>
        </Asset>
        <Asset xmlns="http://www.scte.org/schemas/236/2017/core" xmlns:XMLSchema-instance="http://www.w3.org/2001/XMLSchema-instance" XMLSchema-instance:type="title:TitleType" uriId="hbo.com/title/HBO" providerVersionNum="20" internalVersionNum="20" creationDateTime="2020-03-30T00:00:00" startDateTime="2020-03-30T18:00:00.000-07:00" endDateTime="2020-03-30T19:01:00.000-07:00" lastModifiedDateTime="2020-03-24T23:21:02.000-07:00">
          <Provider xmlns="http://www.scte.org/schemas/236/2017/core">HBO</Provider>
          <Ext xmlns="http://www.scte.org/schemas/236/2017/core">
            <App_Data name="Season_Number" value="1"></App_Data>
            <App_Data name="Season_Episode_Number" value="3"></App_Data>
            <App_Data name="Series_Episode_Number" value="3"></App_Data>
            <App_Data name="Category" value="Other"></App_Data>
            <App_Data name="Genre" value="Drama"></App_Data>
          </Ext>
          <LocalizableTitle xmlns="http://www.scte.org/schemas/236/2017/title">
            <TitleBrief xmlns="http://www.scte.org/schemas/236/2017/title">Plot Against America</TitleBrief>
            <TitleMedium xmlns="http://www.scte.org/schemas/236/2017/title">The Plot Against America</TitleMedium>
            <TitleLong xmlns="http://www.scte.org/schemas/236/2017/title">The Plot Against America</TitleLong>
            <SummaryShort xmlns="http://www.scte.org/schemas/236/2017/title">Based on Philip Roth&#39;s novel of an alternate American history during World War II.</SummaryShort>
          </LocalizableTitle>
          <Rating xmlns="http://www.scte.org/schemas/236/2017/title" ratingSystem="urn:v-chip">TVMA</Rating>
          <IsClosedCaptioning xmlns="http://www.scte.org/schemas/236/2017/title">true</IsClosedCaptioning>
          <DisplayRunTime xmlns="http://www.scte.org/schemas/236/2017/title">60</DisplayRunTime>
          <Year xmlns="http://www.scte.org/schemas/236/2017/title">2020</Year>
        </Asset>
      </ADI3>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="ScheduledStart" type="timestamp" provider="HBO">2020-03-31T01:00:00Z</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="ScheduledEnd" type="timestamp" provider="HBO">2020-03-31T02:01:00Z</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="StartOver" type="Boolean" provider="HBO">false</MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="StartOverDuration" type="String" provider="HBO"></MetadataDetail>
      <MetadataDetail xmlns="https://lrm.aor.theplatform.com/lrm/schemas/esni/224" name="Lookback" type="Boolean" provider="HBO">false</MetadataDetail>
    </Metadata>
  </MediaPoint>
</Media>`

func TestADIMetadata(t *testing.T) {
	decoder := xml.NewDecoder(strings.NewReader(originalADIMetadata))
	var metaScte Media
	decodeErr := decoder.Decode(&metaScte)
	if nil != decodeErr {
		t.Log(decodeErr)
		t.FailNow()
	}
	if metaScte.Source != "HBOHD" {
		t.Log(metaScte.Source, "should have been \"HBOHD\"")
		t.Fail()
	}
	if 1 != len(metaScte.MediaPoints) {
		t.Log("Expected just one MediaPoint")
		t.Fail()
	}

	pretty, marshalErr := xml.MarshalIndent(metaScte, "", "  ")
	if nil != marshalErr {
		t.Log(marshalErr)
		t.FailNow()
	}
	roundtrip := string(pretty)
	if expectedInlinedADIMetadata != roundtrip {
		t.Log(roundtrip)
		t.Log("did not match")
		t.Log(expectedInlinedADIMetadata)
		t.Fail()
	}
}
