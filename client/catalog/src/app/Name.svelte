<script>
export let placeholder = "Name";
export let value="";

import {restrictToIdent} from './../lib/util.js';
import HLButton from '../ui/form/HLButton.svelte';
import HLRow from '../ui/list/HLRow.svelte';
import HLText from '../ui/form/HLText.svelte';
	
</script>

<style>
.function {
	color:var(--text-color-function);
}
</style>
<div class="function">
	<HLText placeholder={placeholder} bind:value={value} restrict={restrictToIdent}/>
</div>
