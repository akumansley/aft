<script>
export let dt = null;
export let fs = null;
export let storage;

import client from '../../data/client.js';
import { router } from '../router.js';
import { checkSave } from '../save.js';

import { dirtyStore } from '../stores.js';
import { getContext } from 'svelte';

import Name from '../Name.svelte';
import CodeMirror from '../codemirror/CodeMirror.svelte';
import Save from '../Save.svelte';
import Storage from './Storage.svelte';
import HLButton from '../../ui/form/HLButton.svelte';
import HLContent from '../../ui/main/HLContent.svelte';
import HLHeader from '../../ui/main/HLHeader.svelte';

var n;
function isNew() {
	return dt == null || n;
}

if(isNew()) {
	n = true;
	dt = {
		name: "",
		storedAs: storage["string"]["id"],
		validator : {
			create : {
				name : "",
				code: "",
				functionSignature: fs,
			}
		}
	}
} else {
	n = false;
}

var cm; var cmName = "code";
function setUp() {
	cm = getContext(cmName);
	if(isNew()) {
		cm.setValue(`#Code datatypes are used by models.
#Stored As defines how the data is represented in the database.

def main(arg):
    #arg is input to be validated.
    #main should fail() if the input isn't valid.
    return "Return the input unchanged or modified as necessary."`);
	} else {
		cm.setValue(dt.validator.code);	
	}
	cm.setCursor(0,0);
	cm.focus();
}

var name = dt.name; var storedAs = dt.storedAs;
var clean = true;
function checkClean() {
	if(name != dt.name || storedAs != dt.storedAs || !cm.isClean()) {
		dirtyStore.set({'clean' : false});
		clean = false;
	} else {
		dirtyStore.set({'clean' : true});	
		clean = true;
	}
}

async function saveAndNav() {
	var p = await cm.parses();
	if(!p) {
		return;
	}
	await save();
	dirtyStore.set({'clean' : true});
	router.route("/datatypes");
}

async function save() {
	if(isNew()) {
		dt.validator.create.name = dt.name;
		dt.validator.create.code = cm.getValue();
		const d = await client.api.datatype.create({data: dt});
	} else {
		var updateDatatypeOp = {
			name: dt.name,
			storedAs: dt.storedAs,
		}
		var d = await client.api.datatype.update({data: updateDatatypeOp, where : {id: dt.id}});	
		var updateCodeOp = {
			name: dt.name,
			code: cm.getValue()
		}
		var c = await client.api.code.update({data: updateCodeOp, where : {id: dt.validator.id}});
	}
	cm.setClean();
	name = dt.name;
	storedAs = dt.storedAs;
}

function del() {
	console.log("delete goes here");
}
</script>

<style>
.spacer-big {
	width: 2em;
}
.rightAlign {
	margin-left: auto;
}
</style>

<svelte:window on:keyup={checkClean} on:keydown={checkSave(save)} />

<HLHeader>
	<Name id="name" bind:value={dt.name} click={saveAndNav} rightAlignLast={true}>
		<span class="spacer-big"/>
		<Storage bind:storedAs={dt.storedAs} storage={storage} change={checkClean}/>
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
	<CodeMirror name={cmName} on:initialized={setUp}/>
	<Save clean={clean}/>
</HLContent>
