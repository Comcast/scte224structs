package convert

import (
	scte224_2015 "github.com/Comcast/scte224structs/types/scte224v20151115"
	scte224 "github.com/Comcast/scte224structs/types/scte224v20180501"
)

func DowngradeIdentifiableType(identifiableType scte224.IdentifiableType) scte224_2015.IdentifiableType {
	var dst scte224_2015.IdentifiableType
	dst.Id = identifiableType.Id
	for _, altId := range identifiableType.AltIDs {
		if nil != altId {
			dst.AltIDs = append(dst.AltIDs, &scte224_2015.AltID{XMLName: altId.XMLName, Value: altId.Value})
		}
	}
	dst.Description = identifiableType.Description
	if identifiableType.Ext != nil {
		nodes := make([]scte224_2015.Any, len(identifiableType.Ext.Nodes))
		for j, node := range identifiableType.Ext.Nodes {
			nodes[j] = scte224_2015.Any{
				XMLName:    node.XMLName,
				Attributes: node.Attributes,
				Value:      node.Value,
			}
		}
		dst.Ext = &scte224_2015.Ext{
			XMLName: identifiableType.Ext.XMLName,
			Nodes:   nodes,
		}
	}
	if identifiableType.Metadata != nil {
		nodes := make([]scte224_2015.Any, len(identifiableType.Metadata.Nodes))
		for j, node := range identifiableType.Metadata.Nodes {
			nodes[j] = scte224_2015.Any{
				XMLName:    node.XMLName,
				Attributes: node.Attributes,
				Value:      node.Value,
			}
		}

		dst.Metadata = &scte224_2015.Metadata{
			XMLName: identifiableType.Metadata.XMLName,
			Nodes:   nodes,
		}
	}

	dst.LastUpdated = identifiableType.LastUpdated
	dst.XMLBase = identifiableType.XMLBase
	return dst
}

func DowngradeReusableType(reusableType scte224.ReusableType) scte224_2015.ReusableType {
	var dst scte224_2015.ReusableType
	dst.IdentifiableType = DowngradeIdentifiableType(reusableType.IdentifiableType)
	dst.XLinkHRef = reusableType.XLinkHRef
	return dst
}

func DowngradeAudience(audience scte224.Audience) scte224_2015.Audience {
	var dst scte224_2015.Audience
	dst.XMLName = audience.XMLName
	for _, audienceProp := range audience.AudienceProperty {
		dst.AudienceProperty = append(dst.AudienceProperty, scte224_2015.Any{XMLName: audienceProp.XMLName, Value: audienceProp.Value})
	}
	for _, nestedAudience := range audience.Audiences {
		if nil != nestedAudience {
			downer := DowngradeAudience(*nestedAudience)
			dst.Audiences = append(dst.Audiences, &downer)
		}
	}
	dst.Match = scte224_2015.Match(audience.Match)
	dst.ReusableType = DowngradeReusableType(audience.ReusableType)
	return dst
}

func DowngradeViewingPolicy(vp scte224.ViewingPolicy) scte224_2015.ViewingPolicy {
	var dst scte224_2015.ViewingPolicy
	dst.XMLName = vp.XMLName
	if content := vp.Content; nil != content {
		dst.ActionProperty = append(dst.ActionProperty, scte224_2015.Any{
			XMLName: content.XMLName,
			Value:   content.Content,
		})
	}
	if deletion := vp.SignalPointDeletion; nil != deletion {
		dst.ActionProperty = append(dst.ActionProperty, scte224_2015.Any{
			XMLName: deletion.XMLName,
			Value:   deletion.SignalPointDeletion,
		})
	}
	// Downgrading signal point insertion is not supported

	for _, actionProp := range vp.ActionProperty {
		dst.ActionProperty = append(dst.ActionProperty, scte224_2015.Any{XMLName: actionProp.XMLName, Value: actionProp.Value})
	}
	if vp.Audience != nil {
		downgradedAudience := DowngradeAudience(*vp.Audience)
		dst.Audience = &downgradedAudience
	}
	dst.ReusableType = DowngradeReusableType(vp.ReusableType)
	return dst
}

func DowngradePolicy(p scte224.Policy) scte224_2015.Policy {
	var dst scte224_2015.Policy
	dst.XMLName = p.XMLName
	for _, vp := range p.ViewingPolicys {
		downgradedViewingPolicy := DowngradeViewingPolicy(*vp)
		dst.ViewingPolicys = append(dst.ViewingPolicys, &downgradedViewingPolicy)
	}
	dst.ReusableType = DowngradeReusableType(p.ReusableType)
	return dst
}

// stupid helper function to reduce the cut & paste of handling the policy pointers safely
func DowngradePolicyPointer(p *scte224.Policy) *scte224_2015.Policy {
	if nil == p {
		return nil
	}
	downgradedPolicy := DowngradePolicy(*p)
	return &downgradedPolicy
}

func DowngradeMediaPoint(point scte224.MediaPoint) scte224_2015.MediaPoint {
	var dst scte224_2015.MediaPoint
	dst.IdentifiableType = DowngradeIdentifiableType(point.IdentifiableType)
	dst.XMLName = point.XMLName
	dst.Effective = point.Effective
	dst.Expires = point.Expires
	dst.MatchTime = point.MatchTime

	// duration is just a string but with utility functions to aid conversion
	dst.MatchOffset = scte224_2015.Duration(point.MatchOffset)
	dst.Source = point.Source
	dst.ExpectedDuration = scte224_2015.Duration(point.ExpectedDuration)
	dst.Order = point.Order
	dst.Reusable = point.Reusable
	for _, remove := range point.Removes {
		if nil != remove {
			dst.Removes = append(dst.Removes, &scte224_2015.Remove{XMLName: remove.XMLName, Policy: DowngradePolicyPointer(remove.Policy)})
		}
	}
	for _, apply := range point.Applys {
		if nil != apply {
			dst.Applys = append(dst.Applys, &scte224_2015.Apply{XMLName: apply.XMLName, Policy: DowngradePolicyPointer(apply.Policy), Duration: scte224_2015.Duration(apply.Duration), Priority: apply.Priority})
		}
	}
	if nil != point.MatchSignal {
		dst.MatchSignal = &scte224_2015.MatchSignal{XMLName: point.MatchSignal.XMLName, Match: scte224_2015.Match(point.MatchSignal.Match), SignalTolerance: scte224_2015.Duration(point.MatchSignal.SignalTolerance)}
		for _, assertion := range point.MatchSignal.Assertions {
			if nil != assertion {
				dst.MatchSignal.Assertions = append(dst.MatchSignal.Assertions, &scte224_2015.Assert{XMLName: assertion.XMLName, Declaration: assertion.Declaration})
			}
		}
	}
	dst.MediaGuid = point.MediaGuid
	return dst
}

func DowngradeMedia(media scte224.Media) scte224_2015.Media {
	var dst scte224_2015.Media
	dst.ReusableType = DowngradeReusableType(media.ReusableType)
	dst.XMLName = media.XMLName
	dst.Effective = media.Effective
	dst.Expires = media.Expires
	dst.Source = media.Source
	for _, mp := range media.MediaPoints {
		if nil != mp {
			downgradedPoint := DowngradeMediaPoint(*mp)
			dst.MediaPoints = append(dst.MediaPoints, &downgradedPoint)
		}
	}
	return dst
}
