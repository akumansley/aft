---
id: login
title: Login
---

Let's jump back and build our login widget back in `app.js`. We'll also add a greeting for our user once they sign in.

```js title="app.js"
import {html, useState, useCallback} from 'https://unpkg.com/htm/preact/standalone.module.js'
import aft from './aft.js'


export function App (props) {
	const [user, setUser] = useState(null);
	if (user === null) {
		return html`<${Login} setUser=${setUser} />`
	} else {
		return html`<h1>Hello ${user.email}</h1>`
	}
}

function Login({setUser}) {
	const [errorMessage, setErrorMessage] = useState(null);
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	return html`
	<div class="box stack">
		<input type=email placeholder="Email" 
			value=${email} 
			onInput=${(e) => setEmail(e.target.value)}/>
		<input type=password placeholder="Password" 
			value=${password} 
			onInput=${(e) => setPassword(e.target.value)}/>
		${errorMessage && html`<div class=error>${errorMessage}</div>`}
		<button>Sign in</button>
	</div>`
}
```

Refresh and take a look at our login box. Looking good!

Now we'll try and actually connect it to Aft. Add a `submit` callback and connect it to the button's `onClick` property.

Notice how we're invoking the aft "login" RPC. Aft RPCs accept and return a single JSON object. We're able to just call the RPC by name like a native function thanks to the Proxy magic we did earlier in `aft.js`.

```js title="app.js"
...

function Login({setUser}) {
	const [errorMessage, setErrorMessage] = useState(null);
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const submit = useCallback(async () => {
		setErrorMessage(null);
		try {
			const user = await aft.rpc.login({
				email: email,
				password: password,
			});
			setUser(user);
		} catch (e) {
			setErrorMessage(e.message);
		}
	}, [email, password])

	return html`
	<div class="box stack">
		<input type=email placeholder="Email" 
			value=${email} 
			onInput=${(e) => setEmail(e.target.value)}/>
		<input type=password placeholder="Password" 
			value=${password} 
			onInput=${(e) => setPassword(e.target.value)}/>
		${errorMessage && html`<div class=error>${errorMessage}</div>`}
		<button onClick=${submit}>Sign in</button>
	</div>`
}
```

## Adding a user

If you go ahead and try to sign in to our tutorial app, you should of course get an error about the login not working—we haven't added any users yet!

Open up Aft at `http://localhost:8081`, and navigate to the **Terminal**, and we'll create a user.

```python
def main():
    return create("user", {"data": {
		    	"email": "user@example.com", 
		    	"password": "coolpass",
	    	}})
```

Press **Run**, and you should see a JSON representation of the user just created, though the password is salted and hashed. 

Go back to the tutorial app, and try signing in with your new credentials.

In the next section, we'll finish up the login system, using the `me` RPC.

