package scte224v20180501

import (
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const schemaLocation = "http://www.scte.org/schemas/224/SCTE224-20180501.xsd"

// Structs for SCTE 224 2018 ESNI Objects.
// Table 3
type IdentifiableType struct {
	Id          string     `xml:"id,attr,omitempty" json:"id,omitempty"`
	Description string     `xml:"description,attr,omitempty" json:"description,omitempty"`
	LastUpdated *time.Time `xml:"lastUpdated,attr,omitempty" json:"lastUpdated,omitempty"`
	XMLBase     string     `xml:"xml:base,attr,omitempty" json:"-"`
	AltIDs      []*AltID   `xml:"http://www.scte.org/schemas/224 AltID,omitempty" json:"altIDs,omitempty"`
	Metadata    *Metadata  `xml:"http://www.scte.org/schemas/224 Metadata,omitempty" json:"metadata,omitempty"`
	Ext         *Ext       `xml:"http://www.scte.org/schemas/224 Ext,omitempty" json:"ext,omitempty"`
}

//Table 5
type ReusableType struct {
	IdentifiableType
	XLinkHRef string `xml:"http://www.w3.org/1999/xlink href,attr,omitempty" json:"href,omitempty"`
}

//********************* Media Types *************************//
//Table 6
type Media struct {
	ReusableType
	XMLName     xml.Name      `xml:"http://www.scte.org/schemas/224 Media" json:"-"`
	Effective   *time.Time    `xml:"effective,attr,omitempty" json:"effective,omitempty"`
	Expires     *time.Time    `xml:"expires,attr,omitempty" json:"expires,omitempty"`
	Source      string        `xml:"source,attr,omitempty" json:"source,omitempty"`
	MediaPoints []*MediaPoint `xml:"http://www.scte.org/schemas/224 MediaPoint" json:"mediaPoints,omitempty"`
}

// MediaPoint defines an SCTE 224 (ESNI) media point object.
//Table 7
type MediaPoint struct {
	IdentifiableType
	XMLName          xml.Name     `xml:"http://www.scte.org/schemas/224 MediaPoint" json:"-"`
	Effective        *time.Time   `xml:"effective,attr,omitempty" json:"effective,omitempty"`
	Expires          *time.Time   `xml:"expires,attr,omitempty" json:"expires,omitempty"`
	MatchTime        *time.Time   `xml:"matchTime,attr,omitempty" json:"matchTime,omitempty"`
	MatchOffset      Duration     `xml:"matchOffset,attr,omitempty" json:"matchOffset,omitempty"`
	Source           string       `xml:"source,attr,omitempty" json:"source,omitempty"`
	ExpectedDuration Duration     `xml:"expectedDuration,attr,omitempty" json:"expectedDuration,omitempty"`
	Order            *uint        `xml:"order,attr,omitempty" json:"order,omitempty"`
	Reusable         bool         `xml:"reusable,attr,omitempty" json:"reusable,omitempty"`
	Removes          []*Remove    `xml:"http://www.scte.org/schemas/224 Remove" json:"removes,omitempty"`
	Applys           []*Apply     `xml:"http://www.scte.org/schemas/224 Apply" json:"applys,omitempty"`
	MatchSignal      *MatchSignal `xml:"http://www.scte.org/schemas/224 MatchSignal" json:"matchSignal,omitempty"`
	MediaGuid        string       `xml:"-"` // used internally to track which media this point is part of
}

func (mp *MediaPoint) HasExplicitOrder() bool {
	return nil != mp.Order
}

func (mp *MediaPoint) GetOrder() uint {
	if mp.HasExplicitOrder() {
		return *mp.Order
	}
	return 0
}

type Metadata struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Metadata" json:"-"`
	Nodes   []Any    `xml:",any" json:"values,omitempty"`
}

type Ext struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Ext"`
	Nodes   []Any    `xml:",any" json:"values,omitempty"`
}

type Any struct {
	XMLName    xml.Name   `json:"xmlname"`
	Attributes []xml.Attr `xml:",any,attr"`
	Value      string     `xml:",chardata" json:"value"`
}

type Duration string

func (dur Duration) GoDuration() time.Duration {
	return ConvertDuration(string(dur))
}

var durationRegex = regexp.MustCompile("P(-?\\d*)D?T(-?\\d*H?-?\\d*M?-?\\d*\\.?\\d*S?)")

// TODO: I doubt we will see durations longer than days but in theory we should handle them
func ConvertDuration(xmlDuration string) (duration time.Duration) {

	if "" != xmlDuration {
		subMatches := durationRegex.FindStringSubmatch(xmlDuration)
		if len(subMatches) == 3 {
			var hours int
			var nanoseconds time.Duration
			if "" != subMatches[1] {
				days, convErr := strconv.Atoi(subMatches[1])
				if nil == convErr {
					hours = 24 * days
				} else {
					log.Println(convErr)
				}
			}

			if "" != subMatches[2] {
				parsedDur, convErr := time.ParseDuration(strings.ToLower(subMatches[2]))
				if nil == convErr {
					nanoseconds = parsedDur
				} else {
					log.Println(convErr)
				}
			}
			duration = nanoseconds + time.Duration(hours)*time.Hour
		}
	}
	return duration
}

func ToDuration(dur time.Duration) Duration {
	return Duration("P" + strings.ToUpper(dur.Round(time.Second).String()))
}

type AltID struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 AltID" json:"-"`
	Description string   `xml:"description,attr,omitempty" json:"description,omitempty"`
	Value       string   `xml:",chardata" json:"value,omitempty"`
}

//Table 10
type Apply struct {
	XMLName  xml.Name `xml:"http://www.scte.org/schemas/224 Apply" json:"-"`
	Duration Duration `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	Priority *uint    `xml:"priority,attr,omitempty" json:"priority,omitempty"`
	Policy   *Policy  `xml:"http://www.scte.org/schemas/224 Policy,omitempty" json:"policy,omitempty"`
}

func (ap *Apply) HasExplicitPriority() bool {
	return nil != ap.Priority
}

func (ap *Apply) GetPriority() uint {
	if ap.HasExplicitPriority() {
		return *ap.Priority
	}
	return 0
}

//Table 9
type Remove struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Remove" json:"-"`
	Policy  *Policy  `xml:"http://www.scte.org/schemas/224 Policy,omitempty" json:"policy,omitempty"`
}

//Table 8
type MatchSignal struct {
	XMLName         xml.Name  `xml:"http://www.scte.org/schemas/224 MatchSignal" json:"-"`
	Match           Match     `xml:"match,attr,omitempty" json:"match,omitempty"`
	SignalTolerance Duration  `xml:"signalTolerance,attr,omitempty" json:"signalTolerance,omitempty"`
	Assertions      []*Assert `xml:"http://www.scte.org/schemas/224 Assert,omitempty" json:"assertions,omitempty"`
}

type Match string

//	Returns true if the value of this enumerated Match is "ALL".
func (me Match) IsAll() bool { return me == "ALL" }

//	Returns true if the value of this enumerated Match is "ANY".
func (me Match) IsAny() bool { return me == "ANY" }

//	Returns true if the value of this enumerated Match is "NONE".
func (me Match) IsNone() bool { return me == "NONE" }

type Assert struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 Assert" json:"-"`
	Declaration string   `xml:",chardata" json:"declaration,omitempty"`
}

//********************* Audience Types *************************//
//Table 11
type Policy struct {
	ReusableType
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Policy" json:"-"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy,omitempty" json:"viewingPolicys,omitempty"`
}

//Table 12
type ViewingPolicy struct {
	ReusableType
	XMLName        xml.Name      `xml:"http://www.scte.org/schemas/224 ViewingPolicy" json:"-"`
	Audience       *Audience     `xml:"http://www.scte.org/schemas/224 Audience,omitempty" json:"audience,omitempty"`
	ActionProperty []AnyProperty `xml:",any" json:"actionProperty,omitempty"`
}

//Table 13
type Audience struct {
	ReusableType
	XMLName          xml.Name      `xml:"http://www.scte.org/schemas/224 Audience" json:"-"`
	Match            Match         `xml:"match,attr,omitempty" json:"match,omitempty"`
	Audiences        []*Audience   `xml:"http://www.scte.org/schemas/224 Audience,omitempty" json:"audiences,omitempty"`
	AudienceProperty []AnyProperty `xml:",any" json:"audienceProperty,omitempty"`
}

type AnyProperty struct {
	XMLName xml.Name `json:"-"`
	Data    string   `xml:",innerxml" json:"data,omitempty"`
}

//********************* Results Types *************************//
//Table 14
type Results struct {
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Results" json:"-"`
	Size           int              `xml:"size,attr,omitempty" json:"size,omitempty"`
	Medias         []*Media         `xml:"http://www.scte.org/schemas/224 Media" json:"medias,omitempty"`
	MediaPoints    []*MediaPoint    `xml:"http://www.scte.org/schemas/224 MediaPoint" json:"mediaPoints,omitempty"`
	Policys        []*Policy        `xml:"http://www.scte.org/schemas/224 Policy" json:"policys,omitempty"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy" json:"viewingPolicys,omitempty"`
	Audiences      []*Audience      `xml:"http://www.scte.org/schemas/224 Audience" json:"audiences,omitempty"`
	Audits         []*Audit         `xml:"http://www.scte.org/schemas/224 Audit" json:"audits,omitempty"`
}

//********************* Audit Types *************************//
//Table 15
type Audit struct {
	IdentifiableType
	XMLName       xml.Name `xml:"http://www.scte.org/schemas/224 Audit" json:"-"`
	XLinkHRef     string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty" json:"href,omitempty"`
	XLinkRole     string   `xml:"http://www.w3.org/1999/xlink role,attr,omitempty" json:"role,omitempty"`
	Authorization string   `xml:"authorization,attr,omitempty" json:"authorization,omitempty"`
	PolicyMode    string   `xml:"policyMode,attr,omitempty" json:"policyMode,omitempty"`
	Trigger       string   `xml:"trigger,attr,omitempty" json:"trigger,omitempty"`
	Result        string   `xml:"result,attr,omitempty" json:"result,omitempty"`
	Audits        []*Audit `xml:"http://www.scte.org/schemas/224 Audit" json:"audits,omitempty"`
}
