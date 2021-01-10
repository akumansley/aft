import {api} from './api.js'
import {Login} from './login.js'
import {userStore} from './store.js'

const App = {
	components: {
		'login-view': Login,
	},
	template: `
	<div v-if="!loaded">Loading..</div>
	<login-view v-else-if="userStore.value === null"></login-view>
	<div v-else>hello {{userStore.value.email}}</div>
	`,
	created() {
		userStore.load = api.me().then((user) => {
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
