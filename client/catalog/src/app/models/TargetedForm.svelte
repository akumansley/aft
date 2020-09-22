<script>
	export let value;
	import {HLSmallButton} from '../../ui/form/form.js';
	import Box from '../../ui/spacing/Box.svelte';
</script>

<style>
	.spacer {
		width: 1em;
	}

	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
</style>

<div class="hform-row">
	{value.source[0].name}.{value.name}

	<div class="spacer"></div>

	<HLSmallButton on:click>+ add reverse</HLSmallButton>
</div>
