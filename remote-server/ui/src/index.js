
import './app.css'


import { EditorView, basicSetup } from 'codemirror'
import { keymap } from "@codemirror/view"
import {indentWithTab} from "@codemirror/commands"
import {EditorState} from "@codemirror/state"
import { python } from "@codemirror/lang-python"
import { json } from "@codemirror/lang-json"
import { solarizedDark } from 'cm6-theme-solarized-dark'


async function loadProgram(editor) {
  // Load program from API
  const response = await fetch('/program');
  const program = await response.text();
  editor.dispatch({
    changes: {
      from: 0, 
      insert: program,
    },
  });
}

async function saveProgram(editor) {
  // Get program from editor
  const program = editor.state.doc.toString();

  // Save program to API
  await fetch('/program', {
    method: 'POST',
    body: program,
  });
}

function connectResults() {
  const loc = window.location;
  let uri = "ws:";
  if (loc.protocol === "https:") {
    uri = "wss:";
  }
  uri += "//" + loc.host + "/results";
  const socket = new WebSocket(uri);
  return socket;
}

function renderOutput(root, data) {
  const probeId = data.probe_id;
  const timestamp = data.timestamp;
  const content = JSON.stringify(data.result, null, 2);
  
  const div = document.createElement('div');
  div.classList.add('result');

  const header = document.createElement('div');
  header.classList.add('header');


  const timestampView = document.createElement('span');
  timestampView.classList.add('timestamp');
  timestampView.innerText = timestamp;
  header.appendChild(timestampView);

  const tagView = document.createElement('span');
  tagView.classList.add('tag');
  tagView.innerText = data.tag;
  header.appendChild(tagView);


  const probeIdView = document.createElement('span');
  probeIdView.classList.add('probeid');
  probeIdView.innerText = probeId;
  header.appendChild(probeIdView);

  div.appendChild(header);

  const contentView = document.createElement('pre');
  contentView.classList.add('content');
  contentView.innerText = content;

  div.appendChild(contentView);

  root.prepend(div);
}

window.renderOutput = renderOutput;

document.addEventListener("DOMContentLoaded", async () => {


  const editorRoot = document.getElementById('editor');
  let editor = new EditorView({
      doc: '',
      extensions: [
        basicSetup,
        keymap.of([indentWithTab]),
        python(),
        solarizedDark,
      ],
      parent: editorRoot,
  });

  const btnSave = document.getElementById('btn-save');
  btnSave.addEventListener('click', async () => {
    await saveProgram(editor);
    outputRoot.innerHTML = '';
  });

  let outputRoot = document.getElementById('output');

  // connect to websocket
  const resultsSocket = connectResults();
  resultsSocket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    renderOutput(outputRoot, data);
  };

  await loadProgram(editor);
});



