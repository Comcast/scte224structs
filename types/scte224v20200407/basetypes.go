package scte224v20200407

import (
	"encoding/xml"
	"time"

	"github.com/Comcast/scte224structs/convert"
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224_2018 "github.com/Comcast/scte224structs/types/scte224v20180501"
	"github.com/Comcast/scte224structs/types/scte224v20180501/adi30"
)

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

func (idType *IdentifiableType) Get2018() scte224_2018.IdentifiableType {
	destination := scte224_2018.IdentifiableType{}
	if idType == nil {
		return destination
	}

	destination.Id = idType.Id
	destination.Description = idType.Description
	destination.LastUpdated = idType.LastUpdated
	destination.XMLBase = idType.XMLBase

	for _, altID := range idType.AltIDs {
		if altID == nil {
			continue
		}
		altID2018 := &scte224_2018.AltID{
			XMLName:     altID.XMLName,
			Description: altID.Description,
			Value:       altID.Value,
		}
		destination.AltIDs = append(destination.AltIDs, altID2018)
	}

	if idType.Metadata != nil {
		metadata2018 := &scte224_2018.Metadata{
			XMLName: idType.Metadata.XMLName,
			ADI30:   idType.Metadata.ADI30, // pointer copy; this could lead to incoming value being mutated, but caller should expect it since this method has a pointer receiver
		}

		for _, node := range idType.Metadata.Nodes {
			metadata2018.Nodes = append(metadata2018.Nodes, node.Get2018())
		}

		destination.Metadata = metadata2018
	}

	if idType.Ext != nil {
		ext2018 := &scte224_2018.Ext{
			XMLName: idType.Ext.XMLName,
		}

		for _, node := range idType.Metadata.Nodes {
			ext2018.Nodes = append(ext2018.Nodes, node.Get2018())
		}

		destination.Ext = ext2018
	}

	return destination
}

func (idType *IdentifiableType) Get2015() scte224_2015.IdentifiableType {
	idType2018 := idType.Get2018()
	return convert.DowngradeIdentifiableType(idType2018)
}

//Table 5
type ReusableType struct {
	IdentifiableType
	XLinkHRef string `xml:"http://www.w3.org/1999/xlink href,attr,omitempty" json:"href,omitempty"`
}

func (rt ReusableType) Get2018() scte224_2018.ReusableType {
	return scte224_2018.ReusableType{
		IdentifiableType: rt.IdentifiableType.Get2018(),
		XLinkHRef:        rt.XLinkHRef,
	}
}

func (rt ReusableType) Get2015() scte224_2015.ReusableType {
	rt2018 := rt.Get2018()
	return convert.DowngradeReusableType(rt2018)
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

type AltID struct {
	XMLName     xml.Name `xml:"http://www.scte.org/schemas/224 AltID" json:"-"`
	Description string   `xml:"description,attr,omitempty" json:"description,omitempty"`
	Value       string   `xml:",chardata" json:"value,omitempty"`
	Type        string   `xml:"type,attr,omitempty" json:"type,omitempty"`
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

func (any Any) Get2018() scte224_2018.Any {
	node2018 := scte224_2018.Any{
		XMLName:    any.XMLName,
		Namespace:  scte224_2018.NamespaceCleaner(any.Namespace),
		Attributes: any.Attributes,
		Value:      any.Value,
	}

	return node2018
}
