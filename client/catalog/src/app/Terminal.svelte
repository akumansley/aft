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
	const defaultText = `def main():
	# Put code to run in here.
	return "Return what you want to see on the screen."`;

	navStore.set("terminal");

	function setUpTerminal(initialized) {
		terminal = initialized.detail;
		terminal.setOption("lineNumbers", false);
		terminal.setOption("readOnly", "nocursor");
		terminal.setOption("lint", false);
	}

	function setUpCM(initialized) {
		cm = initialized.detail;
		cm.setValue(defaultText);
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
		flex-grow: 0 1;
	}
	.wrap {
		margin-left: var(--box-margin);
		margin-right: var(--box-margin);
		margin-bottom: var(--box-margin);
		margin-top: .5em;
		display:flex;
		align-items:center;
		justify-content: space-between;
	}
	.loopback {
		flex-grow: 1;
		height: 50vh;
		background-color: var(--background-highlight);
		padding-left: var(--box-margin);
	}

	.flex {
		display: flex;
		flex-direction: column;
		height: 100vh;
	}

</style>
<div class="flex">
	<div out:saveCode></div>

	<div class="loopback">
		<CodeMirror name={"terminal"} on:initialized={setUpTerminal} />
	</div>

	<HLContent>
		<HLBorder />
		<div class="v-space"></div>
		<CodeMirror name={"code"} on:initialized={setUpCM} />
	</HLContent>

	<div class="footer">
		<div class="wrap">
			<HLButton on:click={runRepl}>
				Run
			</HLButton>
			<HLButton on:click={clearRepl}>
				Clear
			</HLButton>
		</div>
	</div>

</div>
