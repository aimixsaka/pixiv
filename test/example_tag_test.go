package pixiv_test

import "github.com/yfaimisaka/pixiv/tag"

func ExampleTag() {
	tag.Tag("錦木千束").Num(10).Download() 
	tag.Tag("錦木千束").Num(20).Upload() // upload to minio server
}
