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
	workoutValidator := validator.NewWorkoutValidator()

	userRepository := repository.NewUserRepository(db)
	vdotRepository := repository.NewVdotRepository(db)
	workoutRepository := repository.NewWorkoutRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	vdotUsecase := usecase.NewVdotUsecase(vdotRepository, vdotValidator)
	workoutUsecase := usecase.NewWorkoutUsecase(workoutRepository, workoutValidator)

	userController := controller.NewUserController(userUsecase)
	vdotController := controller.NewVdotController(vdotUsecase)
	workoutController := controller.NewWorkoutController(workoutUsecase)

	e := router.NewRouter(userController, vdotController, workoutController)
	e.Logger.Fatal(e.Start(":8080"))
}
