package pixiv_test
import "github.com/yfaimisaka/pixiv/keyword"

func ExampleKeyWord() {
	keyword.KeyWord("夏娜").Num(10).Download() 
	keyword.KeyWord("优库里伍德").Num(20).Upload() // upload to minio server
}
