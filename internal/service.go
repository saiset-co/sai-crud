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

	collection := req.Prefix
	if s.config.Collection != "" {
		collection += "_" + s.config.Collection
	}

	storageRequest := storageTypes.CreateDocumentsRequest{
		Collection: collection,
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

	collection := req.Prefix
	if s.config.Collection != "" {
		collection += "_" + s.config.Collection
	}

	storageRequest := storageTypes.ReadDocumentsRequest{
		Collection: collection,
		Filter:     req.Filter,
		Limit:      req.Limit,
		Sort:       req.Sort,
		Skip:       req.Skip,
		Count:      req.Count,
		Fields:     req.IncludeFields,
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

	collection := req.Prefix
	if s.config.Collection != "" {
		collection += "_" + s.config.Collection
	}

	storageRequest := storageTypes.UpdateDocumentsRequest{
		Collection: collection,
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

	collection := req.Prefix
	if s.config.Collection != "" {
		collection += "_" + s.config.Collection
	}

	storageRequest := storageTypes.DeleteDocumentsRequest{
		Collection: collection,
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

func (s *Service) Aggregate(ctx *saiTypes.RequestCtx, req types.AggregateRequest) (resp types.AggregateResponse, err error) {
	if err := s.validator.Struct(req); err != nil {
		return types.AggregateResponse{}, saiTypes.WrapError(err, "validation failed")
	}

	collection := req.Prefix
	if s.config.Collection != "" {
		collection += "_" + s.config.Collection
	}

	storageRequest := storageAggregateRequest{
		Collection: collection,
		Pipeline:   req.Pipeline,
		Filter:     req.Filter,
		GroupBy:    req.GroupBy,
		Aggregates: convertAggregateFields(req.Aggregates),
		Sort:       req.Sort,
		Limit:      req.Limit,
		Skip:       req.Skip,
		Fields:     req.Fields,
		Count:      req.Count,
	}

	storageResult, _, err := sai.ClientManager().Call("storage", "POST", "/api/v1/documents/aggregate", storageRequest, nil)
	if err != nil {
		return resp, err
	}

	var result storageAggregateResponse
	err = ctx.Unmarshal(storageResult, &result)
	if err != nil {
		return resp, err
	}

	return types.AggregateResponse{
		Data:  result.Data,
		Total: result.Total,
	}, nil
}

type storageAggregateField struct {
	Field string `json:"field,omitempty"`
	Op    string `json:"op"`
	As    string `json:"as,omitempty"`
}

type storageAggregateRequest struct {
	Collection string                  `json:"collection"`
	Pipeline   types.OrderedPipeline   `json:"pipeline,omitempty"`
	Filter     map[string]interface{}  `json:"filter,omitempty"`
	GroupBy    []string                `json:"group_by,omitempty"`
	Aggregates []storageAggregateField `json:"aggregates,omitempty"`
	Sort       map[string]int          `json:"sort,omitempty"`
	Limit      int                     `json:"limit,omitempty"`
	Skip       int                     `json:"skip,omitempty"`
	Fields     []string                `json:"fields,omitempty"`
	Count      int                     `json:"count,omitempty"`
}

type storageAggregateResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int64                    `json:"total"`
}

func convertAggregateFields(fields []types.AggregateField) []storageAggregateField {
	if len(fields) == 0 {
		return nil
	}
	result := make([]storageAggregateField, 0, len(fields))
	for _, field := range fields {
		result = append(result, storageAggregateField{
			Field: field.Field,
			Op:    field.Op,
			As:    field.As,
		})
	}
	return result
}
