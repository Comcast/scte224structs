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
	"strconv"
	"strings"
	"time"
)


type Service_urls struct {
	Altcontent   string
	AltcontentRO string
	Idm          string
}

var Prod = Service_urls{Altcontent: "https://data.altcontent.tv.theplatform.com/altcontent", AltcontentRO: "https://read.data.altcontent.tv.theplatform.com/altcontent", Idm: "https://identity.auth.theplatform.com/idm"}

var Stage = Service_urls{Altcontent: "http://data.altcontent.tv.sandbox.theplatform.com/altcontent", AltcontentRO: "https://read.data.altcontent.tv.sandbox.theplatform.com/altcontent", Idm: "http://stg-admin.identity.auth.theplatform.com/idm"}

const NOTIFY_PATH = "/notify?block=true&fields=true&clientId=%v&schema=1.1.0&filter={MediaSource},{Audience},{Policy},{ViewingPolicy}"

var token *plclient.IdentityToken
var tokenDate time.Time = time.Now().Add(time.Duration(-2) * time.Hour)
var username string
var password string
var env Service_urls

type AltContentClient struct {
	*plclient.IdentityClient
}

func SetCredentials(user string, pw string, urls Service_urls) AltContentClient {
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
	Type  string
	Entry NotificationEntry
}

type NotificationEntry struct {
	Guid string
}

type TypedGuid struct {
	Guid string
	Type NotificationType
}

type NotificationType string

const (
	MEDIA_SOURCE = "MediaSource"
	POLICY = "Policy"
	VIEWING_POLICY = "ViewingPolicy"
	AUDIENCE = "Audience"
)

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
					if response.StatusCode > 399 {
						var buf bytes.Buffer
						buf.ReadFrom(response.Body)
						log.Println(string(buf.Bytes()))
						err = errors.New(response.Status)
					} else {
						deco := json.NewDecoder(response.Body)
						var decodeErr error
						for nil == decodeErr {
							var noti []Notification
							decodeErr = deco.Decode(&noti)
							if nil == decodeErr && len(noti) > 0 {
								for _, n := range noti {
									lastNotificationId = n.Id
									guid := TypedGuid{Guid:n.Entry.Guid, Type:n.Type}
									if "" != guid.Guid && "" != guid.Type {
										guid_chan <- guid
									}
								}
							}
						}
						if decodeErr != io.EOF {
							err = decodeErr
						}
						response.Body.Close()
					}
				}
			}
		}
	}
	if nil != err {
		log.Println(err)
	}
}

func (client AltContentClient) GetSCTEData(account string, guid string) (buf *bytes.Buffer, err error) {
	var token *plclient.IdentityToken
	token, err = client.GetToken()
	if nil == err {
		get, err := http.NewRequest("GET", env.AltcontentRO+"/data/scte224/"+account+"/"+guid, nil)
		get.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			var response *http.Response
			response, err = http.DefaultClient.Do(get)
			if nil == err {
				buf = &bytes.Buffer{}
				buf.ReadFrom(response.Body)
				if response.StatusCode > 399 {
					log.Printf("Got a %v %v from the DS trying to GET %v", response.StatusCode, response.Status, guid)
					log.Println(string(buf.Bytes()))
				}
				err = response.Body.Close()
			}
		}
	}
	return
}

func (client AltContentClient) PushSCTEData(account string, guid string, content string) {

	token, err := client.GetToken()
	if nil == err {
		put, err := http.NewRequest("PUT", env.Altcontent+"/data/scte224/"+account+"/"+guid, strings.NewReader(content))
		put.Header.Add("Content-Type", "application/xml")
		put.Header.Add("Authorization", token.EncodeBasicAuth(account))
		if nil == err {
			response, err := http.DefaultClient.Do(put)
			if nil == err {
				if response.StatusCode > 399 {
					log.Printf("Got a %v %v from the DS trying to PUT %v", response.StatusCode, response.Status, guid)
					var buf bytes.Buffer
					buf.ReadFrom(response.Body)
					log.Println(string(buf.Bytes()))
				}
				err = response.Body.Close()
			}
		}
	}

	if err != nil {
		log.Panic(err)
	}
}
