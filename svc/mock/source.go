package mock

type Source struct {
	Words string
}

func (s *Source) GetWisdom() string {
	return s.Words
}
