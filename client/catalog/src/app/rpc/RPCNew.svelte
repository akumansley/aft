<script>
import client from '../../data/client.js';
import {restrictToIdent, getEnumsFromObj} from '../util.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import {router} from '../router.js';
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
let load = client.api.datatype.findMany({
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
load.then(obj => {
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"];
	fs = results["fs"];
	breadcrumbStore.set(
		[{
			href: "/rpcs",
			text: "RPCs",
		}, {
			href: "/rpcs/new",
			text: "New",
		}]
	);
});

var cm;
var name = "code";
const newRPCOp = {
	name: "",
	code : {
		create : {
			name : "",
			runtime: "",
			code: "",
			functionSignature: ""
		}
	}
}

function setUpCM() {
	cm = getContext(name);
	cm.setCursor({line: 0, ch: 0});
	cm.focus();
};

async function saveRPC() {
	const parses = await client.rpc.parse({data: {data : cm.getValue()}});	
	if(!parses.parsed) {
		confirm(parses.error);
		return;
	}
	newRPCOp.code.create.name = newRPCOp.name;
	newRPCOp.code.create.code = cm.getValue();
	newRPCOp.code.create.runtime = runtime["starlark"]["id"];
	newRPCOp.code.create.functionSignature = fs["fromJson"]["id"];
	const d = await client.api.rpc.create({data: newRPCOp});
	router.route("/rpcs");
}
</script>

<HLBox>
	{#await load then load}
	<HLTextBig placeholder="Name" bind:value={newRPCOp.name} restrict={restrictToIdent}/>
	<HLTable>
		<h2>RPC</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		<HLRowButton on:click={saveRPC}>
				Save
		</HLRowButton>
	</HLTable>
	{/await}
</HLBox>