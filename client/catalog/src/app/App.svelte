
<script>
	import Nav from './Nav.svelte';
	import Main from './Main.svelte';
</script>
<style>
	:global(body) {
		padding: 0;
		font-size: 18px;
		line-height: 1.5em;
		font-family: Roboto, sans-serif;
		-webkit-font-smoothing: antialiased;
	}
	:global(p) {
		line-height: 1.5;
		margin-top: 1.5em;
		margin-bottom: 1.5em;
	}
	#grid-root {
		position: absolute;
		height: 100%;
		width: 100%;
		display: grid;
		grid-template-columns: 10em 1fr;
		grid-template-rows: 2em 1fr 1em;
		grid-template-areas: "nav head"
		"nav main"
		"nav foot";
	}
	#head {
		grid-area: head;
	}
	#foot {
		grid-area: foot;
	}
	#nav {
		grid-area: nav;
	}
	#main {
		grid-area: main;
	}
</style>
<svelte:head>
	<title>Aft</title>
</svelte:head>

<div id="grid-root">
	<div id="head"></div>
	<div id="nav">
		<Nav/>
	</div>
	<div id="main">
		<Main/>
	</div>
	<div id="foot"></div>
</div>
