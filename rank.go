package pixiv

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type rank struct {
	pixiv
	num int
	cookie string
}

// Rank Constructor of rank per month.
func Rank() *rank {
	r := new(rank)
	r.rname = "rank"
	r.log = myLog.WithField("place", "rank")
	r.baseURL = "https://www.pixiv.net/ranking.php?mode=monthly&p=1&format=json"
	r.savePath = globalConfig.GetString("download.rank.path")
	r.cookie = r.getCookie()
	return r
}

// Cookie Set cookie (necessary).
// 
// cookie -cookie in request header.
func (r *rank) Cookie(cookie string) *rank {
	r.cookie = cookie
	return r
}

// Num num of picture to get.
//
// default 100
func (r *rank) Num(num int) *rank {
	if r.cookie == "" {
		r.log.Fatalln("cookie is null, please use Cookie method to set cookie")
	}
	r.num = num
	return r
}

func (r *rank) DownLoad() {
	if r.num < 0 {
		r.log.Fatalln("Please give a number > 0")
	}
	if r.num == 0 {
		r.num = 100
	}
	r.downLoadImg(r.getImgUrls(r.getIds()))
}

func (r *rank) Upload() {
	if r.num < 0 {
		r.log.Fatalln("Please give a number > 0")
	}
	if r.num == 0 {
		r.num = 100
	}
	r.upLoadImg(r.getImgUrls(r.getIds()))
}

func (r *rank) getIds() chan string {
	ids := make(chan string)
	req, err := http.NewRequest("GET", r.baseURL, nil)
	if err != nil {
		r.log.WithError(err).Fatalf("Fail to create request, URL=%s", r.baseURL)
	}

	// q := req.URL.Query()
	// q.Add("mode", "monthly")
	// q.Add("p", "1")
	// q.Add("format", "json")
	// req.URL.RawQuery = q.Encode()

	setHeader(req)
	req.Header.Set("cookie", r.cookie)

	res, err := client.Do(req)
	if err != nil {
		r.log.WithError(err).Fatalln("Fail to get response")
	}

	if code := res.StatusCode; code != 200 {
		if code == 400 {
			r.log.Fatalln("Cookie Error, please use Cookie to set cookie")
		}
		r.log.Fatalf("Response Code=%d", res.StatusCode)
	}

	reader, _ := gzip.NewReader(res.Body)
	defer res.Body.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		r.log.WithError(err).Fatalln("Fail to read response")
	}

	rankDate := fmt.Sprintf("%d%d%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	r.log.Infof("Rank date is: %s\n", rankDate)
	r.fileDir = rankDate
	contents := jsoniter.Get(content, "contents")
	go func() {
		for s := 0; s < r.num; s++ {
			ids <- contents.Get(s).Get("illust_id").ToString()
		}
		close(ids)
	}()

	return ids
}

func (r *rank) getCookie() string {
	cookieFile := "cookie.txt"
	cookieByte, err := os.ReadFile(cookieFile)
	if err != nil {
		cookieFile = "../cookie.txt"
		cookieByte, err = os.ReadFile(cookieFile)
		if err != nil {
			r.log.WithError(err).Fatalln("Fail to read cookie.txt")	
		}
	}
	cookie := *(*string)(unsafe.Pointer(&cookieByte))
	cookie = strings.TrimSpace(cookie)
	return cookie
}
