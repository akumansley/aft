
<script>
	import Nav from './Nav.svelte';
	import ObjectList from './ObjectList.svelte';
	import ObjectDetail from './ObjectDetail.svelte';
	import ObjectNew from './ObjectNew.svelte';
	import Breadcrumbs from './Breadcrumbs.svelte';

	// "Minimalist" Routing
	import navaid from 'navaid';
	const router = navaid();
	let params = null;
	let page;
	const routes = {
		"/object/:id": ObjectDetail,
		"/objects/new": ObjectNew,
		"/objects": ObjectList,
		"/": ObjectList,
	};
	for (const [route, component] of Object.entries(routes)) {
		router.on(route, (urlps) => {
			page = component;
			if (Object.keys(urlps).length !== 0) {
				params = urlps;
			} else {
				params = null;
			}
		});
	}
	router.listen();
</script>
<style>
	:global(body) {
		padding: 0;
		font-size: 18px;
		line-height: 1.7;
		-webkit-font-smoothing: antialiased;
		background: #0d0a10;
		color: #cac4d1;
		font-family: "Inter";
	}
	:global(p) {
		line-height: 1.5;
		margin-top: 1.5em;
		margin-bottom: 1.5em;
	}
	:global(a) {
		color:inherit;
	}
	:global(a):visited {
		color: inherit;
	}
	#grid-root {
		position: absolute;
		height: 100%;
		width: 100%;
		display: grid;
		grid-template-columns: 10em 1fr;
		grid-template-rows: 3em 1fr;
		grid-template-areas: "nav head"
		"nav main";
	}
	#head {
		padding: .5em 1.5em;
		border-bottom: 1px solid #2b2533;
	}
	#nav {
		grid-area: nav;
		border-right: 1px solid #2b2533;
	}
	#main {
		grid-area: main;
	}
</style>
<svelte:head>
	<title>Aft</title>
	<link href="https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700&display=swap" rel="stylesheet">

</svelte:head>

<div id="grid-root">
	<div id="nav">
		<Nav/>
	</div>
		<div id="head">
			<Breadcrumbs/>
		</div>
	<div id="main">
		{#if params}
			<svelte:component this={page} {params} />
		{:else}
			<svelte:component this={page} />
		{/if}
	</div>
</div>
