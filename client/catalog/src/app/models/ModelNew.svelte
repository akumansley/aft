<script>
import { onMount } from 'svelte';
import client from '../../data/client.js';
import AttributeForm from './AttributeForm.svelte';
import RelationshipForm from './RelationshipForm.svelte';
import HLBox from '../../ui/HLBox.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLButton from '../../ui/HLButton.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLRow from '../../ui/HLRow.svelte';
import { breadcrumbStore } from '../stores.js';
breadcrumbStore.set(
	[{
		href: "/models",
		text: "Models",
	}, {
		href: "/models/new",
		text: "New",
	}]
);
let models=[];
client.model.findMany({}).then((ms) => {
	models = ms;
});
const newModelOp = {
	name: "",
	attributes: {create: []},
	leftRelationships: {create: []},
}

function addAttribute() {
	newModelOp.attributes.create = [...newModelOp.attributes.create, {
		name: "",
		datatypeId: "",
		datatype: { connect: {id: ""}},
	}];
}

function addRelationship() {
	newModelOp.leftRelationships.create = [...newModelOp.leftRelationships.create, {
		leftName: "",
		leftBinding: 0,
		rightName: "",
		rightBinding: 0,
		rightModel: {
			connect: {
				id: "",
			}
		},
	}];
}

import {router} from '../router.js';
async function saveModel() {
	const data = await client.model.create({data: newModelOp});
	router.route("/model/" + data.id);

}
</script>

<HLBox>
	<HLTextBig placeholder="Model name.." bind:value={newModelOp.name}/>
	<h2>Attributes</h2>
	<HLTable>
		{#each newModelOp.attributes.create as attr}
			<AttributeForm bind:attribute={attr}/>
		{/each}
		<div class="v-space"/>
		<HLRowButton on:click={addAttribute}>+add</HLRowButton>
	</HLTable>
	<h2>Relationships</h2>
	<HLTable>
		{#each newModelOp.leftRelationships.create as rel}
			<RelationshipForm modelName={newModelOp.name} bind:relationship={rel} models={models}/>
		{/each}
		<div class="v-space"/>
		<HLRowButton on:click={addRelationship}>+add</HLRowButton>
	</HLTable>

		<div class="v-space-2"></div>
	<HLTable>
		<HLRowButton on:click={saveModel}>
				Save
		</HLRowButton>
	</HLTable>
</HLBox>
