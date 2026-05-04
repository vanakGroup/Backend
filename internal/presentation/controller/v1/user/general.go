package user

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	userService usecase.UserService
	cfg         *bootstrap.Config
	logger      logging.Logger
}

func NewGeneralUserController(
	userService usecase.UserService,
	cfg *bootstrap.Config,
) *GeneralUserController {
	logger := logging.NewLogger(cfg)

	return &GeneralUserController{
		userService: userService,
		cfg:         cfg,
		logger:      logger,
	}
}

func (userController *GeneralUserController) SendOtp(ctx *gin.Context) {
	type sendOtpParams struct {
		Phone string `json:"phone"`
	}

	params := base.Validated[sendOtpParams](ctx)

	err := userController.userService.SendOtp(params.Phone)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, userController.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.SendOtp")
	base.Response(ctx, 200, message, nil)
}
