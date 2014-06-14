package main

import (
	// "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	// "net/url"
	"os"
)

import (
	"github.com/qiniu/api/rs"
)

func decodeUpToken(token string) (data rs.PutPolicy, err error) {
	return
}

func encodeUpToken(jsonData string, ak, sk string) (token string, err error) {
	var policy rs.PutPolicy
	err = json.Unmarshal([]byte(jsonData), &policy)
	return
}

func encodeDownToken(baseUrl string, ak, sk string, seconds uint64) (token string, err error) {
	return
}

func main() {
	upToken := flag.String("uptoken", "", "upload token")
	upPolicy := flag.String("uppolicy", "", "upload policy json string")
	accessKey := flag.String("ak", "", "access key")
	secretKey := flag.String("sk", "", "secret key")

	downUrl := flag.String("downurl", "", "download url for sign")
	expiredDuration := flag.Uint64("expires", 0, "expired duration seconds")

	if *upToken != "" {
		upDecode, err := decodeUpToken(*upToken)
		if err != nil {
			fmt.Println(err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(upDecode, "", "")
		fmt.Println(out)
		return
	}

	if *upPolicy != "" {
		if *accessKey == "" || *secretKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		token, err := encodeUpToken(*upPolicy, *accessKey, *secretKey)
		if err != nil {
			fmt.Println(err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		fmt.Println(token)
		return
	}

	if *downUrl != "" {
		if *accessKey == "" || *secretKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		var dur uint64 = 10 * 365 * 24 * 3600
		if *expiredDuration != 0 {
			dur = *expiredDuration
		}
		token, err := encodeDownToken(*downUrl, *accessKey, *secretKey, dur)
		if err != nil {
			fmt.Println(err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		fmt.Println(token)
		return
	}
	return
}
