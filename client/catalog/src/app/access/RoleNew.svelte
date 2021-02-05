<script>
	import {navStore} from '../stores.js';
	import {router} from '../router.js';
	import RoleForm from './RoleForm.svelte';

	import client from '../../data/client.js';
	import {nonEmpty} from '../../lib/util.js';
	import {ObjectOperation, AttributeOperation, SetOperation, RelationshipOperation, ConnectOperation} from '../../api/object.js';

	import HLButton from '../../ui/form/HLButton.svelte';
	import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js'
	import Name from '../Name.svelte';

	navStore.set("access")
	
	let value = ObjectOperation({
		name: AttributeOperation(""),
		policies: RelationshipOperation(
			ObjectOperation({
				"interface": SetOperation(),
				readWhere: AttributeOperation("{}"),
				createWhere: AttributeOperation("{}"),
				updateWhere: AttributeOperation("{}"),
				allowRead: AttributeOperation(true),
				allowCreate: AttributeOperation(true),
				allowUpdate: AttributeOperation(true),
			}),
			),
		executableFunctions: RelationshipOperation(ConnectOperation()),
		module: ConnectOperation(),
	});

	const saveAndNav = async () => {
		const op = value.op()
		if (nonEmpty(op)) {
			await client.api.role.create(op.create);
		}
		router.route("/roles");
	};

</script>

<RoleForm on:save={saveAndNav} bind:value={value} />
