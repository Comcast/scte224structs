package scte224v20200407

const vp2020Raw string = `<ViewingPolicy xmlns="http://www.scte.org/schemas/224" id="test/program" description="test program" lastUpdated="2021-01-19T18:49:26.298986528Z">
  <Content xmlns="urn:scte:224:action">CONTENT</Content>
  <Allocation xmlns="urn:scte:224:action" ownerType="PROVIDER" ownerName="Stu" duration="PT30S" ads="FreeWheel/MRM">
  	<Slots>
  		<Slot duration="PT15S" offset="PT0S">        
	  		<AdsReferenceId>98765</AdsReferenceId>        
	  		<AdsReferenceId>87654</AdsReferenceId>        
	  		<AdsReferenceId>54321</AdsReferenceId>        
	  		<AdsReferenceId>56343</AdsReferenceId>     
            <SlotRules>
		      <SlotRule rule="exclusionWithInference">
		        <Parameter parameterName="advertiser">advertiser_external_id</Parameter>
		      </SlotRule>
            </SlotRules>
  		</Slot>      
  		<Slot duration="PT30S" offset="PT15S">
       		<AdsReferenceId>123</AdsReferenceId>      
  		</Slot>      		     
  		<Slot duration="PT15S" offset="PT45S"/>   		  
	</Slots>
  </Allocation>
</ViewingPolicy>`

const vpSignalPointInsertion_w_SpliceInfoSection string = `<ViewingPolicy xmlns="http://www.scte.org/schemas/224" id="test.com/viewingpolicy/test/testVP1" description="Test VP with SignalPointInsertion containing a SpliceInfoSection" lastUpdated="2022-04-27T22:09:55Z">
  <Audience xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="test.com/audience/test/testaudience"></Audience>
  <SignalPointInsertion xmlns="urn:scte:224:action" offset="PT1M">
    <SpliceInfoSection xmlns="http://www.scte.org/schemas/35" offset="PT1M">
      <TimeSignal xmlns="http://www.scte.org/schemas/35">
        <SpliceTime xmlns="http://www.scte.org/schemas/35" ptsTime="0"></SpliceTime>
      </TimeSignal>
      <SegmentationDescriptor xmlns="http://www.scte.org/schemas/35" segmentationTypeId="55" segmentsExpected="3" segmentNum="2" segmentationDuration="5400000">
        <SegmentationUpid xmlns="http://www.scte.org/schemas/35" segmentationUpidType="9">SIGNAL:%Base64GUID%</SegmentationUpid>
      </SegmentationDescriptor>
    </SpliceInfoSection>
  </SignalPointInsertion>
  <SignalPointReplacement xmlns="urn:scte:224:action">
    <SpliceInfoSection xmlns="http://www.scte.org/schemas/35" offset="PT1M">
      <TimeSignal xmlns="http://www.scte.org/schemas/35">
        <SpliceTime xmlns="http://www.scte.org/schemas/35" ptsTime="0"></SpliceTime>
      </TimeSignal>
      <SegmentationDescriptor xmlns="http://www.scte.org/schemas/35" segmentationTypeId="54" segmentsExpected="3" segmentNum="2" segmentationDuration="5400000">
        <SegmentationUpid xmlns="http://www.scte.org/schemas/35" segmentationUpidType="9">SIGNAL:%Base64GUID%</SegmentationUpid>
        <SegmentationUpid xmlns="http://www.scte.org/schemas/35" segmentationUpidType="15">urn:comcast:altcon:airdate:1651078800</SegmentationUpid>
      </SegmentationDescriptor>
    </SpliceInfoSection>
  </SignalPointReplacement>
</ViewingPolicy>`

// Same as "vp2020Raw" but without the additional "Allocation" action
const vp2018Raw string = `<ViewingPolicy xmlns="http://www.scte.org/schemas/224" id="test/program" description="test program" lastUpdated="2021-01-19T18:49:26.298986528Z">
    <Content xmlns="urn:scte:224:action">CONTENT</Content>
</ViewingPolicy>`

const aud2020Raw string = `<Audience xmlns="http://www.scte.org/schemas/224" id="foo/audience/any.all" description="Users anywhere on any device" lastUpdated="2021-03-24T03:27:05Z" match="ALL">
    <Audience xmlns="http://www.scte.org/schemas/224" description="location" lastUpdated="2021-03-23T08:16:19.235Z" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="foo/audience/location/any"></Audience>
    <Audience xmlns="http://www.scte.org/schemas/224" description="device" lastUpdated="2021-03-23T08:16:19.235Z" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="foo/audience/device/all"></Audience>
</Audience>`

const media2020Raw = `<Media xmlns="http://www.scte.org/schemas/224" id="test/media/" description="test" lastUpdated="2021-07-27T01:13:25.849Z" source="TEST">
    <MediaPoint xmlns="http://www.scte.org/schemas/224" id="testcontentowner.com/media/test_channel_EAST/resident/testcontentowner" description="Resident MediaPoint for setting channel source" lastUpdated="2021-07-27T01:13:25.839Z" order="0">
        <Apply xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/resident"></Policy>
        </Apply>
    </MediaPoint>
    <MediaPoint xmlns="http://www.scte.org/schemas/224" id="testcontentowner.com/media/test_channel_EAST/program/0000000/testcontentowner/start" description="new show" lastUpdated="2021-07-27T01:13:25.839Z" effective="2021-07-26T08:58:00Z" expires="2021-07-26T09:02:00Z" matchTime="2021-07-26T09:00:00Z" source="test_channel_EAST" order="1">
        <AltID xmlns="http://www.scte.org/schemas/224" description="programid" type="CallSign">12345</AltID>
        <AltID xmlns="http://www.scte.org/schemas/224" description="seriesid" type="CallSign">67890</AltID>
        <Apply xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/cleared"></Policy>
        </Apply>
        <MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY" schema="http://www.scte.org/schemas/35">
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;12345&#39;)]</Assert>
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;34567&#39;)]</Assert>
        </MatchSignal>
    </MediaPoint>
</Media>`

// Same as "media2020Raw" but without these additional elements
// AltID.type attribute
// MatchSignal.schema attribute
const media2018Raw = `<Media xmlns="http://www.scte.org/schemas/224" id="test/media/" description="test" lastUpdated="2021-07-27T01:13:25.849Z" source="TEST">
    <MediaPoint xmlns="http://www.scte.org/schemas/224" id="testcontentowner.com/media/test_channel_EAST/resident/testcontentowner" description="Resident MediaPoint for setting channel source" lastUpdated="2021-07-27T01:13:25.839Z" order="0">
        <Apply xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/resident"></Policy>
        </Apply>
    </MediaPoint>
    <MediaPoint xmlns="http://www.scte.org/schemas/224" id="testcontentowner.com/media/test_channel_EAST/program/0000000/testcontentowner/start" description="new show" lastUpdated="2021-07-27T01:13:25.839Z" effective="2021-07-26T08:58:00Z" expires="2021-07-26T09:02:00Z" matchTime="2021-07-26T09:00:00Z" source="test_channel_EAST" order="1">
        <AltID xmlns="http://www.scte.org/schemas/224" description="programid">12345</AltID>
        <AltID xmlns="http://www.scte.org/schemas/224" description="seriesid">67890</AltID>
        <Apply xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/cleared"></Policy>
        </Apply>
        <MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY">
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;12345&#39;)]</Assert>
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;34567&#39;)]</Assert>
        </MatchSignal>
    </MediaPoint>
</Media>`

const anotherMedia2020Raw = `<Media xmlns="http://www.scte.org/schemas/224" id="test/media/" description="test" lastUpdated="2021-07-27T01:13:25.849Z" source="TEST">
    <MediaPoint xmlns="http://www.scte.org/schemas/224" id="testcontentowner.com/media/test_channel_EAST/program/0000000/testcontentowner/start" description="new show" lastUpdated="2021-07-27T01:13:25.839Z" effective="2021-07-26T08:58:00Z" expires="2021-07-26T09:02:00Z" matchTime="2021-07-26T09:00:00Z" source="test_channel_EAST" order="1">
        <Apply xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/cleared"></Policy>
        </Apply>
        <Remove xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/cleared"/>
        </Remove>
        <Remove xmlns="http://www.scte.org/schemas/224">
            <Policy xmlns="http://www.scte.org/schemas/224" xmlns:xlink="http://www.w3.org/1999/xlink" xlink:href="testcontentowner.com/policy/test_channel_EAST/cleared"/>
        </Remove>
        <MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY" schema="http://www.scte.org/schemas/35">
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;12345&#39;)]</Assert>
            <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),&#39;34567&#39;)]</Assert>
        </MatchSignal>
    </MediaPoint>
</Media>`
