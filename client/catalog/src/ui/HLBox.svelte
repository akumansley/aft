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
	.v-space{
		height: .5em;
	}

	.v-space-2{
		height: 2em;
	}
	
	.hform-row {
		display: flex; 
		flex-direction: row;
	}
	.col {

	}
	.spacer {
		width: 1em;
		height: 0;
	}

	.hl-row-header {
		font-size: var(--scale--2);
		text-transform: uppercase;
		font-weight: 600;
	}

</style>
<div class="box">
<slot />
</div>