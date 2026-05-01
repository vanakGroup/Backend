package mathService

import (
	addDto "github.com/Vanak/Backend/internal/application/dto/math"
	mathdto "github.com/Vanak/Backend/internal/application/dto/math"
)

type MathService struct{}

func NewMathService() *MathService {
	return &MathService{}
}

func (ms *MathService) Add(numbers mathdto.AddRequest) (mathdto.AddResponse, error) {
	return addDto.AddResponse{
		Answer: numbers.FirstNumber + numbers.SecondNumber,
	}, nil
}
