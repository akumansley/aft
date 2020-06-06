<script>
export let params;
import client from '../../data/client.js';
import {Runtime, Storage} from '../../data/enums.js';
import { breadcrumbStore } from '../stores.js';
import {cap} from '../util.js';
import { getContext } from 'svelte'
import HLRowButton from '../../ui/HLRowButton.svelte';
import HLTable from '../../ui/HLTable.svelte';
import HLCodeMirror from '../../ui/HLCodeMirror.svelte';
import HLButton from '../../ui/HLButton.svelte';
import HLRow from '../../ui/HLRow.svelte';
import HLTextBig from '../../ui/HLTextBig.svelte';
import HLSelect from '../../ui/HLSelect.svelte';

let id = params.id;
let load = client.datatype.findOne({where: {id: id}, include: {validator: true}});
var cm;
var name = "code";
var runtime;
var code = "";

load.then(obj => {
	breadcrumbStore.set(
		[{
			href: "/datatypes",
			text: "Datatypes",
		}, {
			href: "/datatypes/" + id,
			text: cap(obj.name),
		}]
	);
	runtime = obj.validator.runtime;
	code = obj.validator.code;
});

function setUpCM() {
	cm = getContext(name);
	cm.setValue(code);
	cm.setOption("readOnly", true);
}
</script>

<style>
	.box {
		margin: 1em 1.5em;
	}
	h1 {
		font-size: var(--scale-3);
		font-weight: 600;
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
</style>

<div class="box">
	{#await load}
		&nbsp;
	{:then datatype}
	<h1>{cap(datatype.name)}</h1>
	<HLTable>
		<h2>Runtime: {Runtime[datatype.validator.runtime]}</h2>
	</HLTable>
	{#if runtime == 2}
	<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
	{/if}
	<HLTable>
		<h2>Stored As: {Storage[datatype.storedAs]}</h2>
	</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>

