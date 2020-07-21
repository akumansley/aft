<script>
export let params = null;
import { onMount } from 'svelte';
import { navStore } from '../stores.js';
import {router} from '../router.js';

import client from '../../data/client.js';
import AttributeForm from './AttributeForm.svelte';
import RelationshipForm from './RelationshipForm.svelte';

import HLRowButton from '../../ui/list/HLRowButton.svelte';
import HLButton from '../../ui/form/HLButton.svelte';
import HLRow from '../../ui/list/HLRow.svelte';
import HLHeader from '../../ui/main/HLHeader.svelte';
import HLContent from '../../ui/main/HLContent.svelte';
import Name from '../Name.svelte';

navStore.set("model");
function isNew() {
	return params == null || params.id == "new";
}

let models=[]; let model;
var attributes = [];
var relationships = [];
let load = client.api.model.findMany({
	include: {
		relationships: true, 
		attributes: true
	}
});
load.then((ms) => {
	models = ms;
	if(isNew()) {
		model = {
			name: "",
			attributes: {create: []},
			relationships: {create: []},
		}	
	} else {
		for(let i = 0; i < models.length; i++) {
			if(params.id === models[i].id) {
				model = models[i];
				attributes = model.attributes;
				leftRelationships = model.relationships;
				return
			}
		}	
	}
});

function addAttribute() {
	attributes = [...attributes, {
		name: "",
		datatypeId: "",
		datatype: { connect: {id: ""}},
	}];
}

function addRelationship() {
	leftRelationships = [...leftRelationships, {
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

async function saveAndNav() {
	await save();
	router.route("/models");
}

async function save() {
	if(isNew()) {
		model.attributes.create = attributes;
		model.leftRelationships.create = leftRelationships;
		const data = await client.api.model.create({data: model});
	} else {
		var updateModelOp = {
			name: model.name
		}
		await client.api.model.update({data: updateModelOp, where : {id: model.id}});
		for(var i = 0; i < attributes.length; i++) {
			var updateAttributeOp = {
				name: attributes[i].name,
			}
			await client.api.attribute.update({data: updateAttributeOp, where : {id: attributes[i].id}});			
		}
		for(var i = 0; i < leftRelationships.length; i++) {
			var updateRelationshipOp = {
				leftName: leftRelationships[i].leftName,
				rightName: leftRelationships[i].rightName,
				leftBinding: leftRelationships[i].leftBinding,
				rightBinding: leftRelationships[i].rightBinding,
			}
			await client.api.relationship.update({data: updateRelationshipOp, where : {id: leftRelationships[i].id}});			
		}
	}
}

function del() {
	console.log("delete goes here");
}
</script>

<style>
.rightAlign {
	margin-left: auto;
}
.v-space {
	height: .5em;
}
</style>

{#await load then load}
<HLHeader>
	<Name id="name" placeholder="Model name.." bind:value={model.name} click={saveAndNav} rightAlignLast={true}>
		<div class="rightAlign">
			<HLButton on:click={del}>Delete</HLButton>
		</div>
	</Name>
</HLHeader>
<HLContent>
	<h2>Attributes</h2>
	{#each attributes as attr}
		<AttributeForm bind:attribute={attr}/>
	{/each}
	<div class="v-space"/>
	<HLRowButton on:click={addAttribute}>+add</HLRowButton>

	<h2>Relationships</h2>
	{#each leftRelationships as rel}
		<RelationshipForm modelName={model.name} bind:relationship={rel} models={models}/>
	{/each}
	<div class="v-space"/>
	<HLRowButton on:click={addRelationship}>+add</HLRowButton>
</HLContent>
{/await}
