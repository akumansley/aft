<script >
	export let name = "";
	export let href = null;
	export let click = null;
	import {cap} from '../../lib/util.js';
</script>
<style>
	a.object-box {
		display: flex;
		flex-direction: column;
		color: inherit;
		width: 150px;
		padding: var(--box-margin);
	}
	.obj-title {
		border-top: 1px solid var(--border-color);
		padding-top: .25em;
		transition: margin-top .1s, border-top-width .1s, padding-bottom .1s;
	}

	a.object-box:hover > div.obj-title {
		margin-top: -8px;
		padding-bottom: 5px;
		border-top-width: 4px;
		transition: margin-top .1s, border-top-width .1s, padding-bottom .1s;
	}
</style>
<a href="{href}" class="object-box" on:click={click}>
	<div class="obj-title">{cap(name)}</div>
	<slot />
</a>
