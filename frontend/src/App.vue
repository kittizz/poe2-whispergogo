<script lang="ts" setup>
import { Notivue, Notification, darkTheme, NotificationProgress } from "notivue"
</script>
<template>
	<v-app>
		<v-main>
			<router-view />
		</v-main>
	</v-app>
	<Notivue v-slot="item">
		<Notification :item="item" :theme="darkTheme">
			<NotificationProgress :item="item" />
		</Notification>
	</Notivue>
</template>
