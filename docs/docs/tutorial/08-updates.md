---
id: updates
title: Updating
---

In this section we're going to add the ability to check off completed Todos.

First, we'll add a hook to toggle a Todo, calling the Aft API's `update` method.

```js title="app.js"
function Todos({user, setUser}) {
	...
}

function useToggle(setTodos, todo) {
	const toggle = async () => {
		const updated = await aft.api.todo.update({
			where: {id: todo.id},
			data: {done: !todo.done},
		})
		setTodos(todos => {
			return todos.map(t => t.id === todo.id? updated: t)
		})
	}
	return toggle
}
```

And then update our Todos template to pass the hook to the todo component, and the Todo component to call it.

```js title="app.js"
function Todos({user, setUser}) {
	...
	return html`
	<div class="box">
		<div class="row">
			<div><b>Todos</b></div><a onClick=${signout}>Sign out</a>
		</div>
		${todos.map(todo => {
			const toggle = useToggle(setTodos, todo);
			return html`<${Todo} key=${todo.id} toggle=${toggle} todo=${todo}/>`
		})}
		<${AddTodo} user=${user} setTodos=${setTodos} todos=${todos}/>
	</div>`
}

function Todo({todo, toggle}) {
	return html`
	<div class="row">
		<div>${todo.text}</div>
		<input type=checkbox onClick=${toggle} checked=${todo.done} />
	</div>`
}
```

Try it out! Todos should stick around when checked, but not get reloaded on a page refresh.

In our last step, we'll take a look at access controls so we make sure we're only letting users see their own todos.