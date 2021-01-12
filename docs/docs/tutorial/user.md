---
id: user
title: User
---

When a user signs into Aft successfully, Aft sets an cookie authenticating the user.

Upon starting up, we'll invoke the `me` RPC, which will return a user object if the login cookie is present.

```js title="app.js"
import {api} from './api.js'
import {userStore} from './store.js'
import {Login} from './login.js'

const App = {
	components: {
		'login-view': Login,
	},
	created() {
		api.me().then((user) => {
			userStore.value = user;
			this.loaded = true
		}, () => {
			this.loaded = true
		});
	},
	template: `
	<div v-if="!loaded"></div>
	<login-view v-else-if="userStore.value === null" />
	<div v-else>Hello {{userStore.value.email}}!</div>
	`,
	data() {
		return {
			userStore: userStore,
			loaded: false,
		}
	}
}

export const app = Vue.createApp(App)
```

That's it! We're ready to give our users some data.