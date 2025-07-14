package internal

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/saiset-co/sai-crud/types"
	saiTypes "github.com/saiset-co/sai-service/types"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(ctx *saiTypes.RequestCtx) {
	var req types.CreateRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Error(err, http.StatusInternalServerError)
		return
	}

	response, err := h.service.Create(ctx, req)
	if err != nil {
		ctx.Error(err, http.StatusInternalServerError)
		return
	}

	ctx.SuccessJSON(response)
}

func (h *Handler) Read(ctx *saiTypes.RequestCtx) {
	req := types.ReadRequest{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Error(saiTypes.WrapError(err, "Invalid JSON in request body"), fasthttp.StatusBadRequest)
		return
	}

	response, err := h.service.Read(ctx, req)
	if err != nil {
		ctx.Error(err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SuccessJSON(response)
}

func (h *Handler) Update(ctx *saiTypes.RequestCtx) {
	var req types.UpdateRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Error(saiTypes.WrapError(err, "Invalid JSON in request body"), fasthttp.StatusInternalServerError)
		return
	}

	response, err := h.service.Update(ctx, req)
	if err != nil {
		ctx.Error(err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SuccessJSON(response)
}

func (h *Handler) Delete(ctx *saiTypes.RequestCtx) {
	var req types.DeleteRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Error(saiTypes.WrapError(err, "Invalid JSON in request body"), fasthttp.StatusBadRequest)
		return
	}

	response, err := h.service.Delete(ctx, req)
	if err != nil {
		ctx.Error(err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SuccessJSON(response)
}
