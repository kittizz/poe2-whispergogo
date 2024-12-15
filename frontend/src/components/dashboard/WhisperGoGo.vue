<script setup lang="ts">
import { ref, watch } from "vue"
import PoE2Logo from "@/assets/poe2-logo.png"
import VueQrcode from "@chenfengyuan/vue-qrcode"
import kebabCase from "kebab-case"

const prefixTopics = "whispergogo"
const deviceName = ref("kittizzPC")
const ntfyTopics = ref("xxxxxxxxxxxx")
const gameStatus = ref(false)

const generateNtfyLink = (): string => {
	return `https://ntfy.sh/${prefixTopics}-${kebabCase(ntfyTopics.value)}`
}

const ntfyLink = ref(generateNtfyLink())

const resetTopics = () => {
	ntfyTopics.value = deviceName.value
}
watch(ntfyTopics, () => {
	ntfyLink.value = generateNtfyLink()
})
</script>

<template>
	<v-card class="pa-4">
		<v-card-title class="text-h5"> WhisperGoGo</v-card-title>

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
			:prefix="prefixTopics + '-'"
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
			readonly
		/>

		<div class="d-flex justify-center align-center my-4">
			<vue-qrcode
				:value="ntfyLink"
				:options="{
					width: 150,
				}"
			/>
		</div>
		<v-card-text class="text-center">
			scan qr code for subscribe notification on phone
		</v-card-text>
	</v-card>
</template>
