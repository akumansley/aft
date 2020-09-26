<script>
	import {navStore} from '../stores.js';
	import {router} from '../router.js';

	import client from '../../data/client.js';
	import {nonEmpty} from '../../lib/util.js';
	import {ObjectOperation, AttributeOperation} from '../../api/object.js';

	import HLRowButton from '../../ui/list/HLRowButton.svelte';
	import HLButton from '../../ui/form/HLButton.svelte';
	import HLRow from '../../ui/list/HLRow.svelte';
	import {HLHeader, HLContent, HLHeaderItem} from '../../ui/page/page.js'
	import Name from '../Name.svelte';


	const role = ObjectOperation({
		name: AttributeOperation(""),
	});

	const saveAndNav = async () => {
		const op = role.op()
		if (nonEmpty(op)) {
			await client.api.role.create(op.create);
		}
		router.route("/roles");
	};

</script>

<HLHeader>
	<HLHeaderItem>	
		<Name id="name" placeholder="Role name.." bind:value={role.name} />
	</HLHeaderItem>	

	<HLHeaderItem>	
		<HLButton on:click={saveAndNav}>Save</HLButton>
	</HLHeaderItem>
</HLHeader>
