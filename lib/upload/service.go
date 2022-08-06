package upload

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
)

//NewService ...
func NewService(env *lib.Env) (lib.UploadService, error) {

	service := new(Service)

	if err := service.validate(env); err != nil {
		return nil, err
	}

	return service, nil

}

//Service our generic upload service that serves as a gateway to one of
//the platforms that we support
type Service struct {
	*lib.Env
	linodeenv *lib.LinodeEnv
}

//Save allows us to save file
func (s *Service) Save(ctx context.Context, file *lib.File) (url string, err error) {
	switch s.Env.Upload.Type {
	case lib.UploadTypeMemory:
		return s.memory(ctx, file)
	case lib.UploadTypeLinode:
		return s.linode(ctx, file)
	}

	return file.Name, nil
}

func (s *Service) memory(ctx context.Context, file *lib.File) (url string, err error) {
	return
}

func (s *Service) linode(ctx context.Context, file *lib.File) (url string, err error) {
	return
}

//validate checks to see whether or not the required configuration
//is present for the specified upload type. If it isn't it raises an
//expection other wise it sets the required config in the service
//itself.
func (s *Service) validate(env *lib.Env) (err error) {
	switch env.Upload.Type {
	case lib.UploadTypeLinode:
		if err := env.Upload.Linode.Validate(); err != nil {
			return err
		}
		s.linodeenv = env.Upload.Linode
		return
	}

	return
}
