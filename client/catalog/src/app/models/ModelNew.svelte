<script>
import { onMount } from 'svelte';
import client from '../../data/client.js';
import AttributeForm from './AttributeForm.svelte';
import RelationshipForm from './RelationshipForm.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLButton from '../../ui/HLButton.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLRow from '../../ui/HLRow.svelte';
import { breadcrumbStore } from '../stores.js';
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}, {
		href: "/objects/new",
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
	router.route("/object/" + data.id);

}
function restrict(s) {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal;
}
</script>

<style>
	.box {
		margin: 1em 1.5em; 
	}
	.v-space{
		height: .5em;
	}
	.v-space-2{
		height: 2em;
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
	.hform-row {
		display: flex; 
		flex-direction: row;
	}
	.col {

	}
	.spacer {
		width: 1em;
		height: 0;
	}

	.hl-row-header {
		font-size: var(--scale--2);
		text-transform: uppercase;
		font-weight: 600;
	}

</style>

<div class="box">
	<HLTextBig placeholder="Object name.." bind:value={newModelOp.name}/>
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
</div>
