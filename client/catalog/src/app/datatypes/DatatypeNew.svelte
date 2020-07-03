<script>
import client from '../../data/client.js';
import {restrictToIdent} from '../util.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import { cap, getEnumsFromObj } from '../util.js';
import {router} from '../router.js';
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';

let load = client.api.datatype.findMany({
	where: {
		OR :[
			{name: "storedAs"}, 
			{name: "runtime"},
			{name: "functionSignature"}
		]
	}, 
	include: {enumValues: true}
});
var cm;
var name = "code";
const newDatatypeOp = {
	name: "",
	storedAs: "",
	validator : {
		create : {
			name : "",
			runtime: "",
			code: "",
			functionSignature: ""
		}
	}
}
var runtime = {};
var storage = {};
var fs = {};
load.then(obj => {
	var results = getEnumsFromObj(obj);
	runtime = results["runtime"];
	storage = results["storage"];
	fs = results["fs"];
	breadcrumbStore.set(
		[{
			href: "/datatypes",
			text: "Datatypes",
		}, {
			href: "/datatypes/new",
			text: "New",
		}]
	);
});

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

async function saveDatatype() {
	const parses = await client.rpc.parse({data: {data : cm.getValue()}});
	if(!parses.parsed) {
		confirm(parses.error);
		return;
	}
	newDatatypeOp.validator.create.name = newDatatypeOp.name;
	newDatatypeOp.validator.create.code = cm.getValue();
	newDatatypeOp.validator.create.runtime = runtime["starlark"]["id"];
	newDatatypeOp.validator.create.functionSignature = fs["fromJson"]["id"];
	const d = await client.api.datatype.create({data: newDatatypeOp});
	router.route("/datatypes");
}
</script>

<style>
.spacer {
	width: 1em;
}
</style>

<HLBox>
	{#await load then load}
	<HLTextBig placeholder="Name" bind:value={newDatatypeOp.name} restrict={restrictToIdent}/>
	<HLTable>
		<h2>Validator Function</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		<HLRow>
			Stored As: <span class="spacer"/>
			<HLSelect bind:value={newDatatypeOp.storedAs}>
				{#each Object.entries(storage) as it, ix}
				<option value={it[1]["id"]}>
					{cap(it[1]["name"])}
				</option>
				{/each}
			</HLSelect>
		</HLRow>
		<HLRowButton on:click={saveDatatype}>
				Save
		</HLRowButton>
	</HLTable>
	{/await}
</HLBox>