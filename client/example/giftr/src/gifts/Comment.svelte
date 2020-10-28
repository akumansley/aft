<script>
	export let value;
</script>
<style>
	.box {
		padding: .5em .75em;
	}
	.comment-box {
		background: var(--color-1-darker);
		color: white;
	}
</style>
<div class="box comment-box">
	{value.text}
</div>
