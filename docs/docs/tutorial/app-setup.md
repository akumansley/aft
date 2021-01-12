---
id: app-setup
title: App setup
---

First, let's make a new file, `app.js`, and create a root component for our app.

```js title="app.js"
const App = {
	template: `
	<div>Hello App!</div>
	`,
}

export const app = Vue.createApp(App)
```

And then let's bootstrap it in our `index.html` file:

```html title="index.html"
<script type=module>
	import {app} from './app.js';
	window.onload = () => {
		app.mount("#app");
	}
</script>
```

Hit refresh on the client—you should see your app rendering its greeting.

We're using JS modules to keep our code clean and separated, which is lovely, but it's worth noting that modules are only supported in modern browsers. 

## API client

Now we'll add a small API client—some objects that will make it easy for us to talk to Aft.

Make a new file, `api.js`:

```js title="api.js"
export class APIError extends Error {
	constructor(message) {
		super(message);
		this.name = "APIError"; 
	}
}

async function makeRPC(name, args) {
	try {
		const result = await fetch('/rpc/' + name, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({args: args}),
		})
		const response = await result.json();
		if (response.code) {
			// server error
			throw new APIError(response.message);
		}
		return response.data;
	} catch (e) {
		// client error
		throw new APIError(e.message);
	}
}
```

The function `makeRPC` is a helper function that handles the "transport" aspect of talking to Aft, and armed with it, we can make local methods for the first two RPCs we're going to use: `me` and `login`.


```js title="api.js"
export const api = {
	async me() {
		return makeRPC('me', {})
	},
	async login(email, password) {
		return makeRPC('login', {
			email: email,
			password: password,
		})
	},
}
```

Okay, nice work! Up next, we'll make our login UI and sign in to our app.