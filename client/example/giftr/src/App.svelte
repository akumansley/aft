<script>
	import { Router, Link, Route, navigate } from "svelte-routing";

	import Login from './Login.svelte';
	import Index from './Index.svelte';
	import List from './List.svelte';
	import Signup from './Signup.svelte';

	import client from './client.js';

	import user from './user.js';

	async function getUser() {
		try {
			let userResp = await client.rpc.me({});
			user.set(userResp);
		} catch (err) {
			navigate("/login", {replace: true});
			console.log("error", err);
		}
	}
	let load = getUser();

	function logout() {
		document.cookie = "tok= ; expires = Thu, 01 Jan 1970 00:00:00 GMT";
		navigate("/login", {replace: true});
		user.set(null);
	}
	export let url = "";


</script>


<style>
	:global(:root) {
		--color-1: #522D6D;
		--color-1-darker: #361c4a;
		--color-3: #e799cd;

		--scale-4: 2.074em;
		--scale-3: 1.728em;
		--scale-2: 1.44em;
		--scale-1: 1.2em;
		--scale-0: 1em;
		--scale--1: .833em;
		--scale--2: .694em;
		--scale--3: .579em;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
	}
	:global(body) {
		background: var(--color-1);
	}
	:global(a) {
		color: var(--color-3);
	}

	main {
		font-size: 18px;
		line-height: 1.7;
		-webkit-font-smoothing: antialiased;
		text-rendering: optimizeLegibility;
		color: #fff;
		padding: .5em 1em;
		border-radius: .5em;
		max-width: 600px;
		margin: 0 auto;
		display: grid;
		grid-template-columns: 1fr 1fr;
	}

	.account {
		grid-column: 2;
		font-size: var(--scale--1);
		align-self: end;
		line-height: calc(var(--scale-2) * 1.7);
	}
	.content {
		grid-column-start: 1;
		grid-column-end: 3;
		margin-top: 1em;
	}


	h1 {
		font-size: var(--scale-2);
		margin: 0;
		grid-column: 1;
	}
	.link-button {
		appearance: none;
		color: var(--color-3);
		background: none;
		border: none;
		display: inline;
		padding: 0;
		cursor: pointer;
	}

</style>
<Router url="{url}">

	<main>
		<h1>Giftr</h1>
		{#await load then _}

		{#if $user !== null} 
		<div class="account">
			{$user.email} <button class="link-button" on:click={logout}>(sign out)</button>
		</div>
		{/if}

		<div class="content">
			<Route path="login">
				<Login/>
			</Route>
			<Route path="signup" >
				<Signup/>
			</Route>
			<Route path="list/:listUserId" let:params>
				<List listUserId={params.listUserId}/>
			</Route>
			<Route path="/" >
				<Index/>
			</Route>
		</div>
		{/await}
	</main>
</Router>
