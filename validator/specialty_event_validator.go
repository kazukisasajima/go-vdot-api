package validator

import (
	"errors"
	"go_vdot_api/model"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// 選択肢リスト
var EventChoices = []string{
	"800m", "1500m", "1mile", "3000m", "3000mSC",
	"2mile", "5000m", "10000m", "ハーフマラソン", "フルマラソン",
}

// []string → []interface{} に変換
func toInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

type ISpecialtyEventValidator interface {
	SpecialtyEventValidate(event model.SpecialtyEvent) error
}

type specialtyEventValidator struct{}

func NewSpecialtyEventValidator() ISpecialtyEventValidator {
	return &specialtyEventValidator{}
}

func (v *specialtyEventValidator) SpecialtyEventValidate(event model.SpecialtyEvent) error {
	err := validation.ValidateStruct(&event,
		validation.Field(&event.EventName,
			validation.Required,
			validation.In(toInterfaceSlice(EventChoices)...).Error("invalid event name"),
		),
		validation.Field(&event.BestTime, validation.Required),
	)

	if err != nil {
		return err
	}

	// best_time の形式チェック
	formats := []string{
		`^\d{1,2}:\d{2}:\d{2}$`,      // 例: 2:22:25
		`^\d{1,2}'\d{2}"(\d{1,2})?$`, // 例: 4'12"11
	}
	matched := false
	for _, f := range formats {
		if matched, _ = regexp.MatchString(f, event.BestTime); matched {
			break
		}
	}
	if !matched {
		return errors.New("invalid time format. Use hh:mm:ss or m'ss\"SS")
	}

	return nil
}
