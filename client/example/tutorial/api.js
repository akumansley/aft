export class APIError extends Error {
	constructor(message) {
		super(message);
		this.name = "APIError"; 
	}
}

async function makeRPC(name, args) {
	try {
		const result = await fetch('http://localhost:8080/rpc/' + name, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({args: args}),
		})
		const response = await result.json();

		// TODO remove this hack
		if (response.code) {
			throw new APIError(response.message);
		}
		if (response.data.code) {
			throw new APIError(response.data.message);
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
	}
}
