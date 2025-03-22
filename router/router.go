package router

import (
	"go_vdot_api/controller"
	"net/http"
	"os"

	mymiddleware "go_vdot_api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, vc controller.IVdotController) *echo.Echo {
	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteNoneMode,
		CookieSameSite: http.SameSiteDefaultMode,
		//CookieMaxAge:   60,
	}))
	// ユーザー用エンドポイント
	auth := router.Group("/api/auth")
	auth.POST("/signup", uc.SignUp)
	auth.POST("/login", uc.LogIn)
	auth.POST("/logout", uc.LogOut)
	auth.GET("/csrf", uc.CsrfToken)

	vdot := router.Group("/api/vdots")
	vdot.Use(mymiddleware.JWTMiddleware())
	vdot.POST("", vc.CreateVdot)
	vdot.GET("/:id", vc.GetVdotByID)
	vdot.PUT("/:id", vc.UpdateVdot)

	return router
}