package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitedData := strings.Split(data, ",")

	if len(splitedData) != 3 {
		return 0, "", time.Duration(0), errors.New("wrong len of splited data")
	}

	stepCount, err := strconv.Atoi(splitedData[0])

	if err != nil {
		joinedError := errors.Join(errors.New("step count parsing error"), err)
		return 0, "", time.Duration(0), joinedError
	}

	if stepCount <= 0 {
		return 0, "", time.Duration(0), errors.New("wrong step count")
	}

	activityDuration, err := time.ParseDuration(splitedData[2])

	if err != nil {
		joinedError := errors.Join(errors.New("activity duration parsing error"), err)
		return 0, "", time.Duration(0), joinedError
	}

	if activityDuration <= time.Duration(0) {
		return 0, "", time.Duration(0), errors.New("wrong step count")
	}

	return stepCount, splitedData[1], activityDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := stepLengthCoefficient * height

	length := stepLength * float64(steps) / mInKm

	return length
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= time.Duration(0) {
		return 0
	}

	dist := distance(steps, height)

	result := dist / duration.Hours()

	return result
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepCount, trainingType, trainingDuration, err := parseTraining(data)

	if err != nil {
		return "", err
	}

	var (
		distance      = distance(stepCount, height)
		meanSpeed     = meanSpeed(stepCount, height, trainingDuration)
		spentCalories float64
	)

	switch trainingType {
	case "Бег":
		spentCalories, err = RunningSpentCalories(stepCount, weight, height, trainingDuration)

	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(stepCount, weight, height, trainingDuration)

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		joinedError := errors.Join(errors.New("spent calories calculation error"), err)
		return "", joinedError
	}

	resultFormat := "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n"

	resultParams := []any{trainingType, trainingDuration.Hours(), distance, meanSpeed, spentCalories}

	result := fmt.Sprintf(resultFormat, resultParams...)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= time.Duration(0) {
		return 0, errors.New("wrong argument value")
	}

	speed := meanSpeed(steps, height, duration)

	spentCalories := (weight * speed * duration.Minutes()) / minInH

	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= time.Duration(0) {
		return 0, errors.New("wrong argument value")
	}

	speed := meanSpeed(steps, height, duration)

	spentCalories := (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient

	return spentCalories, nil
}
