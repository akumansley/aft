<script>
	import { getContext, onMount } from 'svelte';
	import { key } from '../../api/api.js';

	export let init = null;

	let value = {
		name: "",
		id: "",
	};

	function clone(v) {
		return JSON.parse(JSON.stringify(v));
	}

	onMount(() => {
		if (init) {
			value = clone(init);
		}
		setOp();
	})


	export let op = null;

	import DatatypeSelect from './DatatypeSelect.svelte';
	import HLSelect from '../../ui/form/HLSelect.svelte';
	import HLText from '../../ui/form/HLText.svelte';
	import {restrictToIdent, cap, isObject, isEmptyObject} from '../../lib/util.js';

	let operation = getContext(key);

	let datatype = {};
	let datatypeOp;

	let data = {
		datatype: datatypeOp,
		name: "",
	};

	function setOp() {
		if (init !== null) {
			if (isEmptyObject(data)) {
				op = {};
			} else {
				op = {
					where: {id: value.id},
					data: data, 
				}
			}
		} else {
			op = data;
		}
		console.log("setOp", op);
	}

	function bindOp(data, init, key) {
		return function(e) {
			let newVal = e.detail.target.value;
			if (init !== null) {
				let initVal = init[key];

				console.log("bindOp", newVal, init);
				if (newVal === initVal) {
					delete data[key];
				} else {
					data[key] = newVal;
				}
			} else {
				data[key] = newVal;
				}
			console.log("bindOp", data);
			setOp();
		}
	}

</script>

<style>
	.hform-row {
		display: flex; 
		flex-direction: row;
		padding: calc(var(--box-margin)/ 2) var(--box-margin);
	}
	.spacer {
		width: 1em;
		height: 0;
	}
</style>

<div class="hform-row">
	
	<HLText placeholder="Attribute name.." on:input={bindOp(data, init, "name")} value={value.name} restrict={restrictToIdent}/>

	<div class="spacer"/>

	<DatatypeSelect bind:value={datatype} bind:op={datatypeOp}/>
</div>