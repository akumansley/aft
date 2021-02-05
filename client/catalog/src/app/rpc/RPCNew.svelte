<script>
import client from '../../data/client.js';
import {nonEmpty} from '../../lib/util.js';
import {navStore} from '../stores.js';
import {router} from '../router.js';
import {AttributeOperation, ObjectOperation, SetOperation, ConnectOperation} from '../../api/object.js';

import RPCForm from './RPCForm.svelte';

const RPC = "4b8db42e-d084-4328-a758-a76939341ffa";
navStore.set("rpcs");

const defaultText = `# Run function from the api via client.rpc.[name]({args : [json_object]})

def main(args):
    # args can be any valid json object.
    return "Return json back to the client here."`

let value = ObjectOperation({
	name: AttributeOperation(""),
	funcType: AttributeOperation(RPC),
	code: AttributeOperation(defaultText),
	role: SetOperation(),
	module: ConnectOperation(),
}); 

async function saveAndNav() {
	const op = value.op();
	if (nonEmpty(op)) {
		await client.api.starlarkFunction.create(op.create);
	}
	router.route("/rpcs");
}
</script>

<RPCForm bind:value={value} on:save={saveAndNav} />
