<script>
	export let checked;
</script>
<style>
	.hl-check {
		margin-right:.3em;
		margin-bottom: 0;
		margin-left: 0;
		margin-top: 0;
	}
	.hl-label {
		display: flex;
		flex-direction: row;
		align-items: center;
	}

</style>
<label class="hl-label">
<input class="hl-check" type=checkbox bind:checked={checked} />

<slot/>
</label>
