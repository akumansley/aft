<script>
	export let name;
	export let header;
	import {cap} from '../app/util.js';
	import HLRow from './HLRow.svelte';
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
	{cap(name)}
	<div class="v-space"/>
	<dl>
		<dt>{header}</dt>
	<dd>
		<slot />
	</dd>
	</dl>
	</div>
</HLRow>