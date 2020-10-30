<script>
	import UpdateOp from './UpdateOp.svelte';
	import CreateOp from './CreateOp.svelte';
	import ConnectOp from './ConnectOp.svelte';
	import {HLButton} from '../../ui/form/form.js';

	export let value;
	let ops = value.Ops;
	let opMap = {
		0: "Connect",
		1: "Disconnect",
		2: "Create",
		3: "Update",
		4: "Delete",
	}
	let open = false;
	function toggle () {
		open = !open;
	}
</script>
<style>
	.op-closed {
		overflow: hidden;
		text-overflow: ellipsis;
	}
	.tx-entry {
		padding: 1em;
	}
	.optype {
		font-size: var(--scale--2);
		text-transform: uppercase;
		margin-top: .5em;
	}
	.tx-header {
		cursor: pointer;
		color: white;
	}
</style>

<div class="tx-entry">
	{#if open}

	<div on:click={toggle} class="tx-header">
		&#8722; Transaction with {ops.length} operations
	</div>

	{#each ops as opEntry}
	{#if opEntry.OpType === 3}
	<div class="optype">Update</div>
	<UpdateOp value={opEntry.Op} />
	{:else if opEntry.OpType === 2}
	<div class="optype">Create</div>
	<CreateOp value={opEntry.Op} />
	{:else if opEntry.OpType === 0}
	<div class="optype">Connect</div>
	<ConnectOp value={opEntry.Op} />
	{:else}
	{opMap[opEntry.OpType]}
	<div class="op-closed">{JSON.stringify(opEntry.Op, null, 2)}</div>
	{/if}
	{/each}
	{:else}

	<div on:click={toggle} class="tx-header">
		&#43; Transaction with {ops.length} operations
	</div>

	{/if}

</div>
