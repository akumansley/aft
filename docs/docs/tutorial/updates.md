---
id: updates
title: Updating
---

And again we'll add to our API to include two methods for updating todos.

```js title="api.js"
export const api = {
	...
	async updateTodoText(id, text) {
		return call('api', 'todo.update', 
		{
			"where": {
				"id": id,
			},
			"data": {
				"text": text, 
			},
		});
	},
	async updateTodoDone(id, done) {
		return call('api', 'todo.update', 
		{
			"where": {
				"id": id,
			},
			"data": {
				"done": done, 
			},
		});
	}
	...
}
```

Then call them from a new component `todo.js`.

```js title="todo.js"
import {api} from './api.js'
import {todoStore, userStore} from './store.js'

export const Todo = {
	props: ['todo'],
	data() {
		return {
			errorMessage: null,
		}
	},
	methods: {
		async saveText() {
			this.errorMessage = null;
			try {
				await api.updateTodoText(this.todo.id, this.todo.text);
			} catch (err) {
				this.errorMessage = err.message;
			}
		},
		async saveDone() {
			this.errorMessage = null;
			try {
				await api.updateTodoDone(this.todo.id, this.todo.done);
			} catch (err) {
				this.errorMessage = err.message;
			}
		}
	},
	template: `
		<div class="list-group-item todo-item">
			<div class="input-group">
				<input @change="saveText" type="text" v-model="todo.text" class="form-control" />
				<div class="input-group-append">
					<div class="input-group-text">
						<input @change="saveDone" type="checkbox" v-model="todo.done" />
					</div>
				</div>
			</div>
			<div v-if="errorMessage" class="alert alert-danger">{{errorMessage}}</div>
		</div>
	`,
}
```

And lastly, we'll add it to the todo list in `todos.js`, by registering it the `components` key, and then adding it to our template at the end of the list.

```js title="todos.js"
import {api} from './api.js'
import {AddTodo} from './add-todo.js'
import {Todo} from './todo.js'
import {todoStore} from './store.js'

export const Todos = {
	components: {
		'add-todo': AddTodo,
		'todo': Todo,
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
		<todo v-for="todo in todoStore.todos" :todo="todo"/>
		<add-todo />
	</div>
	`,
}
```

That's our app—looking great!

In our last step, we'll take a look at access controls so we make sure we're only letting users see their own todos.