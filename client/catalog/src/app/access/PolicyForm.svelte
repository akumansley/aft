<script>
	export let value; 

	import client from '../../data/client.js';
	import {HLButton, HLCheckbox, HLSelect} from '../../ui/form/form.js';
	import {HalfBox} from '../../ui/spacing/spacing.js';
	import {ConnectSelect} from '../../api/api.js';
	import {HLSectionTitle} from '../../ui/page/page.js';
	import CodeMirror from '../codemirror/CodeMirror.svelte';


	let showOperations = false;
	const toggleOperations = () => { showOperations = !showOperations };

	let showWhere = false;
	const toggleWhere = () => { showWhere = !showWhere };

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
	.where {
		border: 1px solid var(--border-color);
		display: flex;
		flex-direction: column;
		height: 20em;
		overflow: hidden;
	}
</style>

<div class="hform-row">
	<ConnectSelect bind:value={value['interface']} iface={"interface"} />
	<div class="spacer" />
	<HLButton on:click={toggleOperations}>Operations</HLButton>
	<div class="spacer" />
	<HLButton on:click={toggleWhere}>Where</HLButton>
</div>


{#if showOperations}
<HalfBox>
	<h2>Operations</h2>
	<HLCheckbox bind:checked={value.read}>Read</HLCheckbox>
	<HLCheckbox bind:checked={value.write}>Write</HLCheckbox>
</HalfBox>
{/if}
{#if showWhere}
<HalfBox>
	<h2>Where</h2>
	<div class="where">
		<CodeMirror bind:value={value.text}></CodeMirror>
	</div>
</HalfBox>
{/if}




