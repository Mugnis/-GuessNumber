package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type GameResult struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	IsWin    bool      `json:"iswin"`
	Attempts int       `json:"attempts"`
}

func main() {
	for {
		var again string
		attempts, randrange := selectDifficulty()
		fmt.Printf("Игра 'Угадай число' - от 1 до %d началась!\nУгадайте число за %d попыток!\n", randrange-1, attempts)
		win, numAttempts := game(attempts, randrange)
		saveResult(win, numAttempts)
		fmt.Print("Сыграть ещё раз? (y-да) ")
		fmt.Scanln(&again)
		if again != "y" {
			return
		}
	}
}

func game(attempts int, randrange int) (win bool, x int) {
	random := rand.IntN(randrange)
	fmt.Println(random)
	nums := make([]int, 0, attempts)
	for x < attempts {
		x++
		num := getUserInput()
		nums = append(nums, num)
		if random == num {
			win = true
			break
		}
		printHint(random, num)
		fmt.Printf("Ранее введённые числа: ")
		for _, x := range nums {
			fmt.Printf("%d ", x)
		}
	}
	if win {
		color.Green("\nВы выиграли!!!!, количество потраченных попыток: %d\n", x)
	} else {
		color.Red("\nВы проиграли😢\nСекретное число было: %d\n", random)
	}
	return win, x
}

func printHint(random, num int) {
	dif := random - num
	if dif <= 5 && dif >= 0 {
		color.Yellow("🔥Горячо, Секретное число больше👆")
	} else if dif >= -5 && dif <= 0 {
		color.Yellow("🔥Горячо, Секретное число меньше👇")
	} else if dif <= 15 && dif >= 0 {
		color.Yellow("🙂Тепло, Секретное число больше👆")
	} else if dif >= -15 && dif <= 0 {
		color.Yellow("🙂Тепло, Секретное число меньше👇")
	} else if random > num {
		color.Yellow("❄️Холодно, Секретное число больше👆")
	} else {
		color.Yellow("❄️Холодно, Секретное число меньше👇")
	}
}

func getUserInput() (num int) {
	var input string
	var err error
	for {
		fmt.Print("\nВведите число ")
		fmt.Scanln(&input)
		num, err = strconv.Atoi(input)
		if err == nil {
			break
		} else {
			fmt.Println("\nВведите число!!")
		}
	}
	return num
}

func selectDifficulty() (attempts, randrange int) {
	var difficulty string
	for {
		fmt.Print("Выберите сложность! (1-Easy, 2-Medium, 3-Hard)\n")
		fmt.Scanln(&difficulty)
		switch difficulty {
		case "1":
			attempts = 15
			randrange = 51
		case "2":
			attempts = 10
			randrange = 101
		case "3":
			attempts = 5
			randrange = 151
		default:
			fmt.Println("Введена неверная сложность")
		}
		if attempts == 5 || attempts == 10 || attempts == 15 {
			break
		}
	}
	return attempts, randrange
}

func saveResult(win bool, numAttempts int) {
	var id int
	var history []GameResult
	file, err := os.ReadFile("stats.json")
	if err == nil && len(file) > 0 {
		err = json.Unmarshal(file, &history)
		if err != nil {
			fmt.Println("Ошибка чтения старой истории:", err)
		}
		id = history[len(history)-1].ID + 1
	} else {
		id = 1
	}

	result := GameResult{
		ID:       id,
		Date:     time.Now(),
		IsWin:    win,
		Attempts: numAttempts,
	}
	history = append(history, result)
	jsonData, err := json.MarshalIndent(history, "", "    ")
	if err != nil {
		fmt.Println("Ошибка кодирования в JSON:", err)
		return
	}

	err = os.WriteFile("stats.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}

}
