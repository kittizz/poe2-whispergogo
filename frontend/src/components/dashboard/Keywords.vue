<script setup lang="ts">
import { ref } from "vue"

const keywords = ref([
	{ text: "xxxxxx", checked: true },
	{ text: "xxxxxx", checked: true },
	{ text: "xxxxxx", checked: true },
])
</script>

<template>
	<v-card>
		<v-card-title class="text-h5 pa-4">Keywords</v-card-title>

		<v-card-text>
			<v-checkbox
				v-for="(keyword, index) in keywords"
				:key="index"
				v-model="keyword.checked"
				:label="keyword.text"
				hide-details
				class="mb-2"
			></v-checkbox>
		</v-card-text>
	</v-card>
</template>
