<script>
	import Login from './Login.svelte';
	import user from './user.js';
	import client from './client.js';

	function logout() {
		user.set(null);
		document.cookie = "tok= ; expires = Thu, 01 Jan 1970 00:00:00 GMT";
	}

</script>


<style>
	:global(:root) {
		--color-1: #522D6D;
		--color-2: #B578A3;
		--color-3: #C681AF;
		--color-4: #462960;
		--color-5: #7C3C7C;
		--color-6: #644CD8;

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

	}

	h1 {
		font-size: var(--scale-2);
		margin: 0;
	}
	.logout {
		margin-top: .5em;
	}

</style>

<main>
	<h1>Giftr</h1>
	{#if $user === null} 
	<Login />
	{:else} 

	<div>
		Welcome to Giftr {$user.email}!
	</div>

	<div class="logout">
		<button on:click={logout}>logout</button>
	</div>
	{/if}

</main>
