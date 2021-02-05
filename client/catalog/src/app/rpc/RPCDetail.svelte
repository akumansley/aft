<script>
	export let params = null;
	import client from '../../data/client.js';
	import {nonEmpty} from '../../lib/util.js';
	import { navStore } from '../stores.js';
	import {router} from '../router.js';
	import {AttributeOperation, ObjectOperation, TypeSpecifier, SetOperation} from '../../api/object.js';

	import RPCForm from './RPCForm.svelte';
	import NativeRPC from './NativeRPC.svelte';

	const RPC = "4b8db42e-d084-4328-a758-a76939341ffa";
	navStore.set("rpcs");

	let starlarkFunction = ObjectOperation({
		name: AttributeOperation(""),
		type: TypeSpecifier("starlarkFunction"),
		code: AttributeOperation(""),
		funcType: AttributeOperation(RPC),
		role: SetOperation(),
		module: SetOperation(),
	}); 
	let value;

	let load = client.api.function.findOne({
		where: {id: params.id},
		include: {
			role: true,
			module: true
		}
	}).then((v) => {
		value = v;
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