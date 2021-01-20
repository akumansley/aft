<script>
	import UpdateOp from './UpdateOp.svelte';
	import CreateOp from './CreateOp.svelte';
	import ConnectOp from './ConnectOp.svelte';
	import {HLButton} from '../../ui/form/form.js';

	export let value;
	let ops = value;
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
	{#if opEntry.opType === "update"}
	<div class="optype">Update</div>
	<UpdateOp value={opEntry} />
	{:else if opEntry.opType === "create"}
	<div class="optype">Create</div>
	<CreateOp value={opEntry} />
	{:else if opEntry.opType === "connect"}
	<div class="optype">Connect</div>
	<ConnectOp value={opEntry} />
	{:else}
	{opMap[opEntry.opType]}
	<div class="op-closed">{JSON.stringify(opEntry.Op, null, 2)}</div>
	{/if}
	{/each}
	{:else}

	<div on:click={toggle} class="tx-header">
		&#43; Transaction with {ops.length} operations
	</div>

	{/if}

</div>
