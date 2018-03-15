package main

import (
	"encoding/xml"
	"flag"
	"github.com/metaleap/go-xsd/types"
	"code.comcast.com/jcolwe200/scte224/altcon_ds_client"
	"code.comcast.com/jcolwe200/scte224/go-xsd-generated-types/www.scte.org/schemas/224/2015/SCTE224.xsd_go"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"
	"code.comcast.com/jcolwe200/scte224/go-xsd-generated-types/www.scte.org/schemas/224/2018/SCTE224-2018.xsd_go"
)

type MediaPayload struct {
	XMLName xml.Name
	go_Scte2242018.TMediaType
}

var scte224Template string
var username string
var password string
var account string
var sourceGUID string
var mediaDurationStr string
var mediaPointDurationStr string
var mediaDuration, mediaPointDuration time.Duration

const SCTE_TIME_FMT = "2006-01-02T15:04:05.999Z"

func init() {
	flag.StringVar(&scte224Template, "t", "", "SCTE-224 Media template path")
	flag.StringVar(&username, "u", "", "MPX username")
	flag.StringVar(&password, "p", "", "MPX password")
	flag.StringVar(&account, "a", "", "MPX account (your account will need role http://access.auth.theplatform.com/data/Role/496696 to work)")
	flag.StringVar(&sourceGUID, "g", "superflaco.com/media/CSNPH", "Media Source GUID")
	flag.StringVar(&mediaDurationStr, "d", "336h", "total Media duration")
	flag.StringVar(&mediaPointDurationStr, "m", "30m", "MediaPoint window duration")
}

func main() {
	flag.Parse()

	var failedToStart bool = false
	var merr, mperr error
	mediaDuration, merr = time.ParseDuration(mediaDurationStr)
	if nil != merr {
		failedToStart = true
		log.Println("media point duration could not be parsed: ", mediaDurationStr)
	}
	mediaPointDuration, mperr = time.ParseDuration(mediaPointDurationStr)
	if nil != mperr {
		failedToStart = true
		log.Println("media point duration could not be parsed: ", mediaPointDurationStr)
	}
	if "" == scte224Template {
		failedToStart = true
		log.Println("a media template is required")
	}

	if "" == username {
		failedToStart = true
		log.Println("an MPX username is required")
	}

	if "" == password {
		failedToStart = true
		log.Println("an MPX password is required")
	}

	if "" == account {
		failedToStart = true
		log.Println("an MPX account is required")
	}

	if failedToStart {
		flag.Usage()
	} else {
		data, err := ioutil.ReadFile(scte224Template)

		if nil == err {
			//fileStr := string(data)
			//log.Println(fileStr)
			var media = &MediaPayload{}
			err2 := xml.Unmarshal(data, media)
			now := time.Now().UTC().Truncate(time.Hour).Add(time.Duration(-1) * time.Hour)
			updateMedia(media, now)
			updatePoints(media, now)
			if nil == err2 {
				display, err3 := xml.MarshalIndent(media, "", "  ")
				if nil == err3 {
					preamblePattern := regexp.MustCompile("^<Media xmlns=\"http://www.scte.org/schemas/224/2015\"")
					preambleReplacement := "<Media xmlns:xlink=\"http://www.w3.org/1999/xlink\" xmlns=\"http://www.scte.org/schemas/224/2015\""
					display = preamblePattern.ReplaceAllLiteral(display, []byte(preambleReplacement))
					hrefPattern := regexp.MustCompile(" href=")
					display = hrefPattern.ReplaceAllLiteral(display, []byte(" xlink:href="))
					output := string(display)

					log.Println(output)
					client := scte224DSClient.SetCredentials(username, password, scte224DSClient.Prod)
					client.PushSCTEData(account, sourceGUID, output)
				} else {
					log.Println(err3)
				}
			} else {
				log.Println(err2)
			}
		} else {
			log.Println(err)
		}
	}
}

func formatXMLTime(when time.Time) go_Scte2242018.ConvertibleDateTime {
	return go_Scte2242018.ConvertibleDateTime(when.Format(SCTE_TIME_FMT))
}

func updateMedia(payload *MediaPayload, when time.Time) {

	twoWeeksOut := when.Add(mediaDuration)

	nowfmt := formatXMLTime(when)
	payload.LastUpdated = nowfmt
	payload.Effective = nowfmt
	payload.Expires = formatXMLTime(twoWeeksOut)
}

func updatePoints(payload *MediaPayload, when time.Time) {

	updatedTimeStamp := formatXMLTime(when)
	expiresTimeStamp := formatXMLTime(when.Add(mediaPointDuration))

	startPoint := resetPointTimes(payload.MediaPoints[0], updatedTimeStamp, expiresTimeStamp)
	originalStartId := startPoint.Id

	endPoint := resetPointTimes(payload.MediaPoints[1], updatedTimeStamp, expiresTimeStamp)
	originalEndId := endPoint.Id

	// make enough media points to fill the media duration with both a waxon and a waxoff point
	pointCount := 2 * int(mediaDuration.Seconds()) / int(mediaPointDuration.Seconds())
	pointList := make([]*go_Scte2242018.TMediaPointType, 0, pointCount)
	for j := 0; j < (pointCount / 2); j++ {
		// intentionally dereferencing the pointer to force a copy so we clone the points as we increment the fields
		startPoint = incrementPoint(*startPoint, originalStartId+xsdt.AnyURI("/"+strconv.Itoa(j)))
		endPoint = incrementPoint(*endPoint, originalEndId+xsdt.AnyURI("/"+strconv.Itoa(j)))
		pointList = append(pointList, startPoint, endPoint)
	}
	payload.MediaPoints = pointList
}

func resetPointTimes(point *go_Scte2242018.TMediaPointType, effective go_Scte2242018.ConvertibleDateTime, expires go_Scte2242018.ConvertibleDateTime) *go_Scte2242018.TMediaPointType {
	point.LastUpdated = effective
	point.Effective = effective
	point.Expires = expires
	return point
}

func incrementPoint(point go_Scte2242018.TMediaPointType, uri xsdt.AnyURI) *go_Scte2242018.TMediaPointType {
	// passing in the object is intended to force it to create a new copy which gets modified
	point.Effective = formatXMLTime(point.Effective.Time().Add(mediaPointDuration))
	point.Expires = formatXMLTime(point.Expires.Time().Add(mediaPointDuration))
	point.Id = uri
	return &point
}
