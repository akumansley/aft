<script>

	import { router } from '../router.js';
	import client from '../../data/client.js';
	import {ObjectOperation, AttributeOperation, RelationshipOperation} from '../../api/object.js';

	import EnumForm from './EnumForm.svelte';

	async function saveAndNav() {
		router.route("/datatypes");
	}

	let value = ObjectOperation({
		name: AttributeOperation(""),
		enumValues: RelationshipOperation(
			ObjectOperation({
				name: AttributeOperation(""),
			})
			)
	});

</script>

<EnumForm on:save={saveAndNav} bind:value={value} />

