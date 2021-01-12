import {api} from './api.js'
import {userStore} from './store.js'

export const Login = {
	data() {
		return {
			email: "",
			password: "",
			errorMessage: null,
		}
	},
	methods: {
		async login() {
			this.errorMessage = null;
			try {
				userStore.value = await api.login(this.email, this.password)
			} catch (err) {
				this.errorMessage = err.message;
			}
		}
	},
	template: `
	<div class="card card-body" style="width: 20rem;">
		<div class="form-group">
			<label>Email</label>
			<input class="form-control" v-model=email type="email" />
		</div>
		<div class="form-group">
			<label>Password</label>
			<input class="form-control" v-model=password type="password"/>
		</div>
		<div v-if="errorMessage" class="alert alert-danger">{{errorMessage}}</div>
		<button @click="login" class="btn btn-primary">Login</button>
	</div>
	`,
}