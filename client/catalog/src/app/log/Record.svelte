<script>
	export let value;
	let keys = Object.keys(value);
</script>

<style>
	.key {
		font-weight: 600;

	}
	.value {

	}
	.record {
		display: flex;
		flex-direction: row;
		font-size: var(--scale--1);
	}
	.entry {
		display: flex;
		flex-direction: column;
		border-left: 1px solid var(--border-color);
		border-top: 1px solid var(--border-color);
		border-bottom: 1px solid var(--border-color);

		padding: .25em .5em;
	}
	.entry:last-child {
		border-right: 1px solid var(--border-color);
	}
</style>

<div class="record">
	{#each keys as key}
	<div class="entry">
		<div class="key">{key}</div>
		<div class="value">{value[key]}</div>
	</div>
	{/each}
</div>
