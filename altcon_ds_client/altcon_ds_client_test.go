package scte224DSClient

import (
	"log"
	"testing"
)

func TestGetToken(t *testing.T) {

	SetCredentials("mpx/stuart_kurkowski@comcast.com", "$tuartRulz!", Prod)
	token, err := GetToken()
	if nil != err {
		t.Error(err)
	} else {
		log.Println(*token)
	}
}
