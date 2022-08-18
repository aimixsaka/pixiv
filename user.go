package pixiv

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

const (
	proxy   = "http://127.0.0.1:10809"
	imgPath = "C:/ELOI/pixiv/user/"
)

var (
	proxyURL, _ = url.Parse(proxy)
	trans       = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client = &http.Client{
		Transport: trans,
	}
	
)

func getIds(userId string) (ids []string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pixiv.net/ajax/user/%s/profile/all", userId), nil)
	if err != nil {
		log.Fatalln("构造请求失败", err)
	}

	q := req.URL.Query()
	q.Add("lang", "zh")
	req.URL.RawQuery = q.Encode()
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")
	req.Header.Set("Referer", "https://www.pixiv.net/")
	req.Header.Set("Cookie", "first_visit_datetime_pc=2022-07-23+12%3A38%3A44; p_ab_id=4; p_ab_id_2=5; p_ab_d_id=112262462; yuid_b=kpRAgQA; PHPSESSID=79545833_qTF9sOptluFXa3wDDyTSaWC4maFBqPfT; device_token=1b27312ca01e70d8e2a3511955340658; c_type=19; privacy_policy_agreement=0; privacy_policy_notification=0; a_type=0; b_type=0; QSI_S_ZN_5hF4My7Ad6VNNAi=v:0:0; p_b_type=1; tag_view_ranking=q_J28dYJ9d~UBwhLy7Ngq~UD-63UkJba~QYP1NVhSHo~WlM_PwZpNM~P2_D2lj6Ce~50qaEUZGRV~wqBB0CzEFh~yqXYmaGSd-~ZZltVrbyeV~gllivtPcvC~hRkZZnS6_e~MvmmIlzoCJ~yoeOInrCjz~e-qa30G05S~1G1bsV2xcg~3W4zqr4Xlx~Lt-oEicbBr~2u-0Jtvqqd~Zoyka_WNME; __cf_bm=7bS5f9jbYXH3IlcImvwLzjb.j0KNpwvF0G.SAXkutP8-1658722130-0-AedKqorbZ6Ucup7yAiVD57Xo3mYElNW7+pIm7egXSvRKjOFMeMGdDmZdIMgBGEMiyTjiK6/YZo5zvhnta1n2d6JP9rZP7Ms8aesqdpYB4AN5")
	
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("客户端执行请求失败", err)
	}
	
	if code := res.StatusCode; code != 200 {
		log.Fatalf("请求失败，%s",res.Status)
	}

	reader, _ := gzip.NewReader(res.Body)
	defer res.Body.Close()
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln("响应读取失败")
	}
	
	idNum := jsoniter.Get(content, "body").Get("illusts").Size()
	if idNum == 0 {
		log.Fatalln("获取的id列表为空")
	}
	
	log.Printf("获取的id总数为%d", idNum)
	ids = jsoniter.Get(content, "body").Get("illusts").Keys()
	return
	
}

func getImgUrls(ids []string) (imgUrls []string) {
	var wp sync.WaitGroup
	for index, id := range(ids) {
		log.Printf("尝试获取id为%s的作品", id)
		wp.Add(1)
		go func(id string) {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.pixiv.net/artworks/%s", id), nil)
			if err != nil {
				log.Fatalln("请求创建失败", err)
			}
			
			// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
			// req.Header.Set("Accept", "*/*")
			// req.Header.Set("Accept-Encoding", "gzip,deflate,br")
			// req.Header.Set("Referer", "https://www.pixiv.net/")
			// req.Header.Set("Cookie", "first_visit_datetime_pc=2022-07-23+12%3A38%3A44; p_ab_id=4; p_ab_id_2=5; p_ab_d_id=112262462; yuid_b=kpRAgQA; PHPSESSID=79545833_qTF9sOptluFXa3wDDyTSaWC4maFBqPfT; device_token=1b27312ca01e70d8e2a3511955340658; c_type=19; privacy_policy_agreement=0; privacy_policy_notification=0; a_type=0; b_type=0; QSI_S_ZN_5hF4My7Ad6VNNAi=v:0:0; p_b_type=1; tag_view_ranking=q_J28dYJ9d~UBwhLy7Ngq~UD-63UkJba~QYP1NVhSHo~WlM_PwZpNM~P2_D2lj6Ce~50qaEUZGRV~wqBB0CzEFh~yqXYmaGSd-~ZZltVrbyeV~gllivtPcvC~hRkZZnS6_e~MvmmIlzoCJ~yoeOInrCjz~e-qa30G05S~1G1bsV2xcg~3W4zqr4Xlx~Lt-oEicbBr~2u-0Jtvqqd~Zoyka_WNME; __cf_bm=7bS5f9jbYXH3IlcImvwLzjb.j0KNpwvF0G.SAXkutP8-1658722130-0-AedKqorbZ6Ucup7yAiVD57Xo3mYElNW7+pIm7egXSvRKjOFMeMGdDmZdIMgBGEMiyTjiK6/YZo5zvhnta1n2d6JP9rZP7Ms8aesqdpYB4AN5")

			res, err := client.Do(req)
			if err != nil {
				log.Fatalln("无法得到响应", err)
			}
			if res.StatusCode != 200 {
				log.Fatalf("请求错误， code=%d", res.StatusCode)
			}	
			defer res.Body.Close()

			htmlByte, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("读取响应失败")
			}
			html := string(htmlByte)
			
			reg := regexp.MustCompile(`(?s)"original":"(.*?)"}`)
			url := reg.FindStringSubmatch(html)[1]
			if url == "" {
				log.Fatalf("获取id为%s的imgurl失败", id)
			}

			log.Printf("获取到第%d张图片", index)
			imgUrls = append(imgUrls, url)
			wp.Done()
		}(id)
		wp.Wait()
	}	
	
	log.Printf("图片url总数为%d", len(imgUrls))
	return
}

func downLoadImg(imgUrls []string, userId string) {
	if ok, _ := PathExists(imgPath + userId); ok {
	} else {
		er := os.Mkdir(imgPath + userId, 0644)
		if er != nil {
			log.Fatalln("创建目录失败", er)
		}
	}
	var wg sync.WaitGroup
	for index, url := range(imgUrls) {
		wg.Add(1)
		go func(url string, index int) {
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

			imgByte, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("读取响应体失败", err)
			}

			fileName := imgPath + userId + "/" + strconv.Itoa(index) + url[len(url)-4:]
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

				wg.Done()
			}
		}(url, index)
		wg.Wait()
	}
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func main() {
	downLoadImg(getImgUrls(getIds("27517")), "27517")
	//getIds()
}