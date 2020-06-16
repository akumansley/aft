<script>
import client from '../../data/client.js';
import {restrictToIdent} from '../util.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import { Storage } from '../../data/enums.js';
import HLButton from '../../ui/HLButton.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';

breadcrumbStore.set(
	[{
		href: "/datatypes",
		text: "Datatypes",
	}, {
		href: "/Datatpes/new",
		text: "New",
	}]
);

const newDatatypeOp = {
	name: "",
	storedAs: 0,
	validator : {
		create : {
			name : "",
			runtime: 2,
			code: "",
			functionSignature: 0,
		}
	}
}
var cm;
var name = "code";
function setUpCM() {
	cm = getContext(name);
	cm.setValue(
`# Compile Regular Expression for valid US Phone Numbers
phone = re.Compile(r"^\\D?(\\d{3})\\D?\\D?(\\d{3})\\D?(\\d{4})$")

def validator(input):
    # Ensure input is a string
    ps = str(input)
    # Check if input matches the regex
    if phone.Match(ps):
        # If yes, return it striped of formatting
        clean = ps.replace(" ","").replace("-","")
        return clean.replace("(","").replace(")","")
    # Otherwise, raise an error
    error("Invalid phone number: %s", input)
`);
	cm.setCursor({line: 0, ch: 0});
	cm.focus();
};

import {router} from '../router.js';
async function saveDatatype() {
	newDatatypeOp.validator.create.name = newDatatypeOp.name;
	newDatatypeOp.validator.create.code = cm.getValue();
	const d = await client.datatype.create({data: newDatatypeOp});
	router.route("/datatypes");
}
</script>

<style>
	.box {
		margin: 1em 1.5em;
	}
	h1 {
		font-size: var(--scale-3);
		font-weight: 600;
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
	.v-space{
		height: .5em;
	}

</style>

<div class="box">
	<HLTextBig placeholder="Name" bind:value={newDatatypeOp.name} restrict={restrictToIdent}/>
	<h2>Validator function</h2>
	<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
	<h2>Stored as</h2>
	<HLRow>
		<HLSelect bind:value={newDatatypeOp.storedAs}>
			{#each Object.entries(Storage) as it, ix}
			<option value={ix}>
				{it[1]}
			</option>
			{/each}
		</HLSelect>
	</HLRow>
	<HLRowButton on:click={saveDatatype}>
			Save
	</HLRowButton>
</div>