<script>
import {router} from './router.js';

function login() {
	router.route("/models");
}
</script>

<style>
img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    z-index:1;
}

.hero {
  position: absolute;
  top: 35%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index:2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width:40%;
}

.text {
  text-align: center;
  text-shadow: 3px 3px 0px rgba(0,0,0,0.2);
  font-size: var(--scale-4);
  font-weight: 600;
  -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Old versions of Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none;
   margin-bottom:1em;
}

.small {
  text-align: center;
  text-shadow: 3px 3px 0px rgba(0,0,0,0.2);
  font-weight: 400;
  -webkit-touch-callout: none; /* iOS Safari */
    -webkit-user-select: none; /* Safari */
     -khtml-user-select: none; /* Konqueror HTML */
       -moz-user-select: none; /* Old versions of Firefox */
        -ms-user-select: none; /* Internet Explorer/Edge */
            user-select: none;
   margin-top:1em;
}

button {
	-webkit-appearance: none;
	-moz-appearance: none;
	appearance: none;
	background-color: var(--background);
	border-radius: 0;
	margin: 0;
	padding: 1em 2em;
	color: var(--text-color);
	border-color: var(--border-color);
    font-weight: 500;
}

button:active {
	-webkit-appearance: none;
	-moz-appearance: none;
	appearance: none;
	background-color: var(--background-highlight);
	border-radius: 0;
	margin: 0;
	color: inherit;
	border-color: var(--border-color);
}

</style>
<img src="./shutterstock_663523441.jpg" alt="boat background"/>
<div class="hero">
		<div class="text" style="">Aft</div>
		<button on:click={login}>Start</button>
		<div class="small">Serenity is not freedom from the storm, but peace amid it.</div>
</div>