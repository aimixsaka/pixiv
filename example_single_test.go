package pixiv_test

import "github.com/yfaimisaka/pixiv"

func ExampleSingle() {
	pixiv.Single("100510734").Name("蒜").Download() // or .Upload()
}
