// Package pixiv
//
// a pixiv spider in go
package pixiv

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

type pixiv struct {
	baseURL  string
	num      int
	log      *logrus.Entry
	savePath string
	fileDir  string
	// real struct name
	rname string
}

func (p *pixiv) getImgUrls(ids chan string) chan string {
	imgUrls := make(chan string, p.num)
	var urlsCount int
	urlsChan := make(chan int)
	for id := range ids {
		urlsCount++
		go func(id string) {
			URL := fmt.Sprintf("https://www.pixiv.net/artworks/%s", id)
			req, err := http.NewRequest("GET", URL, nil)
			if err != nil {
				p.log.WithError(err).Fatalf("Fail to create request, URL=%s", URL)
			}

			res, err := client.Do(req)
			if err != nil {
				p.log.WithError(err).Fatalln("Fail to get response")
			}
			if res.StatusCode != 200 {
				p.log.Fatalf("Response Code=%d", res.StatusCode)
			}
			defer res.Body.Close()

			htmlByte, err := io.ReadAll(res.Body)
			if err != nil {
				p.log.WithError(err).Fatalln("Fail to read response")
			}
			html := string(htmlByte)

			reg := regexp.MustCompile(`(?s)"original":"(.*?)"}`)
			u := reg.FindStringSubmatch(html)[1]
			if u == "" {
				p.log.WithError(err).Fatalf("Fail to get url of id=%s", id)
			}

			p.log.Infof("Got work that id is %s", id)
			imgUrls <- u

			urlsChan <- 1
		}(id)
	}

	go func(urlsCount int, urlsChan chan int) {
		for ; urlsCount > 0; urlsCount-- {
			<-urlsChan
		}
		close(imgUrls)
	}(urlsCount, urlsChan)
	return imgUrls
}

func (p *pixiv) upLoadImg(imgUrls chan string) {
	minioClient := initMinio()
	bucketName := globalConfig.GetString("upload.bucketName")
	fmt.Println(bucketName)
	var count int
	var imgsCount int
	imgsChan := make(chan int)
	for imgUrl := range imgUrls {
		imgsCount++
		go func(url string) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				p.log.WithError(err).Fatalf("Fail to send request, URL=%s", url)
			}
			setHeader(req)

			res, _ := client.Do(req)
			if res.StatusCode != 200 {
				p.log.Fatalf("URL Code=%d", res.StatusCode)
			}
			defer res.Body.Close()

			name :=  "/" + "pixiv" + "/" + p.rname + "/" + strconv.Itoa(count) + url[len(url)-4:]
			contentType := res.Header.Get("content-type")
			n, err := minioClient.PutObject(bucketName, name, res.Body, res.ContentLength, minio.PutObjectOptions{ContentType: contentType})
			if n == 0 || err != nil {
				p.log.WithError(err).Fatalln("Fail to upload to minio")
			}
			p.log.Infoln(name + "upload succeded")
			lock := sync.Mutex{}
			lock.Lock()
			count++
			lock.Unlock()
			imgsChan <- 1
		}(imgUrl)
	}
	for ; imgsCount > 0; imgsCount-- {
		<-imgsChan
	}
	p.log.Infof("Total Uploaded: %d pictures", count)
}

func (p *pixiv) downLoadImg(imgUrls chan string) {
	var wg sync.WaitGroup
	var count int
	var sep string // seperator -depends on os(windows/linux)
	sysType := runtime.GOOS
	if sysType == "linux" {
		sep = "/"
	} else if sysType == "windows" {
		sep = "\\"	
	}
	if ok, _ := pathExists(p.savePath + sep + p.fileDir); !ok {
		err := os.Mkdir(p.savePath + sep + p.fileDir, 0766)
		if err != nil {
			p.log.WithError(err).Fatalf("Fail to create dir: %s", p.savePath+p.fileDir)
		}
	}
	for imgUrl := range imgUrls {
		wg.Add(1)
		count++
		go func(url string, count int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				p.log.WithError(err).Fatalf("Fail to create request, URL=%s", url)
			}
			setHeader(req)

			res, _ := client.Do(req)
			if res.StatusCode != 200 {
				p.log.Fatalf("Response Code=%d", res.StatusCode)
			}
			defer res.Body.Close()

			imgByte, err := io.ReadAll(res.Body)
			if err != nil {
				p.log.WithError(err).Fatalln("Fail to read response")
			}

			
			fileName := p.savePath + sep + p.fileDir + sep + strconv.Itoa(count) + url[len(url)-4:]
			if ok, _ := pathExists(fileName); !ok {
				file, err := os.Create(fileName)
				if err != nil {
					p.log.WithError(err).Fatalf("Fail to create file: %s", fileName)
				}

				writer := bufio.NewWriter(file)
				nn, err := writer.Write(imgByte)
				if nn == 0 || err != nil {
					p.log.WithError(err).Fatalf("Fail to write picture: %s", fileName)
				}
				p.log.Infoln(fileName + " save successfully")

			}
			mu := sync.Mutex{}
			mu.Lock()
			count++
			mu.Unlock()
		}(imgUrl, count)
		wg.Wait()
		p.log.Infof("Total: %d", count)
	}
}
