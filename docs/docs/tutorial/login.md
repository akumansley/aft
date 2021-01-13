---
id: login
title: Login
---

First, let's create a store for tracking the signed in user.

```js title="store.js"
export const userStore = Vue.reactive({value: null})
```

We'll start off by building our small login widget in a new file, `login.js`.

```js title="login.js"
import {api} from './api.js'
import {userStore} from './store.js'

export const Login = {
	data() {
		return {
			email: "",
			password: "",
			errorMessage: null,
		}
	},
	methods: {
		async login() {
			this.errorMessage = null;
			try {
				userStore.value = await api.login(this.email, this.password)
			} catch (err) {
				this.errorMessage = err.message;
			}
		}
	},
	template: `
	<div class="card card-body" style="width: 20rem;">
		<div class="form-group">
			<label>Email</label>
			<input class="form-control" v-model=email type="email" />
		</div>
		<div class="form-group">
			<label>Password</label>
			<input class="form-control" v-model=password type="password"/>
		</div>
		<div v-if="errorMessage" class="alert alert-danger">{{errorMessage}}</div>
		<button @click="login" class="btn btn-primary">Login</button>
	</div>
	`,
}
```

It's not wired up to anything yet, so let's go add it to `app.js`, replacing our "hello" message.

```js title="app.js"
import {Login} from './login.js'

const App = {
	components: {
		'login-view': Login,
	},
	template: `
	<login-view />
	`,
}

export const app = Vue.createApp(App)
```


## Adding a user

If you go ahead and try to sign in to our tutorial app, you should of course get an error about the login not working—we haven't added any users yet!

Open up Aft, and navigate to the **Terminal**, and we'll create a user.

```python
def main(aft):
    return aft.api.create("user", {"data":{"email":"user@example.com", "password":"coolpass"}})
```

Press **Run**, and you should see a JSON representation of the user just created, though the password is salted and hashed. 

Go back to the tutorial app, and try signing in with your new credentials. If you open DevTools, you should see a user object in the result of the `console.log` statement.

In the next section, we'll finish up the login system, using the `me` RPC.

