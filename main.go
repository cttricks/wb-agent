package main

import (
	"fmt"
	"log"
	"os"
	"wb-agent/controllers"

	"github.com/joho/godotenv"
	"golang.design/x/hotkey"
)

func main() {
	// Optional: Load .env (OPENAI_API_KEY=...)
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	dbName := os.Getenv("DATABASE_NAME")
	if apiKey == "" {
		log.Fatal("[ERROR] OPENAI_API_KEY not set (put it in .env or env vars)")
	}

	// Register Alt + A
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModAlt}, hotkey.KeyA)
	if err := hk.Register(); err != nil {
		log.Fatalf("[ERROR] HotKey register failed: %v", err.Error())
	}
	defer hk.Unregister()

	fmt.Println("Registered Alt+A — press it in Workbench to convert NL -> SQL")

	for {
		<-hk.Keydown()                               // wait for key down
		go controllers.OnHotkeyPress(apiKey, dbName) // handle asynchronously so registration keeps working
		<-hk.Keyup()                                 // wait for key up (optional)
	}
}
