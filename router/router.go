package router

import (
	"go_vdot_api/controller"
	"net/http"
	"os"
	"go_vdot_api/pkg/logger"

	mymiddleware "go_vdot_api/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, vc controller.IVdotController, wc controller.IWorkoutController, sec controller.ISpecialtyEventController) *echo.Echo {
	router := echo.New()
	router.Use(logger.RequestLogger())
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
	// CSRFトークンを取得するためのエンドポイント
	auth := router.Group("/api/auth")
	auth.POST("/signup", uc.SignUp)
	auth.POST("/login", uc.LogIn)
	auth.POST("/logout", uc.LogOut)
	auth.GET("/csrf", uc.CsrfToken)
	auth.Use(mymiddleware.JWTMiddleware())

	// ログイン確認用エンドポイント
	authCheck := auth.Group("")
	authCheck.Use(mymiddleware.JWTMiddleware())
	authCheck.GET("/check", mymiddleware.CheckAuth)

	// ユーザー情報取得用エンドポイント
	user := router.Group("/api/user")
	user.Use(mymiddleware.JWTMiddleware())
	user.PATCH("", uc.UpdateUser)
	user.DELETE("", uc.DeleteUser)
	
	// Vdot関連のエンドポイント
	vdot := router.Group("/api/vdots")
	vdot.Use(mymiddleware.JWTMiddleware())
	vdot.POST("", vc.CreateVdot)
	vdot.GET("", vc.GetVdot)
	vdot.PATCH("/:id", vc.UpdateVdot)
	vdot.GET("/value", vc.GetUserVdotValue)

	// Workout関連のエンドポイント
	workout := router.Group("/api/workouts")
	workout.Use(mymiddleware.JWTMiddleware())
	workout.POST("", wc.CreateWorkout)
	workout.GET("", wc.GetWorkoutPerMonth)
	workout.PATCH("/:id", wc.UpdateWorkout)

	// SpecialtyEvent関連のエンドポイント
	specialtyEvent := router.Group("/api/specialty_events")
	specialtyEvent.Use(mymiddleware.JWTMiddleware())
	specialtyEvent.POST("", sec.CreateSpecialtyEvent)
	specialtyEvent.GET("", sec.GetSpecialtyEvent)
	specialtyEvent.PATCH("/:id", sec.UpdateSpecialtyEvent)

	return router
}
