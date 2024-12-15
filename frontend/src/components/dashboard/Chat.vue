<script setup lang="ts">
import { ref } from "vue"

type ChatType = "dungeon" | "public" | "party" | "whisper" | "trade" | "guild"

interface ChatButton {
	enable: boolean
	type: ChatType
}
interface ChatMessage {
	content: string
	time: string
	type?: ChatType
	sender?: string
}

const selectedButtons = ref<ChatType[]>([])

const buttons: ChatButton[] = [
	{ enable: true, type: "dungeon" },
	{ enable: true, type: "public" },
	{ enable: true, type: "party" },
	{ enable: true, type: "whisper" },
	{ enable: true, type: "trade" },
	{ enable: true, type: "guild" },
]

const messages = ref<ChatMessage[]>([
	{
		content:
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quam velit, vulputate eu pharetra nec.",
		time: "1 minute ago",
	},
])
</script>

<template>
	<v-card class="pa-2">
		<v-card-title class="text-h5">Chat</v-card-title>

		<v-card-text>
			<v-card class="py-4">
				<v-btn-toggle
					v-model="selectedButtons"
					multiple
					class="d-flex flex-wrap gap-2"
				>
					<v-btn
						v-for="button in buttons"
						:key="button.type"
						:value="button.type"
						variant="outlined"
						class="flex-grow-1"
						:class="{
							'bg-teal': selectedButtons.includes(button.type),
						}"
					>
						{{ button.type }}
					</v-btn>
				</v-btn-toggle>
			</v-card>

			<v-card variant="outlined" class="pa-4">
				<div v-for="(message, index) in messages" :key="index">
					<div class="d-flex">
						<p class="text-start">{{ message.content }}</p>
						<p class="text-end">
							{{ message.time }}
						</p>
					</div>
				</div>
			</v-card>
		</v-card-text>
	</v-card>
</template>
