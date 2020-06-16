<script >
	export let name;
	import { onMount } from 'svelte';
	import { setContext } from 'svelte';
	import { createEventDispatcher } from 'svelte';
	import CodeMirror from 'codemirror';
	import 'codemirror/mode/python/python.js';
	import 'codemirror/addon/selection/active-line.js';

	const dispatch = createEventDispatcher();
	onMount(async () => {
	var cm = CodeMirror.fromTextArea(document.querySelector("#" + name), {
		mode: {name:"python"},
		lineNumbers: true,
		indentUnit: 4,
		theme: "duotone-dark",
	});
	setContext(name, cm);
	dispatch('initialized');
});

</script>

<textarea id={name} name={name} style="display: none;"></textarea>