package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	V1Domains "github.com/kelvinator07/go-rest-template/internal/business/domains/v1"
	"github.com/kelvinator07/go-rest-template/internal/constants"
	"github.com/kelvinator07/go-rest-template/internal/datasources/caches"
	"github.com/kelvinator07/go-rest-template/internal/http/datatransfers/requests"
	"github.com/kelvinator07/go-rest-template/internal/http/datatransfers/responses"
	"github.com/kelvinator07/go-rest-template/pkg/jwt"
	"github.com/kelvinator07/go-rest-template/pkg/validators"
)

type UserHandler struct {
	usecase        V1Domains.UserUsecase
	redisCache     caches.RedisCache
	ristrettoCache caches.RistrettoCache
}

func NewUserHandler(usecase V1Domains.UserUsecase, redisCache caches.RedisCache, ristrettoCache caches.RistrettoCache) UserHandler {
	return UserHandler{
		usecase:        usecase,
		redisCache:     redisCache,
		ristrettoCache: ristrettoCache,
	}
}

func (userH UserHandler) Register(ctx *gin.Context) {
	var UserRegisterRequest requests.UserRequest
	if err := ctx.ShouldBindJSON(&UserRegisterRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(UserRegisterRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := UserRegisterRequest.ToV1Domain()
	userDomainn, statusCode, err := userH.usecase.Store(ctx.Request.Context(), userDomain)

	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go userH.ristrettoCache.Del("users")

	NewSuccessResponse(ctx, statusCode, "registration user success", map[string]interface{}{
		"user": responses.FromV1Domain(userDomainn),
	})
}

func (userH UserHandler) Login(ctx *gin.Context) {
	var UserLoginRequest requests.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(UserLoginRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain, statusCode, err := userH.usecase.Login(ctx.Request.Context(), UserLoginRequest.ToV1Domain())
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	NewSuccessResponse(ctx, statusCode, "login success", responses.FromV1Domain(userDomain))
}

func (userH UserHandler) SendOTP(ctx *gin.Context) {
	var userOTP requests.UserSendOTPRequest

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(userOTP); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	otpCode, statusCode, err := userH.usecase.SendOTP(ctx.Request.Context(), userOTP.Email)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	go userH.redisCache.Set(otpKey, otpCode)

	NewSuccessResponse(ctx, statusCode, fmt.Sprintf("otp code has been send to %s", userOTP.Email), nil)
}

func (userH UserHandler) VerifyOTP(ctx *gin.Context) {
	var userOTP requests.UserVerifyOTPRequest

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validators.ValidatePayloads(userOTP); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	otpRedis, err := userH.redisCache.Get(otpKey)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	statusCode, err := userH.usecase.VerifyOTP(ctx.Request.Context(), userOTP.Email, userOTP.Code, otpRedis)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	statusCode, err = userH.usecase.ActivateUser(ctx.Request.Context(), userOTP.Email)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go userH.redisCache.Del(otpKey)
	go userH.ristrettoCache.Del("users")

	NewSuccessResponse(ctx, statusCode, "otp verification success", nil)
}

func (u UserHandler) GetUserData(ctx *gin.Context) {
	// get authenticated user from context
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)
	if val := u.ristrettoCache.Get(fmt.Sprintf("user/%s", userClaims.Email)); val != nil {
		NewSuccessResponse(ctx, http.StatusOK, "user data fetched successfully", map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	userDom, statusCode, err := u.usecase.GetByEmail(ctxx, userClaims.Email)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.FromV1Domain(userDom)

	go u.ristrettoCache.Set(fmt.Sprintf("user/%s", userClaims.Email), userResponse)

	NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"user": userResponse,
	})
}

func (u UserHandler) GetAllUsers(ctx *gin.Context) {
	if val := u.ristrettoCache.Get("users"); val != nil {
		NewSuccessResponse(ctx, http.StatusOK, "users data fetched successfully", map[string]interface{}{
			"users": val,
		})
		return
	}

	usersDom, statusCode, err := u.usecase.GetAllUsers(ctx.Request.Context())
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	usersResponse := responses.ToResponseList(usersDom)

	go u.ristrettoCache.Set("users", usersResponse)

	NewSuccessResponse(ctx, statusCode, "users data fetched successfully", map[string]interface{}{
		"users": usersResponse,
	})

}
