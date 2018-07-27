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
	Id          string     `xml:"id,attr,omitempty"`
	Description string     `xml:"description,attr,omitempty"`
	LastUpdated *time.Time `xml:"lastUpdated,attr,omitempty"`
	XMLBase     string     `xml:"xml:base,attr,omitempty"`
	AltIDs      []*AltID   `xml:"http://www.scte.org/schemas/224 AltID,omitempty"`
	Metadata    *Metadata  `xml:"http://www.scte.org/schemas/224 Metadata,omitempty"`
	Ext         *Metadata  `xml:"http://www.scte.org/schemas/224 Ext,omitempty"`
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
	XMLName     xml.Name      `xml:"http://www.scte.org/schemas/224 Media"`
	Effective   *time.Time    `xml:"effective,attr,omitempty"`
	Expires     *time.Time    `xml:"expires,attr,omitempty"`
	Source      string        `xml:"source,attr,omitempty"`
	MediaPoints []*MediaPoint `xml:"http://www.scte.org/schemas/224 MediaPoint"`
}

// MediaPoint defines an SCTE 224 (ESNI) media point object.
//Table 7
type MediaPoint struct {
	IdentifiableType
	XMLName          xml.Name     `xml:"http://www.scte.org/schemas/224 MediaPoint"`
	Effective        *time.Time   `xml:"effective,attr,omitempty"`
	Expires          *time.Time   `xml:"expires,attr,omitempty"`
	MatchTime        *time.Time   `xml:"matchTime,attr,omitempty"`
	MatchOffset      Duration     `xml:"matchOffset,attr,omitempty"`
	Source           string       `xml:"source,attr,omitempty"`
	ExpectedDuration Duration     `xml:"expectedDuration,attr,omitempty"`
	Order            uint         `xml:"order,attr,omitempty"`
	Reusable         bool         `xml:"reusable,attr,omitempty"`
	Removes          []*Remove    `xml:"http://www.scte.org/schemas/224 Remove"`
	Applys           []*Apply     `xml:"http://www.scte.org/schemas/224 Apply"`
	MatchSignal      *MatchSignal `xml:"http://www.scte.org/schemas/224 MatchSignal"`
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
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 AltID"`
	Description string   `xml:"description,attr,omitempty"`
	Value       string   `xml:",chardata"`
}

//Table 10
type Apply struct {
	XMLName  xml.Name `xml:"http://www.scte.org/schemas/224 Apply"`
	Duration Duration `xml:"duration,attr,omitempty"`
	Priority uint     `xml:"priority,attr,omitempty"`
	Policy   *Policy  `xml:"http://www.scte.org/schemas/224 Policy,omitempty"`
}

//Table 9
type Remove struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Remove"`
	Policy  *Policy  `xml:"http://www.scte.org/schemas/224 Policy",omitempty`
}

//Table 8
type MatchSignal struct {
	XMLName         xml.Name  `xml:"http://www.scte.org/schemas/224 MatchSignal"`
	Match           Match     `xml:"match,attr,omitempty"`
	SignalTolerance Duration  `xml:"signalTolerance,attr,omitempty"`
	Assertions      []*Assert `xml:"http://www.scte.org/schemas/224 Assert,omitempty"`
}

type Match string

//	Returns true if the value of this enumerated Match is "ALL".
func (me Match) IsAll() bool { return me == "ALL" }

//	Returns true if the value of this enumerated Match is "ANY".
func (me Match) IsAny() bool { return me == "ANY" }

//	Returns true if the value of this enumerated Match is "NONE".
func (me Match) IsNone() bool { return me == "NONE" }

type Assert struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 Assert"`
	Declaration string   `xml:",chardata"`
}

//********************* Audience Types *************************//
//Table 11
type Policy struct {
	ReusableType
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Policy"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy,omitempty"`
}

//Table 12
type ViewingPolicy struct {
	ReusableType
	XMLName        xml.Name  `xml:"http://www.scte.org/schemas/224 ViewingPolicy"`
	Audience       *Audience `xml:"http://www.scte.org/schemas/224 Audience,omitempty"`
	ActionProperty string    `xml:",any,omitempty"`
}

//Table 13
type Audience struct {
	ReusableType
	XMLName          xml.Name    `xml:"http://www.scte.org/schemas/224 Audience"`
	Match            Match       `xml:"match,attr,omitempty"`
	Audiences        []*Audience `xml:"http://www.scte.org/schemas/224 Audience,omitempty"`
	AudienceProperty string      `xml:",any,omitempty"`
}

//********************* Results Types *************************//
//Table 14
type Results struct {
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Results"`
	Size           int              `xml:"size,attr,omitempty"`
	Medias         []*Media         `xml:"http://www.scte.org/schemas/224 Media"`
	MediaPoints    []*MediaPoint    `xml:"http://www.scte.org/schemas/224 MediaPoint"`
	Policys        []*Policy        `xml:"http://www.scte.org/schemas/224 Policy"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy"`
	Audiences      []*Audience      `xml:"http://www.scte.org/schemas/224 Audience"`
	Audits         []*Audit         `xml:"http://www.scte.org/schemas/224 Audit"`
}

//********************* Audit Types *************************//
//Table 15
type Audit struct {
	IdentifiableType
	XMLName       xml.Name `xml:"http://www.scte.org/schemas/224 Audit"`
	XLinkHRef     string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XLinkRole     string   `xml:"http://www.w3.org/1999/xlink role,attr,omitempty"`
	Authorization string   `xml:"authorization,attr,omitempty"`
	PolicyMode    string   `xml:"policyMode,attr,omitempty"`
	Trigger       string   `xml:"trigger,attr,omitempty"`
	Result        string   `xml:"result,attr,omitempty"`
	Audits        []*Audit `xml:"http://www.scte.org/schemas/224 Audit"`
}
