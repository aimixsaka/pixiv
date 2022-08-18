package pixiv

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
)

type single struct {
	singleURL string
}

func (s *single) DownLoadAndSave(name string) {
	imgPath := globalConfig.GetString("download.singlePath")
	fileName := imgPath + "/" + name + s.singleURL[len(s.singleURL)-4:]
	if ok, _ := PathExists(fileName); !ok {
		file, err := os.Create(fileName)
		if err != nil {
			myLog.WithField("place", "single").WithError(err).Fatal("Create file FAILED.")
		}

		writer := bufio.NewWriter(file)
		nn, err := writer.Write(getBytes(s.singleURL))
		if nn == 0 || err != nil {
			myLog.WithField("place", "single").WithError(err).Fatal("Write Picture to file FAILED.")
		}
		myLog.WithField("place", "single").Info(fileName + "Saved.")
	}
}

func (s *single) Upload(name string) {
	minioClient := initMinio()
	minioClient
}

func getBytes(singleURL string) []byte {
	req, err := http.NewRequest("GET", singleURL, nil)
	if err != nil {
		myLog.WithField("place", "single").WithError(err).Fatalln("Fail to request")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.pixiv.net/")
	req.Header.Add("Accept-Charset", "utf-8")

	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		myLog.WithField("place", "single").Fatalf("Response StatusCode: %d!!", res.StatusCode)
	}
	defer res.Body.Close()

	imgByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		myLog.WithField("place", "single").WithError(err).Fatalln("Read Response Body FAILED", err)
	}
	return imgByte
}

func Single() *single {
	return &single{}
}

func (s *single) URL(u string) *single {
	s.singleURL = u
	return s
}