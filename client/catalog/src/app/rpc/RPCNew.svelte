<script>
export let params = null;
import client from '../../data/client.js';
import {nonEmpty} from '../../lib/util.js';
import {navStore} from '../stores.js';
import {router} from '../router.js';
import {AttributeOperation, ObjectOperation} from '../../api/object.js';

import RPCForm from './RPCForm.svelte';

navStore.set("rpc");

const defaultText = `# Run function from the api via client.rpc.[name]({args : [json_object]})

def main(args):
    # args can be any valid json object.
    return "Return json back to the client here."`

let value = ObjectOperation({
	name: AttributeOperation(""),
	code: AttributeOperation(defaultText),
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
