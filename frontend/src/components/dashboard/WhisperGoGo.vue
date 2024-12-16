<script setup lang="ts">
import { ref, watch, onMounted } from "vue"
import PoE2Logo from "@/assets/poe2-logo.png"
import TGQrcode from "@/assets/tgqrcode.png"
import { push } from "notivue"
import { useAppStore } from "@/stores/app"

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

		<v-card-text class="text-center">
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

			<div class="d-flex justify-center align-center my-4">
				<v-img :src="TGQrcode" height="200" />
			</div>
			<div>
				Scan QR code with your mobile phone to subscribe for
				notifications
			</div>
		</v-card-text>
	</v-card>
</template>
