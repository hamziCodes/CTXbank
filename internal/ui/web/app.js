// CTXbank Interactive Dashboard Application Logic (VERTEX Design)

document.addEventListener('DOMContentLoaded', () => {
  // Navigation State
  const navItems = document.querySelectorAll('.nav-item[data-view]');
  const viewPanels = document.querySelectorAll('.view-panel');
  let currentView = 'graph-view';

  navItems.forEach(item => {
    item.addEventListener('click', () => {
      const viewId = item.getAttribute('data-view');
      if (!viewId) return;

      navItems.forEach(i => i.classList.remove('active'));
      viewPanels.forEach(p => p.classList.remove('active'));

      item.classList.add('active');
      const targetPanel = document.getElementById(viewId);
      if (targetPanel) targetPanel.classList.add('active');
      currentView = viewId;

      if (viewId === 'graph-view') loadGraph();
      if (viewId === 'memory-view') loadFiles();
      if (viewId === 'checkpoints-view') loadCheckpoints();
    });
  });

  // Theme Toggle
  const btnTheme = document.getElementById('btn-theme-toggle');
  btnTheme.addEventListener('click', () => {
    const isDark = document.documentElement.classList.toggle('dark');
    localStorage.setItem('ctx_theme', isDark ? 'dark' : 'light');
    loadGraph(); // re-render graph colors
  });

  // Initial Data Load
  loadStatus();
  loadGraph();

  // Polling project health every 10 seconds
  setInterval(loadStatus, 10000);

  // 1. Status Loading
  async function loadStatus() {
    try {
      const res = await fetch('/api/status');
      const data = await res.json();

      document.getElementById('repo-name').textContent = data.project_name;
      document.getElementById('branch-chip').textContent = 'branch: ' + (data.branch || 'main');

      const dirtyChip = document.getElementById('dirty-chip');
      if (data.dirty_count > 0) {
        dirtyChip.textContent = `${data.dirty_count} files dirty`;
        dirtyChip.className = 'chip dirty';
      } else {
        dirtyChip.textContent = 'clean';
        dirtyChip.className = 'chip clean';
      }

      // Active Focus card
      document.getElementById('focus-content-text').textContent = data.active_focus;
      const budgetBadge = document.getElementById('focus-budget-badge');
      budgetBadge.textContent = `${data.active_lines} / 150 lines`;

      const meterBar = document.getElementById('focus-meter-bar');
      const pct = Math.min(100, Math.round((data.active_lines / 150) * 100));
      meterBar.style.width = pct + '%';
      meterBar.className = 'meter-fill ' + data.budget_status;

    } catch (err) {
      console.error('Failed to load status:', err);
    }
  }

  // 2. Interactive SVG Architecture Graph
  async function loadGraph() {
    const svg = document.getElementById('graph-svg');
    svg.innerHTML = '<text x="50%" y="50%" text-anchor="middle" fill="currentColor">Loading architecture tree...</text>';

    try {
      const res = await fetch('/api/graph');
      const data = await res.json();
      renderGraph(data.nodes, data.edges);
      document.getElementById('graph-node-count').textContent = `${data.nodes.length} nodes`;
    } catch (err) {
      svg.innerHTML = '<text x="50%" y="50%" text-anchor="middle" fill="red">Failed to load graph.</text>';
    }
  }

  function renderGraph(nodes, edges) {
    const svg = document.getElementById('graph-svg');
    svg.innerHTML = '';
    const width = svg.clientWidth || 800;
    const height = svg.clientHeight || 520;

    // Node positioning layout
    const coords = {
      projectbrief:   { x: width * 0.15, y: height * 0.25 },
      productContext: { x: width * 0.45, y: height * 0.15 },
      systemPatterns: { x: width * 0.15, y: height * 0.65 },
      techContext:    { x: width * 0.45, y: height * 0.75 },
      activeContext:  { x: width * 0.50, y: height * 0.45 },
      progress:       { x: width * 0.82, y: height * 0.30 },
      decisionLog:    { x: width * 0.82, y: height * 0.65 },
    };

    // Checkpoint nodes placed along right flow
    let ckptOffset = 0;
    nodes.forEach(n => {
      if (n.type === 'checkpoint') {
        coords[n.ID] = { x: width * 0.65 + (ckptOffset * 30), y: height * 0.85 };
        ckptOffset++;
      }
    });

    const isDark = document.documentElement.classList.contains('dark');
    const edgeColor = isDark ? '#30363d' : '#c2d0dd';
    const textColor = isDark ? '#f0f6fc' : '#000000';
    const accentColor = isDark ? '#609abe' : '#042940';

    // Render Edges
    edges.forEach(edge => {
      const src = coords[edge.Source];
      const tgt = coords[edge.Target];
      if (!src || !tgt) return;

      const line = document.createElementNS('http://www.w3.org/2000/svg', 'line');
      line.setAttribute('x1', src.x);
      line.setAttribute('y1', src.y);
      line.setAttribute('x2', tgt.x);
      line.setAttribute('y2', tgt.y);
      line.setAttribute('stroke', edgeColor);
      line.setAttribute('stroke-width', '1.5');
      line.setAttribute('stroke-dasharray', '4 2');
      svg.appendChild(line);
    });

    // Render Nodes (VERTEX Solid Cards)
    nodes.forEach(node => {
      const pos = coords[node.ID] || { x: width * 0.5, y: height * 0.5 };
      const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
      group.setAttribute('class', 'node-group');
      group.setAttribute('transform', `translate(${pos.x - 75}, ${pos.y - 25})`);

      const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
      rect.setAttribute('width', '150');
      rect.setAttribute('height', '50');
      rect.setAttribute('rx', '4');

      if (node.ID === 'activeContext') {
        rect.style.stroke = accentColor;
        rect.style.strokeWidth = '2';
      }

      const title = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      title.setAttribute('x', '10');
      title.setAttribute('y', '22');
      title.setAttribute('class', 'title');
      title.textContent = node.Label;

      const sub = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      sub.setAttribute('x', '10');
      sub.setAttribute('y', '38');
      sub.setAttribute('class', 'subtitle');
      sub.textContent = node.Subtitle + (node.Lines ? ` (${node.Lines}l)` : '');

      group.appendChild(rect);
      group.appendChild(title);
      group.appendChild(sub);

      group.addEventListener('click', () => {
        if (node.Type !== 'checkpoint') {
          openFileEditor(node.Label);
        }
      });

      svg.appendChild(group);
    });
  }

  // 3. Memory Bank Explorer & File Editor
  async function loadFiles() {
    const grid = document.getElementById('memory-cards-grid');
    grid.innerHTML = '<div class="state-loading">Loading memory bank files...</div>';

    try {
      const res = await fetch('/api/files');
      const files = await res.json();
      grid.innerHTML = '';

      files.forEach(f => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
          <div class="card-header">
            <span class="card-title">${f.Name}</span>
            <span class="chip ${f.BudgetStatus}">${f.LineCount} lines</span>
          </div>
          <p style="font-size: 12px; color: var(--app-ink-muted);">
            Volatility: <strong>${f.Volatility}</strong> • ${f.ByteSize} bytes
          </p>
          <div style="display: flex; justify-content: flex-end; margin-top: auto;">
            <button class="btn" onclick="openFileEditor('${f.Name}')">Edit File</button>
          </div>
        `;
        grid.appendChild(card);
      });
    } catch (err) {
      grid.innerHTML = '<div class="state-error">Failed to load memory bank files.</div>';
    }
  }

  window.openFileEditor = async function(filename) {
    const modal = document.getElementById('editor-modal');
    document.getElementById('editor-filename').textContent = filename;
    const textarea = document.getElementById('editor-textarea');
    textarea.value = 'Loading...';
    modal.classList.add('active');

    try {
      const res = await fetch('/api/files');
      const files = await res.json();
      const target = files.find(f => f.Name === filename);
      if (target) {
        textarea.value = target.Content;
        updateEditorIndicator();
      }
    } catch (err) {
      textarea.value = 'Error reading file.';
    }
  };

  const editorTextarea = document.getElementById('editor-textarea');
  editorTextarea.addEventListener('input', updateEditorIndicator);

  function updateEditorIndicator() {
    const lines = editorTextarea.value.split('\n').length;
    const indicator = document.getElementById('editor-budget-indicator');
    const filename = document.getElementById('editor-filename').textContent;

    if (filename === 'activeContext.md') {
      indicator.textContent = `${lines} / 150 lines`;
      indicator.className = 'chip ' + (lines >= 150 ? 'danger' : (lines > 120 ? 'warning' : 'clean'));
    } else {
      indicator.textContent = `${lines} lines`;
      indicator.className = 'chip';
    }
  }

  document.getElementById('btn-close-editor').addEventListener('click', () => {
    document.getElementById('editor-modal').classList.remove('active');
  });

  document.getElementById('btn-save-file').addEventListener('click', async () => {
    const filename = document.getElementById('editor-filename').textContent;
    const content = editorTextarea.value;
    const statusMsg = document.getElementById('editor-status-msg');

    statusMsg.textContent = 'Saving atomically...';
    try {
      const res = await fetch('/api/file/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ filename, content })
      });
      const data = await res.json();
      if (data.error) {
        statusMsg.textContent = 'Error: ' + data.error;
      } else {
        statusMsg.textContent = 'Saved successfully!';
        loadStatus();
        setTimeout(() => {
          document.getElementById('editor-modal').classList.remove('active');
          statusMsg.textContent = 'Atomic write ready';
        }, 600);
      }
    } catch (err) {
      statusMsg.textContent = 'Failed to save: ' + err.message;
    }
  });

  // 4. Research Ingestion Drag & Drop
  const dropzone = document.getElementById('ingest-dropzone');
  const fileInput = document.getElementById('ingest-file-input');
  let currentProposal = null;

  dropzone.addEventListener('click', () => fileInput.click());
  fileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) handleIngestFile(e.target.files[0]);
  });

  dropzone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropzone.classList.add('dragover');
  });
  dropzone.addEventListener('dragleave', () => dropzone.classList.remove('dragover'));
  dropzone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropzone.classList.remove('dragover');
    if (e.dataTransfer.files.length > 0) handleIngestFile(e.dataTransfer.files[0]);
  });

  async function handleIngestFile(file) {
    const text = await file.text();
    try {
      const res = await fetch('/api/ingest/preview', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ filename: file.name, content: text })
      });
      currentProposal = await res.json();

      document.getElementById('ingest-preview-container').style.display = 'flex';
      document.getElementById('ingest-target-file').textContent = currentProposal.target_file;
      document.getElementById('ingest-dedup-badge').textContent = `${currentProposal.deduplicated_count} duplicates suppressed`;

      const diffViewer = document.getElementById('ingest-diff-viewer');
      diffViewer.textContent = currentProposal.proposed_diff || 'All content near-duplicate. No addition proposed.';

    } catch (err) {
      alert('Ingestion analysis failed: ' + err.message);
    }
  }

  document.getElementById('btn-cancel-ingest').addEventListener('click', () => {
    document.getElementById('ingest-preview-container').style.display = 'none';
    currentProposal = null;
  });

  document.getElementById('btn-commit-ingest').addEventListener('click', async () => {
    if (!currentProposal) return;
    try {
      const res = await fetch('/api/ingest/commit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(currentProposal)
      });
      const data = await res.json();
      if (data.status === 'committed') {
        alert('Ingested and committed to ' + data.target);
        document.getElementById('ingest-preview-container').style.display = 'none';
        currentProposal = null;
        loadStatus();
      }
    } catch (err) {
      alert('Commit failed: ' + err.message);
    }
  });

  // 5. Checkpoints List & Creation
  async function loadCheckpoints() {
    const list = document.getElementById('checkpoints-list');
    list.innerHTML = '<div class="state-loading">Loading checkpoint history...</div>';

    try {
      const res = await fetch('/api/checkpoints');
      const data = await res.json();

      if (!data || data.length === 0) {
        list.innerHTML = '<div class="state-empty">No checkpoints recorded yet. Click "New Snapshot" to create one.</div>';
        return;
      }

      list.innerHTML = '';
      data.reverse().forEach(c => {
        const item = document.createElement('div');
        item.className = 'card';
        item.style.padding = '14px 18px';
        item.innerHTML = `
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <strong style="font-size: 14px;">${c.id}</strong>
              <span class="chip" style="margin-left: 8px;">branch: ${c.branch || 'main'}</span>
            </div>
            <span style="font-size: 12px; color: var(--app-ink-muted); font-variant-numeric: tabular-nums;">
              ${new Date(c.timestamp).toLocaleString()}
            </span>
          </div>
          <p style="font-size: 13px; color: var(--app-ink-soft); margin-top: 4px;">${c.active_focus || 'No focus summary'}</p>
          ${c.diff_stat ? `<pre style="font-size: 11px; color: var(--app-ink-muted); margin-top: 4px;">${c.diff_stat}</pre>` : ''}
        `;
        list.appendChild(item);
      });
    } catch (err) {
      list.innerHTML = '<div class="state-error">Failed to load checkpoints.</div>';
    }
  }

  // Checkpoint Modal
  const ckptModal = document.getElementById('checkpoint-modal');
  document.getElementById('btn-quick-checkpoint').addEventListener('click', () => ckptModal.classList.add('active'));
  document.getElementById('btn-add-checkpoint').addEventListener('click', () => ckptModal.classList.add('active'));
  document.getElementById('btn-close-modal').addEventListener('click', () => ckptModal.classList.remove('active'));
  document.getElementById('btn-modal-cancel').addEventListener('click', () => ckptModal.classList.remove('active'));

  document.getElementById('btn-modal-save').addEventListener('click', async () => {
    const focus = document.getElementById('ckpt-focus-input').value;
    const notes = document.getElementById('ckpt-notes-input').value;

    try {
      const res = await fetch('/api/checkpoint/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ focus, notes })
      });
      const data = await res.json();
      if (data.id) {
        ckptModal.classList.remove('active');
        document.getElementById('ckpt-focus-input').value = '';
        document.getElementById('ckpt-notes-input').value = '';
        loadStatus();
        if (currentView === 'checkpoints-view') loadCheckpoints();
        if (currentView === 'graph-view') loadGraph();
      }
    } catch (err) {
      alert('Failed to create checkpoint: ' + err.message);
    }
  });

  // Action Buttons: Audit and Lint
  document.getElementById('btn-run-audit').addEventListener('click', async () => {
    const btn = document.getElementById('btn-run-audit');
    btn.textContent = 'Running...';
    try {
      const res = await fetch('/api/audit', { method: 'POST' });
      const report = await res.json();
      alert(`Reconnaissance Complete!\n\nDetected Language: ${report.manifest_scan.language}\nComponents: ${report.structural_heuristics.length}\nSymbols: ${report.symbols_summary.exported_symbols_count}`);
    } catch (err) {
      alert('Audit failed: ' + err.message);
    } finally {
      btn.textContent = 'Run Audit';
    }
  });

  document.getElementById('btn-run-lint').addEventListener('click', async () => {
    const btn = document.getElementById('btn-run-lint');
    btn.textContent = 'Linting...';
    try {
      const res = await fetch('/api/lint?fix=true');
      const data = await res.json();
      if (data.passed) {
        alert('All memory bank token and format budgets satisfied (< 150 lines)!');
      } else {
        alert('Linter Warnings/Violations:\n' + data.violations.join('\n'));
      }
      loadStatus();
    } catch (err) {
      alert('Lint failed: ' + err.message);
    } finally {
      btn.textContent = 'Lint Memory';
    }
  });

});
