<script>
	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, ConnectOperation, TypeSpecifier, ReadOnly} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';
	
	import ModelForm from './ModelForm.svelte';

	navStore.set("schema");

	let model = ObjectOperation({
		name: AttributeOperation(""),
		implements: RelationshipOperation(ConnectOperation()),
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
		targeted: ReadOnly([]),
		module: ConnectOperation(),
	});

	async function saveAndNav() {
		if (nonEmpty(model.op())) {
			let createOp = model.op().create;
			const data = await client.api.model.create(createOp);
			router.route("/model/" + data.id);
		}
	}

</script>

<ModelForm bind:value={model} on:save={saveAndNav} />
