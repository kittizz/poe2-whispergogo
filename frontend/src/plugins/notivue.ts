import { createNotivue } from "notivue"

import "notivue/notification.css" // Only needed if using built-in <Notification />
import "notivue/animations.css" // Only needed if using default animations
import "notivue/notification-progress.css"

export default createNotivue({
	position: "top-center",
	limit: 4,
	enqueue: true,
	avoidDuplicates: true,
	pauseOnHover: true,
	notifications: {
		global: {
			duration: 2000,
		},
	},
})
