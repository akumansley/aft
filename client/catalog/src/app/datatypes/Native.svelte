<script>
export let dt;
import client from '../../data/client.js';
import { router } from '../router.js';
import { checkSave } from '../save.js';
import { dirtyStore } from '../stores.js';
import { cap } from '../../lib/util.js';

import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js';
import {Box} from '../../ui/spacing/spacing.js';
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
	var d = await client.api.datatype.update({data: updateDatatypeOp, where : {id: dt.id}});
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
	<HLHeaderItem>
		<Name id="name" bind:value={dt.name} click={saveAndNav}></Name>
	</HLHeaderItem>
</HLHeader>
<HLContent>
	<Box>
		{cap(dt.name)} is implemented in native code.
	</Box>
</HLContent>