package main

import (
	"code.comcast.com/jcolwe200/scte224/altcon_ds_client"
	"code.comcast.com/jcolwe200/scte224/go-xsd-generated-types/www.scte.org/schemas/224/2015/SCTE224.xsd_go"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

var mpx_user, mpx_password, account string
var notification int
var s3svc *s3.S3

func main() {
	flag.Parse()
	if "" == mpx_user || "" == mpx_password {
		flag.Usage()
	} else {
		guid_spew := make(chan scte224DSClient.TypedGuid)
		client := scte224DSClient.SetCredentials(mpx_user, mpx_password, scte224DSClient.Prod)
		go client.PollForNotifications(account, "Mirror", guid_spew, notification)
		for {
			select {
			case tg := <-guid_spew:
				fmt.Println("received notification for ", tg)
				go mirrorGuid(client, account, tg)
			}
		}
	}
}

func init() {
	flag.StringVar(&mpx_user, "u", "", "MPX Username")
	flag.StringVar(&mpx_password, "p", "", "MPX Password")
	flag.StringVar(&account, "a", "2687225861", "Account")
	flag.IntVar(&notification, "n", 0, "Starting Notification")
}

func mirrorGuid(client scte224DSClient.AltContentClient, account string, tg scte224DSClient.TypedGuid) {
	switch tg.Type {
	case scte224DSClient.MEDIA_SOURCE:
		mirrorMedia(client, account, tg.Guid)
	case scte224DSClient.POLICY:
		mirrorPolicy(client, account, tg.Guid)
	case scte224DSClient.VIEWING_POLICY:
		//mirrorViewingPolicy(client, account, tg.Guid)
	case scte224DSClient.AUDIENCE:
		//mirrorAudience(client, account, tg.Guid)
	}
}

func mirrorPolicy(client scte224DSClient.AltContentClient, account, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		var policy = &go_Scte224.PolicyPayload{}
		err = xml.NewDecoder(scteData).Decode(policy)
		if nil == err {
			log.Println("Found Policy: ", policy.Id.String())
			scanForViewingPolicies(policy)
		}
	} else {
		log.Println(err)
	}
}

func scanForViewingPolicies(policy go_Scte224.PolicyPayload)  {

}

func mirrorMedia(client scte224DSClient.AltContentClient, account, guid string) {
	scteData, err := client.GetSCTEData(account, guid)
	if nil == err {
		var media = &go_Scte224.MediaPayload{}
		err = xml.NewDecoder(scteData).Decode(media)
		if nil == err {
			log.Println("Found Media: ", media.Id.String())
			scanForPolicies(media)
		}
	} else {
		log.Println(err)
	}
}

func scanForPolicies(payload *go_Scte224.MediaPayload) {
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
	var j = 0
	policies := make([]string, len(policiesMap))
	for policy := range policiesMap {
		policies[j] = policy
		j++
	}
	log.Println("Found media: ", payload.Id.String(), " containing Policies: ", policies)
}


/*
func initS3() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	))

	s3svc = s3.New(sess)
}


func stashInS3(path string, reader io.Reader) {


input := s3.PutObjectInput{Bucket: aws.String(BUCKET), Key: mpKey, Body: mpValReader, ACL: aws.String("public-read"), ContentType: aws.String("application/xml; charset=UTF-8"), CacheControl: aws.String("public proxy-revalidate s-maxage=60")}
putResp, err := s3svc.PutObject(&input)
if nil == err {
log.Println("PUT ", *mpKey, " containing ", mpVal)
log.Println(putResp)
} else {
log.Println(err)
}
}*/
