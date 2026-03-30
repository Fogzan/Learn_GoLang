package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Читаем входные данные
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")

	n := len(parts)
	nums := make([]int64, n)
	for i := 0; i < n; i++ {
		nums[i], _ = strconv.ParseInt(parts[i], 10, 64)
	}

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Пытаемся найти решение методом случайного поиска
	for attempt := 0; attempt < 1000000; attempt++ {
		// Создаем два случайных подмножества
		mask1 := rand.Int63n(1<<40 - 1)
		mask2 := rand.Int63n(1<<40 - 1)

		// Проверяем, что они не пересекаются и не пустые
		if mask1&mask2 != 0 {
			continue
		}
		if mask1 == 0 || mask2 == 0 {
			continue
		}

		// Считаем суммы
		sum1 := int64(0)
		sum2 := int64(0)

		for i := 0; i < n; i++ {
			if mask1&(1<<uint(i)) != 0 {
				sum1 += nums[i]
			}
			if mask2&(1<<uint(i)) != 0 {
				sum2 += nums[i]
			}
		}

		// Если суммы равны, нашли решение
		if sum1 == sum2 {
			// Собираем индексы
			indices1 := []int{}
			indices2 := []int{}

			for i := 0; i < n; i++ {
				if mask1&(1<<uint(i)) != 0 {
					indices1 = append(indices1, i+1)
				}
				if mask2&(1<<uint(i)) != 0 {
					indices2 = append(indices2, i+1)
				}
			}

			// Выводим результат
			fmt.Println(len(indices1))
			for _, idx := range indices1 {
				fmt.Printf("%d ", idx)
			}
			fmt.Println()
			fmt.Println(len(indices2))
			for _, idx := range indices2 {
				fmt.Printf("%d ", idx)
			}
			fmt.Println()
			return
		}
	}

	// Если не нашли случайным поиском, используем meet-in-the-middle для 20 чисел
	// Разбиваем на две группы по 20
	half := 20

	// Храним сумму и маску для первой группы
	sums := make(map[int64]int)
	for mask := 1; mask < (1 << half); mask++ {
		sum := int64(0)
		for i := 0; i < half; i++ {
			if mask&(1<<i) != 0 {
				sum += nums[i]
			}
		}
		sums[sum] = mask
	}

	// Ищем во второй группе
	secondHalf := n - half
	for mask := 1; mask < (1 << secondHalf); mask++ {
		sum := int64(0)
		for i := 0; i < secondHalf; i++ {
			if mask&(1<<i) != 0 {
				sum += nums[half+i]
			}
		}

		if firstMask, ok := sums[sum]; ok {
			// Нашли решение
			indices1 := []int{}
			indices2 := []int{}

			for i := 0; i < half; i++ {
				if firstMask&(1<<i) != 0 {
					indices1 = append(indices1, i+1)
				}
			}

			for i := 0; i < secondHalf; i++ {
				if mask&(1<<i) != 0 {
					indices2 = append(indices2, half+i+1)
				}
			}

			fmt.Println(len(indices1))
			for _, idx := range indices1 {
				fmt.Printf("%d ", idx)
			}
			fmt.Println()
			fmt.Println(len(indices2))
			for _, idx := range indices2 {
				fmt.Printf("%d ", idx)
			}
			fmt.Println()
			return
		}
	}

	fmt.Println("0")
}
