<script setup lang="ts">
import { ref, watch, onMounted, handleError } from "vue"
import PoE2Logo from "@/assets/poe2-logo.png"
import TGQrcode from "@/assets/tgqrcode.png"
import { push } from "notivue"
import { useAppStore } from "@/stores/app"
import { OpenTelegramLink } from "../../wailsjs/go/main/App"

const app = useAppStore()
const telegram_chatid = ref("")
const openSettings = () => {
	app.alertStatus = !app.alertStatus
	if (app.alertStatus) {
		push.success("Alert is turned on.")
	} else {
		push.error("Alert is turned off.")
	}
}

watch(
	() => app.alertStatus,
	async (newStatus) => {
		await app.updateAlertStatus(newStatus)
	}
)
watch(
	// app.alertType
	() => app.alertType,
	async (newType) => {
		await app.updateAlertType(newType!)
	}
)

// Lifecycle Hooks
onMounted(async () => {
	await app.fetchAppData()
	telegram_chatid.value = app.telegram_chatid
})

// Error Handling
const handleTelegramChatID = async () => {
	if (!telegram_chatid.value) {
		push.error("Telegram Chat ID is required.")
	} else {
		const [valid, chat_id] = await app.updateTelegramChatID(
			telegram_chatid.value
		)
		if (valid) {
			push.success("Telegram Chat ID updated successfully.")
			telegram_chatid.value = chat_id
		} else {
			push.error("Telegram Chat ID is invalid.")
			telegram_chatid.value = app.telegram_chatid
		}
	}
}
const handleTelegramLink = () => {
	OpenTelegramLink()
}
</script>

<template>
	<v-card class="pa-2">
		<v-card-title class="text-h5 d-flex justify-space-between align-center">
			<span>WhisperGoGo</span>

			<v-btn
				:color="app.alertStatus ? 'success' : 'error'"
				@click="openSettings"
			>
				<v-icon :icon="app.alertStatus ? 'mdi-eye' : 'mdi-eye-off'" />
				<div>{{ app.alertStatus ? "Alert ON" : "Alert OFF" }}</div>
			</v-btn>
		</v-card-title>

		<v-card-text>
			<div class="d-flex align-center justify-center pa-2">
				<v-img :src="PoE2Logo" max-width="100" class="mr-4" />
				<div>
					<div class="text-h6">Path of Exile 2</div>
					<div
						class="text-subtitle-2"
						:class="app.gameStatus ? 'text-success' : 'text-error'"
					>
						{{ app.gameStatus ? "Running" : "Not Running" }}
					</div>
				</div>
			</div>
			<v-radio-group inline label="Alert Type" v-model="app.alertType">
				<v-radio label="both" value="both"></v-radio>
				<v-radio label="chat filter" value="chat_filter"></v-radio>
				<v-radio label="keyword" value="keyword"></v-radio>
				<p class="text-caption font-italic">
					Select the type of alerts you want to receive: chat filter
					(specific triggers from game chat) or keyword (specific
					words you set).
				</p>
			</v-radio-group>

			<v-text-field
				v-model="telegram_chatid"
				label="Telegram Chat ID"
				variant="outlined"
				density="comfortable"
				class="my-2"
			>
				<template v-slot:prepend>
					<v-icon :color="app.telegram_chatid ? 'green' : 'red'">
						{{
							app.telegram_chatid
								? "mdi-link-variant"
								: "mdi-link-variant-off"
						}}
					</v-icon>
				</template>

				<template #append>
					<v-btn color="primary" @click="handleTelegramChatID"
						>Save</v-btn
					>
				</template>
			</v-text-field>

			<p class="text-caption font-italic">
				Enter your Telegram Chat ID to link the notifications to this
				Telegram chat. Use the instructions below to retrieve your Chat
				ID.
			</p>

			<div class="justify-center align-center text-center my-4">
				<v-img :src="TGQrcode" height="150" />
				<br />
				<a
					class="text-decoration-none"
					href="#"
					@click="handleTelegramLink"
				>
					https://t.me/whispergogo_bot
				</a>
			</div>

			<p class="text-caption font-italic">
				Scan the QR code or search for @WHISPERGOGO_BOT on Telegram.
				After adding the bot, use
				<b class="text-blue">/start</b> to receive your Chat ID.
			</p>
		</v-card-text>
	</v-card>
</template>
