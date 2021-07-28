package scte224v20200407

import (
	"encoding/xml"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Comcast/scte224structs/convert"
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224_2018 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

const schemaLocation = "https://www.scte.org/standards-development/library/standards-catalog/ansiscte-224-2018r1/"

// Structs for SCTE 224 2020 ESNI Objects.
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

func (m *Media) Get2018() scte224_2018.Media {
	destination := scte224_2018.Media{}
	if m == nil {
		return destination
	}

	destination.ReusableType = m.ReusableType.Get2018()
	destination.XMLName = m.XMLName
	destination.Effective = m.Effective
	destination.Expires = m.Expires
	destination.Source = m.Source

	for _, mp := range m.MediaPoints {
		if mp == nil {
			continue
		}

		mp2018 := mp.Get2018()
		destination.MediaPoints = append(destination.MediaPoints, &mp2018)
	}

	return destination
}

func (m *Media) Get2015() scte224_2015.Media {
	media2018 := m.Get2018()
	return convert.DowngradeMedia(media2018)
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

func (mp *MediaPoint) Get2018() scte224_2018.MediaPoint {
	destination := scte224_2018.MediaPoint{}
	if mp == nil {
		return destination
	}

	destination.IdentifiableType = mp.IdentifiableType.Get2018()
	destination.XMLName = mp.XMLName
	destination.Effective = mp.Effective
	destination.Expires = mp.Expires
	destination.MatchTime = mp.MatchTime
	destination.MatchOffset = scte224_2018.Duration(mp.MatchOffset)
	destination.Source = mp.Source
	destination.ExpectedDuration = scte224_2018.Duration(mp.ExpectedDuration)
	destination.Order = mp.Order
	destination.Reusable = mp.Reusable
	destination.MediaGuid = mp.MediaGuid

	if mp.MatchSignal != nil {
		matchSignal2018 := &scte224_2018.MatchSignal{}
		matchSignal2018.XMLName = mp.MatchSignal.XMLName
		matchSignal2018.Match = scte224_2018.Match(mp.MatchSignal.Match)
		matchSignal2018.SignalTolerance = scte224_2018.Duration(mp.MatchSignal.SignalTolerance)

		for _, assertion := range mp.MatchSignal.Assertions {
			if assertion == nil {
				continue
			}

			matchSignal2018.Assertions = append(matchSignal2018.Assertions, &scte224_2018.Assert{
				XMLName:     assertion.XMLName,
				Declaration: assertion.Declaration,
			})
		}

		destination.MatchSignal = matchSignal2018
	}

	for _, apply := range mp.Applys {
		if apply == nil {
			continue
		}

		apply2018 := scte224_2018.Apply{
			XMLName:  apply.XMLName,
			Duration: scte224_2018.Duration(apply.Duration),
			Priority: apply.Priority,
		}
		applyPolicy2018 := apply.Policy.Get2018()
		apply2018.Policy = &applyPolicy2018

		destination.Applys = append(destination.Applys, &apply2018)
	}

	for _, remove := range mp.Removes {
		if remove == nil {
			continue
		}

		remove2018 := scte224_2018.Remove{
			XMLName: remove.XMLName,
		}
		removePolicy2018 := remove.Policy.Get2018()
		remove2018.Policy = &removePolicy2018

		destination.Removes = append(destination.Removes, &remove2018)
	}

	return destination
}

func (mp *MediaPoint) Get2015() scte224_2015.MediaPoint {
	mp2018 := mp.Get2018()
	return convert.DowngradeMediaPoint(mp2018)
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

const matchSignalSchemaDefault = "http://www.scte.org/schemas/35"

//Table 8
type MatchSignal struct {
	XMLName         xml.Name  `xml:"http://www.scte.org/schemas/224 MatchSignal" json:"-"`
	Match           Match     `xml:"match,attr,omitempty" json:"match,omitempty"`
	SignalTolerance Duration  `xml:"signalTolerance,attr,omitempty" json:"signalTolerance,omitempty"`
	Assertions      []*Assert `xml:"http://www.scte.org/schemas/224 Assert,omitempty" json:"assertions,omitempty"`
	Schema          string    `xml:"schema,attr,omitempty" json:"schema,omitempty"`
}

// Custom unmarshalling to provide default value to "schema" attribute, when unavailable
func (ms *MatchSignal) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	/**
	* https://pkg.go.dev/encoding/xml#Unmarshaler
	* One common implementation strategy is to unmarshal into a separate value with a layout matching the expected XML using d.DecodeElement,
	* and then to copy the data from that value into the receiver.
	* this is to prevent recursive unmarshaling ending in an infinite loop
	 */
	type dupMatchSignalType MatchSignal
	matchSignal := dupMatchSignalType{}

	err := d.DecodeElement(&matchSignal, &start)
	if err != nil {
		return err
	}

	if strings.Trim(matchSignal.Schema, " ") == "" {
		// set the default value
		matchSignal.Schema = matchSignalSchemaDefault
	}

	*ms = MatchSignal(matchSignal)
	return nil
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
