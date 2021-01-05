async function makeRPC(name, args) {
	const result = await fetch('http://localhost:8080/rpc/' + name, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({args: args}),
	})
	const response = await result.json();
	return response.data;
}

const api = {
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



