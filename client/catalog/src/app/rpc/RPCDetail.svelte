<script>
export let params;
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import {cap, restrictToIdent} from '../util.js';
import { getContext } from 'svelte';
import {router} from '../router.js';
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';

let id = params.id;
let load = client.rpc.findOne({where: {id: id}, include: {code: true}});
var cm;
var name = "code";
var rpc = {
	name : "",
	id : "",
	code : {
		id : "",
		code : "",
		runtime : 2,
		functionSignature : 2
	}
}

load.then(obj => {
	breadcrumbStore.set(
		[{
			href: "/rpcs",
			text: "RPCs",
		}, {
			href: "/rpc/" + id,
			text: cap(obj.name),
		}]
	);
	rpc.name = obj.name;
	rpc.id = obj.id;
	rpc.code.id = obj.code.id;
	rpc.code.code = obj.code.code;
	if (cm != null) {
		cm.setValue(rpc.code.code);
	}
});

function setUpCM() {
	cm = getContext(name);
	cm.setValue(rpc.code.code);
}

async function updateRPC() {
	var updateRPCOp = {
		name: rpc.name
	}
	var d = await client.rpc.update({data: updateRPCOp, where : {id: id}});
	var updateCodeOp = {
		name: rpc.name,
		code: cm.getValue()
	}
	var c = await client.code.update({data: updateCodeOp, where : {id: rpc.code.id}});

	router.route("/rpcs");
}
</script>

<HLBox>
	{#await load}
		&nbsp;
	{:then}
	<HLTextBig placeholder="Name" bind:value={rpc.name} restrict={restrictToIdent}/>
	<HLTable>
		<h2>RPC</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		<HLRowButton on:click={updateRPC}>
				Update
		</HLRowButton>
	</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLBox>
