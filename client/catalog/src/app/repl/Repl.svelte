<script>
import client from '../../data/client.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLButton from '../../ui/HLButton.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';
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
}

function setUpCM() {
	cm = getContext("code");
	cm.focus();
}

async function runRepl() {
	const d = await client.repl({input: cm.getValue().trim()});
	if(repl.getValue() == "") {
		repl.setValue(">>> " + cm.getValue().trim() + "\n" + d.output);
		repl.setOption("styleActiveLine", true);
	} else {
		repl.setValue(repl.getValue() + "\n>>> " + cm.getValue().trim() + "\n" + d.output);
	} 
	repl.setCursor(repl.lastLine(), 0);
	cm.focus();
}
</script>

<style>
	.box {
		margin: 1em 1.5em;
	}
	.spacer {
		height: .5em;
	}
	h1 {
		font-size: var(--scale-3);
		font-weight: 600;
	}
	
</style>

<div class="box">
	<h1>Repl</h1>
	<HLCodeMirror name={"repl"} on:initialized={setUpREPL}></HLCodeMirror>
	<div class="spacer"></div>
	<HLCodeMirror name={"code"} on:initialized={setUpCM}></HLCodeMirror>
	<HLRowButton on:click={runRepl}>
		Run
	</HLRowButton>
</div>