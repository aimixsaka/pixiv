package pixiv

import (
	"net/http"
	"os"

	"github.com/minio/minio-go"
)

func setHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.pixiv.net/")
	req.Header.Add("Accept-Charset", "utf-8")
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

func initMinio() *minio.Client {
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
