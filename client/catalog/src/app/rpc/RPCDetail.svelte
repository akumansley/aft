<script>
	export let params = null;
	import client from '../../data/client.js';
	import {nonEmpty} from '../../lib/util.js';
	import { navStore } from '../stores.js';
	import {router} from '../router.js';
	import {AttributeOperation, ObjectOperation, TypeSpecifier} from '../../api/object.js';

	import RPCForm from './RPCForm.svelte';
	import NativeRPC from './NativeRPC.svelte';

	navStore.set("rpc");

	let starlarkFunction = ObjectOperation({
		name: AttributeOperation(""),
		type: TypeSpecifier("starlarkFunction"),
		code: AttributeOperation(""),
	}); 
	let value;

	let load = client.api.rpc.findOne({where: {id: params.id}, include: {function: true}}).then((v) => {
		value = v['function'];
		if (value.type === "starlarkFunction") {
			starlarkFunction.initialize(value);
			starlarkFunction = starlarkFunction;

		}
	});

	async function saveAndNav() {
		if (value.type === "starlarkFunction") {
			const op = starlarkFunction.op();
			if (nonEmpty(op)) {
				await client.api.starlarkFunction.update(op.update);
			}
		}
		router.route("/rpcs");
	}
</script>

{#await load then _}
{#if value.type === "starlarkFunction"}
<RPCForm bind:value={starlarkFunction} on:save={saveAndNav} />
{:else if value.type === "nativeFunction"}
<NativeRPC value={value} on:save={saveAndNav} />
{/if}

{/await}