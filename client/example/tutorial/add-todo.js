import {api} from './api.js'
import {todoStore, userStore} from './store.js'

export const AddTodo = {
	data() {
		return {
			text: "",
			submitting: false,
			errorMessage: null,
		}
	},
	methods: {
		async submit() {
			this.submitting = true;
			try {
				this.errorMessage = null;
				let todo = await api.createTodo(this.text, userStore.value);
				todoStore.todos.push(todo);
				this.submitting = false;
				this.text = "";
			} catch (err) {
				this.submitting = false;
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