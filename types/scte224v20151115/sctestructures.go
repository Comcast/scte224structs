package scte224v20151115

import (
	"encoding/xml"
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
	AltIds      []AltId    `xml:"http://www.scte.org/schemas/224/2015 AltID,omitempty"`
	Metadata    *Metadata  `xml:"http://www.scte.org/schemas/224/2015 Metadata,omitempty"`
	Ext         string     `xml:"##other Ext,omitempty"`
}

//Table 5
type ReusableableType struct {
	IdentifiableType
	XLinkHRef string `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	//XLinkHRef string `xml:"xlink:href,attr,omitempty"`
}

//********************* Media Types *************************//
//Table 6
type Media struct {
	ReusableableType

	XMLName     xml.Name     `xml:"http://www.scte.org/schemas/224/2015 Media"`
	Effective   time.Time    `xml:"effective,attr,omitempty"`
	Expires     time.Time    `xml:"expires,attr,omitempty"`
	Source      string       `xml:"source,attr,omitempty"`
	MediaPoints []MediaPoint `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
}

// MediaPoint defines an SCTE 224 (ESNI) media point object.
//Table 7
type MediaPoint struct {
	IdentifiableType
	XMLName      xml.Name      `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
	Effective    time.Time     `xml:"effective,attr,omitempty"`
	Expires      time.Time     `xml:"expires,attr,omitempty"`
	MatchTime    time.Time     `xml:"matchTime,attr,omitempty"`
	MatchOffset  time.Duration `xml:"matchOffset,attr,omitempty"`
	Source       string        `xml:"source,attr,omitempty"`
	Removes      []Remove      `xml:"http://www.scte.org/schemas/224/2015 Remove"`
	Applys       []Apply       `xml:"http://www.scte.org/schemas/224/2015 Apply"`
	MatchSignals []MatchSignal `xml:"http://www.scte.org/schemas/224/2015 MatchSignal"`
	MediaGuid    string        `xml:"-"` // used internally to track which media this point is part of
}
type AltId struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224/2015 AltID"`
	Value   string   `xml:",chardata"`
}

type Metadata struct {
	XMLName        xml.Name         `xml:"http://www.scte.org/schemas/224/2015 Metadata"`
	MetadataDetail []MetadataDetail `xml:"http://ctsrmm.com/ctsesni MetadataDetail,omitempty"`
	ScheduledStart *time.Time       `xml:"http://thistech.com/esni ScheduledStart,omitempty"`
	ScheduledEnd   *time.Time       `xml:"http://thistech.com/esni ScheduledEnd,omitempty"`

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
	AScheduledStart    *time.Time `xml:"http://www.scte.org/schemas/224/2015 ScheduledStart,omitempty"`
	AScheduledEnd      *time.Time `xml:"http://www.scte.org/schemas/224/2015 ScheduledEnd,omitempty"`
	ProgramTitle       string     `xml:"http://www.scte.org/schemas/224/2015 ProgramTitle,omitempty"`
	ProgramDescription string     `xml:"http://www.scte.org/schemas/224/2015 ProgramDescription,omitempty"`
	EpisodeTitle       string     `xml:"http://www.scte.org/schemas/224/2015 EpisodeTitle,omitempty"`
	EpisodeDescription string     `xml:"http://www.scte.org/schemas/224/2015 EpisodeDescription,omitempty"`
}

//Normalized stuff
type MetadataDetail struct {
	Name           string           `xml:"name,attr,omitempty"`
	Type           string           `xml:"type,attr,omitempty"`
	Provider       string           `xml:"provider,attr,omitempty"`
	Value          string           `xml:",chardata"`
	MetadataDetail []MetadataDetail `xml:"http://ctsrmm.com/ctsesni MetadataDetail,omitempty"`
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
	XMLName          xml.Name      `xml:"http://schema.foxinc.com/esni MediaMetaData"`
	Network          FoxAttributes `xml:"http://schema.foxinc.com/esni Network,omitempty"`
	Region           FoxAttributes `xml:"http://schema.foxinc.com/esni Region,omitempty"`
	DistributionArea FoxAttributes `xml:"http://schema.foxinc.com/esni DistributionArea"`
}

type ChannelMetadata struct {
	XMLName          xml.Name      `xml:"http://schema.foxinc.com/esni ChannelMetadata"`
	AudienceType     string        `xml:"http://schema.foxinc.com/esni AudienceType,omitempty"`
	LocationType     string        `xml:"http://schema.foxinc.com/esni LocationType,omitempty"`
	Network          FoxAttributes `xml:"http://schema.foxinc.com/esni Network,omitempty"`
	Region           FoxAttributes `xml:"http://schema.foxinc.com/esni Region,omitempty"`
	DistributionArea FoxAttributes `xml:"http://schema.foxinc.com/esni DistributionArea,omitempty"`
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
	XMLName            xml.Name           `xml:"http://schema.foxinc.com/esni DeliveryRestrictions,omitempty"`
	WebEnabled         bool               `xml:"http://schema.foxinc.com/esni WebEnabled,omitempty"`
	DeviceRestrictions DeviceRestrictions `xml:"http://schema.foxinc.com/esni DeviceRestrictions,omitempty"`
}

type DeviceRestrictions struct {
	XMLName xml.Name `xml:"http://schema.foxinc.com/esni DeviceRestrictions,omitempty"`
	Match   string   `xml:"http://schema.foxinc.com/esni match,attr,omitempty"`
	Devices []Item   `xml:"http://schema.foxinc.com/esni Device,omitempty"`
}

//Table 10
type Apply struct {
	XMLName  xml.Name `xml:"http://www.scte.org/schemas/224/2015 Apply"`
	Duration string   `xml:"duration,attr,omitempty"`
	Policys  []Policy `xml:"http://www.scte.org/schemas/224/2015 Policy,omitempty"`
}

//Table 9
type Remove struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224/2015 Remove"`
	Policys []Policy `xml:"http://www.scte.org/schemas/224/2015 Policy",omitempty`
}

//Table 8
type MatchSignal struct {
	XMLName         xml.Name `xml:"http://www.scte.org/schemas/224/2015 MatchSignal"`
	Match           string   `xml:"match,attr,omitempty"`
	SignalTolerance string   `xml:"signalTolerance,attr,omitempty"`
	Assertions      []Assert `xml:"http://www.scte.org/schemas/224/2015 Assert,omitempty"`
}

type Assert struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224/2015 Assert"`
	Declaration string   `xml:",chardata"`
}

//********************* Audience Types *************************//
//Table 11
type Policy struct {
	ReusableableType
	XMLName        xml.Name        `xml:"http://www.scte.org/schemas/224/2015 Policy"`
	ViewingPolicys []ViewingPolicy `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy,omitempty"`
	Metadata       string          `xml:"http://www.scte.org/schemas/224/2015 Metadata,omitempty"`
}

//Table 12
type ViewingPolicy struct {
	ReusableableType
	XMLName  xml.Name   `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy"`
	Audience *Audience  `xml:"http://www.scte.org/schemas/224/2015 Audience,omitempty"`
	Anys     []Any      `xml:"http://www.scte.org/schemas/224/2015 Any,omitempty"`
	Contents []AContent `xml:"urn:scte:224:action Content,omitempty"`
	FF       []Item     `xml:"urn:scte:224:action FastForward,omitempty"`
	Capture  []Capture  `xml:"urn:scte:224:action Capture,omitempty"`
	Metadata string     `xml:"http://www.scte.org/schemas/224/2015 Metadata,omitempty"`
}

//Table 13
type Audience struct {
	ReusableableType
	XMLName        xml.Name   `xml:"http://www.scte.org/schemas/224/2015 Audience"`
	Match          string     `xml:"match,attr,omitempty"`
	Metadata       *Metadata  `xml:"http://www.scte.org/schemas/224/2015 Metadata,omitempty"`
	Audiences      []Audience `xml:"http://www.scte.org/schemas/224/2015 Audience,omitempty"`
	Virds          []Item     `xml:"urn:scte:224:audience Vird,omitempty"`
	Zips           []Item     `xml:"urn:scte:224:audience Zip,omitempty"`
	Features       []Item     `xml:"urn:scte:224:audience Feature,omitempty"`
	DeviceFeatures []Item     `xml:"urn:scte:224:audience DeviceFeature,omitempty"`
	Devices        []Item     `xml:"urn:scte:224:audience Device,omitempty"`
	FoxDevices     []Item     `xml:"http://schema.foxinc.com/esni Device,omitempty"`
	Networks       []Item     `xml:"urn:scte:224:audience Network,omitempty"`
	Default        []Item     `xml:"urn:scte:224:audience Default,omitempty"`
}

type Capture struct {
	XMLName     xml.Name     `xml:"urn:scte:224:action Capture"`
	StartWindow *StartWindow `xml:"urn:scte:224:action StartWindow,omitempty"`
	StopWindow  *StopWindow  `xml:"urn:scte:224:action StopWindow,omitempty"`
	MidrollDAI  *Item        `xml:"urn:scte:224:action MidrollDAI,omitempty"`
}

type StartWindow struct {
	Percentage Item `xml:"urn:scte:224:action Percentage,omitempty"`
}

type StopWindow struct {
	Offset     Item `xml:"urn:scte:224:action Offset,omitempty"`
	Percentage Item `xml:"urn:scte:224:action Percentage,omitempty"`
}

type AContent struct {
	XMLName xml.Name `xml:"urn:scte:224:action Content"`
	Value   string   `xml:",chardata"`
}
type Item struct {
	Value string `xml:",chardata"`
}

type Any struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/224/2015 Any"`
	Value   string   `xml:",chardata"`
}

//********************* Results Types *************************//
//Table 14
type Results struct {
	XMLName        xml.Name        `xml:"http://www.scte.org/schemas/224/2015 Results"`
	Size           int             `xml:"size,attr,omitempty"`
	Medias         []Media         `xml:"http://www.scte.org/schemas/224/2015 Media"`
	MediaPoints    []MediaPoint    `xml:"http://www.scte.org/schemas/224/2015 MediaPoint"`
	Policys        []Policy        `xml:"http://www.scte.org/schemas/224/2015 Policy"`
	ViewingPolicys []ViewingPolicy `xml:"http://www.scte.org/schemas/224/2015 ViewingPolicy"`
	Audiences      []Audience      `xml:"http://www.scte.org/schemas/224/2015 Audience"`
	Audits         []Audit         `xml:"http://www.scte.org/schemas/224/2015 Audit"`
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
	Audits        []Audit  `xml:"http://www.scte.org/schemas/224/2015 Audit"`
}
