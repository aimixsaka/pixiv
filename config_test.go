package pixiv

import "testing"

func TestConfig(t *testing.T) {
	configs := []string{
		"download.single.path",
		"download.keyword.path",
		"download.tag.path",
		"download.rank.path",
		"download.user.path",
		"upload.endPoint",
 		"upload.accessKeyID",
 		"upload.secretAccessKey",
 		"upload.useSSL",
 		"upload.bucketName",
	}
	
	for _, config := range configs {
		if globalConfig.GetString(config) == "" {
			t.Errorf("Config: %s is null\n", config)
		}
	}
	
}