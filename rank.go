package pixiv

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type rank struct {
	pixiv
	num int
}

func Rank() *rank {
	r := new(rank)
	r.log = myLog.WithField("place", "rank")
	r.baseURL = "https://www.pixiv.net/ajax/top/illust"
	r.savePath = globalConfig.GetString("download.rank.path")
	return r
}

func (r *rank) Num(num int) *rank {
	r.num = num
	return r
}

func (r *rank) DownLoad() {
	r.downLoadImg(r.getImgUrls(r.getIds()))
}

func (r *rank) Upload() {
	r.upLoadImg(r.getImgUrls(r.getIds()))
}

func (r *rank) getIds() chan string {
	ids := make(chan string)
	req, err := http.NewRequest("GET", r.baseURL, nil)
	if err != nil {
		r.log.WithError(err).Fatalf("Fail to create request, URL=%s", r.baseURL)
	}

	q := req.URL.Query()
	q.Add("mode", "all")
	q.Add("lang", "zh")
	req.URL.RawQuery = q.Encode()

	setHeader(req)

	res, err := client.Do(req)
	if err != nil {
		r.log.WithError(err).Fatalln("Fail to get response")
	}

	if code := res.StatusCode; code != 200 {
		r.log.Fatalf("Response Code=%d", res.StatusCode)
	}

	reader, _ := gzip.NewReader(res.Body)
	defer res.Body.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		r.log.WithError(err).Fatalln("Fail to read response")
	}

	rankDate := jsoniter.Get(content, "body").Get("page").Get("ranking").Get("date").ToString()
	if rankDate == "" {
		log.Fatalln("Fail to get today's rank date")
	}

	r.log.Infof("Rank date is: %s\n", rankDate)
	r.fileDir = rankDate

	go func() {
		for s := 0; s < r.num; s++ {
			ids <- jsoniter.Get(content, "body").Get("page").Get("ranking").Get("items", s).Get("id").ToString()
		}
		close(ids)
	}()

	return ids
}
