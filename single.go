package pixiv

type single struct {
	pixiv
	workId chan string
}

// Constructor of single picture.
// workId -the id of the single picture.
func Single(workId string) *single {
	s := new(single)
	s.rname = "single"
	s.log = myLog.WithField("place", "single")
	s.savePath = globalConfig.GetString("download.single.path")
	s.num = 1
	s.workId = make(chan string, 1)
	s.defaultSingleDir()
	go func() {
		s.workId <- workId
		close(s.workId)
	}()
	return s
}

// Set single dir name.
// Default is single
func (s *single) Name(dirName string) *single {
	s.fileDir = dirName
	return s
}

func (s *single) Download() {
	s.downLoadImg(s.getImgUrls(s.workId))
}

func (s *single) Upload() {
	s.upLoadImg(s.getImgUrls(s.workId))
}

func (s *single) defaultSingleDir() {
	s.fileDir = "single"
}