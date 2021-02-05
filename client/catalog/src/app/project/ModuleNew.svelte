<script>
	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, SetOperation, ConnectOperation, TypeSpecifier, ReadOnly, Case} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';
	import ModuleForm from './ModuleForm.svelte';

	navStore.set("modules");

	let mod = ObjectOperation({
		name: AttributeOperation(""),
		functions: ReadOnly([]),
		interfaces: ReadOnly([]),
		datatypes: ReadOnly([]),
		roles: ReadOnly([]),
	});

	async function saveAndNav() {
		await save();
		router.route("/modules");
	}

	async function save() {
		if (nonEmpty(mod.op())) {
			return client.api.module.create(mod.op().create);
		}
		return
	}
</script>
	
<ModuleForm bind:value={mod} on:save={saveAndNav}/>
