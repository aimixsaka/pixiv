package pixiv_test

import "github.com/yfaimisaka/pixiv"

func ExampleKeyWord() {
	pixiv.KeyWord("夏娜").Num(10).Download() 
	pixiv.KeyWord("优库里伍德").Num(20).Upload() // upload to minio server
}
