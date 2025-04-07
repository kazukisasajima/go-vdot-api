package usecase

import (
	"fmt"
	"go_vdot_api/model"
	"go_vdot_api/repository"
	"go_vdot_api/validator"
	"math"
	"strconv"
	"strings"

	"go_vdot_api/pkg/logger"
)

type IVdotUsecase interface {
	CreateVdot(vdot model.Vdot) (model.VdotResponse, error)
	GetVdot(userId uint) (model.VdotResponse, error)
	UpdateVdot(vdot model.Vdot, userId uint, vdotId uint) (model.VdotResponse, error)
	GetUserVdotValue(userId uint) (map[string]interface{}, error)
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

func (vu *vdotUsecase) GetVdot(userId uint) (model.VdotResponse, error) {
	vdot := model.Vdot{}
	if err := vu.vr.GetVdot(&vdot, userId); err != nil {
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

const COEFF1 float64 = 0.1894393
const COEFF2 float64 = -0.012788
const COEFF3 float64 = 0.2989558
const COEFF4 float64 = -0.1932605

func (vu *vdotUsecase) GetUserVdotValue(userId uint) (map[string]interface{}, error) {
	// ユーザーIDに基づいてVDOT値を取得
	vdot := model.Vdot{}
	if err := vu.vr.GetVdot(&vdot, userId); err != nil {
		return nil, fmt.Errorf("vdot data not found: %v", err)
	}
	logger.Info("vdot", "vdot", vdot)

	// 距離と時間の変換
	distance, err := DistanceUnitConvert(vdot)
	if err != nil {
		return nil, fmt.Errorf("failed to convert distance: %v", err)
	}
	logger.Info("distance", "distance", distance)

	timeInMinutes, err := TimeUnitConvert(vdot)
	if err != nil {
		return nil, fmt.Errorf("failed to convert time: %v", err)
	}
	logger.Info("timeInMinutes", "timeInMinutes", timeInMinutes)

	// 各種計算
    velocity := CalculateVelocity(distance, timeInMinutes)
    vo2max := CalculateVo2Max(timeInMinutes)
    vdotValue := CalculateVdot(vo2max, velocity)
    paceZones := CalculatePaceZones(velocity)
    raceTimes := PredictRaceTimes(vdot)
	
	// 結果をマップにまとめる
	data := map[string]interface{}{
        "id":            vdot.ID,
        "distanceValue": vdot.DistanceValue,
        "distanceUnit":  vdot.DistanceUnit,
        "time":          vdot.Time,
        "elevation":     vdot.Elevation,
        "temperature":   vdot.Temperature,
        "pace_zones":    paceZones,
        "VDOT":          vdotValue,
        "race_times":    raceTimes,		
	}

	return data, nil
}

func TimeUnitConvert(vdot model.Vdot) (float64, error) {
	tokens := strings.Split(vdot.Time, ":")
	if len(tokens) != 3 {
		logger.Error("Time parsing failed", "value", vdot.Time)
		return 0, fmt.Errorf("invalid time format: %s", vdot.Time)
	}
	hh, err1 := strconv.Atoi(tokens[0])
	mm, err2 := strconv.Atoi(tokens[1])
	ss, err3 := strconv.Atoi(tokens[2])
	if err1 != nil || err2 != nil || err3 != nil {
		logger.Error("Time conversion failed", "hh", tokens[0], "mm", tokens[1], "ss", tokens[2])
		return 0, fmt.Errorf("invalid time values")
	}

	totalMinutes := float64(hh)*60 + float64(mm) + float64(ss)/60
	return totalMinutes, nil
}


func DistanceUnitConvert(vdot model.Vdot) (float64, error) {
	distance_value := vdot.DistanceValue
	if distance_value < 0 {
		return 0, fmt.Errorf("invalid distance value")
	}
	distance_unit := vdot.DistanceUnit

	if distance_unit == "km" {
		return distance_value * 1000, nil
	} else if distance_unit == "mile" {
		return distance_value * 1609.34, nil
	} else if distance_unit == "m" {
		return distance_value / 1000, nil
	}

	return distance_value, nil
}

func CalculateVelocity(distance float64, timeInMinutes float64) float64 {
	velocity := distance / timeInMinutes
	return velocity
}

func CalculateVo2Max(time_in_minutes float64) float64 {
	VO2max_percentage := 0.8 + COEFF1 * math.Exp(COEFF2 * time_in_minutes) + COEFF3 * math.Exp(COEFF4 * time_in_minutes)
    return VO2max_percentage
}

func CalculateVdot(vo2max float64, velocity float64) float64 {
	vo2 := -4.6 + (0.182258 * velocity) + (0.000104 * math.Pow(velocity, 2))
	vdot := math.Round(vo2 / vo2max)
	return vdot
}

func CalculatePaceZones(velocity float64) []map[string][]map[string]map[string]string {
    zoneOrder := []string{"E", "M", "T", "I", "R"}
    distanceOrder := []string{"1mi", "1Km", "1200m", "800m", "600m", "400m", "300m", "200m"}

    zones := map[string][2]float64{
        "E": {70, 77},
        "M": {88, 0},
        "T": {92.5, 0},
        "I": {100.5, 0},
        "R": {108.25, 0},
    }

    distances := map[string]float64{
        "1mi":   1609.34,
        "1Km":   1000,
        "1200m": 1200,
        "800m":  800,
        "600m":  600,
        "400m":  400,
        "300m":  300,
        "200m":  200,
    }

    var orderedZones []map[string][]map[string]map[string]string

    for _, zone := range zoneOrder {
        bounds := zones[zone]
        lowerBound := bounds[0]
        upperBound := bounds[1]

        // 距離ごとのデータを順序付きで保持
        var orderedDistances []map[string]map[string]string

        for _, distance := range distanceOrder {
            distanceM := distances[distance]
            lowerPace := CalculatePace(velocity, lowerBound, distanceM)
            var upperPace float64
            if upperBound != 0 {
                upperPace = CalculatePace(velocity, upperBound, distanceM)
            }

            // このdistanceだけのデータを map としてまとめて配列に追加
            paceData := map[string]map[string]string{
                distance: {
                    "lower_pace": FormatPace(lowerPace),
                    "upper_pace": FormatPace(upperPace),
                },
            }
            orderedDistances = append(orderedDistances, paceData)
        }

        orderedZones = append(orderedZones, map[string][]map[string]map[string]string{
            zone: orderedDistances,
        })
    }

    return orderedZones
}

func CalculatePace(velocity, vo2maxPercentage, distance float64) float64 {
	pace := distance / (velocity * (vo2maxPercentage / 100))
	return pace
}

func FormatPace(pace float64) string {
	if pace > 0 {
		// pace（分）をmm:ss形式に変換
		minutes := int(pace)
		seconds := int((pace - float64(minutes)) * 60)
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return ""
}

func PacePerKm(timeMinutes, distance float64) string {
	paceMinutes := timeMinutes / (distance / 1000)
	minutes := int(paceMinutes)
	seconds := int((paceMinutes - float64(minutes)) * 60)
	return fmt.Sprintf("%02d:%02d /km", minutes, seconds)
}


type RaceTime struct {
	Race          string `json:"race"`
	PredictedTime string `json:"predicted_time"`
	PacePerKm     string `json:"pace_per_km"`
}

func PredictRaceTimes(vdot model.Vdot) []RaceTime {
	targetTime, _ := TimeUnitConvert(vdot)
	targetDistance, _ := DistanceUnitConvert(vdot)

	type Distance struct {
		Race     string
		Distance float64
	}

	distances := []Distance{
		{"マラソン", 42195},
		{"ハーフマラソン", 21097.5},
		{"30Km", 30000},
		{"10Mile", 16093.4},
		{"15Km", 15000},
		{"10Km", 10000},
		{"8Km", 8000},
		{"6Km", 6000},
		{"5Km", 5000},
		{"2Mile", 3218.69},
		{"3200m", 3200},
		{"3Km", 3000},
		{"1Mile", 1609.34},
		{"1600m", 1600},
		{"1500m", 1500},
	}

	var result []RaceTime

	for _, d := range distances {
		var predictedTimeMinutes float64
		if d.Distance == targetDistance {
			predictedTimeMinutes = targetTime
		} else {
			predictedTimeMinutes = targetTime * math.Pow(d.Distance/targetDistance, 1.06)
		}

		hours := int(predictedTimeMinutes / 60)
		remainder := math.Mod(predictedTimeMinutes, 60)
		minutes := int(remainder)
		seconds := int(math.Round((remainder - float64(minutes)) * 60))

		if seconds >= 60 {
			seconds -= 60
			minutes++
		}
		if minutes >= 60 {
			minutes -= 60
			hours++
		}

		formattedTime := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
		pacePerKm := predictedTimeMinutes / (d.Distance / 1000)
		formattedPace := fmt.Sprintf("%02d:%02d /km", int(pacePerKm), int((pacePerKm-float64(int(pacePerKm)))*60))

		result = append(result, RaceTime{
			Race:          d.Race,
			PredictedTime: formattedTime,
			PacePerKm:     formattedPace,
		})
	}

	return result
}
