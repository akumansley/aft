<script>
	import { setContext } from 'svelte';
	import { key } from './context.js';

	setContext(key, "update");
</script>

<slot/>