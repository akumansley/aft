<script>
export let params;
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import {cap, getEnumsFromObj} from '../util.js';
import { getContext } from 'svelte'
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLDetailRow from '../../ui/HLDetailRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import {router} from '../router.js';

let id = params.id;
let load = client.api.datatype.findMany({
	where: {
		OR :[
			{id: id}, 
			{name: "storedAs"}, 
			{name: "runtime"}
		]
	}, 
	include: {validator: true, enumValues: true}
});

var cm;
var name = "code";
var dt = {};
var runtime = {};
var storage = {};
load.then(obj => {
	for (var i = 0; i < obj.length; i++) {
		var name = obj[i]["name"];
		if(obj[i]["id"] == id){
			dt = obj[i];
		}
	}
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"];
	storage = results["storage"];
	breadcrumbStore.set(
		[{
			href: "/datatypes",
			text: "Datatypes",
		}, {
			href: "/datatype/" + id,
			text: cap(dt.name),
		}]
	);
});


function setUpCM() {
	cm = getContext(name);
	cm.setValue(dt.validator.code);
}

async function updateDatatype() {
	if(dt.enum == false && cm != null) {
		const parses = await client.rpc.parse({data: {data : cm.getValue()}});
		if(!parses.parsed) {
			confirm(parses.error);
			return;
		}	
	}
		
	var updateDatatypeOp = {
		name: dt.name,
		storedAs: dt.storedAs,
	}
	var d = await client.api.datatype.update({data: updateDatatypeOp, where : {id: id}});
	if(dt.enum == false) {
		var code = dt.validator.code;
		if(cm != null) {
		  code = cm.getValue();
		}
		var updateCodeOp = {
			name: dt.name,
			runtime: dt.validator.runtime,
			code: code
		}
		var c = await client.api.code.update({data: updateCodeOp, where : {id: dt.validator.id}});
	}
	router.route("/datatypes");
}

</script>

<style>
.spacer {
	width: 1em;
}
</style>

<HLBox>
	{#await load then load}
	<HLTextBig placeholder="Name" bind:value={dt.name}/>
	<HLTable>
		{#if dt.enum == false}
		{#if dt.validator.runtime == runtime["starlark"].id}
		<h2>Validator function</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		{#if dt.native == false}		
		<HLRow>
			Stored as: <span class="spacer"/>
			<HLSelect bind:value={dt.storedAs}>
				{#each Object.entries(storage) as it, ix}
				<option value={it[1]["id"]}>
					{cap(it[1]["name"])}
				</option>
				{/each}
			</HLSelect>
		</HLRow>
		{/if}
		{/if}
		{:else}
		{#each dt.enumValues as enumValue}
		<HLDetailRow name={enumValue.name} header={"id"}>
			{enumValue.id}
		</HLDetailRow>
		{/each}
		{/if}
		<HLRowButton on:click={updateDatatype}>
				Update
		</HLRowButton>
	</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLBox>
