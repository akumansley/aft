<script>
export let params = null;
import aft from '../../data/aft.js';
import { navStore, dirtyStore } from '../stores.js';

import Native from './Native.svelte';
import Enum from './Enum.svelte';
import Starlark from './Starlark.svelte';
import Type from './Type.svelte';

navStore.set("datatype");
function isNew() {
	return params == null || params.id == "new";
}
let load;
if(!isNew()) {
	load = aft.api.datatype.findOne({
		where: {id: params.id}, 
		case: {
			coreDatatype: { include: {validator: true}, },
			enum: { include: { enumValues: true } }
		}
	});
}
var dt = null;

var storage = {enumValues:[]};
aft.api.enum.findOne({where: {name: "storedAs"}, include: {enumValues: true}}).then(s => storage = s);

var fs = {};
aft.api.enum.findOne({where: {name: "functionSignature"}}).then(f => fs = f);

load.then(obj => {
	if(!isNew()) {
		dt = obj;
	}
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
	{#if !isNew() && dt.type === "coreDatatype" && dt.validator.type === "nativeFunction"}
		<Native dt={dt} />
	{:else if (isNew() && type == "code") || (!isNew() && dt.type !== "enum")}
		<Starlark dt={dt} fs={fs} storage={storage} >
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
