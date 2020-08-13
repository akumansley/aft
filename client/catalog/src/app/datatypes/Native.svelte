<script>
export let dt;
import aft from '../../data/aft.js';
import { router } from '../router.js';
import { checkSave } from '../save.js';
import { dirtyStore } from '../stores.js';
import { cap } from '../util.js';

import HLHeader from '../../ui/main/HLHeader.svelte';
import HLContent from '../../ui/main/HLContent.svelte';
import Name from '../Name.svelte';
import Save from '../Save.svelte';

var name = dt.name;
var clean = true;
function checkClean() {
	if(name != dt.name) {
		dirtyStore.set({'clean' : false});
		clean = false;
	} else {
		dirtyStore.set({'clean' : true});
		clean = true;
	}
}

async function saveAndNav() {
	await save();
	router.route("/datatypes");
}

async function save() {
	var updateDatatypeOp = {
		name: dt.name
	}
	var d = await aft.api.datatype.update({data: updateDatatypeOp, where : {id: dt.id}});
	name = dt.name;
}

</script>

<style>
.spacer-small {
	width: .1em;
}
</style>

<svelte:window on:keyup={checkClean} on:keydown={checkSave(save)}/>

<HLHeader>
	<Name id="name" bind:value={dt.name} click={saveAndNav}></Name>
</HLHeader>
<HLContent>
	<Save bind:clean={clean} />
	{cap(dt.name)} comes prepackaged with Aft.
</HLContent>