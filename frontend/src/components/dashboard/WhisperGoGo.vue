<script setup lang="ts">
import { ref, watch, onMounted } from "vue"
import PoE2Logo from "@/assets/poe2-logo.png"
import VueQrcode from "@chenfengyuan/vue-qrcode"
import kebabCase from "kebab-case"
import { push } from "notivue"
import { useNtfyStore } from "@/stores/ntfy"
import { useAppStore } from "@/stores/app"

const ntfy = useNtfyStore()
const app = useAppStore()

// Reactive State
const gameStatus = ref(false)

// Validation
const topicErrors = ref<string[]>([])
const isValidTopic = computed(() => ntfy.topics?.trim().length > 0)

const topicRules = [
	(v: string) => !!v?.trim() || "Please enter a Topic name",
	(v: string) =>
		v?.trim().length <= 16 ||
		"The Topic name must not exceed 16 characters",
]
// Validation
const validateAndUpdateTopic = (value: string) => {
	// ตรวจสอบ rules ทั้งหมด
	const errors = topicRules
		.map((rule) => rule(value))
		.filter((result) => result !== true)

	topicErrors.value = errors

	// ถ้าไม่มี error จึงอัพเดท store
	if (errors.length === 0) {
		ntfy.updateTopics(value)
	}
}

const clearTopicErrors = () => {
	topicErrors.value = []
}

// Methods
const resetTopics = async (): Promise<void> => {
	if (!app.deviceName) {
		topicErrors.value = [
			"Cannot reset because the device name is not found",
		]
		push.error("Cannot reset because the device name is not found.")
		return
	}
	await ntfy.updateTopics(app.deviceName)
	clearTopicErrors()
}

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

// Lifecycle Hooks
onMounted(async () => {
	ntfy.fetchTopics()
	app.fetchAppData()
})

// Error Handling
const handleQRError = (error: Error): void => {
	console.error("QR Code generation error:", error)
	push.error("Failed to generate QR Code.")
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
						:class="gameStatus ? 'text-success' : 'text-error'"
					>
						{{ gameStatus ? "Running" : "Not Running" }}
					</div>
				</div>
			</div>

			<v-text-field
				v-model="app.deviceName"
				label="Device Name"
				readonly
				variant="outlined"
				density="comfortable"
			/>

			<v-text-field
				v-model="ntfy.topics"
				label="Ntfy.sh topics"
				variant="outlined"
				density="comfortable"
				:prefix="`${ntfy.NTFY_PREFIX_TOPICS}-`"
				:rules="topicRules"
				:error-messages="topicErrors"
				@input="clearTopicErrors"
				@update:model-value="validateAndUpdateTopic"
				class="my-2"
			>
				<template #append>
					<v-btn color="primary" @click="resetTopics">reset</v-btn>
				</template>
			</v-text-field>

			<v-text-field
				v-model="ntfy.ntfyLink"
				label="Ntfy.sh link"
				variant="outlined"
				density="comfortable"
				class="my-2"
				readonly
			/>

			<div class="d-flex justify-center align-center my-4">
				<vue-qrcode
					v-if="ntfy.ntfyLink && isValidTopic"
					:value="ntfy.ntfyLink"
					:options="{
						width: 150,
					}"
					@error="handleQRError"
				/>
				<div v-else class="text-caption text-grey">
					Please enter a valid topic to generate QR Code
				</div>
			</div>
			<div>
				Scan QR code with your mobile phone to subscribe for
				notifications
			</div>
		</v-card-text>
	</v-card>
</template>
