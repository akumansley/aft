<script>
import client from '../../data/client.js';
import {restrictToIdent} from '../util.js';
import { breadcrumbStore } from '../stores.js';
import { getContext } from 'svelte'
import {router} from '../router.js';
import HLBox from '../../ui/HLBox.svelte';
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';

var cm;
var name = "code";

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
			functionSignature: 2
		}
	}
}

function setUpCM() {
	cm = getContext(name);
	cm.setCursor({line: 0, ch: 0});
	cm.focus();
};

async function saveRPC() {
	newRPCOp.code.create.name = newRPCOp.name;
	newRPCOp.code.create.code = cm.getValue();
	const d = await client.rpc.create({data: newRPCOp});
	router.route("/rpcs");
}
</script>

<HLBox>
	<HLTextBig placeholder="Name" bind:value={newRPCOp.name} restrict={restrictToIdent}/>
	<HLTable>
		<h2>RPC</h2>
		<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
		<HLRowButton on:click={saveRPC}>
				Save
		</HLRowButton>
	</HLTable>
</HLBox>