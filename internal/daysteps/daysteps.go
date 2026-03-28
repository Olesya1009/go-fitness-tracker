package daysteps

import (
	"fmt"
	"time"
    "strings"
	"strconv"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage разбирает строку формата "шаги,длительность"
func parsePackage(data string) (int, time.Duration, error) {
	
	parts := strings.Split(data, ",")
	
	// Проверяем что получили ровно 2 части
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])  
    if err != nil {
        return 0, 0, fmt.Errorf("не могу прочитать шаги: %v", err)
    }

    if steps <= 0 {                        
        return 0, 0, fmt.Errorf("шагов должно быть больше нуля")
    }

    duration, err := time.ParseDuration(parts[1])  
    if err != nil {
        return 0, 0, fmt.Errorf("ошибка времени: %v", err)
    }

    return steps, duration, nil            
}

// DayActionInfo возвращает итоговую строку с результатами за день
func DayActionInfo(data string, weight, height float64) string {
	 // Вызываем parsePackage
    steps, duration, err := parsePackage(data)
    if err != nil {
        fmt.Println(err)  
        return ""        
    }

    // Проверяем что шаги > 0
    if steps <= 0 {
        return ""
    }

    // Считаем дистанцию
    distanceMeters := float64(steps) * stepLength
    distanceKm := distanceMeters / mInKm

    // Калории считает готовая функция из другого пакета
    calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
    if err != nil {
        fmt.Println(err)
        return ""
    }

    return fmt.Sprintf(
        "Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
        steps, distanceKm, calories,
    )
}
