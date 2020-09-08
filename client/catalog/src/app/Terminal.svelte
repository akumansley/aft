<script>
import client from '../data/client.js';
import { terminalStore, navStore } from './stores.js';
import { getContext } from 'svelte';

import HLContent from '../ui/page/HLContent.svelte';
import HLHeader from '../ui/page/HLHeader.svelte';
import HLRow from '../ui/list/HLRow.svelte';
import HLBorder from '../ui/page/HLBorder.svelte';
import HLButton from '../ui/form/HLButton.svelte';
import CodeMirror from './codemirror/CodeMirror.svelte';

var terminal;
var cm;

navStore.set("terminal");

function setUpTerminal() {
	terminal = getContext("terminal");
	terminal.setOption("lineNumbers", false);
	terminal.setOption("readOnly", "nocursor");
	terminal.setSize(null, 250);
	terminal.setOption("lint", false);
}

function setUpCM() {
	cm = getContext("code");
	cm.setValue(
`def main():
    #Put code to run in here.
    return "Return what you want to see on the screen."`);
	cm.setCursor(cm.lastLine(), 1000);
	terminalStore.subscribe(value => {
		if ('code' in value) {
			cm.setValue(value["code"]);
		}
		if ('history' in value) {
			cm.setHistory(value["history"]);
		}
		if ('cursor' in value) {
			cm.setCursor(value["cursor"]);
		}
		if ('terminal' in value) {
			terminal.setValue(value["terminal"]);
		}
	});
	cm.focus();
}

async function runRepl() {
	const result = await client.rpc.terminal({args: {data : cm.getValue().trim()}});
	if(terminal.getValue() == "") {
		if (result == "") {
			terminal.setValue(">>> " + cm.getValue().trim());
		} else {
			terminal.setValue(">>> " + cm.getValue().trim() + "\n" + result);
		}
		terminal.setOption("styleActiveLine", true);
	} else {
		if (result == "") {
			terminal.setValue(terminal.getValue() + "\n>>> " + cm.getValue().trim());
		} else {
			terminal.setValue(terminal.getValue() + "\n>>> " + cm.getValue().trim() + "\n" + result);
		}
	} 
	terminal.setCursor(terminal.lastLine(), 0);
}

function clearRepl() {
	terminal.setValue("");
	terminal.setOption("styleActiveLine", false);
}

function saveCode(){
  terminalStore.set({
  	"code" : cm.getValue(), 
  	"history" : cm.getHistory(),
  	"cursor"  : cm.getCursor(),
  	"terminal"    : terminal.getValue()
  });
};

</script>
<style>
	.v-space {
		height: var(--box-margin);
	}
	.footer {
  		flex: 0 1 40px;
	}
	.wrap {
		margin-left: var(--box-margin);
		margin-right: var(--box-margin);
		margin-bottom: var(--box-margin);
		margin-top: .5em;
		display:flex;
		align-items:center;
		width:100%;
		justify-content: space-between;
	}
	.wrap div:last-child {
		margin-left: auto;
	}
	.highlight {
		background-color: var(--background-highlight);
	}
	.flex {
           display: flex;
           flex-flow: column;
           height: 100vh;
	}

</style>
<div class="flex">
<div out:saveCode></div>
<HLHeader>
	<div class="highlight">
		<CodeMirror name={"terminal"} on:initialized={setUpTerminal} />
	</div>
</HLHeader>
<HLContent>
	<HLBorder />
	<div class="v-space"></div>
	<CodeMirror name={"code"} on:initialized={setUpCM} />
</HLContent>

<div class="footer">
	<HLRow>
		<div class="wrap">
			<HLButton on:click={runRepl} style="padding-left: 6em; padding-right: 6em;">
				Run
			</HLButton>
			<HLButton on:click={clearRepl} style="padding-left: 1.5em; padding-right: 1.5em;">
				Clear
			</HLButton>
		</div>
	</HLRow>
</div>

</div>
