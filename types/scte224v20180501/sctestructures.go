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
	Ext         string     `xml:"##other Ext,omitempty"`
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

type Metadata struct {
	XMLName        xml.Name          `xml:"http://www.scte.org/schemas/224 Metadata"`
	MetadataDetail []*MetadataDetail `xml:"http://ctsrmm.com/ctsesni MetadataDetail,omitempty"`
	ScheduledStart *time.Time        `xml:"http://thistech.com/esni ScheduledStart,omitempty"`
	ScheduledEnd   *time.Time        `xml:"http://thistech.com/esni ScheduledEnd,omitempty"`

	//TODO: get this customer specific stuff out of here
	//fox stuff
	FoxMediaMetadata      *FoxMediaMetadata      `xml:"http://schema.foxinc.com/esni MediaMetaData,omitempty"`
	ChannelMetadata       *ChannelMetadata       `xml:"http://schema.foxinc.com/esni ChannelMetadata,omitempty"`
	DeviceMetadata        *DeviceMetadata        `xml:"http://schema.foxinc.com/esni DeviceMetadata,omitempty"`
	LocationMetadata      *LocationMetadata      `xml:"http://schema.foxinc.com/esni LocationMetadata,omitempty"`
	FoxMediaPointMetadata *FoxMediaPointMetadata `xml:"http://schema.foxinc.com/esni MediaPointMetaData,omitempty"`

	//nbc stuff
	HomeTeam   string `xml:"http://xml.nbcuni.com/listing/field homeTeam,omitempty"`
	HomeTeamId string `xml:"http://xml.nbcuni.com/listing/field homeTeamId,omitempty"`
	AwayTeam   string `xml:"http://xml.nbcuni.com/listing/field awayTeam,omitempty"`
	AwayTeamId string `xml:"http://xml.nbcuni.com/listing/field awayTeamId,omitempty"`
	Sport      string `xml:"http://xml.nbcuni.com/listing/field sport,omitempty"`

	//A&E stuff
	AScheduledStart    *time.Time `xml:"http://www.scte.org/schemas/224 ScheduledStart,omitempty"`
	AScheduledEnd      *time.Time `xml:"http://www.scte.org/schemas/224 ScheduledEnd,omitempty"`
	ProgramTitle       string     `xml:"http://www.scte.org/schemas/224 ProgramTitle,omitempty"`
	ProgramDescription string     `xml:"http://www.scte.org/schemas/224 ProgramDescription,omitempty"`
	EpisodeTitle       string     `xml:"http://www.scte.org/schemas/224 EpisodeTitle,omitempty"`
	EpisodeDescription string     `xml:"http://www.scte.org/schemas/224 EpisodeDescription,omitempty"`
}

//Normalized stuff
type MetadataDetail struct {
	Name           string            `xml:"name,attr,omitempty"`
	Type           string            `xml:"type,attr,omitempty"`
	Provider       string            `xml:"provider,attr,omitempty"`
	Value          string            `xml:",chardata"`
	MetadataDetail []*MetadataDetail `xml:"http://ctsrmm.com/ctsesni MetadataDetail,omitempty"`
}

//FOX stuff
type FoxMediaPointMetadata struct {
	ScheduledAiringId    string                `xml:"http://schema.foxinc.com/esni ScheduledAiringId,omitempty"`
	ScheduledStart       *time.Time            `xml:"http://schema.foxinc.com/esni ScheduledStart,omitempty"`
	ScheduledEnd         *time.Time            `xml:"http://schema.foxinc.com/esni ScheduledEnd,omitempty"`
	Category             string                `xml:"http://schema.foxinc.com/esni Category,omitempty"`
	ProgramId            string                `xml:"http://schema.foxinc.com/esni ProgramId,omitempty"`
	ProgramTitle         string                `xml:"http://schema.foxinc.com/esni ProgramTitle,omitempty"`
	ProgramDescription   string                `xml:"http://schema.foxinc.com/esni ProgramDescription,omitempty"`
	SeasonId             string                `xml:"http://schema.foxinc.com/esni SeasonId,omitempty"`
	SeasonTitle          string                `xml:"http://schema.foxinc.com/esni SeasonTitle,omitempty"`
	EpisodeId            string                `xml:"http://schema.foxinc.com/esni EpisodeId,omitempty"`
	EpisodeTitle         string                `xml:"http://schema.foxinc.com/esni EpisodeTitle,omitempty"`
	EpisodeDescription   string                `xml:"http://schema.foxinc.com/esni EpisodeDescription,omitempty"`
	MaterialId           string                `xml:"http://schema.foxinc.com/esni MaterialId,omitempty"`
	Rating               string                `xml:"http://schema.foxinc.com/esni Rating,omitempty"`
	Live                 string                `xml:"http://schema.foxinc.com/esni Live,omitempty"`
	StartOver            string                `xml:"http://schema.foxinc.com/esni StartOver,omitempty"`
	LookBack             string                `xml:"http://schema.foxinc.com/esni LookBack,omitempty"`
	DeliveryRestrictions *DeliveryRestrictions `xml:"http://schema.foxinc.com/esni DeliveryRestrictions,omitempty"`
}

type FoxMediaMetadata struct {
	XMLName          xml.Name       `xml:"http://schema.foxinc.com/esni MediaMetaData"`
	Network          *FoxAttributes `xml:"http://schema.foxinc.com/esni Network,omitempty"`
	Region           *FoxAttributes `xml:"http://schema.foxinc.com/esni Region,omitempty"`
	DistributionArea *FoxAttributes `xml:"http://schema.foxinc.com/esni DistributionArea"`
}

type ChannelMetadata struct {
	XMLName          xml.Name       `xml:"http://schema.foxinc.com/esni ChannelMetadata"`
	AudienceType     string         `xml:"http://schema.foxinc.com/esni AudienceType,omitempty"`
	LocationType     string         `xml:"http://schema.foxinc.com/esni LocationType,omitempty"`
	Network          *FoxAttributes `xml:"http://schema.foxinc.com/esni Network,omitempty"`
	Region           *FoxAttributes `xml:"http://schema.foxinc.com/esni Region,omitempty"`
	DistributionArea *FoxAttributes `xml:"http://schema.foxinc.com/esni DistributionArea,omitempty"`
}

type DeviceMetadata struct {
	XMLName      xml.Name `xml:"http://schema.foxinc.com/esni DeviceMetadata"`
	AudienceType string   `xml:"http://schema.foxinc.com/esni AudienceType,omitempty"`
}

type LocationMetadata struct {
	XMLName      xml.Name `xml:"http://schema.foxinc.com/esni LocationMetadata"`
	AudienceType string   `xml:"http://schema.foxinc.com/esni AudienceType,omitempty"`
}

type FoxAttributes struct {
	Name         string `xml:"name,attr,omitempty"`
	Call_sign    string `xml:"call_sign,attr,omitempty"`
	Network_type string `xml:"network_type,attr,omitempty"`
}

type DeliveryRestrictions struct {
	XMLName            xml.Name            `xml:"http://schema.foxinc.com/esni DeliveryRestrictions,omitempty"`
	WebEnabled         bool                `xml:"http://schema.foxinc.com/esni WebEnabled,omitempty"`
	DeviceRestrictions *DeviceRestrictions `xml:"http://schema.foxinc.com/esni DeviceRestrictions,omitempty"`
}

type DeviceRestrictions struct {
	XMLName xml.Name `xml:"http://schema.foxinc.com/esni DeviceRestrictions,omitempty"`
	Match   string   `xml:"http://schema.foxinc.com/esni match,attr,omitempty"`
	Devices []*Item  `xml:"http://schema.foxinc.com/esni Device,omitempty"`
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
	Metadata       string           `xml:"http://www.scte.org/schemas/224 Metadata,omitempty"`
}

//Table 12
type ViewingPolicy struct {
	ReusableType
	XMLName  xml.Name    `xml:"http://www.scte.org/schemas/224 ViewingPolicy"`
	Audience *Audience   `xml:"http://www.scte.org/schemas/224 Audience,omitempty"`
	Anys     []*Any      `xml:"http://www.scte.org/schemas/224 Any,omitempty"`
	Contents []*AContent `xml:"urn:scte:224:action Content,omitempty"`
	FF       []*Item     `xml:"urn:scte:224:action FastForward,omitempty"`
	Capture  []*Capture  `xml:"urn:scte:224:action Capture,omitempty"`
	Metadata string      `xml:"http://www.scte.org/schemas/224 Metadata,omitempty"`
}

//Table 13
type Audience struct {
	ReusableType
	XMLName        xml.Name    `xml:"http://www.scte.org/schemas/224 Audience"`
	Match          string      `xml:"match,attr,omitempty"`
	Metadata       *Metadata   `xml:"http://www.scte.org/schemas/224 Metadata,omitempty"`
	Audiences      []*Audience `xml:"http://www.scte.org/schemas/224 Audience,omitempty"`
	Virds          []*Item     `xml:"urn:scte:224:audience Vird,omitempty"`
	Zips           []*Item     `xml:"urn:scte:224:audience Zip,omitempty"`
	Features       []*Item     `xml:"urn:scte:224:audience Feature,omitempty"`
	DeviceFeatures []*Item     `xml:"urn:scte:224:audience DeviceFeature,omitempty"`
	Devices        []*Item     `xml:"urn:scte:224:audience Device,omitempty"`
	FoxDevices     []*Item     `xml:"http://schema.foxinc.com/esni Device,omitempty"`
	Networks       []*Item     `xml:"urn:scte:224:audience Network,omitempty"`
	Default        []*Item     `xml:"urn:scte:224:audience Default,omitempty"`
}

type Capture struct {
	XMLName     xml.Name     `xml:"urn:scte:224:action Capture"`
	StartWindow *StartWindow `xml:"urn:scte:224:action StartWindow,omitempty"`
	StopWindow  *StopWindow  `xml:"urn:scte:224:action StopWindow,omitempty"`
	MidrollDAI  *Item        `xml:"urn:scte:224:action MidrollDAI,omitempty"`
}

type StartWindow struct {
	Percentage *Item `xml:"urn:scte:224:action Percentage,omitempty"`
}

type StopWindow struct {
	Offset     *Item `xml:"urn:scte:224:action Offset,omitempty"`
	Percentage *Item `xml:"urn:scte:224:action Percentage,omitempty"`
}

type AContent struct {
	XMLName xml.Name `xml:"urn:scte:224:action Content"`
	Value   string   `xml:",chardata"`
}
type Item struct {
	Value string `xml:",chardata"`
}

type Any struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224 Any"`
	Value   string   `xml:",chardata"`
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
