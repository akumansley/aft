<script>
	export let params = null;

	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, ConnectOperation, TypeSpecifier} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';

	import ModelForm from './ModelForm.svelte';

	navStore.set("schema");

	let model = ObjectOperation({
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
			}))
	});

	let load = client.api.model.findOne({
		where: {id: params.id},
		include: {
			attributes: {
				include: {datatype: true},
			},
			relationships: {
				case: {
					concreteRelationship: {
						include: {target: true},
					}
				}
			}
		},
	}).then(m => { 
		try {
			model.initialize(m);
			model = model;
		} catch (e) {
			console.log(e);
		}
	});

	async function saveAndNav() {
		await save();
		router.route("/schema");
	}

	async function save() {
		if (nonEmpty(model.op())) {
			return client.api.model.update(model.op().update);
		}
		return
	}

</script>
	
{#await load then _}
<ModelForm bind:value={model} on:save={saveAndNav}/>
{/await}