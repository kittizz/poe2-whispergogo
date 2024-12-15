<script setup lang="ts">
import { ref } from "vue"

interface ChatMessage {
	label: string
	content: string
	time: string
}

const messages = ref<ChatMessage[]>([
	{
		label: "Label",
		content:
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quam velit, vulputate eu pharetra nec.",
		time: "1 minute ago",
	},
	{
		label: "Label",
		content:
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla quam velit, vulputate eu pharetra nec.",
		time: "1 minute ago",
	},
])
</script>

<template>
	<v-card>
		<v-card-title class="text-h5 pa-4">Chat</v-card-title>

		<v-card-text>
			<v-card variant="outlined" class="pa-4">
				<div v-for="(message, index) in messages" :key="index">
					<div class="d-flex">
						<p class="text-start">{{ message.content }}</p>
						<p class="text-end">
							Right aligned text on LTR and left aligned on RTL.
						</p>
					</div>
				</div>
			</v-card>
		</v-card-text>
	</v-card>
</template>
