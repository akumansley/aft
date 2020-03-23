<script>
	import { onMount } from 'svelte';

	let objects = [];

	onMount(async () => {
		const res = await fetch(`http://localhost:8080/objects`);
		objects = await res.json();
	});
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
	{#each objects as object}
		<div>{object.name}</div>
		{#each object.fields as field}
			<div>{field.name}</div>
		{/each}
	{/each}
</div>
