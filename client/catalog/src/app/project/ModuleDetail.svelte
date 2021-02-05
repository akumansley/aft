<script>
	export let params = null;

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

	let load = client.api.module.findOne({
		where: {id: params.id},
		include: {
			functions: {select: {name: true, id: true}},
			interfaces: {select: {name: true, id: true}},
			datatypes: {select: {name: true, id:true}},
			roles: {select: {name: true, id:true}},
		},
	}).then(m => { 
		try {
			mod.initialize(m);
			mod = mod;
		} catch (e) {
			console.log(e);
		}
	});

	async function saveAndNav() {
		await save();
		router.route("/modules");
	}

	async function save() {
		if (nonEmpty(mod.op())) {
			return client.api.module.update(mod.op().update);
		}
		return
	}
</script>
	
{#await load then _}
<ModuleForm bind:value={mod} on:save={saveAndNav}/>
{/await}