---
id: creates
title: Creating
---

Let's start by extending our API to include a `createTodo` method, very similar to the code we used in the **Terminal** to add a todo before.

```js title="api.js"
export const api = {
	...
	async createTodo(text, user) {
		return call('api', 'todo.create', 
		{
			"data": {
				"text": text, 
				"user": {"connect": {"id": user.id}}
			}
		});
	},
	...
}
```

Then let's call it from a new component `add-todo.js`.

```js title="add-todo.js"
import {api} from './api.js'
import {todoStore, userStore} from './store.js'

export const AddTodo = {
	data() {
		return {
			text: "",
			errorMessage: null,
		}
	},
	methods: {
		async submit() {
			try {
				this.errorMessage = null;
				let todo = await api.createTodo(this.text, userStore.value);
				todoStore.todos.push(todo);
				this.text = "";
			} catch (err) {
				this.errorMessage = err.message;
			}
		}
	},
	template: `
	<div class="list-group-item">
		<div class="form-group">
			<input class="form-control" v-model="text" type="text"/> 
		</div>
		<div v-if="errorMessage" class="alert alert-danger">{{errorMessage}}</div>
		<button @click="submit" class="btn btn-primary">Add</button>
	</div>
	`,
}
```

And lastly, let's add it to the todo list in `todos.js`, by registering it the `components` key, and then adding it to our template at the end of the list.

```js title="todos.js"
import {api} from './api.js'
import {AddTodo} from './add-todo.js'
import {todoStore} from './store.js'

export const Todos = {
	components: {
		'add-todo': AddTodo,
	},
	data() {
		return {
			loaded: false,
			todoStore: todoStore,
		}
	},
	created() {
		api.loadTodos().then((todos) => {
			todoStore.todos = todos;
			this.loaded = true;
		}, (err) => {
			this.error = err;
			this.loaded = true;
		});
	},
	template: `
	<div v-if="loaded" class="list-group" style="width: 20rem;">
		<div v-for="todo in todoStore.todos">{{todo.text}} - {{todo.done}}</div>
		<add-todo />
	</div>
	`,
}
```

Great! Now we can write our todos. In the next section, we'll look at updating existing todos.