<script>
export let params;
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import {cap, restrictToIdent, getEnumsFromObj} from '../util.js';
import { getContext } from 'svelte';
import {router} from '../router.js';
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';

let id = params.id;
let load = client.api.rpc.findOne({where: {id: id}, include: {code: true}});
let enms = client.api.datatype.findMany({
	where: {
		OR :[
			{name: "runtime"},
			{name: "functionSignature"}
		]
	}, 
	include: {enumValues: true}
});

var runtime = {};
var fs = {};
enms.then(obj => {
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"];
	fs = results["fs"];
});

var cm;
var name = "code";
var rpc = {
	name : "",
	id : "",
	code : {
		id : "",
		code : "",
		runtime : "",
		functionSignature : ""
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
	rpc = obj;
	if (cm != null) {
		cm.setValue(rpc.code.code);
	}
});

function setUpCM() {
	cm = getContext(name);
	cm.setValue(rpc.code.code);
}

async function updateRPC() {
	const parses = await client.rpc.parse({data: {data : cm.getValue()}});	
	if(!parses.parsed) {
		confirm(parses.error);
		return;
	}
	rpc.code.runtime = runtime["starlark"]["id"];
	rpc.code.functionSignature = fs["fromJson"]["id"];
	var updateRPCOp = {
		name: rpc.name
	}
	var d = client.api.rpc.update({data: updateRPCOp, where : {id: id}});
	var updateCodeOp = {
		name: rpc.name,
		code: cm.getValue()
	}
	var c = await client.api.code.update({data: updateCodeOp, where : {id: rpc.code.id}});

	router.route("/rpcs");
}
</script>

<HLBox>
	{#await load then load}
	<HLTextBig placeholder="Name" bind:value={rpc.name} restrict={restrictToIdent}/>
	<HLTable>
		<h2>RPC</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
			{#await enms then enms}
			<HLRowButton on:click={updateRPC}>
					Update
			</HLRowButton>
			{/await}
	</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLBox>
