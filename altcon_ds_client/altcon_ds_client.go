package scte224DSClient

import (
	"bytes"
	"code.comcast.com/svp/plclient"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataServiceGuidResponse struct {
	Entries []DataServiceGuidEntry
}

type DataServiceUpdatedResponse struct {
	Entries []DataServiceUpdatedEntry
}

type DataServiceGuidEntry struct {
	Id   string
	Guid string
}

type DataServiceUpdatedEntry struct {
	DataServiceGuidEntry
	Updated int64
}

type Service_urls struct {
	Altcontent   string
	AltcontentRO string
	Idm          string
}

var Prod = Service_urls{Altcontent: "https://data.altcontent.tv.theplatform.com/altcontent", AltcontentRO: "https://read.data.altcontent.tv.theplatform.com/altcontent", Idm: "https://identity.auth.theplatform.com/idm"}

var Stage = Service_urls{Altcontent: "http://data.altcontent.tv.sandbox.theplatform.com/altcontent", AltcontentRO: "https://read.data.altcontent.tv.sandbox.theplatform.com/altcontent", Idm: "http://stg-admin.identity.auth.theplatform.com/idm"}

const NOTIFY_PATH = "/notify?block=true&fields=true&clientId=%v&schema=1.1.0&filter={MediaSource|put,post},{Audience|put,post},{Policy|put,post},{ViewingPolicy|put,post}"

const GuidLookupTemplate = "/data/MediaPoint/%v?schema=1.3.0&form=cjson&fields=guid,id"

var token *plclient.IdentityToken
var tokenDate time.Time = time.Now().Add(time.Duration(-2) * time.Hour)
var username string
var password string
var env Service_urls
var logger *log.Logger

type AltContentClient struct {
	*plclient.IdentityClient
}

func milliTime(millis int64) time.Time {
	sec := millis / 1000
	nsec := (millis % 1000) * (int64(time.Millisecond) / int64(time.Nanosecond))
	return time.Unix(sec, nsec)
}

func timeToMillis(tempo time.Time) int64 {
	return tempo.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func SetCredentials(user string, pw string, urls Service_urls) AltContentClient {
	logger = log.New(os.Stdout, "AltConDS Client", log.LstdFlags|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	username = user
	password = pw
	env = urls
	return AltContentClient{IdentityClient: plclient.NewIdentityClient(env.Idm, "1.0")}
}

func (client AltContentClient) GetToken() (*plclient.IdentityToken, error) {

	var err error = nil
	if tokenDate.Before(time.Now().Add(time.Duration(-1) * time.Hour)) {
		if nil != token {
			client.DeleteToken(token)
		}
		token, err = client.CreateToken(username, password)
		if nil == err {
			tokenDate = time.Now()
		}
	}
	return token, err
}

type Notification struct {
	Id    int
	Type  DataType
	Entry NotificationEntry
}

type NotificationEntry struct {
	Guid string
}

type TypedGuid struct {
	Guid string
	Type DataType
}

type DataType string

const (
	MEDIA_SOURCE   DataType = "MediaSource"
	POLICY         DataType = "Policy"
	VIEWING_POLICY DataType = "ViewingPolicy"
	AUDIENCE       DataType = "Audience"
)

func (client AltContentClient) LoadMediaPointGuids(account string, mediaPointIds []string) (guidMapping map[string]string, err error) {
	var token *plclient.IdentityToken
	token, err = client.GetToken()
	guidMapping = make(map[string]string, len(mediaPointIds))
	if nil == err {
		get, err := http.NewRequest("GET", env.AltcontentRO + fmt.Sprintf(GuidLookupTemplate, strings.Join(mediaPointIds, ",")), nil)
		get.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			var response *http.Response
			response, err = http.DefaultClient.Do(get)
			if nil == err {
				defer response.Body.Close()
				if response.StatusCode != 200 {
					errMsg := fmt.Sprintf("Got a %v from the DS trying to GET MediaPoint Guids", response.Status)
					logger.Println(errMsg)
					err = errors.New(errMsg)
				} else {
					var decodeErr error
					decoder := json.NewDecoder(response.Body)
					for nil == decodeErr {
						guidLookupPayload := &DataServiceGuidResponse{}
						decodeErr = decoder.Decode(guidLookupPayload)
						  if nil == decodeErr {
							  for _,guidResult := range guidLookupPayload.Entries {
								  // map both directions since I don't yet know which I want
								  guidMapping[guidResult.Guid] = guidResult.Id
								  guidMapping[guidResult.Id] = guidResult.Guid
							  }
						  }
					}
					if decodeErr != io.EOF {
						err = decodeErr
					}
				}
			}
		}
	}
	return
}

func (client AltContentClient) PollForNotifications(account, client_id string, guid_chan chan TypedGuid, lastNotificationId int) {
	var token *plclient.IdentityToken
	var err error
	var happy = true

	var baseNotifyPath string = env.AltcontentRO + fmt.Sprintf(NOTIFY_PATH, client_id)

	for nil == err && happy {
		var path string
		if lastNotificationId > 0 {
			path = baseNotifyPath + "&since=" + strconv.Itoa(lastNotificationId)
		} else {
			path = baseNotifyPath
		}
		token, err = client.GetToken()
		if nil == err {
			get, err := http.NewRequest("GET", path, nil)
			get.Header.Add("Authorization", token.EncodeBasicAuth(account))
			if nil == err {
				var response *http.Response
				response, err = http.DefaultClient.Do(get)
				if nil == err {
					defer response.Body.Close()
					if response.StatusCode != 200 {
						var buf bytes.Buffer
						buf.ReadFrom(response.Body)
						err = errors.New(response.Status + " " + buf.String())
					} else {
						deco := json.NewDecoder(response.Body)
						var decodeErr error
						for nil == decodeErr {
							var noti []Notification
							decodeErr = deco.Decode(&noti)
							if nil == decodeErr && len(noti) > 0 {
								for _, n := range noti {
									lastNotificationId = n.Id
									guid := TypedGuid{Guid: n.Entry.Guid, Type: n.Type}
									if "" != guid.Guid && "" != guid.Type {
										guid_chan <- guid
									}
								}
							}
						}
						if decodeErr != io.EOF {
							err = decodeErr
						}
					}
				}
			}
		}
	}
	if nil != err {
		logger.Println(err)
	}
}

func (client AltContentClient) GetRecentlyUpdated(account string, updatedSince time.Time) (updatedMap map[TypedGuid]time.Time, err error) {

	updatedMap = make(map[TypedGuid]time.Time)
	err = client.queryUpdatedSince(account, updatedMap, AUDIENCE, updatedSince)
	if nil == err {
		err = client.queryUpdatedSince(account, updatedMap, VIEWING_POLICY, updatedSince)
		if nil == err {
			err = client.queryUpdatedSince(account, updatedMap, POLICY, updatedSince)
			if nil == err {
				err = client.queryUpdatedSince(account, updatedMap, MEDIA_SOURCE, updatedSince)
			}
		}
	}
	return
}

func (client AltContentClient) GetUpdatedTimestamps(account string, guids ...TypedGuid) (updatedMap map[TypedGuid]time.Time, err error) {

	var mediaGuids, policyGuids, viewingPolicyGuids, audienceGuids []string
	for _, guid := range guids {
		switch guid.Type {
		case AUDIENCE:
			audienceGuids = append(audienceGuids, guid.Guid)
		case VIEWING_POLICY:
			viewingPolicyGuids = append(viewingPolicyGuids, guid.Guid)
		case POLICY:
			policyGuids = append(policyGuids, guid.Guid)
		case MEDIA_SOURCE:
			mediaGuids = append(mediaGuids, guid.Guid)
		}
	}
	updatedMap = make(map[TypedGuid]time.Time)
	err = client.queryForUpdated(account, updatedMap, AUDIENCE, audienceGuids)
	if nil == err {
		err = client.queryForUpdated(account, updatedMap, VIEWING_POLICY, viewingPolicyGuids)
		if nil == err {
			err = client.queryForUpdated(account, updatedMap, POLICY, policyGuids)
			if nil == err {
				err = client.queryForUpdated(account, updatedMap, MEDIA_SOURCE, mediaGuids)
			}
		}
	}
	return
}

func (client AltContentClient) queryUpdatedSince(account string, updatedMap map[TypedGuid]time.Time, dt DataType, updatedSince time.Time) (err error) {
	var token *plclient.IdentityToken
	token, err = client.GetToken()
	if nil == err {
		get, err := http.NewRequest("GET", env.AltcontentRO+"/data/"+string(dt)+"?schema=1.3.0&fields=id,guid&form=cjson&byUpdated="+strconv.Itoa(int(timeToMillis(updatedSince)))+"~", nil)
		get.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			var response *http.Response
			response, err = http.DefaultClient.Do(get)
			if nil == err {
				defer response.Body.Close()
				if response.StatusCode != 200 {
					errMsg := fmt.Sprintf("Got a %v from the DS trying to GET recently updated %v items", response.Status, dt)
					logger.Println(errMsg)
					err = errors.New(errMsg)
				} else {
					updatePayload := &DataServiceUpdatedResponse{}
					err = json.NewDecoder(response.Body).Decode(updatePayload)
					if nil == err {
						for _, entry := range updatePayload.Entries {
							updatedMap[TypedGuid{Guid: entry.Guid, Type: dt}] = milliTime(entry.Updated)
						}
					}

				}
			}
		}
	}
	return
}

func (client AltContentClient) queryForUpdated(account string, updatedMap map[TypedGuid]time.Time, dt DataType, guids []string) (err error) {
	if len(guids) > 0 {
		var token *plclient.IdentityToken
		token, err = client.GetToken()
		if nil == err {
			get, err := http.NewRequest("GET", env.AltcontentRO+"/data/"+string(dt)+"?schema=1.3.0&form=cjson&byGuid="+strings.Join(guids, "|"), nil)
			get.Header.Add("Authorization", token.EncodeBasicAuth(account))
			if nil == err {
				var response *http.Response
				response, err = http.DefaultClient.Do(get)
				if nil == err {
					defer response.Body.Close()
					if response.StatusCode != 200 {
						errMsg := fmt.Sprintf("Got a %v from the DS trying to GET updated times for: %v", response.Status, guids)
						logger.Println(errMsg)
						err = errors.New(errMsg)
					} else {
						updatePayload := &DataServiceUpdatedResponse{}
						err = json.NewDecoder(response.Body).Decode(updatePayload)
						if nil == err {
							for _, entry := range updatePayload.Entries {
								updatedMap[TypedGuid{Guid: entry.Guid, Type: dt}] = milliTime(entry.Updated)
							}
						}
					}
				}
			}
		}
	}
	return
}

func (client AltContentClient) GetSCTEData(account, guid string) (buf *bytes.Buffer, err error) {
	var token *plclient.IdentityToken
	token, err = client.GetToken()
	if nil == err {
		get, err := http.NewRequest("GET", env.AltcontentRO+"/data/scte224/"+account+"/"+guid, nil)
		get.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			var response *http.Response
			response, err = http.DefaultClient.Do(get)
			if nil == err {
				defer response.Body.Close()
				buf = &bytes.Buffer{}
				buf.ReadFrom(response.Body)
				if response.StatusCode != 200 {
					errMsg := fmt.Sprintf("Got a %v from the DS trying to GET %v: %v", response.Status, guid, buf.String())
					logger.Println(errMsg)
					err = errors.New(errMsg)
				}
			}
		}
	}
	return
}

func (client AltContentClient) PushSCTEData(account string, guid string, content string) error {

	token, err := client.GetToken()
	if nil == err {
		put, err := http.NewRequest("PUT", env.Altcontent+"/data/scte224/"+account+"/"+guid, strings.NewReader(content))
		put.Header.Add("Content-Type", "application/xml")
		put.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			response, err := http.DefaultClient.Do(put)
			if nil == err {
				defer response.Body.Close()
				if response.StatusCode != 200 && response.StatusCode != 204 {
					var buf bytes.Buffer
					buf.ReadFrom(response.Body)
					errMsg := fmt.Sprintf("Got a %v from the DS trying to PUT %v: %v", response.Status, guid, buf.String())
					logger.Println(errMsg)
					err = errors.New(errMsg)
				}
			}
		}
	}
	return err
}
