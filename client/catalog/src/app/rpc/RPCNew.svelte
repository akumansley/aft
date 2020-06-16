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
		href: "/rpcs",
		text: "RPCs",
	}, {
		href: "/rpcs/new",
		text: "New",
	}]
);

const newRPCOp = {
	name: "",
	code : {
		create : {
			name : "",
			runtime: 2,
			code: "",
			functionSignature: 1,
		}
	}
}
var cm;
var name = "code";
function setUpCM() {
	cm = getContext(name);
	cm.setCursor({line: 0, ch: 0});
	cm.focus();
};

import {router} from '../router.js';
async function saveRPC() {
	newRPCOp.code.create.name = newRPCOp.name;
	newRPCOp.code.create.code = cm.getValue();
	const d = await client.rpc.create({data: newRPCOp});
	router.route("/rpcs");
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
	<HLTextBig placeholder="Name" bind:value={newRPCOp.name} restrict={restrictToIdent}/>
	<h2>RPC</h2>
	<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
	<HLRowButton on:click={saveRPC}>
			Save
	</HLRowButton>
</div>