package usecase

import (
	"go_vdot_api/model"
	"go_vdot_api/repository"
	"go_vdot_api/validator"
)

type IVdotUsecase interface {
	CreateVdot(vdot model.Vdot) (model.VdotResponse, error)
	GetVdotByID(userId uint, vdotId uint) (model.VdotResponse, error)
	UpdateVdot(vdot model.Vdot, userId uint, vdotId uint) (model.VdotResponse, error)
}

type vdotUsecase struct {
	vr repository.IVdotRepository
	vv validator.IVdotValidator
}

func NewVdotUsecase(vr repository.IVdotRepository, vv validator.IVdotValidator) IVdotUsecase {
	return &vdotUsecase{vr, vv}
}

func (vu *vdotUsecase) CreateVdot(vdot model.Vdot) (model.VdotResponse, error) {
	if err := vu.vv.VdotValidate(vdot); err != nil {
		return model.VdotResponse{}, err
	}

	if err := vu.vr.CreateVdot(&vdot); err != nil {
		return model.VdotResponse{}, err
	}
	
	resVdot := model.VdotResponse{
		ID:            vdot.ID,
		DistanceValue: vdot.DistanceValue,
		DistanceUnit:  vdot.DistanceUnit,
		Time:          vdot.Time,
		Elevation:     vdot.Elevation,
		Temperature:   vdot.Temperature,
	}
	return resVdot, nil
}

func (vu *vdotUsecase) GetVdotByID(userId uint, vdotId uint) (model.VdotResponse, error) {
	vdot := model.Vdot{}
	if err := vu.vr.GetVdotByID(&vdot, userId, vdotId); err != nil {
		return model.VdotResponse{}, err
	}
	resVdot := model.VdotResponse{
		ID:            vdot.ID,
		DistanceValue: vdot.DistanceValue,
		DistanceUnit:  vdot.DistanceUnit,
		Time:          vdot.Time,
		Elevation:     vdot.Elevation,
		Temperature:   vdot.Temperature,
	}
	return resVdot, nil
}

func (vu *vdotUsecase) UpdateVdot(vdot model.Vdot, userId uint, vdotId uint) (model.VdotResponse, error) {
	if err := vu.vv.VdotValidate(vdot); err != nil {
		return model.VdotResponse{}, err
	}

	if err := vu.vr.UpdateVdot(&vdot, userId, vdotId); err != nil {
		return model.VdotResponse{}, err
	}
	resVdot := model.VdotResponse{
		ID:            vdot.ID,
		DistanceValue: vdot.DistanceValue,
		DistanceUnit:  vdot.DistanceUnit,
		Time:          vdot.Time,
		Elevation:     vdot.Elevation,
		Temperature:   vdot.Temperature,
	}
	return resVdot, nil
}
