package pixiv

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

var (
)

func getHeader() http.Header {
	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	header.Set("Accept", "*/*")
	header.Set("Accept-Encoding", "gzip,deflate,br")
	header.Set("Connection", "keep-alive")
	header.Set("Referer", "https://www.pixiv.net/")
	header.Add("Accept-Charset", "utf-8")
	return header
}

func getWorkId(initURL string, q url.Values, pageUp int, log logrus.Entry) chan string {
	ids := make(chan string, 60*pageUp)
	pageCount := pageUp
	// 等待通道 确保传输全部完成
	idsChan := make(chan int)

	req, err := http.NewRequest("GET", fmt.Sprintf(initURL + "s", name), nil)
	if err != nil {
		log.WithError(err).Fatalln("创建请求失败")
	}
	req.Header = getHeader()

	for ; pageUp > 0; pageUp-- {
		go func(pageUp int, q url.Values) {
			q.Add("p", strconv.Itoa(pageUp))
			req.URL.RawQuery = q.Encode()
			res, err := client.Do(req)
			if err != nil {
				log.Fatalln("客户端执行请求失败", err)
			}

			reader, _ := gzip.NewReader(res.Body)
			content, err := io.ReadAll(reader)
			defer res.Body.Close()
			if err != nil {
				log.Fatalln("响应体读取失败", err)
			}

			idNum := jsoniter.Get(content, "body").Get("illustManga").Get("data").Size()
			fmt.Println("=========", idNum, "==========")
			if idNum == 0 {
				log.Fatalln("获取data错误", err)
			}

			for ; idNum > 0; idNum-- {
				ids <- jsoniter.Get(content, "body").Get("illustManga").Get("data", idNum-1).Get("id").ToString()
			}
			idsChan <- 1

		}(pageUp, q)
	}
	go func(pageCount int, idsChan chan int) {
		for ; pageCount > 0; pageCount-- {
			<-idsChan
		}
		close(ids)
	}(pageCount, idsChan)
	return ids
}

func getImgUrls(pageUp int, ids chan string, log logrus.Entry) chan string {
	imgUrls := make(chan string, 60*pageUp)
	var urlsCount int
	urlsChan := make(chan int)
	for id := range ids {
		urlsCount++
		go func(id string) {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pixiv.net/artworks/%s", id), nil)
			if err != nil {
				log.Fatalln("请求创建失败", err)
			}

			res, err := client.Do(req)
			if err != nil {
				log.Fatalln("无法得到响应", err)
			}
			if res.StatusCode != 200 {
				log.Fatalf("请求错误， code=%d", res.StatusCode)
			}
			defer res.Body.Close()

			htmlByte, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("读取响应失败")
			}
			html := string(htmlByte)

			reg := regexp.MustCompile(`(?s)"original":"(.*?)"}`)
			u := reg.FindStringSubmatch(html)[1]
			if u == "" {
				log.Fatalf("获取id为%s的imgurl失败", id)
			}

			log.Printf("====已获取到id为%s的作品", id)
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


func upLoadImg(imgUrls chan string, log *logrus.Entry) {
	minioClient := initMinio()
	bucketName := globalConfig.GetString("upload.bucketName")
	var count int
	var imgsCount int
	imgsChan := make(chan int)
	for imgUrl := range imgUrls {
		imgsCount++
		go func(url string) {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalln("发起图片请求失败")
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Encoding", "gzip,deflate,br")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Referer", "https://www.pixiv.net/")
			req.Header.Add("Accept-Charset", "utf-8")

			res, _ := client.Do(req)
			if res.StatusCode != 200 {
				log.Fatalf("获取图片响应码:%d", res.StatusCode)
			}
			defer res.Body.Close()

			name := url[len(url)-15:len(url)-7] + url[len(url)-4:] 
			contentType := res.Header.Get("content-type")
			n, err := minioClient.PutObject(bucketName, name, res.Body, res.ContentLength, minio.PutObjectOptions{ContentType: contentType})		
			if n == 0 || err != nil {
				log.Fatalln("minio上传错误", err)
			}
			log.Println(name + "上传成功")
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
}

func downLoadImg(tagOrName string, imgUrls chan string, log logrus.Entry) {
	if ok, _ := PathExists(imgPath + tagOrName); !ok {
		err := os.Mkdir(imgPath + tagOrName, 0644)
		if err != nil {
			log.Fatalln("创建目录失败", err)	
		}
	}
	var wg sync.WaitGroup
	var count int
	for imgUrl := range imgUrls {
		wg.Add(1)
		count++
		go func(url string, count int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalln("发起图片请求失败")
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Encoding", "gzip,deflate,br")
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("Referer", "https://www.pixiv.net/")
			req.Header.Add("Accept-Charset","utf-8")
			
			res, _ := client.Do(req)
			if res.StatusCode != 200 {
				log.Fatalf("获取图片响应码:%d", res.StatusCode)
			}
			defer res.Body.Close()

			imgByte, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("读取响应体失败", err)
			}

			
			fileName := imgPath + tagOrName + "/" + strconv.Itoa(count) + url[len(url)-4:]
			if ok, _ := PathExists(fileName); !ok {
				file, err := os.Create(fileName)	
				if err != nil {
					log.Fatalln("创建文件失败", err)
				}

				writer := bufio.NewWriter(file)		
				nn, err := writer.Write(imgByte)
				if nn == 0 || err != nil {
					log.Fatalln("图片写入失败", err)
				}
				log.Println(fileName + "已成功保存")

			}
			mu := sync.Mutex{}
			mu.Lock()
			count++
			mu.Unlock()
		}(imgUrl, count)
		wg.Wait()
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


func initMinio() *minio.Client{
	var err error
	minioClient, err := minio.New(
		globalConfig.GetString("upload.endPoint"), 
		globalConfig.GetString("upload.accessKeyID"), 
		globalConfig.GetString("secretAccessKey"),
	    globalConfig.GetBool("useSSL"),
	)
	if err != nil {
		myLog.WithField("place", "config").WithError(err).Fatal("Fail to Create minio Client, please check your config")
	}
	myLog.WithField("place", "config").Info("minio init succeed")
	return minioClient
}
