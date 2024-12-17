// Utilities
import { defineStore } from "pinia"
import { main } from "../../wailsjs/go/models"
import {
	GetKeywords,
	SetKeywords,
	ResetKeywords,
} from "../../wailsjs/go/main/App"
import type { IKeyword } from "@/types/keyword"

export const useKeywordStore = defineStore("keyword", {
	state: () => ({
		keywords: [] as IKeyword[],
		ntfyTopics: [] as string[],
	}),

	actions: {
		async fetchKeywords() {
			const keywords = await GetKeywords()
			this.keywords = keywords.map(
				(keyword) =>
					({
						text: keyword.Keyword,
						enable: keyword.Enable,
					} as IKeyword)
			)
		},

		async updateKeywords(keywords: IKeyword[]) {
			const mainKeywords = keywords.map((ikeyword) => {
				return {
					Keyword: ikeyword.text,
					Enable: ikeyword.enable,
				} as main.Keyword
			})
			await SetKeywords(mainKeywords)

			await this.fetchKeywords()
		},

		async resetKeywords() {
			await ResetKeywords()
			await this.fetchKeywords()
		},
	},
})
