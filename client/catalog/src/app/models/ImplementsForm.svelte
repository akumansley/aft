<script>
	import { ConnectSelect } from '../../api/api.js';
	import {HLSmallButton} from '../../ui/form/form.js';

	export let value = [];
	function add() {
		value = value.add();
	}
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
	{#each value as v, ix}
	<ConnectSelect bind:value={v} except={value.slice(0, ix)} iface={"concreteInterface"} />
	<div class="spacer"></div>
	{/each}
	<HLSmallButton on:click={add}>+add</HLSmallButton>
</div>
