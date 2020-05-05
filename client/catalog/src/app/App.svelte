
<script>
	import Nav from './Nav.svelte';
	import ObjectList from './ObjectList.svelte';
	import ObjectDetail from './ObjectDetail.svelte';
	import ObjectNew from './ObjectNew.svelte';
	import Breadcrumbs from './Breadcrumbs.svelte';
	import LogList from './LogList.svelte';

	// "Minimalist" Routing
	import navaid from 'navaid';
	const router = navaid();
	let params = null;
	let page;
	const routes = {
		"/object/:id": ObjectDetail,
		"/objects/new": ObjectNew,
		"/objects": ObjectList,
		"/log": LogList,
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
	:global(:root) {
		--background: #0d0a10;
		--background-highlight: #130f17;
		--text-color: #e4e1e8;
		--text-color-darker: #635b6d;
		--border-color: #2b2533;
		--highlight-color: #543c6c;

		--scale-4: 2.074em;
		--scale-3: 1.728em;
		--scale-2: 1.44em;
		--scale-1: 1.2em;
		--scale-0: 1em;
		--scale--1: .833em;
		--scale--2: .694em;
		--scale--3: .579em;
		

	}
	:global(body) {
		padding: 0;
		font-size: 18px;
		line-height: 1.7;
		-webkit-font-smoothing: antialiased;
		text-rendering: optimizeLegibility;
		background: var(--background);
		color: var(--text-color);
		font-family: "Inter",sans-serif;
	}
	:global(::selection) {
		background: #0041c6;
		color: #f5f6ff;
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
		border-bottom: 1px solid var(--border-color);
	}
	#nav {
		grid-area: nav;
		border-right: 1px solid var(--border-color);
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
