package pixiv

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// global config
	globalConfig *viper.Viper
	
	client *http.Client
	myLog = &logrus.Logger{
		Level: logrus.DebugLevel,
	}
)

func init() {
	readConfig()
	proxyConfig()
	initMinio()
}


func readConfig() {
	globalConfig.SetConfigName("config")
	globalConfig.SetConfigType("yaml")
	globalConfig.AddConfigPath(".")
	globalConfig.AddConfigPath("./config")
	globalConfig.SetDefault("logOutPut", os.Stdout)
	myLog.Out = globalConfig.Get("logOutPut").(io.Writer)
	err := globalConfig.ReadInConfig()
	if err != nil {
		myLog.WithField("place", "config").Errorf("read config file %s failed", "config.yaml Or config/*.yaml")	
	}
}

func proxyConfig() {
	proxyHost := globalConfig.GetString("proxy.Host")
	proxyPort := globalConfig.GetString("proxy.Port")
	rawURL := "http://" + proxyHost + ":" + proxyPort
	proxyURL, err := url.Parse(rawURL)
	if err != nil {
		myLog.WithField("place", "config").Fatal("URL: [%s] Parse FAILED")
	}
	trans       := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client = &http.Client{
		Transport: trans,
	}
}
