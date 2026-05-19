package user

import (
	"github.com/Vanak/Backend/bootstrap"
	userdto "github.com/Vanak/Backend/internal/application/dto/user"
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

// SendOtp godoc
// @Summary      Send OTP Verification Code
// @Description  Validates the Iranian phone number and sends a temporary verification SMS
// @Tags         Authentication
// @Param        request  body      object{phone=string}  true  "User's mobile phone number"  example(object{phone="09123456789"})
// @Success      200      {object}  base.SingleMessageResponse
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /auth/send-otp [post]
func (userController *GeneralUserController) SendOtp(ctx *gin.Context) {
	type sendOtpParams struct {
		Phone string `json:"phone" validate:"required,ir_phone"`
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

// VerifyOTP godoc
// @Summary      Verify OTP Code
// @Description  Verifies the temporary OTP code sent to the user's mobile phone number
// @Tags         Authentication
// @Param        request  body      object{phone=string,otp=string}  true  "Phone number and received OTP code"  example(object{phone="09123456789",otp="123456"})
// @Success      200      {object}  base.SingleMessageResponse   "Phone number verified successfully"
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /auth/verify-otp [post]
func (userController *GeneralUserController) VerifyOTP(ctx *gin.Context) {
	type VerifyPhoneOtp struct {
		Phone string `json:"phone" validate:"required,ir_phone"`
		Otp   string `json:"otp" validate:"required"`
	}

	params := base.Validated[VerifyPhoneOtp](ctx)

	err := userController.userService.VerifyOTP(userdto.VerifyPhoneOtp(params))

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, userController.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.VerifyOtp")
	base.Response(ctx, 200, message, nil)

}

// Login godoc
// @Summary      User Login with Phone
// @Description  Authenticate user into the Vanak Parking System using phone number and returns user session info
// @Tags         Authentication
// @Param        request  body      object{phone=string}  true  "User's mobile phone number"  example(object{phone="09123456789"})
// @Success      200      {object}  base.SingleMessageResponse "Successfully logged in and session created"
// @Failure      400      {object}  base.MultipleMessageResponse
// @Router       /auth/login [post]
func (userController *GeneralUserController) Login(ctx *gin.Context) {
	type LoginWithPhone struct {
		Phone string `json:"phone" validate:"required,ir_phone"`
	}

	params := base.Validated[LoginWithPhone](ctx)

	userInfo, err := userController.userService.Login(params.Phone)

	if err != nil {
		panic(err)
	}

	trans := base.GetTranslator(ctx, userController.cfg.Constant.Context.Translator)
	message, _ := trans.Translate("successMessage.login")
	base.Response(ctx, 200, message, userInfo)
}
