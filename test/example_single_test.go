package pixiv_test

import "github.com/yfaimisaka/pixiv/single"

func ExampleSingle() {
	single.Single("100510734").Name("蒜").Download() // or .Upload()
}
