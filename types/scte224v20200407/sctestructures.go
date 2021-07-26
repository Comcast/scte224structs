package scte224v20200407

import (
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Comcast/scte224structs/convert"
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224_2018 "github.com/Comcast/scte224structs/types/scte224v20180501"
	"github.com/Comcast/scte224structs/types/scte224v20180501/adi30"
)

const schemaLocation = "https://www.scte.org/standards-development/library/standards-catalog/ansiscte-224-2018r1/"

// Structs for SCTE 224 2020 ESNI Objects.
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

func (m *Media) Get2018() (*scte224_2018.Media, error) {
	encodedMedia, err := xml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Unmarshal into 2018 struct
	var media2018 *scte224_2018.Media
	err = xml.Unmarshal(encodedMedia, &media2018)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Apply the MediaPoint transformations
	mediaPoints2018 := make([]*scte224_2018.MediaPoint, len(m.MediaPoints))
	for _, mp := range m.MediaPoints {
		mp2018, err := mp.Get2018()
		if err != nil {
			return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
		}

		mediaPoints2018 = append(mediaPoints2018, mp2018)
	}
	media2018.MediaPoints = mediaPoints2018

	return media2018, nil
}

func (m *Media) Get2015() (*scte224_2015.Media, error) {
	// First downgrade it to 2018
	media2018, err := m.Get2018()
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Then downgrade it to 2015, with the existing utility
	media2015 := convert.DowngradeMedia(*media2018)
	return &media2015, nil
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

func (mp *MediaPoint) Get2018() (*scte224_2018.MediaPoint, error) {
	encodedMP, err := xml.Marshal(mp)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Unmarshal into 2018 struct
	var mp2018 *scte224_2018.MediaPoint
	err = xml.Unmarshal(encodedMP, &mp2018)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Apply Policy transformations
	applyPolicies2018 := make([]*scte224_2018.Apply, len(mp.Applys))
	for _, applyElement := range mp.Applys {
		applyPolicy2018 := scte224_2018.Apply{
			XMLName:  applyElement.XMLName,
			Duration: scte224_2018.Duration(applyElement.Duration),
			Priority: applyElement.Priority,
		}
		policy2018, err := applyElement.Policy.Get2018()
		if err != nil {
			return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
		}
		applyPolicy2018.Policy = policy2018

		applyPolicies2018 = append(applyPolicies2018, &applyPolicy2018)
	}

	removePolicies2018 := make([]*scte224_2018.Remove, len(mp.Removes))
	for _, removeElement := range mp.Removes {
		removePolicy2018 := scte224_2018.Remove{
			XMLName: removeElement.XMLName,
		}
		policy2018, err := removeElement.Policy.Get2018()
		if err != nil {
			return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
		}
		removePolicy2018.Policy = policy2018

		removePolicies2018 = append(removePolicies2018, &removePolicy2018)
	}

	mp2018.Applys = applyPolicies2018
	mp2018.Removes = removePolicies2018

	return mp2018, nil
}

func (mp *MediaPoint) Get2015() (*scte224_2015.MediaPoint, error) {
	// First downgrade it to 2018
	mp2018, err := mp.Get2018()
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Then downgrade it to 2015, with the existing utility
	mp2015 := convert.DowngradeMediaPoint(*mp2018)
	return &mp2015, nil
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
	XMLName xml.Name     `xml:"http://www.scte.org/schemas/224 Metadata" json:"-"`
	ADI30   *adi30.ADI30 `xml:"http://www.scte.org/schemas/236/2017/core ADI3" json:"-"`
	Nodes   []Any        `xml:",any" json:"values,omitempty"`
}

type Ext struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Ext"`
	Nodes   []Any    `xml:",any" json:"values,omitempty"`
}

type NamespaceCleaner string

func (nc NamespaceCleaner) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{}, nil
}

type Any struct {
	XMLName xml.Name `json:"xmlname"`
	// mapping xmlns to a field that will avoid marshalling a duplicate namespace
	Namespace  NamespaceCleaner `xml:"xmlns,attr"`
	Attributes []xml.Attr       `xml:",any,attr"`
	Value      string           `xml:",innerxml" json:"value"`
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

type AltID struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 AltID" json:"-"`
	Description string   `xml:"description,attr,omitempty" json:"description,omitempty"`
	Value       string   `xml:",chardata" json:"value,omitempty"`
	Type        string   `xml:"type,attr,omitempty" json:"type,omitempty"`
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

//********************* Audience Types *************************//
//Table 11
type Policy struct {
	ReusableType
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224 Policy" json:"-"`
	ViewingPolicys []*ViewingPolicy `xml:"http://www.scte.org/schemas/224 ViewingPolicy,omitempty" json:"viewingPolicys,omitempty"`
}

func (p *Policy) Get2018() (*scte224_2018.Policy, error) {
	encodedPolicy, err := xml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Unmarshal into 2018 struct
	var policy2018 *scte224_2018.Policy
	err = xml.Unmarshal(encodedPolicy, &policy2018)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Apply ViewingPolicy transformations
	viewingPolicies2018 := make([]*scte224_2018.ViewingPolicy, len(p.ViewingPolicys))
	for _, vp := range p.ViewingPolicys {
		vp2018, err := vp.Get2018()
		if err != nil {
			return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
		}
		viewingPolicies2018 = append(viewingPolicies2018, vp2018)
	}
	policy2018.ViewingPolicys = viewingPolicies2018

	return policy2018, nil
}

func (p *Policy) Get2015() (*scte224_2015.Policy, error) {
	// First downgrade it to 2018
	policy2018, err := p.Get2018()
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Then downgrade it to 2015, with the existing utility
	policy2015 := convert.DowngradePolicy(*policy2018)
	return &policy2015, nil
}

//Table 13
type Audience struct {
	ReusableType
	XMLName          xml.Name    `xml:"http://www.scte.org/schemas/224 Audience" json:"-"`
	Match            Match       `xml:"match,attr,omitempty" json:"match,omitempty"`
	Audiences        []*Audience `xml:"http://www.scte.org/schemas/224 Audience,omitempty" json:"audiences,omitempty"`
	AudienceProperty []Any       `xml:",any" json:"audienceProperty,omitempty"`
}

func (aud *Audience) Get2018() (*scte224_2018.Audience, error) {
	encodedAud, err := xml.Marshal(aud)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Unmarshal into 2018 struct
	var aud2018 *scte224_2018.Audience
	err = xml.Unmarshal(encodedAud, &aud2018)
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	return aud2018, nil
}

func (aud *Audience) Get2015() (*scte224_2015.Audience, error) {
	// First downgrade it to 2018
	aud2018, err := aud.Get2018()
	if err != nil {
		return nil, fmt.Errorf("error attempting to downgrade the 2020 element, %+v", err)
	}

	// Then downgrade it to 2015, with the existing utility
	aud2015 := convert.DowngradeAudience(*aud2018)
	return &aud2015, nil
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
