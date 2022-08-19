package pixiv

type single struct {
	pixiv
	workId chan string
}

func Single(workId string) *single {
	s := new(single)
	s.log = myLog.WithField("place", "single")
	s.savePath = globalConfig.GetString("download.single.path")
	go func() {
		s.workId <- workId
		close(s.workId)
	}()
	return s
}

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
