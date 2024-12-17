<script setup lang="ts">
import { ref, watch } from "vue"
import { useKeywordStore } from "@/stores/keyword"
import { push } from "notivue"
import type { IKeyword } from "@/types/keyword"

const keyword = useKeywordStore()
const keywords = ref<IKeyword[]>([])

const enableAll = ref(keywords.value.every((keyword) => keyword.enable))
const isAdding = ref(false)

onMounted(async () => {
	await keyword.fetchKeywords()
	keywords.value = keyword.keywords
})

watch(
	keywords,
	() => {
		enableAll.value = keywords.value.every((keyword) => keyword.enable)

		const updateKeywords = keywords.value.filter(
			(keyword) => !keyword.isNew && !keyword.isEditing
		)
		keyword.updateKeywords(updateKeywords)
	},
	{ deep: true }
)

const toggleEnableAll = () => {
	enableAll.value = !enableAll.value
	keywords.value.forEach((keyword) => (keyword.enable = enableAll.value))
}

const startAddNew = () => {
	isAdding.value = true
	keywords.value.unshift({
		text: "",
		enable: false,
		isEditing: true,
		isNew: true,
	})
}

const saveNew = (keyword: IKeyword) => {
	if (keyword.text.trim()) {
		keyword.isEditing = false
		keyword.isNew = false
		keyword.enable = true
	} else {
		keywords.value = keywords.value.filter((k) => !k.isNew)
	}
	isAdding.value = false
}

const cancelNew = () => {
	keywords.value = keywords.value.filter((k) => !k.isNew)
	isAdding.value = false
}

const startEditing = (keyword: IKeyword) => {
	keywords.value.forEach((k) => (k.isEditing = false))
	keyword.originalText = keyword.text
	keyword.isEditing = true
}

const saveEdit = (keyword: IKeyword) => {
	if (keyword.text.trim()) {
		keyword.isEditing = false
		delete keyword.originalText
	}
}

const cancelEdit = (keyword: IKeyword) => {
	if (keyword.originalText !== undefined) {
		keyword.text = keyword.originalText
		delete keyword.originalText
	}
	keyword.isEditing = false
}

const deleteKeyword = (keywordToDelete: IKeyword) => {
	keywords.value = keywords.value.filter(
		(keyword) => keyword !== keywordToDelete
	)
}
const restore = async () => {
	await keyword.resetKeywords()
	keywords.value = keyword.keywords
	push.warning("Keywords restored to default")
}
</script>
<template>
	<v-card class="pa-2">
		<v-card-title class="text-h6">Keywords</v-card-title>

		<v-card-text>
			<p class="text-caption font-italic mb-2">
				Set the keywords you want to monitor in chat. If a set keyword
				is detected, an alert will be sent to your Telegram chat.
			</p>
			<v-table density="compact">
				<thead>
					<tr>
						<th class="text-left px-2" width="40px">
							<div
								class="d-flex align-center justify-space-between"
							>
								<v-checkbox-btn
									v-model="enableAll"
									@click="toggleEnableAll"
									density="compact"
								></v-checkbox-btn>
							</div>
						</th>
						<th class="text-left">Keyword</th>
						<th class="text-left" width="80px">
							<v-btn-group>
								<v-btn
									size="small"
									icon="mdi-plus"
									@click="startAddNew"
									:disabled="isAdding"
								></v-btn>

								<v-btn
									size="small"
									icon="mdi-restore"
									@click="restore"
								></v-btn>
							</v-btn-group>
						</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="(keyword, index) in keywords" :key="index">
						<td class="checkbox-column px-2">
							<v-checkbox-btn
								v-model="keyword.enable"
								:disabled="keyword.isNew"
								density="compact"
							></v-checkbox-btn>
						</td>
						<td class="keyword-column">
							<template v-if="!keyword.isEditing">
								<div class="keyword-text">
									{{ keyword.text }}
								</div>
							</template>
							<v-text-field
								v-else
								density="compact"
								v-model="keyword.text"
								variant="underlined"
								autofocus
								hide-details
								@keyup.enter="
									keyword.isNew
										? saveNew(keyword)
										: saveEdit(keyword)
								"
							></v-text-field>
						</td>
						<td class="action-column">
							<template v-if="!keyword.isEditing">
								<v-btn-group>
									<v-btn
										size="small"
										icon="mdi-pencil"
										@click="startEditing(keyword)"
									></v-btn>
									<v-btn
										size="small"
										icon="mdi-delete"
										@click="deleteKeyword(keyword)"
									></v-btn>
								</v-btn-group>
							</template>
							<template v-else>
								<v-btn-group>
									<v-btn
										size="small"
										icon="mdi-check"
										@click="
											keyword.isNew
												? saveNew(keyword)
												: saveEdit(keyword)
										"
										:disabled="!keyword.text.trim()"
									></v-btn>
									<v-btn
										size="small"
										icon="mdi-close"
										@click="
											keyword.isNew
												? cancelNew()
												: cancelEdit(keyword)
										"
									></v-btn>
								</v-btn-group>
							</template>
						</td>
					</tr>
				</tbody>
			</v-table>
		</v-card-text>
	</v-card>
</template>

<style scoped>
.checkbox-column {
	width: 40px;
	min-width: 40px;
	vertical-align: top;
	padding-top: 8px !important;
}

.keyword-column {
	vertical-align: top;
	padding: 8px 12px !important;
}

.action-column {
	width: 80px;
	min-width: 80px;
	vertical-align: top;
	padding-top: 8px !important;
}

.keyword-text {
	word-wrap: break-word;
	white-space: pre-wrap;
	line-height: 1.2;
	min-height: 32px;
	display: flex;
	align-items: center;
}

:deep(.v-table .v-table__wrapper > table > tbody > tr > td) {
	height: auto;
	min-height: 40px;
}
</style>
