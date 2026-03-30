package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readStr(countElem int) []int {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if countElem != -1 {
		if len(parts) != countElem {
			return nil
		}
	}

	var result []int
	for _, elem := range parts {
		a, err := strconv.Atoi(elem)
		if err != nil {
			return nil
		}
		result = append(result, a)
	}
	return result
}

func checkRoute(routeBarrier, typeBarrier []int, lenRoute int) bool {
	barrier := map[int]int{
		1: 1,
		2: 2,
		3: 4,
	}

	if len(routeBarrier) != lenRoute || len(typeBarrier) != lenRoute {
		return false
	}
	for i, elem := range routeBarrier {
		if i > 0 {
			if elem <= routeBarrier[i-1] {
				return false
			}

			if elem <= routeBarrier[i-1]+barrier[typeBarrier[i-1]] {
				// fmt.Printf("\n\ni: %v\nelem: %v\nbar: %v\nres: %v\n\n", i, elem, typeBarrier[i-1], routeBarrier[i-1]+barrier[typeBarrier[i-1]]-1)
				// fmt.Println(routeBarrier[i-1]+barrier[typeBarrier[i-1]]-1, elem)
				return false
			}
		}
	}
	return true
}

func testJumps(routeBarrier, typeBarrier []int, lenRoute int, pointJump, lenJump []int, countJump int) int {
	barrier := map[int]int{
		1: 1,
		2: 2,
		3: 4,
	}
	coin := map[int]int{
		1: 1,
		2: 3,
		3: 5,
	}
	resultCoin := 0
	processed := make([]bool, lenRoute)
	countBarrier := 0
	// Проходимся по всем прыжкам
	for i := 0; i < countJump; i++ {
		startJump := pointJump[i]
		endJump := startJump + lenJump[i]

		// Проходимся по всем препятствиям
		for j := 0; j < lenRoute; j++ {
			// Если препятствие уже обработано, пропускаем
			if processed[j] {
				continue
			}

			startBarrier := routeBarrier[j]
			endBarrier := routeBarrier[j] + barrier[typeBarrier[j]]

			// Проверяем, пересекается ли прыжок с препятствием
			if startJump <= startBarrier && endJump >= startBarrier {
				// fmt.Printf("\n\n-- %v -> %v\n", startJump, endJump)
				// fmt.Printf("-- %v -> %v\n", startBarrier, endBarrier)
				// Помечаем, что препятствие обработано
				processed[j] = true
				countBarrier += 1
				// Проверяем, успешно ли преодолено препятствие
				// Условие: прыжок начался не позже препятствия И покрывает всю длину
				if startJump <= startBarrier && endJump >= endBarrier {
					resultCoin += coin[typeBarrier[j]]
					// fmt.Printf("YES + %v\n\n", coin[typeBarrier[j]])
				} else {
					resultCoin -= 1
					// fmt.Printf("NO = %v\n\n", resultCoin)
				}
			}
		}
	}

	if lenRoute > countBarrier {
		resultCoin -= lenRoute - countBarrier
	}

	if resultCoin < 0 {
		resultCoin = 0
	}
	return resultCoin
}

func checkJump(pointJump, lenJump []int, countJump int) bool {
	if countJump == 0 {
		return true
	}
	for i := 1; i < countJump; i++ {
		if pointJump[i] < pointJump[i-1] {
			return false
		}
	}
	currentPos := 0
	for i, startJump := range pointJump {
		if startJump < currentPos {
			return false
		}
		currentPos = startJump + lenJump[i]
	}

	return true
}

func main() {
	lenRoute := readStr(1)
	if lenRoute == nil {
		return
	}
	routeBarrier := readStr(lenRoute[0])
	if routeBarrier == nil {
		return
	}
	typeBarrier := readStr(lenRoute[0])
	if typeBarrier == nil {
		return
	}

	countJump := readStr(1)
	if countJump == nil {
		return
	}
	pointJump := readStr(countJump[0])
	if pointJump == nil {
		return
	}
	lenJump := readStr(countJump[0])
	if lenJump == nil {
		return
	}
	for _, elem := range lenJump {
		if elem != 1 && elem != 2 && elem != 4 {

			return
		}
	}
	testRoute := checkRoute(routeBarrier, typeBarrier, lenRoute[0])
	if !testRoute {
		fmt.Println("0")
		return
	}
	testJump := checkJump(pointJump, lenJump, countJump[0])
	if !testJump {
		fmt.Println("0")
		return
	}
	resultCoin := testJumps(routeBarrier, typeBarrier, lenRoute[0], pointJump, lenJump, countJump[0])
	fmt.Println(resultCoin)
}
