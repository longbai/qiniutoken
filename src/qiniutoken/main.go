package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/qiniu/api/auth/digest"
	"github.com/qiniu/api/rs"
)

func decodeUpToken(token string) (policy rs.PutPolicy, err error) {
	array := strings.Split(token, ":")
	if len(array) != 3 {
		err = errors.New("invalid token")
		return
	}

	data, err := base64.URLEncoding.DecodeString(array[2])
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &policy)
	return
}

func encodeUpToken(args []string, ak, sk string) (token string, err error) {
	var policy rs.PutPolicy
	for _, v := range args {
		arg := strings.Split(v, "=")
		if len(arg) != 2 {
			err = errors.New("invalid args " + v)
			return
		}
		switch arg[0] {
		case "Scope":
			policy.Scope = arg[1]
		case "Expires":
			data, err1 := strconv.ParseUint(arg[1], 10, 32)
			if err1 != nil {
				err = err1
				return
			}
			policy.Expires = uint32(data)
		}
	}
	mac := digest.Mac{ak, []byte(sk)}
	token = policy.Token(&mac)
	return
}

func encodeDownToken(baseUrl string, ak, sk string, seconds uint64) (downUrl string, err error) {
	var policy rs.GetPolicy
	policy.Expires = uint32(seconds)
	mac := digest.Mac{ak, []byte(sk)}
	downUrl = policy.MakeRequest(baseUrl, &mac)
	return
}

func main() {
	upToken := flag.String("uptoken", "", "upload token")
	upPolicy := flag.Bool("uppolicy", false, "upload policy")
	accessKey := flag.String("ak", "", "access key")
	secretKey := flag.String("sk", "", "secret key")

	downUrl := flag.String("downurl", "", "download url for sign")
	expiredDuration := flag.Uint64("expires", 0, "expired duration seconds")
	flag.Parse()

	if *upToken != "" {
		upDecode, err := decodeUpToken(*upToken)
		if err != nil {
			fmt.Println(err)
			flag.PrintDefaults()
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(upDecode, "", "")
		fmt.Println(string(out))
		fmt.Println(time.Unix(int64(upDecode.Expires), 0))
		return
	}

	if *upPolicy {
		if *accessKey == "" || *secretKey == "" {
			fmt.Println("no ak or sk")
			flag.PrintDefaults()
			os.Exit(1)
		}
		token, err := encodeUpToken(flag.Args(), *accessKey, *secretKey)
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
			fmt.Println("no ak or sk")
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
	flag.PrintDefaults()
	return
}
