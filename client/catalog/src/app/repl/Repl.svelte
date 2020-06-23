<script>
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';

breadcrumbStore.set(
	[{
		href: "/repl",
		text: "Repl",
	}]
);

var repl;
var cm;

function setUpREPL() {
	repl = getContext("repl");
	repl.setOption("lineNumbers", false);
	repl.setOption("readOnly", "nocursor");
	repl.setSize(null, 200);
}

function setUpCM() {
	cm = getContext("code");
	cm.focus();
	cm.setSize(null, 400);
}

async function runRepl() {
	const d = await client.repl({input: cm.getValue().trim()});
	if(repl.getValue() == "") {
		if (d.output == "") {
			repl.setValue(">>> " + cm.getValue().trim());
		} else {
			repl.setValue(">>> " + cm.getValue().trim() + "\n" + d.output);
		}
		repl.setOption("styleActiveLine", true);
	} else {
		if (d.output == "") {
			repl.setValue(repl.getValue() + "\n>>> " + cm.getValue().trim());
		} else {
			repl.setValue(repl.getValue() + "\n>>> " + cm.getValue().trim() + "\n" + d.output);
		}
	} 
	repl.setCursor(repl.lastLine(), 0);
	cm.focus();
}
</script>
<style>
	.v-space{
		height: .5em;
	}
</style>
<HLBox>
	<HLTable>
		<h1>Repl</h1>
		<HLCodeMirror name={"repl"} on:initialized={setUpREPL}></HLCodeMirror>
		<div class="v-space"></div>
		<HLCodeMirror name={"code"} on:initialized={setUpCM}></HLCodeMirror>
		<HLRowButton on:click={runRepl}>
			Run
		</HLRowButton>
	</HLTable>
</HLBox>