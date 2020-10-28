<script>
	export let value;
</script>
<style>
	.box {
		padding: .5em .75em;
	}
	.gift-box {
		background: var(--color-1-darker);
		border-radius: 5px 5px 0 5px;
		color: white;
	}
</style>
<div class="box gift-box">
	<div class="gift">
		{value.description}
	</div>
</div>
