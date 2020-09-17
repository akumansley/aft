<script>
	import { onMount } from 'svelte';
	import HLText from '../../ui/form/HLText.svelte';
	import { ConnectSelect } from '../../api/api.js';
	import {restrictToIdent } from '../../lib/util.js';

	export let value = {};
</script>

<style>
	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
	.spacer {
		width: 1em;
		height: 0;
	}
</style>

<div class="hform-row">
	<HLText placeholder="Attribute name.." bind:value={value.name} restrict={restrictToIdent} />
	<div class="spacer" />
	<ConnectSelect bind:value={value.datatype} iface={"datatype"} />
</div>
