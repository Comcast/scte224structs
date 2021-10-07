package scte224v20200407

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewingPolicyDowngrade(t *testing.T) {
	var vp *ViewingPolicy
	err := xml.Unmarshal([]byte(vp2020Raw), &vp)
	assert.Nil(t, err, "Error unmarshalling viewingpolicy")

	vp2018 := vp.Get2018()
	vp2018Marshaled, err := xml.MarshalIndent(vp2018, "", "    ")
	assert.Nil(t, err, "Error marshaling downgraded viewingpolicy")
	assert.EqualValues(t, vp2018Raw, string(vp2018Marshaled), "Downgrade failed")
}

func TestAudienceDowngrade(t *testing.T) {
	var aud Audience
	err := xml.Unmarshal([]byte(aud2020Raw), &aud)
	assert.Nil(t, err, "Error unmarshalling audience")

	aud2018 := aud.Get2018()
	_, err = xml.MarshalIndent(aud2018, "", "    ")
	assert.Nil(t, err, "Error marshaling downgraded audience")
	assert.Len(t, aud2018.Audiences, 2, "Expected 2 nested Audience")
}

func TestViewingPolicyAllocation(t *testing.T) {
	var vpol ViewingPolicy
	err := xml.Unmarshal([]byte(vp2020Raw), &vpol)
	assert.Nil(t, err, "Error unmarshalling viewingpolicy")

	//var roundtrip []byte
	_, err = xml.Marshal(vpol)
	assert.Nil(t, err, "Error marshalling viewingpolicy")
	//t.Log(string(roundtrip))

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
