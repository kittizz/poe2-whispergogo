// Utilities
import { defineStore } from "pinia"
import { GetChatFilters, SetChatFilters } from "../../wailsjs/go/main/App"
import type { ChatType } from "@/types/chat"

export const useChatStore = defineStore("chat", {
	state: () => ({
		filters: [] as ChatType[],
	}),

	actions: {
		async fetchChatFilters() {
			this.filters = await GetChatFilters()
		},

		async updateChatFilters(filters: ChatType[]) {
			await SetChatFilters(filters)
			this.filters = filters
		},
	},
})
