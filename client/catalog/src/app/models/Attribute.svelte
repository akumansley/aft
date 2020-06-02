<script>
export let attribute
import HLRow from '../../ui/HLRow.svelte';
import { AttrType } from '../../data/enums.js';
</script>
<style>
dl {
	padding: 0; 
	margin:0;
}
dt {
	font-size: var(--scale--2);
	text-transform: uppercase;
	font-weight: 600;
}
dd {
	margin: 0;
}
.v-space{
	height: .5em;
}
.col {
	display: flex;
	flex-direction: column;
}
</style>
<HLRow>
	<div class="col">
	{attribute.name}
	<div class="v-space"/>
	<dl>
		<dt>Type</dt>
	<dd>
		{attribute}
	</dd>
	</dl>
	</div>
</HLRow>

