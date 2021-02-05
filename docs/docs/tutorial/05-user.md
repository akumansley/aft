---
id: user
title: User
---

When a user signs into Aft successfully, Aft sets a cookie authenticating the user.

Upon starting up, we'll invoke the `me` RPC, which will return a user object if the login cookie is present.

And we'll add a link to call the `signout` RPC, which will clear the login cookie. 

We'll use the Preact `useEffect` hook to do this automatically when the app starts.


```js title="app.js"
import {html, useState, useCallback, useEffect} from 'https://unpkg.com/htm/preact/standalone.module.js'
import aft from './aft.js'

export function App (props) {
    const [user, setUser] = useState(null);
    const [loaded, setLoaded] = useState(false);

    useEffect(async () => {
        try {
            setUser(await aft.rpc.me());
        } catch {} finally {
            setLoaded(true);
        }
    }, []);

    const signout = async () => {
        await aft.rpc.logout();
        setUser(null);
    }

    if (!loaded) {
        return html``
    } else if (user === null) {
        return html`<${Login} setUser=${setUser} />`
    } else {
        return html`<h1>Hello ${user.email} - <a onClick=${signout}>Sign out</a></h1>`
    }
}

...
```

That's it! Make sure you can log in and out successfully, and then we're ready to give our users some data.