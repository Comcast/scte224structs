package scte224v20200407

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaDowngrade(t *testing.T) {
	var media *Media
	decoder := xml.NewDecoder(strings.NewReader(media2020Raw))
	err := decoder.Decode(&media)
	assert.Nil(t, err, "Error unmarshalling media")

	media2018Marshaled, err := xml.MarshalIndent(media.Get2018(), "", "    ")
	assert.Nil(t, err, "Error marshaling downgraded media")
	// this assert is not very robust, works for the existing input,
	// but can fail if the input Media has characters that is encoded differently, so watch out while changing input
	assert.EqualValues(t, media2018Raw, string(media2018Marshaled), "Downgrade failed")

	// Another test to specifically check that Apply and Remove policies are correctly downgraded
	// Reset media
	media = nil
	decoder = xml.NewDecoder(strings.NewReader(anotherMedia2020Raw))
	err = decoder.Decode(&media)
	assert.Nil(t, err, "Error unmarshalling media")

	media2018 := media.Get2018()
	if assert.Len(t, media2018.MediaPoints, 1, "Expected 1 Mediapoint") {
		assert.Len(t, media2018.MediaPoints[0].Applys, 1, "Expected 1 Apply policy")
		assert.Len(t, media2018.MediaPoints[0].Removes, 2, "Expected 2 Remove policies")
	}
}

const matchSignalXML = `
<MatchSignal xmlns="http://www.scte.org/schemas/224" match="ANY">
    <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=16]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),'12345')]</Assert>
    <Assert xmlns="http://www.scte.org/schemas/224">/SpliceInfoSection/SegmentationDescriptor[@segmentationTypeId=1]/SegmentationUpid[@segmentationUpidType=1 and contains(text(),'34567')]</Assert>
  </MatchSignal>`

func TestMatchSignal2020(t *testing.T) {
	var matchSignal *MatchSignal
	err := xml.Unmarshal([]byte(matchSignalXML), &matchSignal)
	assert.Nil(t, err, "Error unmarshalling match signal")

	_, err = xml.Marshal(matchSignal)
	assert.Nil(t, err, "Error marshalling match signal")

	assert.Equalf(t, matchSignalSchemaDefault, matchSignal.Schema, "Expected default schema %s but got %s \n", matchSignalSchemaDefault, matchSignal.Schema)
}
