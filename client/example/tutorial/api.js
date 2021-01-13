class APIError extends Error {}

async function call(mount, name, body) {
	const result = await fetch('/' + mount + '/' + name, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(body),
	})
	const response = await result.json();
	if (response.code) {
		throw new APIError(response.message);
	}
	return response.data;
}

export const api = {
	async me() {
		return call('rpc', 'me', {})
	},
	async login(email, password) {
		return call('rpc', 'login', {
			email: email,
			password: password,
		})
	},
	async loadTodos() {
		return call('api', 'todo.findMany', {
			"where": {"done": false}
		});
	},
	async createTodo(text, user) {
		return call('api', 'todo.create', 
		{
			"data": {
				"text": text, 
				"user": {"connect": {"id": user.id}}
			}
		});
	},
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
}
