<script>
	export let show = false;
</script>
<style>
	.detail {
		border-bottom: 1px solid var(--border-color);
		display: flex;
		flex-direction: row;
		background-color: var(--background-soft-highlight);
		align-items: baseline;
		padding: .5em .5em;
	}
	.hidden{
		display: none;
	}
</style>
<div class:hidden={!show} class="detail">
	<slot />
</div>
