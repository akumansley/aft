<script>
export let params = null;
import client from '../../data/client.js';
import { navStore, dirtyStore } from '../stores.js';
import { getEnumsFromObj } from '../util.js';

import Native from './Native.svelte';
import Enum from './Enum.svelte';
import Starlark from './Starlark.svelte';
import Type from './Type.svelte';

navStore.set("datatype");
function isNew() {
	return params == null || params.id == "new";
}
let load;
if(isNew()) {
	load = client.api.datatype.findMany({
		where: {
			OR :[
				{name: "storedAs"}, 
				{name: "runtime"},
				{name: "functionSignature"}
			]
		}, 
	});
} else {
	load = client.api.datatype.findOne({
		where: {id: params.id}, 
	});
}
var dt = null;
var runtime = {};
var storage = {};
var fs = {};
load.then(obj => {
	if(!isNew()) {
		for (var i = 0; i < obj.length; i++) {
			if(obj[i]["id"] == params.id){
				dt = obj[i];
			}
		}	
	}
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"];
	storage = results["storage"];
	fs = results["fs"];
});

const types = ["code", "enum"];
let type = types[0];

function select(e) {
	type = e;
	dirtyStore.set({'clean' : true});
}
</script>

<style></style>
{#await load then load}
	{#if !isNew() && dt.native == true && dt.enum == false}
		<Native dt={dt} />
	{:else if (isNew() && type == "code") || (!isNew() && dt.enum == false)}
	<Starlark dt={dt} storage={storage} runtime={runtime["starlark"]["id"]} fs={fs["fromJson"]["id"]} >
		{#if isNew()}
			<Type types={types} type={type} change={select}/>
		{/if}
	</Starlark>
	{:else}
		<Enum dt={dt}>
			{#if isNew()}
				<Type types={types} type={type} change={select}/>
			{/if}
		</Enum>
	{/if}
{/await}
