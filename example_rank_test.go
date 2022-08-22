package pixiv_test

import "github.com/yfaimisaka/pixiv"

func ExampleRankTest() {
	pixiv.Rank().Num(10).DownLoad()
}
