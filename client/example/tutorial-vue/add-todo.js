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