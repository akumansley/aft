<script >
	export let name = "";
	export let href;
	import {cap} from '../app/util.js';
</script>
<style>
	a.object-box {
		display: flex;
		flex-direction: column;
		color: inherit;
		width: 150px;
		padding: 1em 1.5em;
	}
	a.object-box:hover {
		background: var(--background-highlight);
	}
	a.object-box.center {
		align-items: center;
		justify-content: center;
	}
	.obj-title{
		font-weight: 600;
	}
</style>
<a href="{href}" class="object-box">
	<div class="obj-title">{cap(name)}</div>
	<slot />
</a>