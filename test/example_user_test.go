package pixiv_test

import (
	"github.com/yfaimisaka/pixiv/user"
)

func ExampleUser() {
	user.User("26560096").Num(20).Download()
	user.User("26560096").Name("someone").Upload()
}
