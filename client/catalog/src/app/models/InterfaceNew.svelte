<script>
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
		module: ConnectOperation(),
	});

	async function saveAndNav() {
		if (nonEmpty(iface.op())) {
			const data = await client.api.concreteInterface.create(iface.op().create);
			router.route("/interface/" + data.id);
		}
	}
</script>

<InterfaceForm bind:value={iface} on:save={saveAndNav} />
