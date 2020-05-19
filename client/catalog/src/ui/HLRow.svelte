<script>
export let indent = 0;
</script>
<style>
	.hl-row {
		padding-left: .5em;
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
{#if indent == 0}
<div class="hl-row">
	<slot />
</div>
{/if}

{#if indent == 1}
<div class="hl-row indented-container">
	<div class="iblock"/>
	<div class="indented">
		<slot />
	</div>
</div>
{/if}
