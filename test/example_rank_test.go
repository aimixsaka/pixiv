package pixiv_test

import "github.com/yfaimisaka/pixiv/rank"

// Note, before call Rank(), you need to create cookie.txt
// file in root, and paste request cookie into it.
func ExampleRank() {
	rank.Rank().Num(10).DownLoad()
}
