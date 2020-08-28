<script>
export let rightAlignLast=false;
export let id = "";
export let placeholder = "Name";
export let value="";
export let click=null;

import {restrictToIdent} from './util.js';
import HLButton from '../ui/form/HLButton.svelte';
import HLRow from '../ui/list/HLRow.svelte';
import HLText from '../ui/form/HLText.svelte';

function emptyName(e) {
	if(value == "") {
		alert(placeholder + " can't be blank");
		document.getElementById(id).focus();
		return false;
	}
	click(e);
	return true;
}
	
</script>

<style>
.wrap {
	margin-left: var(--box-margin);
	margin-right: var(--box-margin);
	margin-top: .25em;
	margin-bottom: .25em;
	display:flex;
	align-items:center;
}
.function {
	color:var(--text-color-function);
}
.spacer-small {
	width: 1em;
}
</style>

<div class="header">
		<div class="wrap">
			<div class="function">
				<HLText id={id} placeholder={placeholder} bind:value={value} restrict={restrictToIdent}/>
			</div>
			<div class="spacer-small"></div>
			<HLButton on:click>
					Save
			</HLButton>
			<div class="spacer-small"></div>
			<slot />
		</div>
</div>
