<script>
	import client from './client.js';
	import user from './user.js';

	let email = "";
	let password = "";

	async function login() {
		let userResp = await client.views.login({
			"email": email,
			"password": password,
		});
		user.set(userResp);
	}
	
</script>

<style>
	.form {

	}
	.label {

	}
	.login {
		margin-top: .5em;
	}
</style>

<div class="form">

<label>
<div class="label">Email</div>
<input type=text bind:value={email} />
</label>

<label>
<div class="label">Password</div>
<input type=password bind:value={password} />
</label>

</div>

<div class="login">
<button on:click={login}>login</button>
</div>