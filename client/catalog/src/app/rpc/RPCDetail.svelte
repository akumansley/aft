<script>
export let params = null;
import client from '../../data/client.js';
import { navStore, dirtyStore } from '../stores.js';
import { getContext,setContext } from 'svelte';
import {router} from '../router.js';
import { checkSave } from '../save.js';

import HLContent from '../../ui/main/HLContent.svelte';
import HLHeader from '../../ui/main/HLHeader.svelte';
import HLButton from '../../ui/form/HLButton.svelte';

import Save from '../Save.svelte';
import Name from '../Name.svelte';
import CodeMirror from '../codemirror/CodeMirror.svelte';

navStore.set("rpc");

let name; let rpc;

function isNew() { return params == null; }
async function load() {
	if(isNew()) {
		rpc = {
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
		name = "";
	} else {
		rpc = await client.api.rpc.findOne({where: {id: params.id}, include: {code: true}});
		if (cm != null) {
			cm.setValue(rpc.code.code);
		}
		name = rpc.name;
	}
}
load();

var cm; var cmName = "code";
function setUp() {
	cm = getContext(cmName);
	if(isNew()) {
		cm.setValue(`#Run function from the api via client.rpc.[name]({args : [json_object]})

def main(args):
    #args can be any valid json object.
    return "Return json back to the client here."`);
	} else if(rpc != null) {
		cm.setValue(rpc.code.code);
	}
	cm.setCursor(0,0);
	cm.focus();
}

let enms = client.api.datatype.findMany({
	where: {
		OR :[
			{name: "runtime"},
			{name: "functionSignature"}
		]
	}, 
	include: {enumValues: true}
});

var runtime; var fs;
enms.then(obj => {
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"]["starlark"]["id"];
	fs = results["fs"]["fromJson"]["id"];
});

async function saveAndNav() {
	var p = await cm.parses();
	if(!p) {
		return;
	}
	
	await save();
	dirtyStore.set({'clean' : true});
	router.route("/functions");
}

async function save() {
		rpc.name = name;
	if(isNew()) {
		rpc.code.create.name = rpc.name;
		rpc.code.create.code = cm.getValue();
		rpc.code.create.runtime = runtime;
		rpc.code.create.functionSignature = fs;
		await client.api.rpc.create({data: func});
	} else {
		await client.api.rpc.update({data: {name: rpc.name }, where : {id: rpc.id}});
		var updateCodeOp = {
			name: rpc.name,
			code: cm.getValue()
		}
		await client.api.code.update({data: updateCodeOp, where : {id: rpc.code.id}});	
	} 
	cm.setClean();
}

var clean = true;
function checkClean() {
	if(name != rpc.name || !cm.isClean()) {
		dirtyStore.set({'clean' : false});
		clean=false;
	} else {
		dirtyStore.set({'clean' : true});
		clean=true;	
	}
}
function del() {
	console.log("Delete api calls go here");
}
</script>

<style>
.rightAlign {
	margin-left: auto;
}
</style>

<svelte:window on:keyup={checkClean} on:keydown={checkSave(save)}/>
{#await rpc then rpc}
<HLHeader>
	<Name id="name" bind:value={name} click={saveAndNav} rightAlignLast={!isNew()}>
	{#if !isNew()}
		<div class="rightAlign">
			<HLButton on:click={del}>Delete</HLButton>
		</div>
	{/if}
	</Name>
</HLHeader>
<HLContent>
	<CodeMirror name={cmName} on:initialized={setUp}/>
	<Save bind:clean={clean} />
</HLContent>
{/await}
