package scte224v20200407

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Comcast/scte224structs/convert"
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224_2018 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

//Table 12
type ViewingPolicy struct {
	ReusableType
	XMLName              xml.Name                    `xml:"http://www.scte.org/schemas/224 ViewingPolicy" json:"-"`
	Audience             *Audience                   `xml:"http://www.scte.org/schemas/224 Audience,omitempty" json:"audience,omitempty"`
	SignalPointDeletion  *SignalPointDeletionAction  `xml:"urn:scte:224:action SignalPointDeletion,omitempty" json:"signalPointDeletion,omitempty"`
	SignalPointInsertion *SignalPointInsertionAction `xml:"urn:scte:224:action SignalPointInsertion,omitempty" json:"signalPointInsertion,omitempty"`
	Content              *ContentAction              `xml:"urn:scte:224:action Content,omitempty" json:"content,omitempty"`
	Allocation           *Allocation                 `xml:"urn:scte:224:action Allocation,omitempty" json:"Allocation,omitempty"`
	ActionProperty       []Any                       `xml:",any" json:"actionProperty,omitempty"`
}

func (vp *ViewingPolicy) Get2018() (*scte224_2018.ViewingPolicy, error) {
	encodedVP, err := xml.Marshal(vp)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Unmarshal into 2018 struct
	var vp2018 *scte224_2018.ViewingPolicy
	err = xml.Unmarshal(encodedVP, &vp2018)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Remove any "Allocation" action from the 2018 object, because it is an addition in 2020 spec
	vp2018.RemoveAllocationAction()
	return vp2018, nil
}

func (vp *ViewingPolicy) Get2015() (*scte224_2015.ViewingPolicy, error) {
	// First downgrade it to 2018
	vp2018, err := vp.Get2018()
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Then downgrade it to 2015, with the existing utility
	vp2015 := convert.DowngradeViewingPolicy(*vp2018)
	return &vp2015, nil
}

type Allocation struct {
	XMLName xml.Name `xml:"urn:scte:224:action Allocation" json:"-"`
	Slots   []*Slots `xml:"Slots,omitempty" json:"Slots,omitempty"`
}

type Slots struct {
	XMLName xml.Name `xml:"Slots" json:"-"`
	AdSlots []*Slot  `xml:"Slot,omitempty" json:"Slot,omitempty"`
}

type Slot struct {
	XMLName        xml.Name          `xml:"Slot" json:"-"`
	AdsReferenceId []*AdsReferenceId `xml:"AdsReferenceId,omitempty" json:"AdsReferenceId,omitempty"`
	Duration       Duration          `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	Offset         Duration          `xml:"offset,attr,omitempty" json:"offset,omitempty"`
}

type AdsReferenceId struct {
	XMLName       xml.Name `xml:"AdsReferenceId" json:"-"`
	ID            string   `xml:",chardata" json:"data,omitempty"`
	ReferenceType string   `xml:"referenceType,attr,omitempty" json:"referenceType,omitempty"`
	Exclude       bool     `xml:"exclude,attr,omitempty" json:"exclude,omitempty"`
}

type ContentAction struct {
	XMLName xml.Name `xml:"urn:scte:224:action Content" json:"-"`
	Content string   `xml:",chardata" json:"data,omitempty"`
}

type SignalPointDeletionAction struct {
	XMLName             xml.Name `xml:"urn:scte:224:action SignalPointDeletion" json:"-"`
	SignalPointDeletion string   `xml:",chardata" json:"data,omitempty"`
}

type SignalPointInsertionAction struct {
	SignalPoints []*SignalPoint `xml:"urn:scte:224:action SignalPoint,omitempty" json:"signalPoint,omitempty"`
}

type SignalPoint struct {
	Offset               Duration   `xml:"offset,attr,omitempty" json:"offset,omitempty"`
	SegmentationTypeId   *uint      `xml:"segmentationTypeId,attr,omitempty" json:"segmentationTypeId,omitempty"`
	SegmentationUpidType *uint      `xml:"segmentationUpidType,attr,omitempty" json:"segmentationUpidType,omitempty"`
	SegmentationUpid     string     `xml:"segmentationUpid,attr,omitempty" json:"segmentationUpid,omitempty"`
	RepeatInterval       Duration   `xml:"repeatInterval,attr,omitempty" json:"repeatInterval,omitempty"`
	RepeatStart          *time.Time `xml:"repeatStart,attr,omitempty" json:"repeatStart,omitempty"`
	RepeatStop           *time.Time `xml:"repeatStop,attr,omitempty" json:"repeatStop,omitempty"`
}
