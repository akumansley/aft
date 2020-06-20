<script>
export let params;
import client from '../../data/client.js';
import {Runtime} from '../../data/enums.js';
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
let load = client.rpc.findOne({where: {id: id}, include: {code: true}});
var cm;
var name = "code";
var code = "";

load.then(obj => {
	breadcrumbStore.set(
		[{
			href: "/rpcs",
			text: "RPCs",
		}, {
			href: "/rpc/" + id,
			text: cap(obj.name),
		}]
	);
	code = obj.code.code;
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
	{:then rpc}
	<h1>{cap(rpc.name)}</h1>
	<HLTable>
		<h2>Runtime: {Runtime[rpc.code.runtime]}</h2>
	</HLTable>
	<HLCodeMirror name={name} on:initialized={setUpCM}></HLCodeMirror>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>

