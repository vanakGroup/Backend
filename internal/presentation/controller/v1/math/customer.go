package math

import (
	"github.com/Vanak/Backend/bootstrap"
	addDto "github.com/Vanak/Backend/internal/application/dto/math"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/presentation/controller/base"

	"github.com/gin-gonic/gin"
)

type CustomerMathController struct {
	constant    *bootstrap.Constants
	mathService usecase.MathService
}

func NewCustomerMathController(
	mathService usecase.MathService,
	constant *bootstrap.Constants,
) *CustomerMathController {
	return &CustomerMathController{
		mathService: mathService,
		constant:    constant,
	}
}

func (mathController *CustomerMathController) Adder(ctx *gin.Context) {
	type AddNumber struct {
		FirstNumber  int `uri:"first"  validate:"required,gte=85,lte=100"`
		SecondNumber int `uri:"second" validate:"required,gte=23,lte=200"`
	}

	params := base.Validated[AddNumber](ctx)

	addReqeust := addDto.AddRequest{
		FirstNumber:  params.FirstNumber,
		SecondNumber: params.SecondNumber,
	}

	answer, err := mathController.mathService.Add(addReqeust)

	if err != nil {
		panic(err)
	}
	translator := base.GetTranslator(ctx, mathController.constant.Context.Translator)
	message, _ := translator.Translate("successMessage.addNumber")
	base.Response(ctx, 200, message, answer)
}
