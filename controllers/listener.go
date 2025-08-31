package controllers

import (
	"fmt"
	"strings"
	"time"
	"wb-agent/utils"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func OnHotkeyPress(apiKey string, dbName string) {
	// helper to press a combo: Ctrl+C or Ctrl+V or Home/Shift+End etc.
	press := func(keys []int, ctrl, shift, alt bool) {
		kb, _ := keybd_event.NewKeyBonding()
		kb.SetKeys(keys...)
		kb.HasCTRL(ctrl)
		kb.HasSHIFT(shift)
		kb.HasALT(alt)
		time.Sleep(50 * time.Millisecond)
		kb.Launching()
		time.Sleep(80 * time.Millisecond)
	}

	// Clean Clipboard
	if err := clipboard.WriteAll(""); err != nil {
		fmt.Println("[ERROR] Failed to clear clipboard:", err.Error())
	}

	// Select current line under caret
	press([]int{keybd_event.VK_HOME}, false, false, false) // move caret to start
	time.Sleep(50 * time.Millisecond)

	press([]int{keybd_event.VK_A}, true, false, false) // shift+end -> select line
	time.Sleep(50 * time.Millisecond)

	// Try Ctrl+C first (user should ideally select text or caret should be on line)
	press([]int{keybd_event.VK_C}, true, false, false)
	time.Sleep(150 * time.Millisecond)
	text, _ := clipboard.ReadAll()

	// Unselect All
	press([]int{keybd_event.VK_UP}, false, false, false)

	if strings.TrimSpace(text) == "" {
		fmt.Println("[WARNING] No text found in clipboard after copy. Make sure the editor has focus and text/line is selectable.")
		utils.Notify("User Input Not Found", "No text found in clipboard after copy. Make sure the editor has focus and text/line is selectable.")
		return
	}

	fmt.Println("[IN] ", text)

	// Call OpenAI to convert NL -> SQL
	sql := ComposeSQL(apiKey, text)
	if sql == "" {
		fmt.Println("[WARNING] OpenAI returned empty SQL")
		utils.Notify("OpenAI Response Is Empty", "OpenAI returned empty SQL query, please try again.")
		return
	}

	sqlQuery := fmt.Sprintf("USE %s;\n%s", dbName, sql)
	fmt.Println("[AI] Generated SQL Query:", strings.ReplaceAll(sql, "\n", " "))

	// Put SQL-Query to clipboard
	if err := clipboard.WriteAll(sqlQuery); err != nil {
		fmt.Println("[ERROR] Failed to paste query from clipboard:", err.Error())
		utils.Notify("Unable to Insert Query", "Failed to paste query from clipboard. Check logs for more info.")
		return
	}
	time.Sleep(80 * time.Millisecond)

	// Clean the Query Page
	press([]int{keybd_event.VK_A}, true, false, false)
	press([]int{keybd_event.VK_BACKSPACE}, false, false, false)

	// Put Query
	time.Sleep(80 * time.Millisecond)
	press([]int{keybd_event.VK_V}, true, false, false)

	// Execute the query
	press([]int{keybd_event.VK_ENTER}, true, true, false)
}
