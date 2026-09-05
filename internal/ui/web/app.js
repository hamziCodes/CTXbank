// CTXbank Interactive Dashboard Application Logic (Smooth UX & Non-Tech Friendly)

document.addEventListener('DOMContentLoaded', () => {
  // Navigation State
  const navItems = document.querySelectorAll('.nav-item[data-view]');
  const viewPanels = document.querySelectorAll('.view-panel');
  let currentView = 'overview-view';

  function switchView(viewId) {
    if (!viewId) return;

    navItems.forEach(i => i.classList.remove('active'));
    viewPanels.forEach(p => p.classList.remove('active'));

    const activeNav = document.querySelector(`.nav-item[data-view="${viewId}"]`);
    if (activeNav) activeNav.classList.add('active');

    const targetPanel = document.getElementById(viewId);
    if (targetPanel) targetPanel.classList.add('active');
    currentView = viewId;

    if (viewId === 'overview-view') loadStatus();
    if (viewId === 'graph-view') loadGraph();
    if (viewId === 'memory-view') loadFiles();
    if (viewId === 'checkpoints-view') loadCheckpoints();
  }

  navItems.forEach(item => {
    item.addEventListener('click', () => {
      const viewId = item.getAttribute('data-view');
      switchView(viewId);
    });
  });

  // Overview Action Buttons
  const btnOverviewOpenMemory = document.getElementById('btn-overview-open-memory');
  if (btnOverviewOpenMemory) {
    btnOverviewOpenMemory.addEventListener('click', () => switchView('memory-view'));
  }

  const btnOverviewSaveCkpt = document.getElementById('btn-overview-save-ckpt');
  if (btnOverviewSaveCkpt) {
    btnOverviewSaveCkpt.addEventListener('click', () => openCheckpointModal());
  }

  const btnEditActiveFocus = document.getElementById('btn-edit-active-focus');
  if (btnEditActiveFocus) {
    btnEditActiveFocus.addEventListener('click', () => openEditor('activeContext.md'));
  }

  // Theme Toggle
  const btnTheme = document.getElementById('btn-theme-toggle');
  btnTheme.addEventListener('click', () => {
    const isDark = document.documentElement.classList.toggle('dark');
    localStorage.setItem('ctx_theme', isDark ? 'dark' : 'light');
    if (currentView === 'graph-view') loadGraph();
  });

  // Initial Data Load
  loadStatus();
  setInterval(loadStatus, 10000);

  // Friendly File Roles Dictionary for Non-Tech Users
  const fileDescriptions = {
    'projectbrief.md': 'The Big Picture — Core mission, problem statement, and goals of this project.',
    'productContext.md': 'User Experience — Why this project exists and how users should experience it.',
    'systemPatterns.md': 'Code Rules & Patterns — The technical design standards and architecture rules.',
    'techContext.md': 'Tech Stack — Languages, frameworks, libraries, and tools used.',
    'activeContext.md': 'Current Task (High Priority) — What you or your AI are working on right now.',
    'progress.md': 'Milestone Tracker — What is finished, what is currently being built, and what is next.',
    'decisionLog.md': 'Decision History — Why key architectural choices were made (stops circular debates).'
  };

  // 1. Status Loading
  async function loadStatus() {
    try {
      const res = await fetch('/api/status');
      const data = await res.json();

      // Header Bar
      document.getElementById('repo-name').textContent = data.project_name || 'Project';
      document.getElementById('branch-name').textContent = data.branch || 'main';

      const gitStatusPill = document.getElementById('git-status-pill');
      const gitStatusText = document.getElementById('git-status-text');
      if (data.dirty_count > 0) {
        gitStatusText.textContent = `${data.dirty_count} files modified`;
        gitStatusPill.className = 'status-pill dirty';
      } else {
        gitStatusText.textContent = 'Clean';
        gitStatusPill.className = 'status-pill clean';
      }

      // Overview Tab Hero & Metric Cards
      document.getElementById('overview-focus-text').textContent = data.active_focus || 'Ready for next task.';
      document.getElementById('overview-branch-val').textContent = data.branch || 'main';
      document.getElementById('overview-tree-val').textContent = data.dirty_count > 0 ? `${data.dirty_count} modified` : 'Clean';

      const budgetLines = data.active_lines || 0;
      document.getElementById('overview-budget-text').textContent = `${budgetLines} / 150 lines`;

      const budgetPct = Math.min(100, Math.round((budgetLines / 150) * 100));
      document.getElementById('hero-budget-pct').textContent = budgetPct + '%';

      const budgetBar = document.getElementById('overview-budget-bar');
      if (budgetBar) {
        budgetBar.style.width = budgetPct + '%';
        budgetBar.className = 'budget-bar ' + (data.budget_status || 'normal');
      }

      // Load snapshot count for hero stat
      fetchCheckpointsCount();

    } catch (err) {
      console.error('Failed to load status:', err);
    }
  }

  async function fetchCheckpointsCount() {
    try {
      const res = await fetch('/api/checkpoints');
      const data = await res.json();
      const count = Array.isArray(data) ? data.length : (data.checkpoints ? data.checkpoints.length : 0);
      const heroStat = document.getElementById('hero-snapshots-count');
      if (heroStat) heroStat.textContent = count;
    } catch (e) {}
  }

  // 2. Interactive SVG Architecture Graph
  async function loadGraph() {
    const svg = document.getElementById('graph-svg');
    svg.innerHTML = '<text x="50%" y="50%" text-anchor="middle" fill="currentColor">Loading architecture tree...</text>';

    try {
      const res = await fetch('/api/graph');
      const data = await res.json();
      renderGraph(data.nodes || [], data.edges || []);
      const countEl = document.getElementById('graph-node-count');
      if (countEl && data.nodes) countEl.textContent = `${data.nodes.length} nodes`;
    } catch (err) {
      svg.innerHTML = '<text x="50%" y="50%" text-anchor="middle" fill="red">Failed to load graph.</text>';
    }
  }

  function renderGraph(nodes, edges) {
    const svg = document.getElementById('graph-svg');
    svg.innerHTML = '';
    const width = svg.clientWidth || 700;
    const height = svg.clientHeight || 520;

    // Node positioning layout
    const coords = {
      projectbrief:   { x: width * 0.16, y: height * 0.25 },
      productContext: { x: width * 0.48, y: height * 0.18 },
      systemPatterns: { x: width * 0.16, y: height * 0.65 },
      techContext:    { x: width * 0.48, y: height * 0.78 },
      activeContext:  { x: width * 0.50, y: height * 0.48 },
      progress:       { x: width * 0.82, y: height * 0.32 },
      decisionLog:    { x: width * 0.82, y: height * 0.68 },
    };

    let ckptOffset = 0;
    nodes.forEach(n => {
      if (n.type === 'checkpoint') {
        coords[n.ID] = { x: width * 0.65 + (ckptOffset * 35), y: height * 0.88 };
        ckptOffset++;
      }
    });

    const isDark = document.documentElement.classList.contains('dark');
    const edgeColor = isDark ? '#33445b' : '#cbd5e1';
    const textColor = isDark ? '#f8fafc' : '#0f172a';
    const subTextColor = isDark ? '#94a3b8' : '#64748b';
    const rectFill = isDark ? '#111720' : '#ffffff';
    const rectStroke = isDark ? '#222e3e' : '#e2e8f0';

    // Render Edges (smooth dashed connection lines)
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
      line.setAttribute('stroke-width', '1.75');
      line.setAttribute('stroke-dasharray', '5 3');
      svg.appendChild(line);
    });

    // Render Nodes (Smooth rounded cards)
    nodes.forEach(node => {
      const pos = coords[node.ID] || { x: width * 0.5, y: height * 0.5 };
      const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
      group.setAttribute('class', 'node-group');
      group.setAttribute('transform', `translate(${pos.x - 75}, ${pos.y - 26})`);

      const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
      rect.setAttribute('width', '150');
      rect.setAttribute('height', '52');
      rect.setAttribute('rx', '10');
      rect.setAttribute('ry', '10');
      rect.setAttribute('fill', rectFill);
      rect.setAttribute('stroke', node.ID === 'activeContext' ? (isDark ? '#609abe' : '#042940') : rectStroke);
      rect.setAttribute('stroke-width', node.ID === 'activeContext' ? '2' : '1.5');
      group.appendChild(rect);

      // Title
      const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      text.setAttribute('x', '14');
      text.setAttribute('y', '22');
      text.setAttribute('fill', textColor);
      text.setAttribute('font-size', '12');
      text.setAttribute('font-weight', '600');
      text.setAttribute('font-family', 'Inter, sans-serif');
      text.textContent = node.Name || node.ID;
      group.appendChild(text);

      // Subtitle / Line Count
      const sub = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      sub.setAttribute('x', '14');
      sub.setAttribute('y', '39');
      sub.setAttribute('fill', subTextColor);
      sub.setAttribute('font-size', '11');
      sub.setAttribute('font-family', 'JetBrains Mono, monospace');
      sub.textContent = node.Lines ? `${node.Lines} lines` : (node.type || 'artifact');
      group.appendChild(sub);

      // Click node to open Inspector
      group.addEventListener('click', () => {
        inspectNode(node);
      });

      svg.appendChild(group);
    });
  }

  function inspectNode(node) {
    const inspector = document.getElementById('node-inspector');
    const nameEl = document.getElementById('inspector-name');
    const typeEl = document.getElementById('inspector-type');
    const descEl = document.getElementById('inspector-desc');
    const detailsEl = document.getElementById('inspector-details');
    const pathEl = document.getElementById('inspector-path');
    const linesEl = document.getElementById('inspector-lines');
    const btnEdit = document.getElementById('btn-inspector-edit');

    const filename = node.ID + '.md';
    nameEl.textContent = node.Name || filename;
    typeEl.textContent = node.type || 'Memory Artifact';
    descEl.textContent = fileDescriptions[filename] || 'A core architectural artifact tracked by CTXbank.';

    detailsEl.style.display = 'block';
    pathEl.textContent = `memory-bank/${filename}`;
    linesEl.textContent = node.Lines ? `${node.Lines} lines` : 'N/A';

    btnEdit.onclick = () => openEditor(filename);
  }

  // 3. Memory Bank Grid
  async function loadFiles() {
    const grid = document.getElementById('memory-cards-grid');
    grid.innerHTML = '<div class="text-sm text-muted">Loading memory bank files...</div>';

    try {
      const res = await fetch('/api/files');
      const files = await res.json();
      grid.innerHTML = '';

      files.forEach(f => {
        const card = document.createElement('div');
        card.className = 'memory-file-card';

        const roleDesc = fileDescriptions[f.filename] || 'Core project knowledge file.';
        const budgetWarning = (f.filename === 'activeContext.md' && f.lines >= 150);

        card.innerHTML = `
          <div class="card-top-row">
            <span class="file-name-title">${f.filename}</span>
            <span class="chip ${budgetWarning ? 'dirty' : 'clean'}">${f.lines} lines</span>
          </div>
          <p class="file-friendly-role">${roleDesc}</p>
          <div class="file-meta-row">
            <span>Size: ${(f.bytes / 1024).toFixed(1)} KB</span>
            <span class="font-mono text-xs">Click to edit &rarr;</span>
          </div>
        `;

        card.addEventListener('click', () => openEditor(f.filename));
        grid.appendChild(card);
      });

    } catch (err) {
      grid.innerHTML = '<div class="text-sm text-muted">Failed to load memory bank files.</div>';
    }
  }

  const btnRefreshFiles = document.getElementById('btn-refresh-files');
  if (btnRefreshFiles) btnRefreshFiles.addEventListener('click', loadFiles);

  // 4. File Editor Modal
  const editorModal = document.getElementById('editor-modal');
  const editorTextarea = document.getElementById('editor-textarea');
  const editorFilename = document.getElementById('editor-filename');
  const editorBudget = document.getElementById('editor-budget-indicator');
  const btnCloseEditor = document.getElementById('btn-close-editor');
  const btnCancelEditor = document.getElementById('btn-cancel-editor');
  const btnSaveFile = document.getElementById('btn-save-file');
  let activeEditingFile = '';

  async function openEditor(filename) {
    activeEditingFile = filename;
    editorFilename.textContent = filename;
    editorTextarea.value = 'Loading file contents...';
    editorModal.classList.add('active');

    try {
      const res = await fetch(`/api/file?name=${encodeURIComponent(filename)}`);
      const data = await res.json();
      editorTextarea.value = data.content || '';
      updateEditorBudget();
    } catch (err) {
      editorTextarea.value = 'Error loading file content.';
    }
  }

  function updateEditorBudget() {
    const lines = editorTextarea.value.split('\n').length;
    editorBudget.textContent = `${lines} lines`;
    if (activeEditingFile === 'activeContext.md') {
      editorBudget.textContent = `${lines} / 150 lines`;
      if (lines > 150) {
        editorBudget.className = 'chip dirty';
      } else {
        editorBudget.className = 'chip clean';
      }
    } else {
      editorBudget.className = 'chip';
    }
  }

  editorTextarea.addEventListener('input', updateEditorBudget);

  if (btnCloseEditor) btnCloseEditor.addEventListener('click', () => editorModal.classList.remove('active'));
  if (btnCancelEditor) btnCancelEditor.addEventListener('click', () => editorModal.classList.remove('active'));

  if (btnSaveFile) {
    btnSaveFile.addEventListener('click', async () => {
      const content = editorTextarea.value;
      btnSaveFile.disabled = true;
      btnSaveFile.textContent = 'Saving...';

      try {
        const res = await fetch('/api/file/save', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ filename: activeEditingFile, content: content })
        });
        const result = await res.json();
        if (result.success) {
          showToast(`Saved ${activeEditingFile} atomically.`);
          editorModal.classList.remove('active');
          loadStatus();
          if (currentView === 'memory-view') loadFiles();
          if (currentView === 'graph-view') loadGraph();
        } else {
          showToast('Failed to save file: ' + (result.error || 'unknown'));
        }
      } catch (err) {
        showToast('Error saving file.');
      } finally {
        btnSaveFile.disabled = false;
        btnSaveFile.textContent = 'Save Changes';
      }
    });
  }

  // 5. Ingest Notes & Dropzone
  const dropzone = document.getElementById('ingest-dropzone');
  const fileInput = document.getElementById('ingest-file-input');
  const previewContainer = document.getElementById('ingest-preview-container');
  const diffViewer = document.getElementById('ingest-diff-viewer');
  const dedupBadge = document.getElementById('ingest-dedup-badge');
  const btnCancelIngest = document.getElementById('btn-cancel-ingest');
  const btnCommitIngest = document.getElementById('btn-commit-ingest');
  let currentRawIngestNote = '';

  if (dropzone) {
    dropzone.addEventListener('click', () => fileInput.click());
    dropzone.addEventListener('dragover', (e) => {
      e.preventDefault();
      dropzone.classList.add('dragover');
    });
    dropzone.addEventListener('dragleave', () => dropzone.classList.remove('dragover'));
    dropzone.addEventListener('drop', (e) => {
      e.preventDefault();
      dropzone.classList.remove('dragover');
      if (e.dataTransfer.files && e.dataTransfer.files[0]) {
        processIngestFile(e.dataTransfer.files[0]);
      }
    });
  }

  if (fileInput) {
    fileInput.addEventListener('change', () => {
      if (fileInput.files && fileInput.files[0]) {
        processIngestFile(fileInput.files[0]);
      }
    });
  }

  async function processIngestFile(file) {
    const text = await file.text();
    currentRawIngestNote = text;

    try {
      const res = await fetch('/api/ingest/preview', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: text, target: 'productContext.md' })
      });
      const data = await res.json();

      previewContainer.style.display = 'block';
      diffViewer.textContent = data.diff || '(No new content detected or all duplicate)';
      dedupBadge.textContent = `${data.duplicate_count || 0} duplicate paragraphs removed`;
      showToast('SimHash deduplication completed.');
    } catch (e) {
      showToast('Failed to preview note ingestion.');
    }
  }

  if (btnCancelIngest) {
    btnCancelIngest.addEventListener('click', () => {
      previewContainer.style.display = 'none';
      diffViewer.textContent = '';
      currentRawIngestNote = '';
    });
  }

  if (btnCommitIngest) {
    btnCommitIngest.addEventListener('click', async () => {
      try {
        const res = await fetch('/api/ingest/commit', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ content: currentRawIngestNote, target: 'productContext.md' })
        });
        const data = await res.json();
        if (data.success) {
          showToast('Successfully committed note to productContext.md!');
          previewContainer.style.display = 'none';
          loadStatus();
        }
      } catch (e) {
        showToast('Failed to commit note.');
      }
    });
  }

  // 6. Checkpoints Timeline
  const ckptModal = document.getElementById('checkpoint-modal');
  const btnCloseModal = document.getElementById('btn-close-modal');
  const btnModalCancel = document.getElementById('btn-modal-cancel');
  const btnModalSave = document.getElementById('btn-modal-save');
  const btnQuickCheckpoint = document.getElementById('btn-quick-checkpoint');
  const btnAddCheckpoint = document.getElementById('btn-add-checkpoint');

  function openCheckpointModal() {
    ckptModal.classList.add('active');
    document.getElementById('ckpt-focus-input').focus();
  }

  if (btnQuickCheckpoint) btnQuickCheckpoint.addEventListener('click', openCheckpointModal);
  if (btnAddCheckpoint) btnAddCheckpoint.addEventListener('click', openCheckpointModal);
  if (btnCloseModal) btnCloseModal.addEventListener('click', () => ckptModal.classList.remove('active'));
  if (btnModalCancel) btnModalCancel.addEventListener('click', () => ckptModal.classList.remove('active'));

  if (btnModalSave) {
    btnModalSave.addEventListener('click', async () => {
      const focus = document.getElementById('ckpt-focus-input').value.trim();
      const notes = document.getElementById('ckpt-notes-input').value.trim();
      if (!focus) {
        showToast('Please provide an active focus summary.');
        return;
      }

      btnModalSave.disabled = true;
      btnModalSave.textContent = 'Saving...';

      try {
        const res = await fetch('/api/checkpoint/create', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ focus: focus, notes: notes })
        });
        const data = await res.json();
        if (data.success) {
          showToast(`Snapshot ${data.checkpoint_id} created safely.`);
          ckptModal.classList.remove('active');
          document.getElementById('ckpt-focus-input').value = '';
          document.getElementById('ckpt-notes-input').value = '';
          loadStatus();
          if (currentView === 'checkpoints-view') loadCheckpoints();
        }
      } catch (e) {
        showToast('Failed to create snapshot.');
      } finally {
        btnModalSave.disabled = false;
        btnModalSave.textContent = 'Save Snapshot';
      }
    });
  }

  async function loadCheckpoints() {
    const list = document.getElementById('checkpoints-list');
    list.innerHTML = '<div class="text-sm text-muted">Loading snapshots...</div>';

    try {
      const res = await fetch('/api/checkpoints');
      const data = await res.json();
      const items = Array.isArray(data) ? data : (data.checkpoints || []);

      if (items.length === 0) {
        list.innerHTML = '<div class="text-sm text-muted">No snapshots created yet. Click "+ New Snapshot" to save your current state.</div>';
        return;
      }

      list.innerHTML = '';
      items.forEach(c => {
        const row = document.createElement('div');
        row.className = 'ckpt-item';
        row.innerHTML = `
          <div>
            <div class="ckpt-id">${c.id || c.ID}</div>
            <div class="ckpt-message">${c.focus || c.Focus || c.message || 'Snapshot save point'}</div>
          </div>
          <div class="ckpt-date">${c.timestamp || c.Created || 'Recently'}</div>
        `;
        list.appendChild(row);
      });
    } catch (e) {
      list.innerHTML = '<div class="text-sm text-muted">Failed to load checkpoints.</div>';
    }
  }

  // 7. Audit & Lint Buttons
  const btnRunAudit = document.getElementById('btn-run-audit');
  if (btnRunAudit) {
    btnRunAudit.addEventListener('click', async () => {
      showToast('Running project audit...');
      try {
        const res = await fetch('/api/audit');
        const data = await res.json();
        showToast(`Audit complete: ${data.packages || 0} packages, ${data.symbols || 0} symbols found.`);
      } catch (e) {
        showToast('Audit failed.');
      }
    });
  }

  const btnRunLint = document.getElementById('btn-run-lint');
  if (btnRunLint) {
    btnRunLint.addEventListener('click', async () => {
      try {
        const res = await fetch('/api/lint');
        const data = await res.json();
        if (data.compliant) {
          showToast('Memory bank budget check PASSED (all < 150 lines).');
        } else {
          showToast(`Budget warning: ${data.warning || 'Exceeds budget'}`);
        }
      } catch (e) {
        showToast('Lint check failed.');
      }
    });
  }

  // Toast Notification System
  function showToast(message) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.textContent = message;
    container.appendChild(toast);

    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateY(10px)';
      toast.style.transition = 'all 0.3s ease';
      setTimeout(() => toast.remove(), 300);
    }, 3200);
  }
});
