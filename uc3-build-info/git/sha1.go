package git

type SHA1 string

func (s SHA1) String() string {
	fullSha := s.Full()
	if FullSHA {
		return fullSha
	}
	return fullSha[0:12]
}

func (s SHA1) Full() string {
	return string(s)
}
