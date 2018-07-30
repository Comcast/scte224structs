package convert

import (
	scte224_2015 "github.comcast.com/jcolwe200/scte224/types/scte224v20151115"
	scte224 "github.comcast.com/jcolwe200/scte224/types/scte224v20180501"
)

func UpgradeIdentifiableType(identifiableType scte224_2015.IdentifiableType) scte224.IdentifiableType {
	var dst scte224.IdentifiableType
	dst.Id = identifiableType.Id
	for _, altId := range identifiableType.AltIDs {
		if nil != altId {
			dst.AltIDs = append(dst.AltIDs, &scte224.AltID{XMLName: altId.XMLName, Value: altId.Value})
		}
	}
	dst.Description = identifiableType.Description
	if identifiableType.Ext != nil {
		dst.Ext = &scte224.Metadata{InnerXml: identifiableType.Ext.InnerXml}
	}
	if identifiableType.Metadata != nil {
		dst.Metadata = &scte224.Metadata{InnerXml: identifiableType.Metadata.InnerXml}
	}
	dst.LastUpdated = identifiableType.LastUpdated
	dst.XMLBase = identifiableType.XMLBase
	return dst
}

func UpgradeReusableType(reusableType scte224_2015.ReusableType) scte224.ReusableType {
	var dst scte224.ReusableType
	dst.IdentifiableType = UpgradeIdentifiableType(reusableType.IdentifiableType)
	dst.XLinkHRef = reusableType.XLinkHRef
	return dst
}

func UpgradeAudience(audience scte224_2015.Audience) scte224.Audience {
	var dst scte224.Audience
	dst.XMLName = audience.XMLName
	dst.AudienceProperty = audience.AudienceProperty
	dst.Match = scte224.Match(audience.Match)
	dst.ReusableType = UpgradeReusableType(audience.ReusableType)
	return dst
}

func UpgradeViewingPolicy(vp scte224_2015.ViewingPolicy) scte224.ViewingPolicy {
	var dst scte224.ViewingPolicy
	dst.XMLName = vp.XMLName
	dst.ActionProperty = vp.ActionProperty
	if vp.Audience != nil {
		upgradedAudience := UpgradeAudience(*vp.Audience)
		dst.Audience = &upgradedAudience
	}
	dst.ReusableType = UpgradeReusableType(vp.ReusableType)
	return dst
}

func UpgradePolicy(p scte224_2015.Policy) scte224.Policy {
	var dst scte224.Policy
	dst.XMLName = p.XMLName
	for _, vp := range p.ViewingPolicys {
		upgradedViewingPolicy := UpgradeViewingPolicy(*vp)
		dst.ViewingPolicys = append(dst.ViewingPolicys, &upgradedViewingPolicy)
	}
	dst.ReusableType = UpgradeReusableType(p.ReusableType)
	return dst
}

// stupid helper function to reduce the cut & paste of handling the policy pointers safely
func UpgradePolicyPointer(p *scte224_2015.Policy) *scte224.Policy {
	if nil == p {
		return nil
	}
	upgradedPolicy := UpgradePolicy(*p)
	return &upgradedPolicy
}

func UpgradeMediaPoint(point scte224_2015.MediaPoint) scte224.MediaPoint {
	var dst scte224.MediaPoint
	dst.IdentifiableType = UpgradeIdentifiableType(point.IdentifiableType)
	dst.XMLName = point.XMLName
	dst.Effective = point.Effective
	dst.Expires = point.Expires
	dst.MatchTime = point.MatchTime

	// duration is just a string but with utility functions to aid conversion
	dst.MatchOffset = scte224.Duration(point.MatchOffset)
	dst.Source = point.Source
	dst.ExpectedDuration = scte224.Duration(point.ExpectedDuration)
	dst.Order = point.Order
	dst.Reusable = point.Reusable
	for _, remove := range point.Removes {
		if nil != remove {
			dst.Removes = append(dst.Removes, &scte224.Remove{XMLName: remove.XMLName, Policy: UpgradePolicyPointer(remove.Policy)})
		}
	}
	for _, apply := range point.Applys {
		if nil != apply {
			dst.Applys = append(dst.Applys, &scte224.Apply{XMLName: apply.XMLName, Policy: UpgradePolicyPointer(apply.Policy), Duration: scte224.Duration(apply.Duration), Priority: apply.Priority})
		}
	}
	if nil != point.MatchSignal {
		dst.MatchSignal = &scte224.MatchSignal{XMLName: point.MatchSignal.XMLName, Match: scte224.Match(point.MatchSignal.Match), SignalTolerance: scte224.Duration(point.MatchSignal.SignalTolerance)}
		for _, assertion := range point.MatchSignal.Assertions {
			if nil != assertion {
				dst.MatchSignal.Assertions = append(dst.MatchSignal.Assertions, &scte224.Assert{XMLName: assertion.XMLName, Declaration: assertion.Declaration})
			}
		}
	}
	dst.MediaGuid = point.MediaGuid
	return dst
}

func UpgradeMedia(media scte224_2015.Media) scte224.Media {
	var dst scte224.Media
	dst.ReusableType = UpgradeReusableType(media.ReusableType)
	dst.XMLName = media.XMLName
	dst.Effective = media.Effective
	dst.Expires = media.Expires
	dst.Source = media.Source
	for _, mp := range media.MediaPoints {
		if nil != mp {
			upgradedPoint := UpgradeMediaPoint(*mp)
			dst.MediaPoints = append(dst.MediaPoints, &upgradedPoint)
		}
	}
	return dst
}
