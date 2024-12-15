<script setup lang="ts">
import { ref, watch } from "vue"
import { ResetKeywords } from "../../../wailsjs/go/main/App"
import { useKeywordStore } from "@/stores/keyword"
import { push } from "notivue"

export interface IKeyword {
	text: string
	enable: boolean
	isEditing?: boolean
	isNew?: boolean
	originalText?: string
}

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
		<v-card-title class="text-h5">Keywords</v-card-title>

		<v-card-text>
			<v-table>
				<thead>
					<tr>
						<th class="text-left" width="56px">
							<div
								class="d-flex align-center justify-space-between"
							>
								<v-checkbox-btn
									v-model="enableAll"
									@click="toggleEnableAll"
								></v-checkbox-btn>
							</div>
						</th>
						<th class="text-left">Keyword</th>
						<th class="text-left" width="100px">
							<v-btn-group>
								<v-btn
									icon="mdi-plus"
									@click="startAddNew"
									:disabled="isAdding"
								></v-btn>

								<v-btn
									icon="mdi-restore"
									@click="restore"
								></v-btn>
							</v-btn-group>
						</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="(keyword, index) in keywords" :key="index">
						<td class="checkbox-column">
							<v-checkbox-btn
								v-model="keyword.enable"
								:disabled="keyword.isNew"
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
								v-model="keyword.text"
								density="compact"
								variant="underlined"
								autofocus
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
										icon="mdi-pencil"
										@click="startEditing(keyword)"
									></v-btn>
									<v-btn
										icon="mdi-delete"
										@click="deleteKeyword(keyword)"
									></v-btn>
								</v-btn-group>
							</template>
							<template v-else>
								<v-btn-group>
									<v-btn
										icon="mdi-check"
										@click="
											keyword.isNew
												? saveNew(keyword)
												: saveEdit(keyword)
										"
										:disabled="!keyword.text.trim()"
									></v-btn>
									<v-btn
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
	width: 56px;
	min-width: 56px;
	vertical-align: top;
	padding-top: 12px !important;
}

.keyword-column {
	vertical-align: top;
	padding: 12px 16px !important;
}

.action-column {
	width: 100px;
	min-width: 100px;
	vertical-align: top;
	padding-top: 12px !important;
}

.keyword-text {
	word-wrap: break-word;
	white-space: pre-wrap;
	line-height: 1.5;
	min-height: 40px;
	display: flex;
	align-items: center;
}

:deep(.v-table .v-table__wrapper > table > tbody > tr > td) {
	height: auto;
	min-height: 48px;
}
</style>
