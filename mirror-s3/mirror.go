package main

import (
	"bytes"
	"code.comcast.com/jcolwe200/scte224/altcon_ds_client"
	"code.comcast.com/jcolwe200/scte224/go-xsd-generated-types/www.scte.org/schemas/224/2015/SCTE224.xsd_go"
	"encoding/xml"
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
	"time"
)

var mpx_user, mpx_password, account, bucket, base_path string
var lastPoll time.Time
var notification int
var s3svc *s3.S3

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "S3 Mirror", log.LstdFlags|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	flag.Parse()
	if "" == mpx_user || "" == mpx_password || "" == bucket {
		flag.Usage()
	} else {
		initS3()
		guid_spew := make(chan scte224DSClient.TypedGuid)
		client := scte224DSClient.SetCredentials(mpx_user, mpx_password, scte224DSClient.Prod)
		go client.PollForNotifications(account, "Mirror", guid_spew, notification)
		go scanForMissedNotifications(client)
		for {
			select {
			case tg := <-guid_spew:
				logger.Println("received notification for ", tg)
				go mirrorGuid(client, tg)
			}
		}
	}
}

func init() {
	flag.StringVar(&mpx_user, "u", "", "MPX Username")
	flag.StringVar(&mpx_password, "p", "", "MPX Password")
	flag.StringVar(&account, "a", "2687225861", "Account")
	flag.IntVar(&notification, "n", 0, "Starting Notification")
	flag.StringVar(&base_path, "f", "", "Base Path for S3 Keys")
	flag.StringVar(&bucket, "b", "", "S3 Bucket")
}

func scanForMissedNotifications(client scte224DSClient.AltContentClient) {
	for {
		// just to ensure a bit of overlap in polling subtracting a second from when we claim this poll was
		thisPoll := time.Now().Add(-time.Second)
		if !lastPoll.IsZero() {
			recentlyUpdated, err := client.GetRecentlyUpdated(account, lastPoll)
			if nil == err {
				for guid, updated := range recentlyUpdated {
					logger.Println("Attempting to re mirror ", guid, updated.Format(time.RFC3339))
					go mirrorGuidIfNewer(client, guid, updated)
				}
			}
		}
		lastPoll = thisPoll
		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func mirrorIfNewer(client scte224DSClient.AltContentClient, tgs ...scte224DSClient.TypedGuid) {
	updatedMap, err := client.GetUpdatedTimestamps(account, tgs...)
	if nil == err {
		for guid, updated := range updatedMap {
			go mirrorGuidIfNewer(client, guid, updated)
		}
	} else {
		logger.Println(err)
	}
}

func mirrorGuidIfNewer(client scte224DSClient.AltContentClient, tg scte224DSClient.TypedGuid, updated time.Time) {
	// find any keys that haven't been updated lately
	head := &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(tg.Guid), IfUnmodifiedSince: aws.Time(updated)}
	_, err := s3svc.HeadObject(head)
	if nil == err {
		// mirror keys that haven't been updated lately
		mirrorGuid(client, tg)
	} else {
		// silently skip if we just didn't get the response due to only wanting responses which haven't been updated lately
		aerr, isawserr := err.(awserr.Error)
		if isawserr {
			switch aerr.Code() {
			case "PreconditionFailed":
			case "NotFound":
				mirrorGuid(client, tg)
			default:
				logger.Println(aerr.Code(), aerr.Message(), aerr.Error())
			}
		} else {
			// log if the error isn't the expected precondition failed error
			logger.Println(err)
		}
	}
}

func mirrorGuid(client scte224DSClient.AltContentClient, tg scte224DSClient.TypedGuid) {
	switch tg.Type {
	case scte224DSClient.MEDIA_SOURCE:
		mirrorMedia(client, tg.Guid)
	case scte224DSClient.POLICY:
		mirrorPolicy(client, tg.Guid)
	case scte224DSClient.VIEWING_POLICY:
		mirrorViewingPolicy(client, tg.Guid)
	case scte224DSClient.AUDIENCE:
		mirrorAudience(client, tg.Guid)
	}
}

func flattenSetToSlice(set map[string]bool) []string {
	var j = 0
	slice := make([]string, len(set))
	for key := range set {
		slice[j] = key
		j++
	}
	return slice
}

func flattenSetToTypedGuidSlice(set map[string]bool, dataType scte224DSClient.DataType) []scte224DSClient.TypedGuid {
	var j = 0
	slice := make([]scte224DSClient.TypedGuid, len(set))
	for key := range set {
		slice[j] = scte224DSClient.TypedGuid{Type: dataType, Guid: key}
		j++
	}
	return slice
}

func mirrorAudience(client scte224DSClient.AltContentClient, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		scteBytes := scteData.Bytes()
		var audience = &go_Scte224.AudiencePayload{}
		err = xml.NewDecoder(scteData).Decode(audience)
		if nil == err {
			nestedAudienceMap := make(map[string]bool)
			for _, nestedAudience := range audience.Audiences {

				href := nestedAudience.Href.String()
				if "" != href && href != audience.Id.String() {
					nestedAudienceMap[href] = true
				} else {
					id := nestedAudience.Id.String()
					if "" != id && id != audience.Id.String() {
						nestedAudienceMap[id] = true
					}
				}
			}
			nestedAudiences := flattenSetToSlice(nestedAudienceMap)
			mirrorIfNewer(client, flattenSetToTypedGuidSlice(nestedAudienceMap, scte224DSClient.AUDIENCE)...)
			logger.Println("Found Audience: ", audience.Id.String(), " containing Nested Audiences: ", nestedAudiences)
			stashInS3(audience.Id.String(), bytes.NewReader(scteBytes))
		}
	} else {
		logger.Println(err)
	}
}

func mirrorViewingPolicy(client scte224DSClient.AltContentClient, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		scteBytes := scteData.Bytes()
		var viewingPolicy = &go_Scte224.ViewingPolicyPayload{}
		err = xml.NewDecoder(scteData).Decode(viewingPolicy)
		if nil == err {
			href := viewingPolicy.Audience.Href.String()
			if "" == href {
				href = viewingPolicy.Audience.Id.String()
			}
			mirrorIfNewer(client, scte224DSClient.TypedGuid{Type: scte224DSClient.AUDIENCE, Guid: href})
			logger.Println("Found Policy: ", viewingPolicy.Id.String(), " containing Audience: ", href)
			stashInS3(viewingPolicy.Id.String(), bytes.NewReader(scteBytes))
		}
	} else {
		logger.Println(err)
	}
}

func mirrorPolicy(client scte224DSClient.AltContentClient, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		scteBytes := scteData.Bytes()
		var policy = &go_Scte224.PolicyPayload{}
		err = xml.NewDecoder(scteData).Decode(policy)
		if nil == err {
			scanForViewingPolicies(client, policy)
			stashInS3(policy.Id.String(), bytes.NewReader(scteBytes))
		}
	} else {
		logger.Println(err)
	}
}

func scanForViewingPolicies(client scte224DSClient.AltContentClient, policy *go_Scte224.PolicyPayload) {
	viewingPoliciesMap := make(map[string]bool)
	for _, viewingPolicy := range policy.ViewingPolicies {

		href := viewingPolicy.Href.String()
		if "" != href {
			viewingPoliciesMap[href] = true
		} else {
			id := viewingPolicy.Id.String()
			if "" != id {
				viewingPoliciesMap[id] = true
			}
		}
	}
	vps := flattenSetToSlice(viewingPoliciesMap)
	logger.Println("Found Policy: ", policy.Id.String(), " containing Viewing Policies: ", vps)
	mirrorIfNewer(client, flattenSetToTypedGuidSlice(viewingPoliciesMap, scte224DSClient.VIEWING_POLICY)...)
}

func mirrorMedia(client scte224DSClient.AltContentClient, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		scteBytes := scteData.Bytes()
		var media = &go_Scte224.MediaPayload{}
		err = xml.NewDecoder(scteData).Decode(media)
		if nil == err {
			scanForPolicies(client, media)
			stashInS3(media.Id.String(), bytes.NewReader(scteBytes))
		}
	} else {
		logger.Println(err)
	}
}

func scanForPolicies(client scte224DSClient.AltContentClient, payload *go_Scte224.MediaPayload) {
	policiesMap := make(map[string]bool)
	for _, point := range payload.MediaPoints {
		for _, apply := range point.Applies {
			href := apply.Policy.Href.String()
			if "" != href {
				policiesMap[href] = true
			} else {
				id := apply.Policy.Id.String()
				if "" != id {
					policiesMap[id] = true
				}
			}
		}
		for _, removes := range point.Removes {
			href := removes.Policy.Href.String()
			if "" != href {
				policiesMap[href] = true
			} else {
				id := removes.Policy.Id.String()
				if "" != id {
					policiesMap[id] = true
				}
			}
		}
	}
	policies := flattenSetToSlice(policiesMap)
	logger.Println("Found media: ", payload.Id.String(), " containing Policies: ", policies)
	mirrorIfNewer(client, flattenSetToTypedGuidSlice(policiesMap, scte224DSClient.POLICY)...)
}

func initS3() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	))

	s3svc = s3.New(sess)
}

func stashInS3(path string, reader io.ReadSeeker) {

	input := s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(path), Body: reader, ACL: aws.String("public-read"), ContentType: aws.String("application/xml; charset=UTF-8"), CacheControl: aws.String("public proxy-revalidate s-maxage=60")}
	putResp, err := s3svc.PutObject(&input)
	if nil == err {
		logger.Println("Mirrored ", path, " to S3")
		logger.Println(putResp)
	} else {
		logger.Println(err)
	}
}
