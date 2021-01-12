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