<script>
export let dt = null;
import client from '../../data/client.js';
import { router } from '../router.js';
import { checkSave } from '../save.js';
import { dirtyStore } from '../stores.js';
import { restrictToIdent, cap } from '../util.js';
import HLGridNew from '../../ui/grid/HLGridNew.svelte';
import HLGridEdit from '../../ui/grid/HLGridEdit.svelte';
import HLGridItem from '../../ui/grid/HLGridItem.svelte';
import HLContent from '../../ui/main/HLContent.svelte';
import HLHeader from '../../ui/main/HLHeader.svelte';
import HLButton from '../../ui/form/HLButton.svelte';

import Name from '../Name.svelte';
import Save from '../Save.svelte';


var name; var n; var ev = [];
function isNew() {
	return dt == null || n;
}
if(isNew()) {
	n = true;
	dt = {
		name: "",
		enum: true,
		enumValues: {create : []}
	}
	name = "";
} else {
	n = false;
	name = dt.name;
	for(var i = 0; i < dt.enumValues.length; i++) {
		ev[i]={"id" : dt.enumValues[i].id, "name" : dt.enumValues[i].name}
	}
}

var clean = true;
function checkClean() {
	if(isNew() && (name != dt.name || newEnumValues.length !== 0)) {
		dirtyStore.set({'clean' : false});
		clean = false;	
	} else if(!isNew() && (name != dt.name || newEnumValues.length !== 0 || !arraysEqual(ev, dt.enumValues))) {
		dirtyStore.set({'clean' : false});
		clean = false;
	} else {
		dirtyStore.set({'clean' : true});
		clean = true;
	}
}

async function saveAndNav() {
	await save();
	dirtyStore.set({'clean' : true});
	router.route("/datatypes");
}

async function save() {
	dt.name = name;
	if(isNew()) {
		dt.enumValues.create = newEnumValues;
		dt = await client.api.datatype.create({data: dt});
		//If they save in the middle of editing, this isn't a new one any more. So reroute.
		router.route("/datatype/" + dt.id);
	} else {
		//I think the below should all be one query. Question for andrew.
		var d = await client.api.datatype.update({data: {name: dt.name}, where : {id: dt.id}});
		//update any changes to old enum values
		for(let i = 0; i < dt.enumValues.length; i++) {
			var value = dt.enumValues[i];
 			await client.api.enumValue.update({data: {name: value.name}, where : {id: value.id}});
		}
		//add any new enum values
		var newVal = [];
		for(let i = 0; i < newEnumValues.length; i++) {
			newEnumValues[i].datatype = {connect: {id : dt.id}};
			newVal.push(await client.api.enumValue.create({data : newEnumValues[i]}));
		}
		//the datatype has new enum values, so add them here
		dt.enumValues = dt.enumValues.concat(newVal);
		for(var i = 0; i < dt.enumValues.length; i++) {
			ev[i]={"id" : dt.enumValues[i].id, "name" : dt.enumValues[i].name}
		}
		newEnumValues=[];
	}
}

let newEnumValues = [];
function addEnum() {
	newEnumValues = [...newEnumValues, {
		name: "",
	}];
	checkClean();
}

function arraysEqual(a, b) {
  if (a === b) return true;

  if (a == null || b == null) return false;
  if (a.length !== b.length) return false;
  for (var i = 0; i < a.length; ++i) {
    if (a[i].name !== b[i].name) return false;
  }
  return true;
}

function del() {
	console.log("delete api call goes here");
}

function removeNew(idx) {
	var out = [];
	newEnumValues.splice(idx, 1);
	//force a deep copy
	for(var i = 0; i < newEnumValues.length; i++) {
		out[i]={"id" : newEnumValues[i].id, "name" : newEnumValues[i].name}
	}
	newEnumValues = out;
	checkClean();
}

function removeOld(i) {
	console.log("delete api call goes here");
}

</script>

<style>
.spacer-small {
	width: .1em;
}
.rightAlign {
	margin-left: auto;
}
.v-space {
	height:1em;
}
</style>

<svelte:window on:keyup={checkClean} on:keydown={checkSave(save)}/>

<HLHeader>
	<Name id="name" bind:value={name} click={saveAndNav} rightAlignLast={true}>
		<div class="rightAlign">
			{#if isNew()}
				<slot />
			{:else}
				<HLButton on:click={del}>Delete</HLButton>
			{/if}
		</div>
	</Name>
</HLHeader>
<HLContent>
	{#if dt.native == true}
	{cap(dt.name)} enum comes prepackaged with Aft.
	<div class="v-space"></div>
	{/if}
	{#if ((isNew() && dt.enumValues.create.length == 0) || (!isNew() && dt.enumValues.length == 0)) && newEnumValues.length == 0}
	Enum is a set of predefined constants to be referenced in models.
	<div class="v-space"></div>
	{/if}
	{#if isNew() == false}
		{#each dt.enumValues as enumValue, ix}
			<HLGridEdit bind:value={enumValue.name} restrict={restrictToIdent} remove={removeOld} idx={ix}/>
			<div class="v-space"></div>
		{/each}
	{/if}
	{#each newEnumValues as enumValue, ix}
		<HLGridEdit bind:value={enumValue.name} restrict={restrictToIdent} remove={removeNew} idx={ix}/>
		<div class="v-space"></div>
	{/each}
	<HLGridNew click={addEnum}/>
	<Save bind:clean={clean} />
</HLContent>

