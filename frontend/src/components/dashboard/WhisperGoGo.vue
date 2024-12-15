<script setup lang="ts">
import { ref, watch, onMounted } from "vue"
import PoE2Logo from "@/assets/poe2-logo.png"
import VueQrcode from "@chenfengyuan/vue-qrcode"
import kebabCase from "kebab-case"
import { push } from "notivue"
import { GetDeviceName } from "../../../wailsjs/go/main/App"
import { main } from "../../../wailsjs/go/models"

// Reactive State
const deviceName = ref("")
const ntfyTopics = ref("")
const gameStatus = ref(false)
const aleartStatus = ref(false)
const ntfyLink = ref("")

// Validation
const topicErrors = ref<string[]>([])
const isValidTopic = computed(() => ntfyTopics.value?.trim().length > 0)

const topicRules = [
	(v: string) => !!v?.trim() || "Please enter a Topic name",
	(v: string) =>
		v?.trim().length <= 16 ||
		"The Topic name must not exceed 16 characters",
]

// Computed
const generateNtfyLink = (): string => {
	const formattedTopic = `${main.Ntfy.NTFY_PREFIX_TOPICS}-${kebabCase(
		ntfyTopics.value
	)}`
	return `${main.Ntfy.NTFY_BASE_URL}/${formattedTopic}`
}
const clearTopicErrors = () => {
	topicErrors.value = []
}

// Methods
const resetTopics = (): void => {
	if (!deviceName.value) {
		topicErrors.value = [
			"Cannot reset because the device name is not found",
		]
		push.error("Cannot reset because the device name is not found.")
		return
	}
	ntfyTopics.value = deviceName.value
	clearTopicErrors()
}
const openSettings = () => {
	aleartStatus.value = !aleartStatus.value
	if (aleartStatus.value) {
		push.success("Alert is turned on.")
	} else {
		push.error("Alert is turned off.")
	}
}

// Watchers
watch(ntfyTopics, () => {
	ntfyLink.value = generateNtfyLink()
})

// Lifecycle Hooks
onMounted(async () => {
	try {
		deviceName.value = await GetDeviceName()
		ntfyTopics.value = deviceName.value
		ntfyLink.value = generateNtfyLink()
	} catch (error) {
		console.error("Failed to get device name:", error)
	}
})

// Error Handling
const handleQRError = (error: Error): void => {
	console.error("QR Code generation error:", error)
	// อาจจะเพิ่ม notification หรือ error state ตามต้องการ
	push.error("Failed to generate QR Code.")
}
</script>

<template>
	<v-card class="pa-4">
		<v-card-title class="text-h5 d-flex justify-space-between align-center">
			<span>WhisperGoGo</span>

			<v-btn
				:color="aleartStatus ? 'success' : 'error'"
				@click="openSettings"
			>
				<v-icon :icon="aleartStatus ? 'mdi-eye' : 'mdi-eye-off'" />
				<div>{{ aleartStatus ? "Alert ON" : "Alert OFF" }}</div>
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
				v-model="deviceName"
				label="Device Name"
				readonly
				variant="outlined"
				density="comfortable"
			/>

			<v-text-field
				v-model="ntfyTopics"
				label="Ntfy.sh topics"
				variant="outlined"
				density="comfortable"
				:prefix="main.Ntfy.NTFY_PREFIX_TOPICS + '-'"
				:rules="topicRules"
				:error-messages="topicErrors"
				@input="clearTopicErrors"
				class="my-2"
			>
				<template #append>
					<v-btn color="primary" @click="resetTopics"> reset </v-btn>
				</template>
			</v-text-field>

			<v-text-field
				v-model="ntfyLink"
				label="Ntfy.sh link"
				variant="outlined"
				density="comfortable"
				class="my-2"
				readonly
			/>

			<div class="d-flex justify-center align-center my-4">
				<vue-qrcode
					v-if="ntfyLink && isValidTopic"
					:value="ntfyLink"
					:options="{
						width: 150,
					}"
					@error="handleQRError"
				/>
				<div v-else class="text-caption text-grey">
					Please enter a valid topic to generate QR Code
				</div>
			</div>
			<div>Scan QR Code on moblie phone to subscribe notification.</div>
		</v-card-text>
	</v-card>
</template>
