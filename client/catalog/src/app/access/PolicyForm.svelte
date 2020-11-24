<script>
	export let value; 

	import client from '../../data/client.js';
	import {HLButton, HLCheckbox, HLSelect} from '../../ui/form/form.js';
	import {HalfBox} from '../../ui/spacing/spacing.js';
	import {ConnectSelect} from '../../api/api.js';
	import {HLSectionTitle} from '../../ui/page/page.js';
	import CodeMirror from '../codemirror/CodeMirror.svelte';


	let showDetail = false;
	const toggleDetail = () => { showDetail = !showDetail };

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
		height: 10em;
		overflow: hidden;
	}
	.ops {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		grid-column-gap: .25em;
	}
	.opBox {
	}
</style>

<div class="hform-row">
	<ConnectSelect bind:value={value['interface']} iface={"interface"} />
	<div class="spacer" />
	<HLButton on:click={toggleDetail}>Detail</HLButton>
</div>


{#if showDetail}
<HalfBox>
	<h2>Operations</h2>
	<div class="ops">

		<div class="opBox">
			<HLCheckbox bind:checked={value.allowRead}>Read</HLCheckbox>
			{#if value.allowRead}
			<h2>Where</h2>
			<div class="where">
				<CodeMirror bind:value={value.readWhere}></CodeMirror>
			</div>
			{/if}
		</div>

		<div class="opBox">
			<HLCheckbox bind:checked={value.allowCreate}>Create</HLCheckbox>
			{#if value.allowCreate}
			<h2>Where</h2>
			<div class="where">
				<CodeMirror bind:value={value.createWhere}></CodeMirror>
			</div>
			{/if}
		</div>

		<div class="opBox">
			<HLCheckbox bind:checked={value.allowUpdate}>Update</HLCheckbox>
			{#if value.allowUpdate}
			<h2>Where</h2>
			<div class="where">
				<CodeMirror bind:value={value.updateWhere}></CodeMirror>
			</div>
			{/if}
		</div>
	</div>
</HalfBox>
{/if}
