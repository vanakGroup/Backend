package base

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type singleMessageResponse struct {
	Status  int         `json:"statusCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type multipleMessageResponse struct {
	Status  int                          `json:"statusCode"`
	Message map[string]map[string]string `json:"message"`
	Data    interface{}                  `json:"data"`
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
			singleMessageResponse{
				Status:  statusCode,
				Message: msgType,
				Data:    data,
			},
		)
	case map[string]map[string]string:
		ctx.JSON(
			statusCode,
			multipleMessageResponse{
				Status:  statusCode,
				Message: msgType,
				Data:    data,
			},
		)
	}
}
