// enum สำหรับประเภทข้อความต่างๆ
export enum ChatType {
	Local = "local", // ข้อความในด่าน
	Global = "global", // # ข้อความ global
	Whisper = "whisper", // @ ข้อความ whisper
	Trade = "trade", // $ ข้อความ trade
	Party = "party", // % ข้อความปาร์ตี้
	Guild = "guild", // & ข้อความกิลด์
	System = "system", // ข้อความระบบ
}

// interface สำหรับข้อมูลที่รับมาจาก backend
export interface ChatMessage {
	Timestamp: string // จาก backend จะส่งมาเป็น ISO string
	ClientID: string
	ServerID: string
	MessageType: ChatType
	Username: string
	Content: string
}
