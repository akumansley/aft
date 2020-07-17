<script>
export let href=null;
</script>
<style>
	.hl-row {
		padding-top: .25em;
		padding-bottom: .25em;
		min-height: 2em;
		display:flex;
		align-items:center;
		flex-direction:row;
	}
	.hl-row:hover {
		background: var(--background-highlight);
	}
	.indented {
		padding-top: .25em;
		padding-bottom: .25em;
		padding-left: .5em;
	}
	.hl-row.indented-container {
		padding-top: 0;
		padding-bottom: 0;
	}
	.iblock {
		width: 2em;
	}
</style>
<a class="hl-row" href="{href}">
	<slot />
</a>
