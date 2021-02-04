---
id: app-setup
title: App Setup
---

First, let's make a new file, `app.js`, and move our App component over, making sure to export it.

```js title="app.js"
import {html} from 'https://unpkg.com/htm/preact/standalone.module.js'

export function App() {
	return html`<h1>Hello Aft!</h1>`
}
```

And then import it in our `index.html` file:

```html title="index.html"
<head>
	<link rel=stylesheet href="./styles.css" />
	<script type=module>
		import {html, render} from 'https://unpkg.com/htm/preact/standalone.module.js'
		import {App} from './app.js'

		render(html`<${App} />`, document.body);
	</script>
</head>
```

Hit refresh on the client—you should still see your app rendering its greeting.

## API client

Now we'll add a small API client—some objects that will make it easy for us to talk to Aft.

Make a new file, `api.js`, and add the following.

```js title="api.js"
async function call(path, body) {
	const result = await fetch(path, {
		method: 'POST',
		headers: {'Content-Type': 'application/json'},
		body: JSON.stringify(body || {}),
	})
	const response = await result.json();
	if (response.code) {
		throw new Error(response.message);
	}
	return response.data;
}

const curryProxy = (inner) => {
	return new Proxy({}, {
		get(_, prop)  { 
			return inner(prop) 
		}
	})
}

export default {
	api: curryProxy((interfaceName) => curryProxy((method) => (params) => {
		return call("api/" + interfaceName + "." + method, params)
	})),
	rpc: curryProxy((rpcName) => (args) => {
		return call("rpc/" + rpcName, args)
	}),
}
```

The use of Proxy isn't really necessary, but it gives us a nice looking syntax for making API calls or RPCs to Aft. This short snippet is all you'll need in your app to use every bit of functionality Aft has to offer.

Okay, nice work! Up next, we'll make our login UI and sign in to our app.
