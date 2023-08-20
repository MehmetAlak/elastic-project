package elastic_operation

import (
	"context"
	"elastic-project/client/elasticsearch"
	"elastic-project/model"
	"github.com/google/uuid"
	"time"
)

type elasticsearchService struct {
	storage elasticsearch.UserInfoStorer
}

type Service interface {
	Create(ctx context.Context, req model.CreateRequest) (model.CreateResponse, error)
	Update(ctx context.Context, userId string, req model.UpdateRequest) error
	Delete(ctx context.Context, req model.DeleteRequest) error
	Find(ctx context.Context, req model.FindRequest) (model.FindResponse, error)
	FindByKeyAndValue(req model.FindByRequest) ([]model.FindResponse, error)
	FindByQuery(query string) ([]model.FindResponse, error)
}

func NewElasticsearchService(storage elasticsearch.UserInfoStorer) Service {
	return &elasticsearchService{storage: storage}
}

func (s elasticsearchService) Create(ctx context.Context, req model.CreateRequest) (model.CreateResponse, error) {
	id := uuid.New().String()
	cr := time.Now().UTC()

	doc := elasticsearch.UserInfo{
		ID:         id,
		Name:       req.Name,
		Job:        req.Job,
		ChildNames: req.ChildNames,
		Comment:    req.Comment,
		CreatedAt:  &cr,
	}

	if err := s.storage.Insert(ctx, doc); err != nil {
		return model.CreateResponse{}, err
	}

	return model.CreateResponse{ID: id}, nil
}

func (s elasticsearchService) Update(ctx context.Context, userId string, req model.UpdateRequest) error {
	doc := elasticsearch.UserInfo{
		ID:         userId,
		Name:       req.Name,
		Job:        req.Job,
		ChildNames: req.ChildNames,
		Comment:    req.Comment,
	}

	if err := s.storage.Update(ctx, doc); err != nil {
		return err
	}

	return nil
}

func (s elasticsearchService) Delete(ctx context.Context, req model.DeleteRequest) error {
	if err := s.storage.Delete(ctx, req.ID); err != nil {
		return err
	}

	return nil
}

func (s elasticsearchService) Find(ctx context.Context, req model.FindRequest) (model.FindResponse, error) {
	userInfo, err := s.storage.FindOne(ctx, req.ID)
	if err != nil {
		return model.FindResponse{}, err
	}

	return model.FindResponse{
		ID:         userInfo.ID,
		Name:       userInfo.Name,
		Job:        userInfo.Job,
		ChildNames: userInfo.ChildNames,
		Comment:    userInfo.Comment,
		CreatedAt:  userInfo.CreatedAt,
	}, nil
}

func (s elasticsearchService) FindByKeyAndValue(req model.FindByRequest) ([]model.FindResponse, error) {
	userInfos, err := s.storage.FindByKeyAndValue(req.QueryType, req.Key, req.Value)
	if err != nil {
		return []model.FindResponse{}, err
	}

	findResponseList := make([]model.FindResponse, 0, len(userInfos))
	for _, userInfo := range userInfos {
		findResponse := model.FindResponse{
			ID:         userInfo.ID,
			Name:       userInfo.Name,
			Job:        userInfo.Job,
			ChildNames: userInfo.ChildNames,
			Comment:    userInfo.Comment,
			CreatedAt:  userInfo.CreatedAt,
		}
		findResponseList = append(findResponseList, findResponse)
	}
	return findResponseList, nil
}

func (s elasticsearchService) FindByQuery(query string) ([]model.FindResponse, error) {
	userInfos, err := s.storage.FindByQuery(query)
	if err != nil {
		return []model.FindResponse{}, err
	}

	findResponseList := make([]model.FindResponse, 0, len(userInfos))
	for _, userInfo := range userInfos {
		findResponse := model.FindResponse{
			ID:         userInfo.ID,
			Name:       userInfo.Name,
			Job:        userInfo.Job,
			ChildNames: userInfo.ChildNames,
			Comment:    userInfo.Comment,
			CreatedAt:  userInfo.CreatedAt,
		}
		findResponseList = append(findResponseList, findResponse)
	}
	return findResponseList, nil
}
