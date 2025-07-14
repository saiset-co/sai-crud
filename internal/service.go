package internal

import (
	"github.com/go-playground/validator/v10"
	"github.com/saiset-co/sai-crud/types"
	"github.com/saiset-co/sai-service/sai"
	saiTypes "github.com/saiset-co/sai-service/types"
	storageTypes "github.com/saiset-co/sai-storage/types"
)

type Service struct {
	validator *validator.Validate
	config    *types.ServiceConfig
}

func NewService(config *types.ServiceConfig) *Service {
	return &Service{
		validator: validator.New(),
		config:    config,
	}
}

func (s *Service) Create(ctx *saiTypes.RequestCtx, req types.CreateRequest) (resp types.CreateResponse, err error) {
	if err := s.validator.Struct(req); err != nil {
		return types.CreateResponse{}, saiTypes.WrapError(err, "validation failed")
	}

	storageRequest := storageTypes.CreateDocumentsRequest{
		Collection: req.Prefix + "_" + s.config.Collection,
		Data:       req.Data,
	}

	storageResult, _, err := sai.ClientManager().Call("storage", "POST", "/api/v1/documents", storageRequest, nil)
	if err != nil {
		return resp, err
	}

	var result storageTypes.CreateDocumentsResponse
	err = ctx.Unmarshal(storageResult, &result)
	if err != nil {
		return resp, err
	}

	return types.CreateResponse{
		Data:    result.Data,
		Created: result.Created,
	}, nil
}

func (s *Service) Read(ctx *saiTypes.RequestCtx, req types.ReadRequest) (resp types.ReadResponse, err error) {
	if err := s.validator.Struct(req); err != nil {
		return types.ReadResponse{}, saiTypes.WrapError(err, "validation failed")
	}

	storageRequest := storageTypes.ReadDocumentsRequest{
		Collection: req.Prefix + "_" + s.config.Collection,
		Filter:     req.Filter,
		Limit:      req.Limit,
		Sort:       req.Sort,
		Skip:       req.Skip,
	}

	storageResult, _, err := sai.ClientManager().Call("storage", "GET", "/api/v1/documents", storageRequest, nil)
	if err != nil {
		return resp, err
	}

	var result storageTypes.ReadDocumentsResponse
	err = ctx.Unmarshal(storageResult, &result)
	if err != nil {
		return resp, err
	}

	return types.ReadResponse{
		Data:  result.Data,
		Total: result.Total,
	}, nil
}

func (s *Service) Update(ctx *saiTypes.RequestCtx, req types.UpdateRequest) (resp types.UpdateResponse, err error) {
	if err := s.validator.Struct(req); err != nil {
		return types.UpdateResponse{}, saiTypes.WrapError(err, "validation failed")
	}

	storageRequest := storageTypes.UpdateDocumentsRequest{
		Collection: req.Prefix + "_" + s.config.Collection,
		Filter:     req.Filter,
		Data:       req.Data,
	}

	storageResult, _, err := sai.ClientManager().Call("storage", "PUT", "/api/v1/documents", storageRequest, nil)
	if err != nil {
		return resp, saiTypes.WrapError(err, "failed to update documents")
	}

	var result storageTypes.UpdateDocumentsResponse
	err = ctx.Unmarshal(storageResult, &result)
	if err != nil {
		return resp, err
	}

	return types.UpdateResponse{
		Updated: result.Updated,
		Data:    result.Data,
	}, nil
}

func (s *Service) Delete(ctx *saiTypes.RequestCtx, req types.DeleteRequest) (resp types.DeleteResponse, err error) {
	if err := s.validator.Struct(req); err != nil {
		return types.DeleteResponse{}, saiTypes.WrapError(err, "validation failed")
	}

	storageRequest := storageTypes.DeleteDocumentsRequest{
		Collection: req.Prefix + "_" + s.config.Collection,
		Filter:     req.Filter,
	}

	storageResult, _, err := sai.ClientManager().Call("storage", "DELETE", "/api/v1/documents", storageRequest, nil)
	if err != nil {
		return resp, err
	}

	var result storageTypes.DeleteDocumentsResponse
	err = ctx.Unmarshal(storageResult, &result)
	if err != nil {
		return resp, err
	}

	return types.DeleteResponse{
		Deleted: result.Deleted,
		Data:    result.Data,
	}, nil
}
