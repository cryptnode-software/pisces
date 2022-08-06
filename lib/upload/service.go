package upload

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
)

func NewService(env *lib.Env) (lib.UploadService, error) {
	return &Service{
		env,
	}, nil
}

type Service struct {
	*lib.Env
}

func (s *Service) Save(ctx context.Context, file *lib.File) (url string, err error) {
	return file.Name, nil
}
