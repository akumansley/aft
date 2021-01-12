export class APIError extends Error {
	constructor(message) {
		super(message);
		this.name = "APIError"; 
	}
}

async function makeRPC(name, args) {
	try {
		const result = await fetch('/rpc/' + name, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({args: args}),
		})
		const response = await result.json();
		if (response.code) {
			// server error
			throw new APIError(response.message);
		}
		return response.data;
	} catch (e) {
		// client error
		throw new APIError(e.message);
	}
}

async function makeAPICall(model, method, body) {
	try {
		const result = await fetch('/api/' + model + '.' + method, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(body),
		})
		const response = await result.json();
		if (response.code) {
			// server error
			throw new APIError(response.message);
		}
		return response.data;
	} catch (e) {
		throw new APIError(e.message);
	}
}

export const api = {
	async me() {
		return makeRPC('me', {})
	},
	async login(email, password) {
		return makeRPC('login', {
			email: email,
			password: password,
		})
	},
	async loadTodos() {
		return makeAPICall('todo', 'findMany', {
			"where": {"done": false}
		});
	},
	async createTodo(text, user) {
		return makeAPICall('todo', 'create', 
			{
				"data": {
					"text": text, 
					"user": {"connect": {"id": user.id}}
				}
			});
	},
	async updateTodoText(id, text) {
		return makeAPICall('todo', 'update', 
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
		return makeAPICall('todo', 'update', 
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
