package utils

import (
	"os"
	"path/filepath"

	"github.com/go-toast/toast"
)

func Notify(title string, message string) {
	notification := toast.Notification{
		AppID:   "WB-AGENT @v1.0.2",
		Title:   title,
		Message: message,
		Icon:    getNotificationIcon(),
	}
	notification.Push()
}

func getNotificationIcon() string {

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	iconPath := filepath.Join(dir, "agent.ico")
	if _, err := os.Stat(iconPath); err != nil {
		return ""
	}

	return iconPath
}
