<script>
	import client from './client.js';
	import user from './user.js';
	import { navigate, Link } from "svelte-routing";

	let email = "";
	let password = "";

	async function login() {
		let userResp = await client.rpc.login({
			"email": email,
			"password": password,
		});
		if (userResp.code){
			console.log(userResp)
		} else {
			user.set(userResp);
			navigate("/", {replace: true});
		}
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
<button on:click={login}>Sign in</button>
</div>

<hr/>
<div>
	New to Giftr?
<Link to="/signup">
		Create an account
</Link>
</div>

