package scte224v20180501

import (
	"encoding/xml"
	"testing"
)

const viewingpolicy string = `
<?xml version="1.0" encoding="utf-8"?>
<ViewingPolicy 
	xmlns:action="urn:scte:224:action" 
  xmlns:xlink="http://www.w3.org/1999/xlink" 
	xmlns:audience="urn:scte:224:audience" 
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
  xmlns:metadata="urn:scte:224:metadata" id="test.com/viewingpolicy/test" lastUpdated="2019-02-22T16:03:19.907Z" 
  xmlns="http://www.scte.org/schemas/224">
  <Audience xlink:href="test.com/audience/all" />
	<action:FastForward>false</action:FastForward>
	<action:Content>TARGETSTREAM</action:Content>
  <action:Capture>
    <action:StartWindow>
      <action:Percentage>0</action:Percentage>
    </action:StartWindow>
    <action:StopWindow>
      <action:Percentage>250</action:Percentage>
    </action:StopWindow>
    <action:MidrollDAI>false</action:MidrollDAI>
  </action:Capture>
</ViewingPolicy>`

type testViewingPolicy struct {
	ReusableType
	ActionProperty []*AnyProperty `xml:"http://www.scte.org/schemas/224 Any,omitempty" json:"anys,omitempty"`
	XMLName        xml.Name       `xml:"http://www.scte.org/schemas/224 ViewingPolicy" json:"-"`
	Audience       *Audience      `xml:"http://www.scte.org/schemas/224 Audience" json:"audience,omitempty"`
	FastForward    string         `xml:"FastForward,omitempty" json:"fastForward,omitempty"`
	Capture        []*testCapture `xml:"Capture,omitempty" "json:"capture,omitempty"`
	Content        string         `xml:"Content,omitempty" "json:"content,omitempty"`
}

//Capture struct
type testCapture struct {
	StartWindow *testStartWindow `xml:"StartWindow,omitempty" json:"startWindow"`
	StopWindow  *testStopWindow  `xml:"StopWindow,omitempty" json:"stopWindow"`
	MidrollDAI  string           `xml:"MidrollDAI,omitempty" json:"midrollDAI"`
}

//StartWindow struct
type testStartWindow struct {
	Percentage string `xml:"Percentage,omitempty" json:"percentage"`
}

//StopWindow struct
type testStopWindow struct {
	Offset     string `xml:"Offset,omitempty" json:"offset"`
	Percentage string `xml:"Percentage,omitempty" json:"percentage"`
}

func TestViewingPolicy2018(t *testing.T) {

	var vpol *ViewingPolicy
	err := xml.Unmarshal([]byte(viewingpolicy), &vpol)

	//mybytes, myerr := xml.Marshal(vpol)
	//if myerr != nil {
	//	t.Errorf("Error marshalling viewingpolicys %v", myerr)
	//	t.FailNow()
	//}
	//
	//mystring := string(mybytes)
	//
	//fmt.Print(mystring)

	if err != nil {
		t.Errorf("Error unmarshalling viewingpolicys %v", err)
		t.FailNow()
	}

	vbyte, marshalErr := xml.Marshal(vpol)
	if marshalErr != nil {
		t.Errorf("Error Marshal %v", err)
		t.FailNow()
	}

	var tvpol *testViewingPolicy
	err = xml.Unmarshal(vbyte, &tvpol)
	if err != nil {
		t.Errorf("Error unmarshalling %v", err)
		t.FailNow()
	}

	if tvpol.FastForward == "" {
		t.Log("action:FastForward is empty")
		t.FailNow()
	}

	if tvpol.Content == "" {
		t.Log("action:Content is empty")
		t.FailNow()
	}

	// Capture
	if len(tvpol.Capture) == 0 {
		t.Log("action:Capture is empty and is supposed to have child nodes")
		t.FailNow()
	}

	for _, capture := range tvpol.Capture {
		// StartWindow
		if capture.StartWindow == nil {
			t.Log("action:Capture StartWindow is nil")
			t.FailNow()
		}
		if capture.StartWindow.Percentage == "" {
			t.Log("action:Capture StartWindow Percentage is empty")
			t.Fail()
		}

		// StopWindow
		if capture.StopWindow == nil {
			t.Log("action:Capture StopWindow is nil")
			t.FailNow()
		}

		if capture.StopWindow.Percentage == "" {
			t.Log("action:Capture StopWindow Percentage is empty")
			t.Fail()
		}

		// MidrollDAI
		if capture.MidrollDAI == "" {
			t.Log("action:Capture MidrollDAI is empty")
			t.Fail()
		}
	}
}

const mediapoint string = `
<MediaPoint xmlns="http://www.scte.org/schemas/224" description="Avail Start" effective="2016-07-05T14:59:50Z"
       expectedDuration="PT30S" expires="2016-07-05T15:10:00Z" id="providerXYZ.com/media/break/1234/start" matchTime="2016-07-05T15:00:00Z">
   <AltID description="Ad-ID">CNPA0484000H</AltID>
   <Metadata>
       <metadata DAIModel="SINGLE_ADVERTISER" duration="PT30S" offset="PT0S" ownerName="XYZ" ownerType="PROVIDER"/>
       <metadata DAIModel="DOUBLE_ADVERTISER" duration="PT30S" offset="PT30S" ownerName="ABCD" ownerType="DISTRIBUTOR"/>
   </Metadata>
   <Apply duration="PT30S">
       <Policy/>
   </Apply>
</MediaPoint>`

func TestMediaPointMetadata(t *testing.T) {

	var mp *MediaPoint
	err := xml.Unmarshal([]byte(mediapoint), &mp)
	if err != nil {
		t.Errorf("Error unmarshalling mediapoint %v", err)
		t.FailNow()
	}

	metadata := mp.Metadata
	if metadata == nil {
		t.Log("unexpected nil metadata")
		t.FailNow()
	}

	nodes := metadata.Nodes
	if nodes == nil {
		t.Log("unexpected nil nodes")
		t.FailNow()
	}

	if len(nodes) != 2 {
		t.Log("unexpected exactly 2 nodes")
		t.FailNow()
	}

	attributes := nodes[0].Attributes
	if attributes == nil {
		t.Log("unexpected nil attributes")
		t.FailNow()
	}

	if len(attributes) != 5 {
		t.Log("unexpected exactly 5 node attributes")
		t.FailNow()
	}

	attr1 := attributes[0]
	if attr1.Name.Local != "DAIModel" {
		t.Errorf("unexpected attribute name %v", attr1.Name.Local)
		t.FailNow()
	}
	if attr1.Value != "SINGLE_ADVERTISER" {
		t.Errorf("unexpected attribute value %v", attr1.Value)
		t.FailNow()
	}

}
