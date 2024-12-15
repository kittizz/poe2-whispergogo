// Utilities
import { defineStore } from "pinia"
import {
	GetNtfyTopics,
	SetNtfyTopics,
	ResetNtfyTopics,
} from "../../wailsjs/go/main/App"
import kebabCase from "kebab-case"
import { main } from "../../wailsjs/go/models"

export const useNtfyStore = defineStore("ntfy", {
	state: () => ({
		topics: "" as string,
		NTFY_PREFIX_TOPICS: main.Ntfy.NTFY_PREFIX_TOPICS as string,
		NTFY_BASE_URL: main.Ntfy.NTFY_BASE_URL as string,
	}),
	getters: {
		ntfyLink(): string {
			const formattedTopic = `${this.NTFY_PREFIX_TOPICS}-${kebabCase(
				this.topics
			)}`
			return `${this.NTFY_BASE_URL}/${formattedTopic}`
		},
	},
	actions: {
		async fetchTopics() {
			const topics = await GetNtfyTopics()
			this.topics = topics
		},

		async updateTopics(topics: string) {
			await SetNtfyTopics(topics)
			await this.fetchTopics()
		},

		async resetTopics() {
			await ResetNtfyTopics()
			await this.fetchTopics()
		},
	},
})
