package errors

type Stream struct {
	err error
}

func (s *Stream) Error(e func() error) *Stream {
	if s.err != nil || e == nil {
		return s
	}

	err := e()
	if err != nil {
		s.err = err
	}
	return s
}

func (s *Stream) GetError() error {
	return s.err
}

func S(e func() error) *Stream {
	s := &Stream{}
	return s.Error(e)
}
