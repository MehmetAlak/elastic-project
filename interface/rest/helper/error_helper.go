package helper

import (
	"elastic-project/model"
	"github.com/gin-gonic/gin"
)

func HandleEndpointError(context *gin.Context, responseError *model.ResponseError) {
	context.JSON(responseError.StatusCode, model.ErrorDto{Message: responseError.Err.Error()})
}
