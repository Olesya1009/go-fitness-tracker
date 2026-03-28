package spentcalories

import (
	"fmt"
	"strconv"
	"time"
    "strings"
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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка шагов: %v", err)
	}
    if steps <= 0 {
    return 0, "", 0, fmt.Errorf("шагов должно быть больше нуля")
    }
	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка времени: %v", err)
	}
	if duration <= 0 {
    return 0, "", 0, fmt.Errorf("продолжительность должна быть больше нуля")
    }
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLen
			return distanceM / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
    if err != nil {
        return "", err
    }

    dist := distance(steps, height)
    speed := meanSpeed(steps, height, duration)

    switch activityType {
    case "Ходьба": 
        calories, err := WalkingSpentCalories(steps, weight, height, duration)
        if err != nil {
            return "", err
        }
        return fmt.Sprintf(
            "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
            activityType, duration.Hours(), dist, speed, calories,
        ), nil

    case "Бег":
        calories, err := RunningSpentCalories(steps, weight, height, duration)
        if err != nil {
            return "", err
        }
        return fmt.Sprintf(
            "Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
            activityType, duration.Hours(), dist, speed, calories,
        ), nil

    default:
        return "", fmt.Errorf("неизвестный тип тренировки")
    }
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, fmt.Errorf("все параметры должны быть больше нуля")
    }
	speed := meanSpeed(steps, height, duration)
	durationMin := duration.Minutes()
	calories := (weight * speed * durationMin) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
        return 0, fmt.Errorf("все параметры должны быть больше нуля")
    }

    speed := meanSpeed(steps, height, duration)       
    durationMin := duration.Minutes()                  
	calories := ((weight * speed * durationMin) / minInH) * walkingCaloriesCoefficient
	return calories, nil
}