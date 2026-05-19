package base

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SingleMessageResponse struct {
	Status  int         `json:"statusCode" example:"200"`
	Message string      `json:"message" example:"Successfully"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
}

type MultipleMessageResponse struct {
	Status  int                          `json:"statusCode" example:"422"`
	Message map[string]map[string]string `json:"message" swaggertype:"object,string"`
	Data    interface{}                  `json:"data" swaggertype:"object"`
}

func Response[T string | map[string]map[string]string](ctx *gin.Context, statusCode int, msg T, data interface{}) {
	statusCodeText := http.StatusText(statusCode)

	if statusCodeText == "" {
		panic(fmt.Errorf("invalid status code"))
	}

	switch msgType := any(msg).(type) {
	case string:
		if msgType == "" {
			msgType = statusCodeText
		}
		ctx.JSON(
			statusCode,
			SingleMessageResponse{
				Status:  statusCode,
				Message: msgType,
				Data:    data,
			},
		)
	case map[string]map[string]string:
		ctx.JSON(
			statusCode,
			MultipleMessageResponse{
				Status:  statusCode,
				Message: msgType,
				Data:    data,
			},
		)
	}
}
