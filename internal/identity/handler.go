package identity

import (
	"div-dash/internal/db"
	"div-dash/internal/httputil"
	"div-dash/internal/logging"
	"div-dash/internal/mail"
	"div-dash/internal/services"
	"div-dash/internal/token"
	"div-dash/util/security"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	handlerDependencies interface {
		logging.LoggerProvider
		db.QueriesProvider
		token.TokenServiceProvider
		mail.MailServiceProvider
	}

	HandlerProvider interface {
		LoginHandler() *Handler
	}

	Handler struct {
		handlerDependencies
	}
)

func NewHandler(h handlerDependencies) *Handler {
	return &Handler{
		handlerDependencies: h,
	}
}

func (h *Handler) RegisterPublicRoutes(api gin.IRoutes) {
	api.POST("/login", h.postLogin)
	api.POST("/register", h.postRegister)
	api.GET("/activate", h.postActivate)
}

func (h *Handler) RegisterPrivateRoutes(api gin.IRoutes) {
	api.POST("/auth/logout", h.getLogout)
	api.POST("/auth/identity", h.getLogout)

	api.GET("/user/:id", h.getUser)
}

func (h *Handler) RegisterMiddleware(api gin.IRoutes) {
	api.Use(h.authRequired)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) postLogin(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		httputil.AbortBadRequest(c, err.Error())
		return
	}

	user, err := h.Queries().FindByEmail(c, loginRequest.Email)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			httputil.Abort(c, http.StatusUnauthorized, "wrong credentials")
			return
		}
		c.Error(err)
		return
	}

	if user.Status != db.UserStatusActivated {
		httputil.Abort(c, http.StatusUnauthorized, "User not activated")
		return
	}

	if !security.VerifyHash(loginRequest.Password, user.PasswordHash) {
		httputil.Abort(c, http.StatusUnauthorized, "wrong credentials")
		return
	}

	token, err := h.TokenService().GenerateToken(user.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetCookie("token", token, int(time.Hour.Seconds())*24, "/api", "localhost", true, true)
	c.Status(http.StatusOK)
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) postRegister(c *gin.Context) {

	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		httputil.AbortBadRequest(c, err.Error())
		return
	}

	exists, err := h.Queries().ExistsByEmail(c, registerRequest.Email)

	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		httputil.Abort(c, http.StatusConflict, "A user with email '"+registerRequest.Email+"' already exists")
		return
	}

	passwordHash, err := security.HashPassword(registerRequest.Password)

	if err != nil {
		c.Error(err)
		return
	}

	registerRequestId, err := services.IdService().NewUUID()

	if err != nil {
		c.Error(err)
		return
	}

	user, err := h.Queries().CreateUser(c, db.CreateUserParams{
		ID:           services.IdService().NewId(16),
		Email:        registerRequest.Email,
		PasswordHash: passwordHash,
		Status:       db.UserStatusRegistered,
	})

	if err != nil {
		c.Error(err)
		return
	}

	createRegistrationParams := db.CreateUserRegistrationParams{
		ID:        registerRequestId,
		UserID:    user.ID,
		Timestamp: time.Now().UTC(),
	}

	registrationId, err := h.Queries().CreateUserRegistration(c, createRegistrationParams)
	if err != nil {
		c.Error(err)
		return
	}

	body := "Please activate your account at localhost:8080/activate?id=" + registrationId.String()

	err = h.MailService().SendMail(user.Email, "no-reply@div-dash.io", "Activate your account", body)

	if err != nil {
		c.Error(err)
		return
	}
	c.Status(200)
}

func (h *Handler) postActivate(c *gin.Context) {
	id := c.Query("id")
	registerRequest, err := uuid.Parse(id)

	if err != nil {
		httputil.AbortBadRequest(c, "Activation id is in wrong format")
		return
	}

	userRegistration, err := h.Queries().GetUserRegistration(c, registerRequest)

	if err != nil {
		httputil.AbortBadRequest(c, "Invalid id")
		return
	}

	if userRegistration.Timestamp.Add(24 * time.Hour).Before(time.Now()) {
		httputil.AbortBadRequest(c, "Registration expired")
		return
	}

	activated, err := h.Queries().IsUserActivated(c, userRegistration.UserID)
	if err != nil {
		c.Error(err)
		return
	}
	if activated {
		httputil.AbortBadRequest(c, "User already activated")
	}

	err = h.Queries().ActivateUser(c, userRegistration.UserID)

	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/api", "localhost", true, true)
	c.Status(http.StatusOK)
}

func (h *Handler) GetAuthIdentity(c *gin.Context) {
	userId := c.GetString("userId")

	user, err := h.Queries().GetUser(c, userId)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
