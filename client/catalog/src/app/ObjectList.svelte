<script>
	import { onMount } from 'svelte';
	import client from '../data/client.js';
	let objects = [];
	let load = client.model.findMany();
	let cap= (s) => { 
		return s.charAt(0).toUpperCase() + s.slice(1)

	};
</script>

<style>
	.box {
		padding-left: 1em;
		margin-top: 1em;
		display: flex;
	}
	.stuff {

	}
	a.object-box {
		color: inherit;
		width: 150px;
		padding: .5em .75em;
		background: rgba(0,0,0,.1);
		border-radius: 3px;
	}
	.spacer {
		width: 1em;
	}

</style>

<div class="box">
	{#await load}
		Loading..
	{:then objects}
		{#each objects as object}
			<a href="/objects/{object.Id}" class="object-box">
				<div><b>{cap(object.Name)}</b></div>
			{#each Object.entries(object.Attributes) as attr}
				<div>{attr[0]}</div>
			{/each}
			</a>
			<div class="spacer"/>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
