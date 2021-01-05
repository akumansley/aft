import {api} from './api.js'
import {Login} from './login.js'
import {userStore} from './store.js'

const App = {
	components: {
		'login-view': Login,
	},
	template: `
	<login-view v-if="userStore.value === null">no user</login-view>
	<div v-else>hello {{userStore.value.email}}</div>
	`,
	data() {
		return {
			userStore: userStore,
		}
	}
}

export const app = Vue.createApp(App)
