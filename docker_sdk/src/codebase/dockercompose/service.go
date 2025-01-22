package dockercompose

import "github.com/compose-spec/compose-go/types"

type Service struct {
	s *types.ServiceConfig
}

func NewService(service *types.ServiceConfig) *Service {
	return &Service{
		s: service,
	}
}

func (s *Service) Name() string {
	return s.s.Name
}

func (s *Service) Image() string {
	return s.s.Image
}

func (s *Service) Workdir() string {
	return s.s.WorkingDir
}
