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