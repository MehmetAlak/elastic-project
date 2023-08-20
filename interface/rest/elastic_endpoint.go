package rest

import (
	"elastic-project/application/elastic_operation"
	"elastic-project/interface/rest/helper"
	"elastic-project/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type elasticsearchEndpoint struct {
	elasticsearchService elastic_operation.Service
}

type ElasticsearchEndpoint interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Find() gin.HandlerFunc
	Delete() gin.HandlerFunc
	FindByKeyAndValue() gin.HandlerFunc
	FindByJsonQuery() gin.HandlerFunc
}

func NewElasticsearchEndpoint(elasticsearchService elastic_operation.Service) ElasticsearchEndpoint {
	return &elasticsearchEndpoint{elasticsearchService: elasticsearchService}
}

// Create godoc
// @Summary create user
// @Description creates
// @Tags elastic
// @Accept json
// @Param body body model.CreateRequest true "CreateRequest"
// @Success 201
// @Router /users [post]
func (endpoint *elasticsearchEndpoint) Create() gin.HandlerFunc {
	return func(context *gin.Context) {
		var requestBody model.CreateRequest

		if err := context.BindJSON(&requestBody); err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: 400,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		createResponse, err := endpoint.elasticsearchService.Create(context, requestBody)

		if err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: 500,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}
		//fmt.Sprintf("response: %v", createResponse)

		//context.Status(model.StatusNoContent)
		context.JSON(http.StatusCreated, createResponse)
	}
}

// Update godoc
// @Summary update user
// @Description update user
// @Tags elastic
// @Accept json
// @Param id path string true "id"
// @Param body body model.UpdateRequest true "UpdateRequest"
// @Success 204
// @Router /users/{id} [put]
func (endpoint *elasticsearchEndpoint) Update() gin.HandlerFunc {
	return func(context *gin.Context) {

		userId := context.Param("id")
		if userId == "" {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: 400,
				Err:        errors.New(fmt.Sprintf("invalid userId")),
			})
			return
		}
		var requestBody model.UpdateRequest

		if err := context.BindJSON(&requestBody); err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: 400,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		err := endpoint.elasticsearchService.Update(context, userId, requestBody)

		if err != nil {
			statusCode := http.StatusInternalServerError
			if model.ErrConflict == err {
				statusCode = http.StatusConflict
			}
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: statusCode,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		context.Status(model.StatusNoContent)
	}
}

// Delete godoc
// @Summary delete user
// @Description delete user
// @Tags elastic
// @Accept json
// @Param id path string true "id"
// @Success 204
// @Router /users/{id} [delete]
func (endpoint *elasticsearchEndpoint) Delete() gin.HandlerFunc {
	return func(context *gin.Context) {

		idParam := context.Param("id")

		err := endpoint.elasticsearchService.Delete(context, model.DeleteRequest{ID: idParam})

		if err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		context.Status(model.StatusNoContent)
	}
}

// Find godoc
// @Summary gets user
// @Description gets user
// @Tags elastic
// @Accept json
// @Param id query string true "id"
// @Success 200 {object} model.FindResponse
// @Router /users [get]
func (endpoint *elasticsearchEndpoint) Find() gin.HandlerFunc {
	return func(context *gin.Context) {
		idParam := context.Query("id")

		response, err := endpoint.elasticsearchService.Find(context, model.FindRequest{ID: idParam})

		if err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		context.JSON(http.StatusOK, response)
	}
}

// FindByKeyAndValue godoc
// @Summary gets user list
// @Description gets user list
// @Tags elastic
// @Accept json
// @Param queryType query string true "queryType" Enums(match, wildcard, match_phrase_prefix, regexp, fuzzy) default(match)
// @Param key query string true "key"
// @Param value query string true "value"
// @Success 200 {object} []model.FindResponse
// @Router /users-by [get]
func (endpoint *elasticsearchEndpoint) FindByKeyAndValue() gin.HandlerFunc {
	return func(context *gin.Context) {
		queryTypeParam := context.Query("queryType")
		keyParam := context.Query("key")
		valueParam := context.Query("value")

		response, err := endpoint.elasticsearchService.FindByKeyAndValue(model.FindByRequest{QueryType: queryTypeParam, Key: keyParam, Value: valueParam})

		if err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		context.JSON(http.StatusOK, response)
	}
}

// FindByJsonQuery godoc
// @Summary gets user list with query
// @Description gets user list with query
// @Tags elastic
// @Accept json
// @Param jsonQuery query string true "jsonQuery"
// @Success 200 {object} []model.FindResponse
// @Router /users-by-query [get]
func (endpoint *elasticsearchEndpoint) FindByJsonQuery() gin.HandlerFunc {
	return func(context *gin.Context) {
		jsonQueryParam := context.Query("jsonQuery")

		response, err := endpoint.elasticsearchService.FindByQuery(jsonQueryParam)

		if err != nil {
			helper.HandleEndpointError(context, &model.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(fmt.Sprintf("invalid request: Error: %v", err.Error())),
			})
			return
		}

		context.JSON(http.StatusOK, response)
	}
}
