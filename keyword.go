package pixiv

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type keyWord struct {
	pixiv
	keyWord string
}

// Constructor of keyword of pictures
// word -keyword to search.
func KeyWord(word string) *keyWord {
	k := new(keyWord)
	k.rname ="keyword"
	k.log = myLog.WithField("place", "keyWord")
	k.baseURL = "https://www.pixiv.net/ajax/search/artworks/%s"
	k.savePath = globalConfig.GetString("download.keyword.path")
	k.keyWord = word
	k.fileDir = word
	return k
}

// Set num of pictures.
func (k *keyWord) Num(num int) *keyWord {
	if num <= 0 {
		k.log.Fatalln("Please give a number > 0")
	}
	k.num = num
	return k
}

func (k *keyWord) Download() {
	if k.num == 0 {
		k.log.Fatalln("Please set num before download")
	}
	k.downLoadImg(k.getImgUrls(k.getIds()))
}

func (k *keyWord) Upload() {
	if k.num == 0 {
		k.log.Fatalln("Please set num before download")
	}
	k.upLoadImg(k.getImgUrls(k.getIds()))
}

func (k *keyWord) getIds() chan string {
	ids := make(chan string, k.num)
	pageUp := k.num/60 + 1
	numLeft := k.num % 60
	pageCount := pageUp
	// wait chan
	idsChan := make(chan int)

	URL := fmt.Sprintf(k.baseURL, k.keyWord)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		k.log.WithError(err).Fatalf("Fail to create request, URL=%s", URL)
	}
	setHeader(req)

	for ; pageUp > 0; pageUp-- {
		go func(pageUp int) {
			idNum := 60
			if pageUp == pageCount {
				idNum = numLeft
			}
			q := req.URL.Query()
			q.Add("word", k.keyWord)
			q.Add("order", "date_d")
			q.Add("mode", "all")
			q.Add("p", strconv.Itoa(pageUp))
			q.Add("s_mode", "s_tag_full")
			q.Add("type", "all")
			q.Add("lang", "zh")
			req.URL.RawQuery = q.Encode()

			res, err := client.Do(req)
			if err != nil {
				k.log.WithError(err).Fatalln("Fail to get response")
			}

			reader, _ := gzip.NewReader(res.Body)
			content, err := io.ReadAll(reader)
			defer res.Body.Close()
			if err != nil {
				k.log.WithError(err).Fatalln("Fail to read response")
			}

			for ; idNum > 0; idNum-- {
				ids <- jsoniter.Get(content, "body").Get("illustManga").Get("data", idNum-1).Get("id").ToString()
			}
			idsChan <- 1

		}(pageUp)
	}
	go func(pageCount int, idsChan chan int) {
		for ; pageCount > 0; pageCount-- {
			<-idsChan
		}
		close(ids)
	}(pageCount, idsChan)
	return ids
}
