package main

import (
	"log"

	"suai-queue/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения:", err)
	}
	defer a.Close()

	a.Run()
}
