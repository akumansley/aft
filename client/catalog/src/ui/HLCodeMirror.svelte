<script >
	export let name;
	import { onMount } from 'svelte';
	import { setContext } from 'svelte';
	import { createEventDispatcher } from 'svelte';
	import CodeMirror from 'codemirror';
	import 'codemirror/mode/python/python.js';
	import 'codemirror/addon/selection/active-line.js';
	import 'codemirror/addon/edit/closebrackets.js';
	import 'codemirror/addon/comment/comment.js';

	var cm;
	const dispatch = createEventDispatcher();
	onMount(async () => {
	cm = CodeMirror.fromTextArea(document.querySelector("#" + name), {
		mode: {name:"python"},
		lineNumbers: true,
		indentUnit: 4,
		theme: "duotone-dark",
	    autoCloseBrackets: true
	});
	cm.setSize(null, 500);
	setContext(name, cm);
	dispatch('initialized');
	
});


function destroyCM(node, params) {
	cm.toTextArea();
	cm = null;
}

</script>
<textarea id={name} name={name} style="display: none;" out:destroyCM></textarea>


