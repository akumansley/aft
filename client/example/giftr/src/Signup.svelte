<script>
	import client from './client.js';
	import { navigate, Link } from "svelte-routing";
	import user from './user.js';

	let email = "";
	let password = "";

	async function signup() {
		let userResp = await client.rpc.signup({
			"email": email,
			"password": password,
		});
		user.set(userResp);
		navigate("/", {replace: true});
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
<input type=email bind:value={email} />
</label>

<label>
<div class="label">Password</div>
<input type=password bind:value={password} />
</label>

</div>

<div class="login">
<button on:click={signup}>Create account</button>
</div>
<hr/>
<div>
	Already have an account?
<Link to="/login">
		Sign in
</Link>
</div>
