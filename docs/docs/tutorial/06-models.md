---
id: models
title: Models
---

Let's add another component in `app.js` to for our Todo data and reference it from our App component so it displays once the user is signed in. We'll also move our signout function into the Todos component, now that it has a better home.

```js title="app.js"
export function App (props) {
	...
	if (!loaded) {
		return html``
	} else if (user === null) {
		return html`<${Login} setUser=${setUser} />`
	} else {
		return html`<${Todos} user=${user} setUser=${setUser}/>`
	}
}

function Todos({user, setUser}) {
	const [todos, setTodos] = useState([]);

	const signout = async () => {
		await aft.rpc.logout();
		setUser(null);
	}

	return html`
	<div class="box">
		<div class="row">
			<div><b>Todos</b></div><a onClick=${signout}>Sign out</a>
		</div>
		${todos.map(todo => {
			return html`<${Todo} key=${todo.id} todo=${todo} />`
		})}
	</div>`
}

function Todo({todo}) {
	return html`
	<div class="row">
		<div>${todo.text}</div>
		<input type=checkbox checked=${todo.done} />
	</div>`
}

```

## Adding Models

Now let's create a basic Todo object on our backend. We'll have two attributes, `text`, a string which is the text of the Todo, and `done` a boolean indicating whether it's done or not. We'll also add a relationship to a `user` object—the owner of the Todo.

Switch back over to Aft, and navigate to the **Schema** section.

Click the **Add Model** button at the top of the screen. 

Fill in the text field that says *Model name..* with `todo`,

Then click the **add** button under Attributes and fill in our first attribute.

For *Attribute name..* type `text`, and select `String` from the dropdown on the right, indicating its type.

Click **add** again, and this time name the attribute `done`, and select `Bool` from the type dropdown.

Now let's add our relationship to `user` by clicking the **add** button under Relationships.

Set the *Relationship name..* to `user` and select `User` from the dropdown. That tells aft that this property points to a User object.

You can leave the multiple box unchecked, since this relationship will only be to a single user, rather than a list of users.

Once you've done that, click **Save** next to the model name at the top of the page, and you're all done!


## Reading data

Aft automatically adds our new Todo model to the API, so lets try and read some data from it.

First, navigate over to **Terminal** in Aft and run the following function to add a Todo.

```python
def main(aft):
    return aft.api.create("todo", {"data": {
    		"text":"connect the backend", 
	    	"done": False, 
	    	"user": {"connect": {"email": "user@example.come"}}
    	}})
```

And finally, let's update our Todos component to call the server!

Here we're making use the Aft API's `findMany` method. You'll note we don't filter on the currently signed in user—we'll be doing that automatically later on when we look at access controls. 

```js title="todos.js"
function Todos({user, setUser}) {
	const [todos, setTodos] = useState([]);
	const [loaded, setLoaded] = useState(false);

	useEffect(async () => {
		try {
			setTodos(await aft.api.todo.findMany({
				where: {
					done: false,
				}
			}));
		} catch {
		} finally {
			setLoaded(true);
		}
	}, []);

	const signout = async () => {
		await aft.rpc.logout();
		setUser(null);
	}

	if (!loaded) {
		return html``
	}

	return html`
	<div class="box">
		<div class="row">
			<div><b>Todos</b></div><a onClick=${signout}>Sign out</a>
		</div>
		${todos.map(todo => {
			return html`<${Todo} key=${todo.id} todo=${todo} />`
		})}
	</div>`
}
```

In the next section, we'll look at adding some data!
