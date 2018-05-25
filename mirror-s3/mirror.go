package main

import (
	"flag"
	"github.comcast.com/jcolwe200/decider/s3"
	"log"
	"os"
)

var mpx_user, mpx_password, account, bucket, base_path string
var notification int

func main() {
	logger := log.New(os.Stdout, "S3 Mirror", log.LstdFlags|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	flag.Parse()
	if "" == mpx_user || "" == mpx_password || "" == bucket {
		flag.Usage()
	} else {
		_, s3err := s3.CreateDefaultS3Mirror(mpx_user, mpx_password, account, bucket, nil, nil, logger, false)
		if nil != s3err {
			logger.Println(s3err)
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
