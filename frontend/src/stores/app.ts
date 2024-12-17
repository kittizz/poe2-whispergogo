// Utilities
import { defineStore } from "pinia"
import {
	GetAlertStatus,
	SetAlertStatus,
	GetGameStatus,
	GetAlertType,
	SetAlertType,
	GetTelegramChatID,
	SetTelegramChatID,
} from "../../wailsjs/go/main/App"
import { AlertType } from "@/types/app"

export const useAppStore = defineStore("app", {
	state: () => ({
		alertStatus: false,
		gameStatus: false,
		alertType: null as AlertType | null,
		telegram_chatid: "",
	}),

	actions: {
		async fetchAlertStatus() {
			this.alertStatus = await GetAlertStatus()
		},
		async fetchAlertType() {
			this.alertType = await GetAlertType()
		},
		async fetchGameStatus() {
			this.gameStatus = await GetGameStatus()
		},
		async startGameStatusWatcher() {
			this.fetchGameStatus()
			setInterval(() => {
				this.fetchGameStatus()
			}, 500)
		},
		async fetchTelegramChatID() {
			this.telegram_chatid = await GetTelegramChatID()
		},

		async updateAlertStatus(status: boolean) {
			await SetAlertStatus(status)
			this.alertStatus = status
		},

		async updateAlertType(type: AlertType) {
			await SetAlertType(type)
			this.alertType = type
		},

		async updateTelegramChatID(chatid: string): Promise<[boolean, string]> {
			const result = await SetTelegramChatID(chatid)
			if (result.valid) this.telegram_chatid = result.chat_id
			return [result.valid, result.chat_id]
		},

		async fetchAppData() {
			await Promise.all([
				this.fetchAlertStatus(),
				this.fetchGameStatus(),
				this.fetchAlertType(),
				this.fetchTelegramChatID(),
			])
		},
	},
})
