package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	config *Config

	gameStatus *bool
	lineStream <-chan string

	whispergogoClient *WhisperGOGO
}

// NewApp creates a new App application struct
func NewApp(config *Config, gameStatus *bool, lineStream <-chan string) *App {
	return &App{
		config:            config,
		gameStatus:        gameStatus,
		lineStream:        lineStream,
		whispergogoClient: NewWhisperGOGO(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.subscribeLineStream()
}
func (a *App) shutdown(_ context.Context) {
	fmt.Println("Shutting down...")
}

func (a *App) onSecondInstanceLaunch(options options.SecondInstanceData) {
	a.BringToFront()
}

func (a *App) subscribeLineStream() {

	for logLine := range a.lineStream {
		msg, err := ParseChatMessage(logLine)
		if err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue
		}

		// filter msg.MessageType not in a.config.ChatFilters
		if !lo.Contains(a.config.ChatFilters, msg.MessageType) &&
			(a.config.AlertType == AlertTypeChatFilter || a.config.AlertType == AlertTypeBoth) {
			continue
		}

		fmt.Printf("%v> Type: %s, User: %s, Content: %s\n",
			msg.Timestamp.Format(time.DateTime),
			msg.MessageType, msg.Username, msg.Content)

		runtime.EventsEmit(a.ctx, "chatMessage", *msg)

		// ตรวจสอบ keyword ถ้า AlertType เป็น keyword หรือ both
		if a.config.AlertType == AlertTypeKeyword || a.config.AlertType == AlertTypeBoth {
			// กรองเฉพาะ keyword ที่เปิดใช้งาน
			enabledKeywords := lo.Filter(a.config.Keywords, func(k Keyword, _ int) bool {
				return k.Enable
			})

			// ตรวจสอบว่ามี keyword ใดอยู่ในข้อความหรือไม่
			hasKeyword := lo.SomeBy(enabledKeywords, func(k Keyword) bool {
				return strings.Contains(
					strings.ToLower(msg.Content),
					strings.ToLower(k.Keyword),
				)
			})

			// ถ้าไม่มี keyword ที่ตรงกัน ให้ข้ามไป
			if !hasKeyword {
				continue
			}
		}

		if !a.config.AlertStatus || a.config.TelegramChatID == "" {
			continue
		}

		messageText := fmt.Sprintf("🎮 Path of Exile 2\n"+
			"📢 %s\n"+ // ประเภทแชท (Local, Global, Trade, etc.)
			"👤 %s\n"+ // ชื่อผู้เล่น
			"💬 %s\n"+ // เนื้อหาข้อความ
			"⏰ %s", // เวลา
			msg.MessageType,
			msg.Username,
			msg.Content,
			msg.Timestamp.Format("15:04:05"),
		)

		err = a.whispergogoClient.SendMessage(a.config.TelegramChatID, messageText)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}

	}

}
func (a *App) BringToFront() {
	runtime.WindowShow(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// GetKeywords returns list of keywords from config
func (a *App) GetKeywords() []Keyword {
	return a.config.Keywords
}

// SetKeywords sets list of keywords in config
func (a *App) SetKeywords(keywords []Keyword) {
	a.config.Keywords = keywords
	a.config.Save()
}

func (a *App) ResetKeywords() {
	a.config.defaultKeywords()
	a.config.Save()
}

func (a *App) GetAlertStatus() bool {
	return a.config.AlertStatus
}
func (a *App) SetAlertStatus(status bool) {
	a.config.AlertStatus = status
	a.config.Save()
}
func (a *App) GetGameStatus() bool {
	return *a.gameStatus
}

func (a *App) GetChatFilters() []ChatType {
	return a.config.ChatFilters
}

func (a *App) SetChatFilters(filters []ChatType) {
	a.config.ChatFilters = filters
	a.config.Save()
}
func (a *App) GetAlertType() AlertType {
	return a.config.AlertType
}
func (a *App) SetAlertType(alertType AlertType) {
	a.config.AlertType = alertType
	a.config.Save()
}

func (a *App) GetTelegramChatID() string {
	return a.config.TelegramChatID
}

func (a *App) SetTelegramChatID(chatid string) *VerifyChatResponse {
	if chatid == "" {
		return nil
	}
	verifyChatRes, err := a.whispergogoClient.VerifyChat(chatid)
	if err != nil {
		return nil
	}
	a.config.TelegramChatID = verifyChatRes.ChatID
	a.config.Save()

	return verifyChatRes
}

func (a *App) OpenTelegramLink() {
	runtime.BrowserOpenURL(a.ctx, "https://t.me/whispergogo_bot")
}
