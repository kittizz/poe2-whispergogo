package main

import (
	"fmt"
	"strings"
	"time"

)

// ChatMessage เก็บข้อมูลของข้อความแชท
type ChatMessage struct {
	Timestamp   time.Time
	ClientID    string
	ServerID    string
	MessageType ChatType // เปลี่ยนเป็น enum
	Username    string
	Content     string
}

// ChatType enum สำหรับประเภทข้อความต่างๆ
type ChatType string

const (
	ChatTypeLocal   ChatType = "local"   // ข้อความในด่าน
	ChatTypeGlobal  ChatType = "global"  // # ข้อความ global
	ChatTypeWhisper ChatType = "whisper" // @ ข้อความ whisper
	ChatTypeTrade   ChatType = "trade"   // $ ข้อความ trade
	ChatTypeParty   ChatType = "party"   // % ข้อความปาร์ตี้
	ChatTypeGuild   ChatType = "guild"   // & ข้อความกิลด์
	ChatTypeSystem  ChatType = "system"  // ข้อความระบบ
)
func ParseChatMessage(logLine string) (*ChatMessage, error) {
    // แยกส่วนของ timestamp และข้อมูลอื่นๆ
    parts := strings.SplitN(logLine, " ", 6)
    if len(parts) < 6 {
        return nil, fmt.Errorf("invalid log format")
    }

    timestamp, err := time.Parse("2006/01/02 15:04:05", parts[0]+" "+parts[1])
    if err != nil {
        return nil, fmt.Errorf("invalid timestamp format: %v", err)
    }

    // แยกส่วน [INFO Client xxxxx] ออกจากข้อความ
    infoAndMessage := strings.SplitN(parts[5], "] ", 2)
    if len(infoAndMessage) != 2 {
        return nil, fmt.Errorf("invalid message format")
    }

    // แยกส่วนของ username และ content
    messageParts := strings.SplitN(infoAndMessage[1], ": ", 2)
    var username, content string
    
    if len(messageParts) == 2 {
        username = messageParts[0]
        content = messageParts[1]
    } else {
        // กรณีเป็นข้อความระบบที่ไม่มี username
        username = ""
        content = messageParts[0]
    }

    // ระบุประเภทของข้อความ
    var messageType ChatType
    switch {
    case strings.HasPrefix(username, "@"):
        messageType = ChatTypeWhisper
        username = strings.TrimPrefix(username, "@")
    case strings.HasPrefix(username, "#"):
        messageType = ChatTypeGlobal
        username = strings.TrimPrefix(username, "#")
    case strings.HasPrefix(username, "$"):
        messageType = ChatTypeTrade
        username = strings.TrimPrefix(username, "$")
    case strings.HasPrefix(username, "%"):
        messageType = ChatTypeParty
        username = strings.TrimPrefix(username, "%")
    case strings.HasPrefix(username, "&"):
        messageType = ChatTypeGuild
        username = strings.TrimPrefix(username, "&")
    case username == "":
        messageType = ChatTypeSystem
    default:
        messageType = ChatTypeLocal
    }

    return &ChatMessage{
        Timestamp:   timestamp,
        ClientID:    parts[2],
        ServerID:    parts[3],
        MessageType: messageType,
        Username:    username,
        Content:     content,
    }, nil
}

