package srv

import (
	"context"
	"github.com/krzysztofzaucha/protobuf-playground/internal/model"
)

type Server struct {
	model.PersonServiceServer
	person model.Person
}

func New() *Server {
	return &Server{}
}

func (s *Server) CreatePerson(ctx context.Context, req *model.CreatePersonRequest) (*model.CreatePersonResponse, error) {
	s.person = model.Person{
		ID:      1,
		Name:    req.GetPerson().GetName(),
		Email:   req.GetPerson().GetEmail(),
		Address: req.GetPerson().GetAddress(),
	}

	return &model.CreatePersonResponse{
		Person: &s.person,
	}, nil
}

func (s *Server) ReadPerson(ctx context.Context, req *model.ReadPersonRequest) (*model.ReadPersonResponse, error) {
	return &model.ReadPersonResponse{
		Person: &s.person,
	}, nil
}
