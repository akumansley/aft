<script>
	export let params = null;

	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, SetOperation, TypeSpecifier, ReadOnly} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';
	
	import InterfaceForm from './InterfaceForm.svelte';

	navStore.set("schema");

	let iface = ObjectOperation({
		name: AttributeOperation(""),
		relationships: RelationshipOperation(
			ObjectOperation({
				type: TypeSpecifier("interfaceRelationship"),
				name: AttributeOperation(""),
				multi: AttributeOperation(false),
				target: SetOperation(),
			})),
		attributes: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
				datatype: SetOperation(),
			})),
		module: SetOperation(),
	});

	let load = client.api.interface.findOne({
		where: {id: params.id},
		include: {
			module: true,
			attributes: {
				include: {datatype: true},
			},
			relationships: {
				case: {
					interfaceRelationship: {
						include: {target: true},
					},
				},
			},
		},
	}).then(i => { 
		try {
			iface.initialize(i);
			iface = iface;
		} catch (e) {
			console.log(e);
		}
	});

	async function saveAndNav() {
		if (nonEmpty(iface.op())) {
			const data = await client.api.concreteInterface.update(iface.op().update);
		}
		router.route("/schema/");
	}
</script>

<InterfaceForm bind:value={iface} on:save={saveAndNav} />
