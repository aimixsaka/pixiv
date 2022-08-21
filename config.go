package pixiv

import (
	"net/http"
	"net/url"
	"os"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
    "github.com/shiena/ansicolor"
)

var (
	// global config
	globalConfig *viper.Viper = viper.New()

	client *http.Client
	myLog  = logrus.New()
)


func init() {
	myLog.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	myLog.SetLevel(logrus.DebugLevel)
	myLog.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		ForceQuote:true,    
		TimestampFormat:"2006-01-02 15:04:05",  
		FullTimestamp:true,    
	})
	readConfig()
	proxyConfig()
}


func readConfig() {
	globalConfig.SetConfigName("config")
	globalConfig.SetConfigType("yaml")
	globalConfig.AddConfigPath(".")
	globalConfig.AddConfigPath("../")
	globalConfig.AddConfigPath("./config")
	err := globalConfig.ReadInConfig()
	if err != nil {
		myLog.WithField("place", "config").WithError(err).Fatalf("read config file %s failed", "config.yaml Or config/*.yaml")
	}
}

func proxyConfig() {
	proxyHost := globalConfig.GetString("proxy.Host")
	proxyPort := globalConfig.GetString("proxy.Port")
	rawURL := "http://" + proxyHost + ":" + proxyPort
	proxyURL, err := url.Parse(rawURL)
	if err != nil {
		myLog.WithField("place", "config").Fatalf("URL: [%s] Parse FAILED", proxyURL)
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
		myLog.WithField("place", "config").WithError(err).Fatalln("Fail to Create minio Client, please check your config")
	}
	myLog.WithField("place", "config").Infoln("minio init successfully")
	return client
}