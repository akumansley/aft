import {api} from './api.js'
import {Login} from './login.js'
import {Todos} from './todos.js'
import {userStore} from './store.js'

const App = {
	components: {
		'login-view': Login,
		'todos': Todos,
	},
	template: `
	<div v-if="!loaded"></div>
	<login-view v-else-if="userStore.value === null" />
	<todos v-else />
	`,
	created() {
		api.me().then((user) => {
			userStore.value = user;
			this.loaded = true
		}, () => {
			this.loaded = true
		});
	},
	data() {
		return {
			userStore: userStore,
			loaded: false,
		}
	}
}

export const app = Vue.createApp(App)
