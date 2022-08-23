package pixiv_test

import "github.com/yfaimisaka/pixiv"

// Note, before call Rank(), you need to create cookie.txt
// file in root, and paste request cookie into it.
func ExampleRank() {
	pixiv.Rank().Num(10).DownLoad()
}
