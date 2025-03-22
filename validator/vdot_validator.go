package validator

import (
	"go_vdot_api/model"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IVdotValidator interface {
	VdotValidate(vdot model.Vdot) error
}

type vdotValidator struct{}

func NewVdotValidator() IVdotValidator {
	return &vdotValidator{}
}

func (vv *vdotValidator) VdotValidate(vdot model.Vdot) error {
	return validation.ValidateStruct(&vdot,
		validation.Field(
			&vdot.DistanceValue,
			validation.Required.Error("distance value is required"),
		),
		validation.Field(
			&vdot.DistanceUnit,
			validation.Required.Error("distance unit is required"),
		),
		validation.Field(
			&vdot.Time,
			validation.Required.Error("time is required"),
			validation.Match(regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)).Error("time must be in HH:MM:SS format"),
		),
	)
}
