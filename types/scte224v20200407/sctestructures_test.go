package scte224v20200407

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

const viewingpolicy string = `
<ViewingPolicy xmlns="http://www.scte.org/schemas/224" id="test/program" description="test program" lastUpdated="2021-01-19T18:49:26.298986528Z">
  <Content xmlns="urn:scte:224:action">CONTENT</Content>
  <Allocation xmlns="urn:scte:224:action">
  	<Slots>
  		<Slot duration="PT15S" offset="PT0S">        
	  		<AdsReferenceId>98765</AdsReferenceId>        
	  		<AdsReferenceId>87654</AdsReferenceId>        
	  		<AdsReferenceId>54321</AdsReferenceId>        
	  		<AdsReferenceId>56343</AdsReferenceId>      
  		</Slot>      
  		<Slot duration="PT30S" offset="PT15S">
       		<AdsReferenceId>123</AdsReferenceId>      
  		</Slot>      		     
  		<Slot duration="PT15S" offset="PT45S"/>   		  
	</Slots>
  </Allocation>
</ViewingPolicy>`

const matchSignalXML = `
<MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY">
    <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),'12345')]</Assert>
    <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),'34567')]</Assert>
  </MatchSignal>`

func TestViewingPolicy2020(t *testing.T) {
	var vpol *ViewingPolicy
	err := xml.Unmarshal([]byte(viewingpolicy), &vpol)
	assert.Nil(t, err, "Error unmarshalling viewingpolicy")

	_, err = xml.Marshal(vpol)
	assert.Nil(t, err, "Error marshalling viewingpolicy")

	// Testing "Allocation" element unmarshalling which is a new addition in 2020/21 version
	defer func() {
		// saves us from having to check that a pointer is not nil each time
		// there is a lot of pointer dereferencing to do with Allocation element
		// we want to know that the test failed, any panicking here in this test does indicate a failure in unmarshalling
		if r := recover(); r != nil {
			t.Errorf("One or more assertions panicked with %v, indicates test failures", r)
		}
	}()

	firstSlot := vpol.Allocation.Slots[0]
	assert.NotNil(t, firstSlot)

	adSlots := firstSlot.AdSlots
	assert.NotNil(t, adSlots)
	assert.Equalf(t, 3, len(adSlots), "Expected 3 but got %d \n", len(adSlots))

	firstAd := adSlots[0].AdsReferenceId[0]
	assert.Equalf(t, "98765", (*firstAd).ID, "Expected 98765 but got %s \n", (*firstAd).ID)
}

func TestMatchSignal2020(t *testing.T) {
	var matchSignal *MatchSignal
	err := xml.Unmarshal([]byte(matchSignalXML), &matchSignal)
	assert.Nil(t, err, "Error unmarshalling match signal")

	_, err = xml.Marshal(matchSignal)
	assert.Nil(t, err, "Error marshalling match signal")

	assert.Equalf(t, matchSignalSchemaDefault, matchSignal.Schema, "Expected default schema %s but got %s \n", matchSignalSchemaDefault, matchSignal.Schema)
}
