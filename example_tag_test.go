package pixiv_test

import "github.com/yfaimisaka/pixiv"

func ExampleTag() {
	pixiv.Tag("錦木千束").Num(10).Download() 
	pixiv.Tag("錦木千束").Num(20).Upload() // upload to minio server
}
