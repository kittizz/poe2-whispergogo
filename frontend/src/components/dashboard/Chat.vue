<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { EventsOn, EventsOff } from "../../../wailsjs/runtime"
import dayjs from "dayjs"

import { ChatType, type ChatMessage } from "@/types/chat"
import { useChatStore } from "@/stores/chat"

interface ChatButton {
	enable: boolean
	type: ChatType
}

const buttons: ChatType[] = [
	ChatType.Local,
	ChatType.Global,
	ChatType.Party,
	ChatType.Whisper,
	ChatType.Trade,
	ChatType.Guild,
]

const chat = useChatStore()

const messages = ref<ChatMessage[]>([])
const chatContainer = ref<HTMLDivElement | null>(null)

const MAX_MESSAGES = 100

const handleChatMessage = (message: ChatMessage) => {
	messages.value.push(message)

	// ถ้าจำนวนข้อความเกิน MAX_MESSAGES ให้ลบข้อความเก่าสุดออก
	if (messages.value.length > MAX_MESSAGES) {
		messages.value = messages.value.slice(-MAX_MESSAGES)
	}
}

const scrollToBottom = () => {
	nextTick(() => {
		if (chatContainer.value) {
			chatContainer.value.scrollTop = chatContainer.value.scrollHeight
		}
	})
}

watch(
	() => messages.value,
	() => {
		scrollToBottom()
	},
	{ deep: true }
)
watch(
	() => chat.filters,
	() => chat.updateChatFilters(chat.filters)
)

onMounted(async () => {
	await chat.fetchChatFilters()
	EventsOn("chatMessage", handleChatMessage)
})

onUnmounted(() => {
	EventsOff("chatMessage")
})
</script>

<template>
	<v-card class="pa-2">
		<v-card-title class="text-h5">Chat</v-card-title>

		<v-card-text>
			<v-btn-toggle
				v-model="chat.filters"
				multiple
				class="d-flex flex-wrap gap-2"
			>
				<v-btn
					v-for="buttontype in buttons"
					:key="buttontype"
					:value="buttontype"
					variant="outlined"
					class="flex-grow-1"
					:class="{
						'bg-teal': chat.filters.includes(buttontype),
					}"
				>
					{{ buttontype }}
				</v-btn>
			</v-btn-toggle>

			<v-card variant="tonal" class="my-2">
				<div
					class="chat-messages px-2 bg-grey-darken-4"
					ref="chatContainer"
				>
					<div v-for="(message, index) in messages" :key="index">
						<div class="d-flex justify-space-between">
							<div>
								<span class="font-weight-bold me-2">
									{{
										message.Username
											? message.Username
											: message.MessageType
									}}:
								</span>
								<span>{{ message.Content }}</span>
							</div>
							<span class="text-caption">
								{{
									dayjs(message.Timestamp).format("HH:mm:ss")
								}}
							</span>
						</div>
					</div>
				</div>
			</v-card>
		</v-card-text>
	</v-card>
</template>
<style scoped>
.chat-messages {
	height: 250px;
	overflow-y: auto;
}
</style>
