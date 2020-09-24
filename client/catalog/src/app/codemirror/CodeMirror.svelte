<script >
	import { onMount } from 'svelte';
	import { createEventDispatcher } from 'svelte';
	import { dirtyStore } from '../../app/stores.js';

	import client from '../../data/client.js';

	import CodeMirror from 'codemirror';
	import 'codemirror/mode/python/python.js';
	import 'codemirror/addon/selection/active-line.js';
	import 'codemirror/addon/edit/closebrackets.js';
	import 'codemirror/addon/comment/comment.js';
	import 'codemirror/addon/edit/matchbrackets.js';
	import 'codemirror/addon/lint/lint.js';
	import "./codemirror.css";
	import "./dracula.css";
	import "./dracula.css";

	export let value = "";

	var cm = {};
	let editor = null;
	const dispatch = createEventDispatcher();
	let textarea;

	function set() {
		if (editor === null) {
			return;
		}
		const val = editor.getValue()
		if (val === value) {
			return;
		}
		const {left, top} = editor.getScrollInfo();
		editor.setValue(value);
		editor.scrollTo(left, top);
	}

	$: set(value);

	onMount(async () => {
		editor = CodeMirror.fromTextArea(textarea, {
			mode: {name:"python"},
			lineNumbers: true,
			indentUnit: 4,
			theme: "dracula",
			autoCloseBrackets: true,
			lineWrapping: true,
			matchBrackets: true,
			gutter: true,
			lint: {
				"getAnnotations": CodeMirror.remoteValidator,
				"async": true,
				"selfContain": true,
				"check_cb": check_syntax
			}
		});

		editor.on('changes', instance => {
			const val = instance.getValue()
			if (val === value) {
				return
			}
			value = val;
		})

		cm.inner = editor;
		
		cm.parses = async function() {
			if(cm.originalCode === cm.getValue()) {
				return true;
			}
			const lint = await client.rpc.lint({args: {data : cm.getValue()}});
			if(!lint.parsed) {
				alert(lint.message + " at line " + lint.line + " char " + lint.start);
				cm.setCursor(lint.line-1, lint.start-1);
				cm.inner.focus();
				return false;
			}
			return true;
		}
		cm.focus = function() {
			cm.inner.focus();
		}
		cm.setValue = function(code) {
			cm.originalCode = code;
			cm.inner.setValue(code);
			cm.inner.doc.clearHistory();
		}

		cm.getValue = function() {
			return cm.inner.getValue();
		}

		cm.setOption = function(a, b) {
			return cm.inner.setOption(a, b);
		}
		cm.setSize = function(a, b) {
			return cm.inner.setSize(a, b);
		}
		
		cm.setCursor = function(a, b) {
			return cm.inner.setCursor(a, b);
		}
		
		cm.getCursor = function() {
			return cm.inner.getCursor();
		}
		
		cm.getHistory = function() {
			return cm.inner.getHistory();
		}
		cm.setHistory = function(a) {
			return cm.inner.setHistory(a);
		}
		
		cm.isClean = function () {
			if(cm.originalCode == cm.getValue()) {
				return true;
			}
			return cm.inner.isClean();
		}
		cm.setClean = function () {
			return cm.inner.doc.markClean();
		}
		cm.lastLine = function() {
			return cm.inner.lastLine();
		}
		
		dispatch('initialized', cm);
		set();
		// this initial set is our first state in undo
		cm.inner.doc.clearHistory();
	});

	var check_syntax = async function (code, result_cb) {
		const lint = await client.rpc.lint({args: {data : code}});
		if(lint.parsed) {
			result_cb([]);
		} else {
			result_cb([{
				line_no: lint.line,
				column_no_start: 0,
				message: lint.message,
				severity: "error"
			}]);
		}
	}

	CodeMirror.remoteValidator = function(text, updateLinting, options) {
		if(text.trim() == "") {
			updateLinting([]);
			return;
		}

		function result_cb(error_list)
		{
			var found = [];
			for(var i in error_list) {
				var error = error_list[i];	
				var line = error.line_no;
				var message = error.message;
				var start_char;
				if(typeof(error.column_no_start) != "undefined") {
					start_char = error.column_no_start - 1;            
				}
				else {
					start_char = 0;            
				}

				var severity;
				if(typeof(error.severity) != "undefined") {
					severity = error.severity;            
				}
				else {
					severity = 'error';            
				}
				found.push({
					from: CodeMirror.Pos(line - 1, start_char),
				//1000 basically sets the to position to infinity. This just highlights the entire line.
				to: CodeMirror.Pos(line - 1, 1000),
				message: message,
				severity: severity
			});
			}
			updateLinting(cm.inner, found);
		}
		options.check_cb(text, result_cb)
	}
</script>

<style>
</style>

<textarea bind:this={textarea} style="display: none;" ></textarea>
