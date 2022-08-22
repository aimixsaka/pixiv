package pixiv_test

import "github.com/yfaimisaka/pixiv"

func ExampleUser() {
	pixiv.User("26560096").Num(20).Download()
	pixiv.User("26560096").Name("someone").Upload()
}
