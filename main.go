package main

import (
	"go_vdot_api/controller"
	"go_vdot_api/model"
	"go_vdot_api/repository"
	"go_vdot_api/router"
	"go_vdot_api/usecase"
	"go_vdot_api/validator"
)

func main() {
	db := model.NewDB()
	userValidator := validator.NewUserValidator()
	vdotValidator := validator.NewVdotValidator()

	userRepository := repository.NewUserRepository(db)
	vdotRepository := repository.NewVdotRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	vdotUsecase := usecase.NewVdotUsecase(vdotRepository, vdotValidator)

	userController := controller.NewUserController(userUsecase)
	vdotController := controller.NewVdotController(vdotUsecase)

	e := router.NewRouter(userController, vdotController)
	e.Logger.Fatal(e.Start(":8080"))
}
