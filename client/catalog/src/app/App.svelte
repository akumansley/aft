<script>
	import Nav from './Nav.svelte';
	import ModelList from './models/ModelList.svelte';
	import ModelDetail from './models/ModelDetail.svelte';
	import ModelNew from './models/ModelNew.svelte';
	import InterfaceNew from './models/InterfaceNew.svelte';
	import InterfaceDetail from './models/InterfaceDetail.svelte';
	import DatatypeList from './datatypes/DatatypeList.svelte';
	import EnumDetail from './datatypes/EnumDetail.svelte';
	import CoreDatatypeDetail from './datatypes/CoreDatatypeDetail.svelte';
	import EnumNew from './datatypes/EnumNew.svelte';
	import RPCList from './rpc/RPCList.svelte';
	import RPCDetail from './rpc/RPCDetail.svelte';
	import RPCNew from './rpc/RPCNew.svelte';
	import Terminal from './Terminal.svelte';
	import LogList from './LogList.svelte';
	import RoleList from './access/RoleList.svelte';
	import RoleDetail from './access/RoleDetail.svelte';
	import RoleNew from './access/RoleNew.svelte';
	import {router, canRoute} from './router.js';
	import {routeStore} from './stores.js';
	import {checkSave} from './save.js';
	
	let params = null;
	let page;
	const routes = {
		"/schema": ModelList,
		"/model/:id": ModelDetail,
		"/models/new": ModelNew,

		"/interfaces/new": InterfaceNew,
		"/interface/:id": InterfaceDetail,

		"/enum/:id": EnumDetail,
		"/enums/new": EnumNew,
		"/core/:id": CoreDatatypeDetail,
		"/datatypes": DatatypeList,

		"/terminal": Terminal,

		"/rpc/:id": RPCDetail,
		"/rpcs": RPCList,
		"/rpcs/new":RPCNew,

		"/log": LogList,

		"/roles": RoleList,
		"/role/:id": RoleDetail,
		"/roles/new": RoleNew,

		"/": ModelList,
	};
	for (const [route, component] of Object.entries(routes)) {
		router.on(route, (urlps) => {
			page = component;

			routeStore.set(route);

			if (Object.keys(urlps).length !== 0) {
				params = urlps;
			} else {
				params = null;
			}
		});
	}
	router.listen();

	//If the user navigates from a page with unsaved changes, then alert them
	window.onbeforeunload = (e) => {
		return canRoute(e);
	}
	
</script>
<style>
	:global(:root) {
		--background: #1e1a23;
		--background-highlight: #302937;
		--text-color: #f4f3f6;
		--text-color-lighter: #fff;
		--text-color-darker: #635b6d;
		--text-color-function: #50fa7b;
		--border-color: #4c4359;
		--highlight-color: #543c6c;

		--scale-4: 2.074em;
		--scale-3: 1.728em;
		--scale-2: 1.44em;
		--scale-1: 1.2em;
		--scale-0: 1em;
		--scale--1: .833em;
		--scale--2: .694em;
		--scale--3: .579em;
		--box-margin: .75em;
		

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
	:global(h1) {
		font-size: var(--scale-3);
		font-weight: 600;
		margin: 0;
	}
	
	:global(h2) {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
	#grid-root {
		position: absolute;
		height: 100%;
		width: 100%;
		display: grid;
		grid-template-columns: 10em 1fr;
		grid-template-rows: 0em 1fr;
		grid-template-areas: "nav head"
		"nav main";
	}
	#nav {
		grid-area: nav;
		border-right: 1px solid var(--border-color);
	}
	#main {
		grid-area: main;
		height: 100vh;
		display: flex;
		flex-direction: column;
	}
</style>

<svelte:head>
	<title>Aft</title>
	<link rel="stylesheet"
		href="https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700&display=swap">
</svelte:head>

<svelte:window on:keydown={checkSave(()=>{})}/>

<div id="grid-root">
	<div id="nav">
		<Nav />
	</div>
	<div id="main">
		{#if params}
			<svelte:component this={page} {params} />
		{:else}
			<svelte:component this={page} />
		{/if}
	</div>
</div>
