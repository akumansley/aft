
<script>
	import Nav from './Nav.svelte';
	import Main from './Main.svelte';
</script>
<style>
	:global(body) {
		padding: 0;
		font-size: 18px;
		line-height: 1.5em;
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
		grid-template-columns: 200px 1fr;
		grid-template-rows: 1fr;
		grid-template-areas: ". .";
	}
</style>
<svelte:head>
	<title>Aft</title>
</svelte:head>

<div id="grid-root">
<Nav/>
<Main/>
</div>
