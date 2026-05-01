package usecase

import (
	mathdto "github.com/Vanak/Backend/internal/application/dto/math"
)

type MathService interface {
	Add(Numbers mathdto.AddRequest) (mathdto.AddResponse, error)
}
