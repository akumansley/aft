<script>
	export let value;
</script>
<style>

	.connect {
		font-size: var(--scale--1);
		display: flex;
		flex-direction: row;
	}
	.box {
		border: 1px solid var(--border-color);
		padding: .25em .5em;
	}
</style>

<div class="connect">
	<div class="box">
	{value.Left} &rarr; {value.Right}
	</div>
</div>

