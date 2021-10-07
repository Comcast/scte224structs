package scte224v20200407

import (
	"encoding/xml"
	"time"

	"github.com/Comcast/scte224structs/convert"
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224_2018 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

//Table 11
type Policy struct {
	ReusableType
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Policy" json:"-"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy,omitempty" json:"viewingPolicys,omitempty"`
}

func (p *Policy) Get2018() scte224_2018.Policy {
	destination := scte224_2018.Policy{}
	if p == nil {
		return destination
	}

	destination.ReusableType = p.ReusableType.Get2018()
	destination.XMLName = p.XMLName

	for _, vp := range p.ViewingPolicys {
		if vp == nil {
			continue
		}

		vp2018 := vp.Get2018()
		destination.ViewingPolicys = append(destination.ViewingPolicys, &vp2018)
	}

	return destination
}

func (p *Policy) Get2015() scte224_2015.Policy {
	policy2018 := p.Get2018()
	return convert.DowngradePolicy(policy2018)
}

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

func (vp *ViewingPolicy) Get2018() scte224_2018.ViewingPolicy {
	destination := scte224_2018.ViewingPolicy{}
	if vp == nil {
		return destination
	}

	destination.ReusableType = vp.ReusableType.Get2018()
	destination.XMLName = vp.XMLName

	if vp.Audience != nil {
		aud2018 := vp.Audience.Get2018()
		destination.Audience = &aud2018
	}

	if vp.SignalPointDeletion != nil {
		destination.SignalPointDeletion = &scte224_2018.SignalPointDeletionAction{
			XMLName:             vp.SignalPointDeletion.XMLName,
			SignalPointDeletion: vp.SignalPointDeletion.SignalPointDeletion,
		}
	}

	if vp.SignalPointInsertion != nil {
		signalPoints2018 := make([]*scte224_2018.SignalPoint, len(vp.SignalPointInsertion.SignalPoints))
		for _, signalPoint := range vp.SignalPointInsertion.SignalPoints {
			if signalPoint == nil {
				continue
			}

			signalPoints2018 = append(signalPoints2018, &scte224_2018.SignalPoint{
				Offset:               scte224_2018.Duration(signalPoint.Offset),
				SegmentationTypeId:   signalPoint.SegmentationTypeId,
				SegmentationUpidType: signalPoint.SegmentationUpidType,
				SegmentationUpid:     signalPoint.SegmentationUpid,
				RepeatInterval:       scte224_2018.Duration(signalPoint.RepeatInterval),
				RepeatStart:          signalPoint.RepeatStart,
				RepeatStop:           signalPoint.RepeatStop,
			})
		}

		destination.SignalPointInsertion = &scte224_2018.SignalPointInsertionAction{
			SignalPoints: signalPoints2018,
		}
	}

	if vp.Content != nil {
		destination.Content = &scte224_2018.ContentAction{
			XMLName: vp.Content.XMLName,
			Content: vp.Content.Content,
		}
	}

	for _, action := range vp.ActionProperty {
		destination.ActionProperty = append(destination.ActionProperty, action.Get2018())
	}

	return destination
}

func (vp *ViewingPolicy) Get2015() scte224_2015.ViewingPolicy {
	vp2018 := vp.Get2018()
	return convert.DowngradeViewingPolicy(vp2018)
}

type Allocation struct {
	XMLName   xml.Name `xml:"urn:scte:224:action Allocation" json:"-"`
	Slots     []*Slots `xml:"Slots,omitempty" json:"Slots,omitempty"`
	OwnerType string   `xml:"ownerType,attr,omitempty" json:"ownerType,attr,omitempty"`
	OwnerName string   `xml:"ownerName,attr,omitempty" json:"ownerName,attr,omitempty"`
	Duration  Duration `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	Ads       string   `xml:"ads,attr,omitempty" json:"ads,attr,omitempty"`
}

type Slots struct {
	XMLName xml.Name `xml:"urn:scte:224:action Slots" json:"-"`
	AdSlots []*Slot  `xml:"Slot,omitempty" json:"Slot,omitempty"`
}

type Slot struct {
	XMLName        xml.Name          `xml:"urn:scte:224:action Slot" json:"-"`
	AdsReferenceId []*AdsReferenceId `xml:"AdsReferenceId,omitempty" json:"AdsReferenceId,omitempty"`
	SlotRules      *SlotRules        `xml:"SlotRules,omitempty" json:"SlotRules,omitempty"`
	Duration       Duration          `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	Offset         Duration          `xml:"offset,attr,omitempty" json:"offset,omitempty"`
}

type AdsReferenceId struct {
	XMLName       xml.Name `xml:"urn:scte:224:action AdsReferenceId" json:"-"`
	ID            string   `xml:",chardata" json:"data,omitempty"`
	ReferenceType string   `xml:"referenceType,attr,omitempty" json:"referenceType,omitempty"`
	Exclude       bool     `xml:"exclude,attr,omitempty" json:"exclude,omitempty"`
}

type SlotRules struct {
	XMLName    xml.Name     `xml:"urn:scte:224:action SlotRules" json:"-"`
	SlotRule []*SlotRule `xml:"SlotRule,omitempty" json:"SlotRule,omitempty"`
}

type SlotRule struct {
	XMLName    xml.Name     `xml:"urn:scte:224:action SlotRule" json:"-"`
	Parameters []*Parameter `xml:"Parameter,omitempty" json:"Parameter,omitempty"`
	Rule       string       `xml:"rule,attr,omitempty" json:"rule,attr,omitempty"`
}

type Parameter struct {
	XMLName       xml.Name `xml:"urn:scte:224:action Parameter" json:"-"`
	ParameterName string   `xml:"parameterName,attr,omitempty" json:"parameterName,attr,omitempty"`
	Value         string   `xml:",chardata" json:"value,omitempty"`
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

//Table 13
type Audience struct {
	ReusableType
	XMLName          xml.Name    `xml:"http://www.scte.org/schemas/224 Audience" json:"-"`
	Match            Match       `xml:"match,attr,omitempty" json:"match,omitempty"`
	Audiences        []*Audience `xml:"http://www.scte.org/schemas/224 Audience,omitempty" json:"audiences,omitempty"`
	AudienceProperty []Any       `xml:",any" json:"audienceProperty,omitempty"`
}

func (aud *Audience) Get2018() scte224_2018.Audience {
	destination := scte224_2018.Audience{}

	if aud == nil {
		return destination
	}

	destination.ReusableType = aud.ReusableType.Get2018()
	destination.XMLName = aud.XMLName
	destination.Match = scte224_2018.Match(aud.Match)

	for _, nestedAud := range aud.Audiences {
		if nestedAud == nil {
			continue
		}

		aud2018 := nestedAud.Get2018()
		destination.Audiences = append(destination.Audiences, &aud2018)
	}

	for _, audProp := range aud.AudienceProperty {
		destination.AudienceProperty = append(destination.AudienceProperty, audProp.Get2018())
	}

	return destination
}

func (aud *Audience) Get2015() scte224_2015.Audience {
	aud2018 := aud.Get2018()
	return convert.DowngradeAudience(aud2018)
}
