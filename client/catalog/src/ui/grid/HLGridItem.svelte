<script >
	export let name = "";
	export let href = null;
	export let click = null;
	import {cap} from '../../app/util.js';
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
	}

	a.object-box:hover > div.obj-title {
		border-top: 6px solid var(--border-color);
		margin-top: -5px;
	}
</style>
<a href="{href}" class="object-box" on:click={click}>
	<div class="obj-title">{cap(name)}</div>
	<slot />
</a>
