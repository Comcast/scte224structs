package scte224DSClient

import (
	"github.comcast.com/svp/plclient"
	"time"
	"net/http"
	"strings"
	"log"
	"bytes"
)

type Service_urls struct {
	Altcontent string
	Idm string
}

var Prod = Service_urls{Altcontent:"https://data.altcontent.tv.theplatform.com/altcontent", Idm:"https://identity.auth.theplatform.com/idm"}

var Stage = Service_urls{Altcontent:"http://data.altcontent.tv.sandbox.theplatform.com/altcontent", Idm:"http://stg-admin.identity.auth.theplatform.com/idm"}

var token *plclient.IdentityToken
var tokenDate time.Time = time.Now().Add(time.Duration(-2) * time.Hour)
var username string
var password string
var env Service_urls
var client *plclient.IdentityClient

func SetCredentials(user string, pw string, urls Service_urls) {
	username = user
	password = pw
	env = urls
	client = plclient.NewIdentityClient(env.Idm, "1.0")
}


func GetToken() (*plclient.IdentityToken, error) {

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


func PushSCTEData(account string, guid string, content string) {

	token, err := GetToken()
	if nil == err {
		put, err := http.NewRequest("PUT", env.Altcontent+"/data/scte224/"+account +"/"+guid, strings.NewReader(content))
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