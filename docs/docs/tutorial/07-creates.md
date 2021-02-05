---
id: creates
title: Creating
---

Let's add a new component, `AddTodo`.

```js title="app.js"
function AddTodo({user, setTodos, todos}) {
	const [text, setText] = useState("");
	const [errorMessage, setErrorMessage] = useState(null);

	const addTodo = useCallback(async () => {
		setErrorMessage(null)
		try {
			const todo = await aft.api.todo.create({
				data: {
					text: text,
					user: {connect: {id: user.id}},
				}
			});
			setText("");
			setTodos((todos) => [...todos, todo]);
		} catch (e) {
			setErrorMessage(e.message);
		}
	}, [todos, user, text])

	return html`
	${errorMessage && html`<div class="row error">${errorMessage}</div>`}
	<div class="row">
		<input type=text
			value=${text} 
			onInput=${e => setText(e.target.value)} />
		<button onClick=${addTodo}>Add</button>
	</div>`
}
```

And lastly, let's add it to the todo list in `app.js` by putting it in our template at the end of the list.

```js title="app.js" {29}
function Todos({user, setUser}) {
	...
	
	return html`
	<div class="box">
		<div class="row">
			<div><b>Todos</b></div><a onClick=${signout}>Sign out</a>
		</div>
		${todos.map(todo => {
			return html`<${Todo} key=${todo.id} todo=${todo} />`
		})}
		<${AddTodo} user=${user} setTodos=${setTodos} todos=${todos}/>
	</div>`
}
```

Try it out! 

Now we can write our todos. In the next section, we'll look at updating existing todos.