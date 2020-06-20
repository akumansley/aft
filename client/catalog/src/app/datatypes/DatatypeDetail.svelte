<script>
export let params;
import client from '../../data/client.js';
import {Runtime, Storage} from '../../data/enums.js';
import { breadcrumbStore } from '../stores.js';
import {cap, restrictToIdent} from '../util.js';
import { getContext } from 'svelte'
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import {router} from '../router.js';

let id = params.id;
let load = client.datatype.findOne({where: {id: id}, include: {validator: true}});
var cm;
var name = "code";
var dt = {
	name : "",
	id : "",
	storedAs : "",
	validator : {
		id: "",
		runtime : "",
		}

};

load.then(obj => {
	breadcrumbStore.set(
		[{
			href: "/datatypes",
			text: "Datatypes",
		}, {
			href: "/datatypes/" + id,
			text: cap(obj.name),
		}]
	);
	dt.name = obj.name;
	dt.id = obj.id;
	dt.storedAs = obj.storedAs;
	dt.validator.id = obj.validator.id;
	dt.validator.runtime = obj.validator.runtime;
	dt.validator.code = obj.validator.code;
});

function setUpCM() {
	cm = getContext(name);
	cm.setValue(dt.validator.code);
}

async function updateDatatype() {
	var updateDatatypeOp = {
		name: dt.name,
		storedAs: dt.storedAs,
	}
	var d = await client.datatype.update({data: updateDatatypeOp, where : {id: id}});
	var updateCodeOp = {
		name: dt.name,
		runtime: dt.validator.runtime,
		code: cm.getValue()
	}
	var c = await client.code.update({data: updateCodeOp, where : {id: dt.validator.id}});

	router.route("/datatypes");
}

</script>

<HLBox>
	{#await load}
		&nbsp;
	{:then}
	<HLTextBig placeholder="Name" bind:value={dt.name} restrict={restrictToIdent}/>
	<HLTable>
		<HLRow>
			<HLSelect bind:value={dt.validator.runtime}>
				{#each Object.entries(Runtime) as it, ix}
				<option value={ix}>
					{it[1]}
				</option>
				{/each}
			</HLSelect>
		</HLRow>
		{#if dt.validator.runtime == 2}
		<h2>Validator function</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		<h2>Stored as</h2>
		<HLRow>
			<HLSelect bind:value={dt.storedAs}>
				{#each Object.entries(Storage) as it, ix}
				<option value={ix}>
					{it[1]}
				</option>
				{/each}
			</HLSelect>
		</HLRow>
		{/if}
		<HLRowButton on:click={updateDatatype}>
				Update
		</HLRowButton>
	</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</HLBox>

