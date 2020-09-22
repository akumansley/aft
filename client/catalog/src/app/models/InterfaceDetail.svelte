<script>
	export let params = null;

	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, ConnectOperation, TypeSpecifier, ReadOnly} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';
	
	import InterfaceForm from './InterfaceForm.svelte';

	navStore.set("schema");

	let iface = ObjectOperation({
		name: AttributeOperation(""),
		relationships: RelationshipOperation(
			ObjectOperation({
				type: TypeSpecifier("concreteRelationship"),
				name: AttributeOperation(""),
				multi: AttributeOperation(false),
				target: ConnectOperation(),
			})),
		attributes: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
				datatype: ConnectOperation(),
			})),
	});

	let load = client.api.interface.findOne({
		where: {id: params.id},
		include: {
			attributes: {
				include: {datatype: true},
			},
			relationships: {
				case: {
					concreteRelationship: {
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
			const data = await client.api.interface.update(iface.op().update);
		}
		router.route("/schema/");
	}
</script>

<InterfaceForm bind:value={iface} on:save={saveAndNav} />
