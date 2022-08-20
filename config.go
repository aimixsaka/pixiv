package pixiv

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// global config
	globalConfig *viper.Viper = viper.New()

	client *http.Client
	myLog  = logrus.New()
)

type myFormatter struct {}

func init() {
	myLog.SetOutput(os.Stdout)
	myLog.SetLevel(logrus.DebugLevel)
	myLog.SetFormatter(&myFormatter{})
	readConfig()
	proxyConfig()
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	newLog := fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func readConfig() {
	globalConfig.SetConfigName("config")
	globalConfig.SetConfigType("yaml")
	globalConfig.AddConfigPath(".")
	globalConfig.AddConfigPath("../")
	globalConfig.AddConfigPath("./config")
	err := globalConfig.ReadInConfig()
	if err != nil {
		myLog.WithField("place", "config").WithError(err).Fatalf("read config file %s failed\n", "config.yaml Or config/*.yaml")
	}
}

func proxyConfig() {
	proxyHost := globalConfig.GetString("proxy.Host")
	proxyPort := globalConfig.GetString("proxy.Port")
	rawURL := "http://" + proxyHost + ":" + proxyPort
	proxyURL, err := url.Parse(rawURL)
	if err != nil {
		myLog.WithField("place", "config").Fatalf("URL: [%s] Parse FAILED\n", proxyURL)
	}
	trans := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client = &http.Client{
		Transport: trans,
	}
}

func initMinio() *minio.Client {
	var err error
	client, err := minio.New(
		globalConfig.GetString("upload.endPoint"),
		globalConfig.GetString("upload.accessKeyID"),
		globalConfig.GetString("secretAccessKey"),
		globalConfig.GetBool("useSSL"),
	)
	if err != nil {
		myLog.WithField("place", "config").WithError(err).Fatal("Fail to Create minio Client, please check your config")
	}
	myLog.WithField("place", "config").Info("minio init successfully")
	return client
}