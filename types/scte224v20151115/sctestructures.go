package scte224v20151115

import (
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const schemaLocation = "http://www.scte.org/schemas/224/SCTE224-20151115.xsd"

// Structs for SCTE 224 2018 ESNI Objects.
// Table 3
type IdentifiableType struct {
	Id          string     `xml:"id,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`
	LastUpdated *time.Time `xml:"lastUpdated,attr,omitempty"`
	XMLBase     string     `xml:"xml:base,attr,omitempty"`
	AltIDs      []*AltID   `xml:"http://www.scte.org/schemas/224/2015 AltID,omitempty"`
	Metadata    *Metadata  `xml:"http://www.scte.org/schemas/224/2015 Metadata,omitempty"`
	Ext         *Metadata  `xml:"http://www.scte.org/schemas/224/2015 Ext,omitempty"`
}

//Table 5
type ReusableType struct {
	IdentifiableType
	XLinkHRef string `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
}

//********************* Media Types *************************//
//Table 6
type Media struct {
	ReusableType
	XMLName     xml.Name      `xml:"http://www.scte.org/schemas/224/2015 Media"`
	Effective   *time.Time    `xml:"effective,attr,omitempty"`
	Expires     *time.Time    `xml:"expires,attr,omitempty"`
	Source      string        `xml:"source,attr,omitempty"`
	MediaPoints []*MediaPoint `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
}

// MediaPoint defines an SCTE 224 (ESNI) media point object.
//Table 7
type MediaPoint struct {
	IdentifiableType
	XMLName          xml.Name     `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
	Effective        *time.Time   `xml:"effective,attr,omitempty"`
	Expires          *time.Time   `xml:"expires,attr,omitempty"`
	MatchTime        *time.Time   `xml:"matchTime,attr,omitempty"`
	MatchOffset      Duration     `xml:"matchOffset,attr,omitempty"`
	Source           string       `xml:"source,attr,omitempty"`
	ExpectedDuration Duration     `xml:"-"` // not in the 2015 XSD
	Order            *uint        `xml:"-"` // used internally for ordering but not in the 2015 XSD
	Reusable         bool         `xml:"-"` // not in the 2015 XSD
	Removes          []*Remove    `xml:"http://www.scte.org/schemas/224/2015 Remove"`
	Applys           []*Apply     `xml:"http://www.scte.org/schemas/224/2015 Apply"`
	MatchSignal      *MatchSignal `xml:"http://www.scte.org/schemas/224/2015 MatchSignal"`
	MediaGuid        string       `xml:"-"` // used internally to track which media this point is part of
}

type Metadata struct {
	InnerXml string `xml:",innerxml"`
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
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224/2015 AltID"`
	Description string   `xml:"-"` // not in 2015 XSD
	Value       string   `xml:",chardata"`
}

//Table 10
type Apply struct {
	XMLName  xml.Name `xml:"http://www.scte.org/schemas/224/2015 Apply"`
	Duration Duration `xml:"duration,attr,omitempty"`
	Priority *uint    `xml:"-"` // not in 2015 XSD
	Policy   *Policy  `xml:"http://www.scte.org/schemas/224/2015 Policy,omitempty"`
}

//Table 9
type Remove struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224/2015 Remove"`
	Policy  *Policy  `xml:"http://www.scte.org/schemas/224/2015 Policy",omitempty`
}

//Table 8
type MatchSignal struct {
	XMLName         xml.Name  `xml:"http://www.scte.org/schemas/224/2015 MatchSignal"`
	Match           Match     `xml:"match,attr,omitempty"`
	SignalTolerance Duration  `xml:"signalTolerance,attr,omitempty"`
	Assertions      []*Assert `xml:"http://www.scte.org/schemas/224/2015 Assert,omitempty"`
}

type Match string

//	Returns true if the value of this enumerated Match is "ALL".
func (me Match) IsAll() bool { return me == "ALL" }

//	Returns true if the value of this enumerated Match is "ANY".
func (me Match) IsAny() bool { return me == "ANY" }

//	Returns true if the value of this enumerated Match is "NONE".
func (me Match) IsNone() bool { return me == "NONE" }

type Assert struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224/2015 Assert"`
	Declaration string   `xml:",chardata"`
}

//********************* Audience Types *************************//
//Table 11
type Policy struct {
	ReusableType
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224/2015 Policy"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy,omitempty"`
}

type AnyProperty struct {
	XMLName xml.Name
	Data    string `xml:",innerxml"`
}

//Table 12
type ViewingPolicy struct {
	ReusableType
	XMLName        xml.Name      `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy"`
	Audience       *Audience     `xml:"http://www.scte.org/schemas/224/2015 Audience,omitempty"`
	ActionProperty []AnyProperty `xml:",any"`
}

//Table 13
type Audience struct {
	ReusableType
	XMLName          xml.Name      `xml:"http://www.scte.org/schemas/224/2015 Audience"`
	Match            Match         `xml:"match,attr,omitempty"`
	Audiences        []*Audience   `xml:"http://www.scte.org/schemas/224/2015 Audience,omitempty"`
	AudienceProperty []AnyProperty `xml:",any"`
}

//********************* Results Types *************************//
//Table 14
type Results struct {
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224/2015 Results"`
	Size           int              `xml:"size,attr,omitempty"`
	Medias         []*Media         `xml:"http://www.scte.org/schemas/224/2015 Media"`
	MediaPoints    []*MediaPoint    `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
	Policys        []*Policy        `xml:"http://www.scte.org/schemas/224/2015 Policy"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy"`
	Audiences      []*Audience      `xml:"http://www.scte.org/schemas/224/2015 Audience"`
	Audits         []*Audit         `xml:"http://www.scte.org/schemas/224/2015 Audit"`
}

//********************* Audit Types *************************//
//Table 15
type Audit struct {
	IdentifiableType
	XMLName       xml.Name `xml:"http://www.scte.org/schemas/224/2015 Audit"`
	XLinkHRef     string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XLinkRole     string   `xml:"http://www.w3.org/1999/xlink role,attr,omitempty"`
	Authorization string   `xml:"authorization,attr,omitempty"`
	PolicyMode    string   `xml:"policyMode,attr,omitempty"`
	Trigger       string   `xml:"trigger,attr,omitempty"`
	Result        string   `xml:"result,attr,omitempty"`
	Audits        []*Audit `xml:"http://www.scte.org/schemas/224/2015 Audit"`
}
