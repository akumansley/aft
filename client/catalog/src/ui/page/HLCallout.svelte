<script >
	import HLPad from './HLPad.svelte'
</script>
<style>
	.callout {
		height: 100%;
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
<HLPad>
<div class="callout"><slot /></div>
</HLPad>