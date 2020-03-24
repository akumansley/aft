<script>
	import { onMount } from 'svelte';
	import client from '../data/client.js';
	let objects = [];
	let load = client.objects.list();
</script>

<style>
	.box {
		padding-left: 1em;
	}
	.stuff {
		margin-top: 1.25em;

	}

</style>

<div class="box">
	<div class="stuff">
		Objects
	</div>
	{#await load}
		Loading..
	{:then objects}
		{#each objects as object}
			<div><a href="/objects/{object.id}">{object.name}</a></div>
			{#each object.fields as field}
				<div>{field.name}</div>
			{/each}
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
