package pixiv

import (
	"net/http"
)

//https://www.pixiv.net/ajax/search/artworks/%
func getQuery(req *http.Request, name string) {
	q := req.URL.Query()
	q.Add("word", name)
	q.Add("order", "date_d")
	q.Add("mode", "all")
	q.Add("s_mode", "s_tag")
	q.Add("type", "all")
	q.Add("lang", "zh")
}


