package adi30

// Package adi11 provides  structs for working with ADI 1.1 metadata files.
// The CableLabs® Asset Distribution Interface Specification Version 1.1
// (MD-SP-ADI1.1-C01-120803) provides a DTD to describe the file format.
//
// Within an ADI 1.1 file, App_Data and Content elements both have an attribute
// named "Value". In these cases, the field is named ValueAttr to prevent duplication
// since the encoding/xml package uses "Value" within the Attr struct to represent
// the "value" in the attribute.

import "encoding/xml"

//ADI30 Top level element of an ADI 1.1 metadata file.
type ADI30 struct {
	XMLName xml.Name `xml:"http://www.scte.org/schemas/236/2017/core ADI3" json:"-"`
	/*
		Offer   string   `xml:"http://www.scte.org/schemas/236/2017/offer offer,attr"`
		Terms   string   `xml:"http://www.scte.org/schemas/236/2017/termmmmms terms,attr" json:"-"`
		Title   string   `xml:"http://www.scte.org/schemas/236/2017/title title,attr"`
		Content string   `xml:"http://www.scte.org/schemas/236/2017/content content,attr" json:"-"`
		Xsi     string   `xml:"http://www.w3.org/2001/XMLSchema-instance xsi,attr" json:"-"`
		//Xmlns   string   `xml:"xmlns,attr" json:"-"`
	*/
	ContentNamespace ContentXSIPrefix `xml:"explicitContentNamespace,attr"`
	TitleNamespace   TitleXSIPrefix   `xml:"explicitTitleNamespace,attr"`
	Asset []*Asset `xml:"Asset,omitempty" json:"asset,omitempty"`
	//Metadata *Metadata `xml:"Metadata"`
}


type TitleXSIPrefix struct{}

func (t TitleXSIPrefix) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name: xml.Name{
			Space: "",
			Local: "xmlns:title",
		},
		Value: "http://www.scte.org/schemas/236/2017/title",
	}, nil
}

type ContentXSIPrefix struct {}

func (c ContentXSIPrefix) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name: xml.Name{
			Space: "",
			Local: "xmlns:content",
		},
		Value: "http://www.scte.org/schemas/236/2017/content",
	}, nil
}

// Metadata elements are containers for a single AMS element and
// zero or more App_Data elements.
type Metadata struct {
	XMLName xml.Name   `xml:"Metadata"`
	Ams     *AMS       `xml:"AMS"`
	AppData []*AppData `xml:"App_Data,omitempty"`
}

// An AMS element typically has a class of package, title, movie or poster (or box-cover).
// The CableLabs specification indicates other @Asset_Class values including
// preview, trickfile, encrypted and barker.
type AMS struct {
	XMLName      xml.Name `xml:"AMS"`
	Product      string   `xml:"Product,attr"`
	AssetID      string   `xml:"Asset_ID,attr"`
	ProviderID   string   `xml:"Provider_ID,attr"`
	Provider     string   `xml:"Provider,attr"`
	AssetClass   string   `xml:"Asset_Class,attr"`
	AssetName    string   `xml:"Asset_Name,attr"`
	Description  string   `xml:"Description,attr"`
	Verb         string   `xml:"Verb,attr,omitempty"`
	VersionMinor string   `xml:"Version_Minor,attr"`
	VersionMajor string   `xml:"Version_Major,attr"`
	CreationDate string   `xml:"Creation_Date,attr"`
	//Value        string   `xml:",chardata"`
}

// AppData (App_Data) specifies additional metadata not included in the AMS element.
// Access to the contents of the "Value" attribute is provided via ValueAttr.
type AppData struct {
	XMLName   xml.Name `xml:"App_Data"`
	App       string   `xml:"App,attr,omitempty"`
	Name      string   `xml:"Name,attr,omitempty"`
	ValueAttr string   `xml:"Value,attr,omitempty"`
	//Value     string   `xml:",chardata"`
}

// The Asset element is a container for Metadata and Content elements and may
// also contain child Asset elements
type Asset struct {
	//Type                 string `xml:"type,attr.omitempty"`
	Type                 string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	UriId                string `xml:"uriId,attr"`
	ProviderVersionNum   string `xml:"providerVersionNum,attr"`
	InternalVersionNum   string `xml:"internalVersionNum,attr"`
	CreationDateTime     string `xml:"creationDateTime,attr"`
	StartDateTime        string `xml:"startDateTime,attr"`
	EndDateTime          string `xml:"endDateTime,attr"`
	LastModifiedDateTime string `xml:"lastModifiedDateTime,attr"`

	AlternateId                *altId          `xml:"AlternateId,omitempty"`
	ProviderQAContact          string          `xml:"ProviderQAContact,omitempty"`
	AssetName                  *deprecAndValue `xml:"AssetName,omitempty"`
	Provider                   string          `xml:"Provider,omitempty"`
	Description                *deprecAndValue `xml:"Description,omitempty"`
	OfrPres                    *presentation   `xml:"http://www.scte.org/schemas/236/2017/offer Presentation,omitempty"`
	PromotionalContentGroupRef *UriId          `xml:"http://www.scte.org/schemas/236/2017/offer PromotionalContentGroupRef,omitempty"`
	ProviderContentTier        string          `xml:"http://www.scte.org/schemas/236/2017/offer ProviderContentTier,omitempty"`
	SourceMetadataSpecVersion  *deprecAndValue `xml:"http://www.scte.org/schemas/236/2017/offer SourceMetadataSpecVersion,omitempty"`
	BillingId                  string          `xml:"http://www.scte.org/schemas/236/2017/offer BillingId,omitempty"`
	TermsRef                   *UriId          `xml:"http://www.scte.org/schemas/236/2017/offer TermsRef,omitempty"`
	ContentGroupRef            *UriId          `xml:"http://www.scte.org/schemas/236/2017/offer ContentGroupRef,omitempty"`

	Ext              *Ext      `xml:"Ext,omitempty"`
	LocalizableTitle *LocTitle `xml:"http://www.scte.org/schemas/236/2017/title LocalizableTitle,omitempty"`

	// LocalizableTitle *title `xml:"LocalizableTitle,omitempty"`

	Rating             *Rating `xml:"http://www.scte.org/schemas/236/2017/title Rating,omitempty"`
	IsClosedCaptioning string  `xml:"http://www.scte.org/schemas/236/2017/title IsClosedCaptioning,omitempty"`
	DisplayRunTime     string  `xml:"http://www.scte.org/schemas/236/2017/title DisplayRunTime,omitempty"`
	Year               int     `xml:"http://www.scte.org/schemas/236/2017/title Year,omitempty"`
	Genre              string  `xml:"http://www.scte.org/schemas/236/2017/title Genre,omitempty"`
	ShowType           string  `xml:"http://www.scte.org/schemas/236/2017/title ShowType,omitempty"`

	AudioType           string              `xml:"http://www.scte.org/schemas/236/2017/content AudioType,omitempty"`
	Language            *Language           `xml:"http://www.scte.org/schemas/236/2017/content Language,omitempty"`
	TrickModeRestricted *Trickmodeexclusion `xml:"http://www.scte.org/schemas/236/2017/content TrickModesRestricted,omitempty"`
	TitleRef            *UriId              `xml:"http://www.scte.org/schemas/236/2017/offer TitleRef,omitempty"`
	MovieRef            *UriId              `xml:"http://www.scte.org/schemas/236/2017/offer MovieRef,omitempty"`
	BillingGracePeriod  string              `xml:"http://www.scte.org/schemas/236/2017/terms BillingGracePeriod,omitempty"`
	SuggestedPrice      string              `xml:"http://www.scte.org/schemas/236/2017/terms SuggestedPrice,omitempty"`
	CategoryPath        string              `xml:"http://www.scte.org/schemas/236/2017/offer CategoryPath,omitempty"`
}

type Ext struct {
	App_Data []*ExtAppData `xml:"App_Data,omitempty"`
}

type presentation struct {
	CategoryRef *[]UriId `xml:"http://www.scte.org/schemas/236/2017/offer CategoryRef,omitempty"`
}

type altId struct {
	IdentifierSystem string `xml:"identifierSystem,attr"`
	Value            string `xml:",chardata"`
}

type deprecAndValue struct {
	Deprecated bool   `xml:"deprecated,attr"`
	Value      string `xml:",chardata"`
}

type UriId struct {
	UriId string `xml:"uriId,attr,omitempty"`
	// Value string `xml:"Value,chardata,omitempty"`
}

type ExtAppData struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type Rating struct {
	RatingSystem string `xml:"ratingSystem,attr"`
	Value        string `xml:",chardata"`
}

type Language struct {
	BitStreamMode string `xml:"bitStreamMode,attr"`
	Value         string `xml:",chardata"`
}

type Trickmodeexclusion struct {
	TrickModeExclusion *trickMode `xml:"TrickModeExclusion,omitempty"`
}

type trickMode struct {
	Type string `xml:"type,attr,omitempty"`
}

// Content elements typically specify a file.
// For example <Content  Value="/pid-fxnetworks.com-aid-DDDE0000103928336978.jpg"/>
// Access to the "Value" attribute is provided via ValueAttr
type Content struct {
	XMLName   xml.Name `xml:"Content"`
	ValueAttr string   `xml:"Value,attr"`
	//Value     string   `xml:",chardata"`
}

type LocTitle struct {
	TitleBrief   string `xml:"http://www.scte.org/schemas/236/2017/title TitleBrief,omitempty"`
	TitleMedium  string `xml:"http://www.scte.org/schemas/236/2017/title TitleMedium,omitempty"`
	TitleLong    string `xml:"http://www.scte.org/schemas/236/2017/title TitleLong,omitempty"`
	SummaryShort string `xml:"http://www.scte.org/schemas/236/2017/title SummaryShort,omitempty"`
	ActorDisplay string `xml:"http://www.scte.org/schemas/236/2017/title ActorDisplay,omitempty"`
}
