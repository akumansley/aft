import navaid from 'navaid';
import { dirtyStore} from './stores';
export const router = navaid();



let clean;
dirtyStore.subscribe(value => { 
		clean = value["clean"];
});

export let canRoute= (e) => {
	if (!clean) {
		let ok = confirm("Are you sure? Page has unsaved Changes");
		if(!ok) {
			e.preventDefault();
			return false;
		}
		dirtyStore.set({'clean' : true});
		return true;
	}
}
