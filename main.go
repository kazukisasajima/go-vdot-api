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
	SpecialtyEventValidator := validator.NewSpecialtyEventValidator()

	userRepository := repository.NewUserRepository(db)
	vdotRepository := repository.NewVdotRepository(db)
	workoutRepository := repository.NewWorkoutRepository(db)
	specialtyEventRepository := repository.NewSpecialtyEventRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	vdotUsecase := usecase.NewVdotUsecase(vdotRepository, vdotValidator)
	workoutUsecase := usecase.NewWorkoutUsecase(workoutRepository, workoutValidator)
	specialtyEventUsecase := usecase.NewSpecialtyEventUsecase(specialtyEventRepository, SpecialtyEventValidator)

	userController := controller.NewUserController(userUsecase)
	vdotController := controller.NewVdotController(vdotUsecase)
	workoutController := controller.NewWorkoutController(workoutUsecase)
	specialtyEventController := controller.NewSpecialtyEventController(specialtyEventUsecase)

	e := router.NewRouter(userController, vdotController, workoutController, specialtyEventController)
	e.Logger.Fatal(e.Start(":8080"))
}
