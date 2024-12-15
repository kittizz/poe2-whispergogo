// Utilities
import { defineStore } from "pinia"
import {
	GetDeviceName,
	GetAlertStatus,
	SetAlertStatus,
} from "../../wailsjs/go/main/App"

export const useAppStore = defineStore("app", {
	state: () => ({
		deviceName: "",
		alertStatus: false,
	}),

	actions: {
		async fetchDeviceName() {
			this.deviceName = await GetDeviceName()
		},

		async fetchAlertStatus() {
			this.alertStatus = await GetAlertStatus()
		},

		async updateAlertStatus(status: boolean) {
			await SetAlertStatus(status)
			this.alertStatus = status
		},

		async fetchAppData() {
			await Promise.all([this.fetchDeviceName(), this.fetchAlertStatus()])
		},
	},
})
