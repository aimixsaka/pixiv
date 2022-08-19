package pixiv

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type user struct {
	pixiv
	userId string
}

// Constructor of userId of pictures.
// userId -id of the user.
func User(userId string) *user {
	u := new(user)
	u.rname = "user"
	u.log = myLog.WithField("place", "user")
	u.baseURL = "https://www.pixiv.net/ajax/user/%s/profile/all"
	u.userId = userId
	u.savePath = globalConfig.GetString("download.user.path")
	u.defaultName()
	return u

}

// Set dir name.
// Default is userid
func (u *user) Name(name string) *user {
	u.fileDir = name
	return u
}

func (u *user) Num(num int) *user {
	u.num = num
	return u
}

func (u *user) Download() {
	u.downLoadImg(u.getImgUrls(u.getIds()))
}

func (u *user) Upload() {
	u.upLoadImg(u.getImgUrls(u.getIds()))
}

func (u *user) defaultName() {
	u.fileDir = u.userId
}

func (u *user) getIds() chan string {
	ids := make(chan string)
	URL := fmt.Sprintf(u.baseURL, u.userId)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		u.log.WithError(err).Fatalf("Fail to create request, URL=%s", URL)
	}

	q := req.URL.Query()
	q.Add("lang", "zh")
	req.URL.RawQuery = q.Encode()
	setHeader(req)

	res, err := client.Do(req)
	if err != nil {
		u.log.WithError(err).Fatalln("Fail to get response")
	}

	if code := res.StatusCode; code != 200 {
		u.log.Fatalf("URL Code=%d", res.StatusCode)
	}

	reader, _ := gzip.NewReader(res.Body)
	defer res.Body.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		u.log.WithError(err).Fatalln("Fail to read response")
	}

	idNum := jsoniter.Get(content, "body").Get("illusts").Size()
	if idNum == 0 {
		u.log.Fatalln("Fail to get ids, ids list is null")
	}
	if u.num > idNum {
		log.Fatalf("Total works in id=%s is: %d, while got %d\n", u.userId, idNum, u.num)
	}
	keys := jsoniter.Get(content, "body").Get("illusts").Keys()

	go func() {
		for i := 0; i < u.num; i++ {
			ids <- keys[i]
		}
		close(ids)
	}()
	return ids
}
