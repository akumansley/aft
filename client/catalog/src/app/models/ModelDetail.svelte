<script>
	export let params = null;

	import client from '../../data/client.js';
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import {ObjectOperation, RelationshipOperation, AttributeOperation, SetOperation, ConnectOperation, TypeSpecifier, ReadOnly, Case} from '../../api/object.js';
	import {nonEmpty} from '../../lib/util.js';

	import ModelForm from './ModelForm.svelte';

	navStore.set("schema");

	let model = ObjectOperation({
		name: AttributeOperation(""),
		implements: RelationshipOperation(ConnectOperation()),
		relationships: RelationshipOperation(
			Case({
				concreteRelationship: ObjectOperation({
					name: AttributeOperation(""),
					type: TypeSpecifier("concreteRelationship"),
					multi: AttributeOperation(false),
					target: SetOperation(),
				}),
				reverseRelationship: ObjectOperation({
					name: AttributeOperation(""),
					type: TypeSpecifier("reverseRelationship"),
					multi: ReadOnly(true),
					referencing: SetOperation(),
				}),
			}),
		),
		attributes: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
				datatype: SetOperation(),
			})),
		targeted: ReadOnly([]),
		module: SetOperation(),
	});

	let load = client.api.model.findOne({
		where: {id: params.id},
		include: {
			attributes: {
				include: {datatype: true},
			},
			implements: true,
			relationships: {
				case: {
					concreteRelationship: {
						include: {target: true},
					},
					reverseRelationship: {
						include: {referencing: {
							case: {
								concreteRelationship: {
									include: {source: true},
								},
								interfaceRelationship: {
									include: {source: true},
								},
							}
						}},
					},
				},
			},
			targeted: {
				include: {source: true},
			},
			implements: true,
			module: true,
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