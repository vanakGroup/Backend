package user

import (
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	userdto "github.com/Vanak/Backend/internal/application/dto/user"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"

	"github.com/gin-gonic/gin"
)

type CustomerUserController struct {
	cfg     *bootstrap.Config
	useCase usecase.UserProfileUseCase
}

func NewCustomerUserController(cfg *bootstrap.Config, useCase usecase.UserProfileUseCase) *CustomerUserController {
	return &CustomerUserController{
		cfg:     cfg,
		useCase: useCase,
	}
}

// UpdateProfile godoc
// @Summary      Update Customer Profile
// @Description  Updates the user's profile information. Phone number cannot be changed here. Requires Bearer Token.
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{firstname=string,lastname=string,national_code=string,shaba_number=string} true "Profile update payload" example(object{firstname="Amir",lastname="Mohammadi",national_code="1234567890",shaba_number="IR123456789012345678901234"})
// @Success      200 {object} base.SingleMessageResponse "Profile updated successfully"
// @Failure      400 {object} base.MultipleMessageResponse "Validation error"
// @Failure      401 {object} base.SingleMessageResponse "Unauthorized - Missing or invalid token"
// @Router       /profile [put]
func (cuc *CustomerUserController) UpdateProfile(ctx *gin.Context) {
	type CustomerFilledParams struct {
		FirstName    *string `json:"firstname"`
		LastName     *string `json:"lastname"`
		NationalCode *string `json:"national_code" validate:"omitempty"`
		ShabaNumber  *string `json:"shaba_number" validate:"omitempty"`
	}

	params := base.Validated[CustomerFilledParams](ctx)
	userIDAny, exists := ctx.Get(cuc.cfg.Constant.Context.UserID)
	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cuc.cfg.Constant.ErrorField.User, cuc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	userInfo := userdto.UserFilled{
		UserID:       userID,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		NationalCode: params.NationalCode,
		ShabaNumber:  params.ShabaNumber,
	}

	err := cuc.useCase.UpdateInformation(ctx.Request.Context(), userInfo)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cuc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.updateProfile")
	base.Response(ctx, 200, message, nil)
}

// GetProfile    godoc
// @Summary      Get User Profile
// @Description  Retrieves current authenticated user's profile info. Requires Bearer Token.
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} base.SingleMessageResponse{data=userdto.ProfileDTO}
// @Router       /profile [get]
func (cuc *CustomerUserController) GetProfile(ctx *gin.Context) {
	userIDAny, exists := ctx.Get(cuc.cfg.Constant.Context.UserID)
	if !exists {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", cuc.cfg.Constant.ErrorField.User, cuc.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	userID := userIDAny.(uint)

	userProfile, err := cuc.useCase.GetInformation(userID)
	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, cuc.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.getProfile")
	base.Response(ctx, 200, message, userProfile)
}
