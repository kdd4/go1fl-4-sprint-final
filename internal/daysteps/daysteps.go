package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	splitedData := strings.Split(data, ",")

	if len(splitedData) != 2 {
		return 0, time.Duration(0), errors.New("wrong len of splited data")
	}

	stepCount, err := strconv.Atoi(splitedData[0])

	if err != nil {
		joinedErr := errors.Join(errors.New("step count parsing error"), err)
		return 0, time.Duration(0), joinedErr
	}

	if stepCount <= 0 {
		return 0, time.Duration(0), errors.New("step count less or equal 0")
	}

	walkDuration, err := time.ParseDuration(splitedData[1])

	if err != nil {
		joinedErr := errors.Join(errors.New("walk duration parsing error"), err)
		return 0, time.Duration(0), joinedErr
	}

	if walkDuration <= time.Duration(0) {
		return 0, time.Duration(0), errors.New("step count less or equal 0")
	}

	return stepCount, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepCount, walkDuration, err := parsePackage(data)

	if err != nil {
		log.Println("parse package error")
		return ""
	}

	if stepCount <= 0 {
		return ""
	}

	length := stepLength * float64(stepCount) / mInKm

	spentCalories, err := spentcalories.WalkingSpentCalories(stepCount, weight, height, walkDuration)

	if err != nil {
		log.Println("walking spent calories calculation error")
		return ""
	}

	formatResult := "Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n"

	result := fmt.Sprintf(formatResult, stepCount, length, spentCalories)

	return result
}
