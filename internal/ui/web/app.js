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
  if (btnTheme) {
    btnTheme.addEventListener('click', () => {
      const isDark = document.documentElement.classList.toggle('dark');
      localStorage.setItem('ctx_theme', isDark ? 'dark' : 'light');
      if (currentView === 'graph-view') loadGraph();
    });
  }

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

  // Initial Data Load
  loadStatus();
  loadFiles();
  loadGraph();
  setInterval(loadStatus, 10000);

  // 1. Status Loading
  async function loadStatus() {
    try {
      const res = await fetch('/api/status');
      const data = await res.json();

      // Header Bar
      const repoNameEl = document.getElementById('repo-name');
      if (repoNameEl) repoNameEl.textContent = data.project_name || 'Project';

      const branchNameEl = document.getElementById('branch-name');
      if (branchNameEl) branchNameEl.textContent = data.branch || 'main';

      const gitStatusPill = document.getElementById('git-status-pill');
      const gitStatusText = document.getElementById('git-status-text');
      if (gitStatusPill && gitStatusText) {
        if (data.dirty_count > 0) {
          gitStatusText.textContent = `${data.dirty_count} modified`;
          gitStatusPill.className = 'status-pill dirty';
        } else {
          gitStatusText.textContent = 'Clean';
          gitStatusPill.className = 'status-pill clean';
        }
      }

      // Check if project has memory bank
      const uninitBanner = document.getElementById('uninit-banner');
      if (uninitBanner) {
        if (data.has_bank === false) {
          uninitBanner.style.display = 'flex';
          const heroFiles = document.getElementById('hero-files-count');
          if (heroFiles) heroFiles.textContent = '0';
        } else {
          uninitBanner.style.display = 'none';
          const heroFiles = document.getElementById('hero-files-count');
          if (heroFiles) heroFiles.textContent = '7';
        }
      }

      // Overview Tab Hero & Metric Cards
      const focusText = document.getElementById('overview-focus-text');
      if (focusText) focusText.textContent = data.active_focus || 'Ready for next task.';

      const branchVal = document.getElementById('overview-branch-val');
      if (branchVal) branchVal.textContent = data.branch || 'main';

      const treeVal = document.getElementById('overview-tree-val');
      if (treeVal) treeVal.textContent = data.dirty_count > 0 ? `${data.dirty_count} modified` : 'Clean';

      const budgetLines = data.active_lines || 0;
      const budgetText = document.getElementById('overview-budget-text');
      if (budgetText) budgetText.textContent = `${budgetLines} / 150 lines`;

      const budgetPct = Math.min(100, Math.round((budgetLines / 150) * 100));
      const heroBudget = document.getElementById('hero-budget-pct');
      if (heroBudget) heroBudget.textContent = budgetPct + '%';

      const budgetBar = document.getElementById('overview-budget-bar');
      if (budgetBar) {
        budgetBar.style.width = budgetPct + '%';
        budgetBar.className = 'budget-bar ' + (data.budget_status || 'normal');
      }

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
    if (!svg) return;
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
    if (!svg) return;
    svg.innerHTML = '';
    const width = svg.clientWidth || 700;
    const height = svg.clientHeight || 520;

    if (!nodes || nodes.length === 0) {
      svg.innerHTML = `
        <text x="50%" y="45%" text-anchor="middle" fill="currentColor" font-size="14" font-weight="600">No architecture nodes yet</text>
        <text x="50%" y="55%" text-anchor="middle" fill="gray" font-size="12">Click "Sync Context" or "Create CTXbank" to generate your architecture tree</text>
      `;
      return;
    }

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
      const nid = n.id || n.ID;
      if (n.type === 'checkpoint') {
        coords[nid] = { x: width * 0.65 + (ckptOffset * 35), y: height * 0.88 };
        ckptOffset++;
      }
    });

    const isDark = document.documentElement.classList.contains('dark');
    const edgeColor = isDark ? '#33445b' : '#cbd5e1';
    const textColor = isDark ? '#f8fafc' : '#0f172a';
    const subTextColor = isDark ? '#94a3b8' : '#64748b';
    const rectFill = isDark ? '#111720' : '#ffffff';
    const rectStroke = isDark ? '#222e3e' : '#e2e8f0';

    // Render Edges
    edges.forEach(edge => {
      const srcId = edge.source || edge.Source;
      const tgtId = edge.target || edge.Target;
      const src = coords[srcId];
      const tgt = coords[tgtId];
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

    // Render Nodes
    nodes.forEach(node => {
      const nodeId = node.id || node.ID;
      const nodeName = node.name || node.Name || node.label || node.Label || nodeId;
      const nodeLines = (node.lines !== undefined) ? node.lines : node.Lines;
      const pos = coords[nodeId] || { x: width * 0.5, y: height * 0.5 };

      const group = document.createElementNS('http://www.w3.org/2000/svg', 'g');
      group.setAttribute('class', 'node-group');
      group.setAttribute('transform', `translate(${pos.x - 75}, ${pos.y - 26})`);

      const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
      rect.setAttribute('width', '150');
      rect.setAttribute('height', '52');
      rect.setAttribute('rx', '10');
      rect.setAttribute('ry', '10');
      rect.setAttribute('fill', rectFill);
      rect.setAttribute('stroke', nodeId === 'activeContext' ? (isDark ? '#609abe' : '#042940') : rectStroke);
      rect.setAttribute('stroke-width', nodeId === 'activeContext' ? '2' : '1.5');
      group.appendChild(rect);

      // Title
      const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      text.setAttribute('x', '14');
      text.setAttribute('y', '22');
      text.setAttribute('fill', textColor);
      text.setAttribute('font-size', '12');
      text.setAttribute('font-weight', '600');
      text.setAttribute('font-family', 'Inter, sans-serif');
      text.textContent = nodeName;
      group.appendChild(text);

      // Subtitle / Line Count
      const sub = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      sub.setAttribute('x', '14');
      sub.setAttribute('y', '39');
      sub.setAttribute('fill', subTextColor);
      sub.setAttribute('font-size', '11');
      sub.setAttribute('font-family', 'JetBrains Mono, monospace');
      sub.textContent = nodeLines !== undefined ? `${nodeLines} lines` : (node.type || 'artifact');
      group.appendChild(sub);

      group.addEventListener('click', () => inspectNode(node));
      svg.appendChild(group);
    });
  }

  function inspectNode(node) {
    const nameEl = document.getElementById('inspector-name');
    const typeEl = document.getElementById('inspector-type');
    const descEl = document.getElementById('inspector-desc');
    const detailsEl = document.getElementById('inspector-details');
    const pathEl = document.getElementById('inspector-path');
    const linesEl = document.getElementById('inspector-lines');
    const btnEdit = document.getElementById('btn-inspector-edit');

    const nodeId = node.id || node.ID;
    const nodeName = node.name || node.Name || node.label || node.Label || nodeId;
    const filename = nodeId.endsWith('.md') ? nodeId : nodeId + '.md';
    const nodeLines = (node.lines !== undefined) ? node.lines : node.Lines;

    if (nameEl) nameEl.textContent = nodeName;
    if (typeEl) typeEl.textContent = node.type || 'Memory Artifact';
    if (descEl) descEl.textContent = fileDescriptions[filename] || 'A core architectural artifact tracked by CTXbank.';

    if (detailsEl) detailsEl.style.display = 'block';
    if (pathEl) pathEl.textContent = `memory-bank/${filename}`;
    if (linesEl) linesEl.textContent = nodeLines !== undefined ? `${nodeLines} lines` : 'N/A';

    if (btnEdit) btnEdit.onclick = () => openEditor(filename);
  }

  // 3. Memory Bank Grid
  async function loadFiles() {
    const grid = document.getElementById('memory-cards-grid');
    if (!grid) return;
    grid.innerHTML = '<div class="text-sm text-muted">Loading memory bank files...</div>';

    try {
      const res = await fetch('/api/files');
      const files = await res.json();
      grid.innerHTML = '';

      if (!files || files.length === 0) {
        grid.innerHTML = `
          <div class="card p-6" style="grid-column: 1 / -1; text-align: center;">
            <div style="font-size: 32px; margin-bottom: 8px;">✨</div>
            <h3 class="font-bold text-sm">No Memory Bank Files Detected</h3>
            <p class="text-xs text-muted" style="margin: 8px 0 16px;">This project has not been initialized with CTXbank yet. Click below to create your 7 core memory bank files and scan your codebase.</p>
            <button class="btn btn-primary" id="btn-empty-create-bank">✨ Create Memory Bank & Scan Code</button>
          </div>
        `;
        const btnEmpty = document.getElementById('btn-empty-create-bank');
        if (btnEmpty) btnEmpty.onclick = () => syncCodebaseContext();
        return;
      }

      files.forEach(f => {
        const fname = f.filename || f.name;
        const lineCount = f.lines !== undefined ? f.lines : (f.line_count || 0);
        const byteSize = f.bytes !== undefined ? f.bytes : (f.byte_size || 0);
        const roleDesc = fileDescriptions[fname] || 'Core project knowledge file.';
        const budgetWarning = (fname === 'activeContext.md' && lineCount >= 150);

        const card = document.createElement('div');
        card.className = 'memory-file-card';
        card.innerHTML = `
          <div class="card-top-row">
            <span class="file-name-title">${fname}</span>
            <span class="chip ${budgetWarning ? 'dirty' : 'clean'}">${lineCount} lines</span>
          </div>
          <p class="file-friendly-role">${roleDesc}</p>
          <div class="file-meta-row">
            <span>Size: ${(byteSize / 1024).toFixed(1)} KB</span>
            <span class="font-mono text-xs">Click to edit &rarr;</span>
          </div>
        `;

        card.addEventListener('click', () => openEditor(fname));
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
    if (editorFilename) editorFilename.textContent = filename;
    if (editorTextarea) editorTextarea.value = 'Loading file contents...';
    if (editorModal) editorModal.classList.add('active');

    try {
      const res = await fetch(`/api/file?name=${encodeURIComponent(filename)}`);
      const data = await res.json();
      if (editorTextarea) editorTextarea.value = data.content || '';
      updateEditorBudget();
    } catch (err) {
      if (editorTextarea) editorTextarea.value = 'Error loading file content.';
    }
  }

  function updateEditorBudget() {
    if (!editorTextarea || !editorBudget) return;
    const lines = editorTextarea.value.split('\n').length;
    editorBudget.textContent = `${lines} lines`;
    if (activeEditingFile === 'activeContext.md') {
      editorBudget.textContent = `${lines} / 150 lines`;
      editorBudget.className = lines > 150 ? 'chip dirty' : 'chip clean';
    } else {
      editorBudget.className = 'chip';
    }
  }

  if (editorTextarea) editorTextarea.addEventListener('input', updateEditorBudget);
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
        if (result.status === 'saved' || result.success) {
          showToast(`Saved ${activeEditingFile} atomically.`);
          editorModal.classList.remove('active');
          loadStatus();
          loadFiles();
          loadGraph();
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

  if (dropzone && fileInput) {
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

      if (previewContainer) previewContainer.style.display = 'block';
      if (diffViewer) diffViewer.textContent = data.diff || '(No new content detected or all duplicate)';
      if (dedupBadge) dedupBadge.textContent = `${data.duplicate_count || 0} duplicate paragraphs removed`;
      showToast('SimHash deduplication completed.');
    } catch (e) {
      showToast('Failed to preview note ingestion.');
    }
  }

  if (btnCancelIngest) {
    btnCancelIngest.addEventListener('click', () => {
      if (previewContainer) previewContainer.style.display = 'none';
      if (diffViewer) diffViewer.textContent = '';
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
        if (data.status === 'committed' || data.success) {
          showToast('Successfully committed note to productContext.md!');
          if (previewContainer) previewContainer.style.display = 'none';
          loadStatus();
          loadFiles();
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
    if (!ckptModal) return;
    ckptModal.classList.add('active');
    const input = document.getElementById('ckpt-focus-input');
    if (input) input.focus();
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
        if (data.id || data.ID || data.success) {
          showToast(`Snapshot created safely.`);
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
    if (!list) return;
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

  // 7. Sync Context & 1-Click Init
  const btnSyncContext = document.getElementById('btn-sync-context');
  if (btnSyncContext) {
    btnSyncContext.addEventListener('click', () => syncCodebaseContext());
  }

  const btnCreateBankOverview = document.getElementById('btn-create-bank-overview');
  if (btnCreateBankOverview) {
    btnCreateBankOverview.addEventListener('click', () => syncCodebaseContext());
  }

  async function syncCodebaseContext() {
    showToast('Syncing codebase with Memory Bank...');
    try {
      const res = await fetch('/api/sync', { method: 'POST' });
      const data = await res.json();
      if (data.success) {
        showToast(`Synced! Detected ${data.packages || 0} packages and ${data.symbols || 0} symbols.`);
        loadStatus();
        loadFiles();
        loadGraph();
      } else {
        showToast('Sync failed: ' + (data.error || 'unknown'));
      }
    } catch (e) {
      showToast('Error syncing context.');
    }
  }

  // 8. Project Selector & Multi-Workspace Manager
  const btnProjectSelector = document.getElementById('btn-project-selector');
  const projectModal = document.getElementById('project-modal');
  const btnCloseProjectModal = document.getElementById('btn-close-project-modal');
  const btnCancelProjectModal = document.getElementById('btn-cancel-project-modal');
  const discoveredProjectsList = document.getElementById('discovered-projects-list');
  const btnInspectCustomFolder = document.getElementById('btn-inspect-custom-folder');
  const inputCustomFolderPath = document.getElementById('input-custom-folder-path');
  const folderInspectResult = document.getElementById('folder-inspect-result');

  if (btnProjectSelector && projectModal) {
    btnProjectSelector.addEventListener('click', () => {
      projectModal.classList.add('active');
      loadProjectsModal();
    });
  }

  if (btnCloseProjectModal && projectModal) {
    btnCloseProjectModal.addEventListener('click', () => projectModal.classList.remove('active'));
  }
  if (btnCancelProjectModal && projectModal) {
    btnCancelProjectModal.addEventListener('click', () => projectModal.classList.remove('active'));
  }

  async function loadProjectsModal() {
    if (!discoveredProjectsList) return;
    try {
      const res = await fetch('/api/projects');
      const data = await res.json();

      const activeProjName = document.getElementById('modal-active-project-name');
      const activeProjPath = document.getElementById('modal-active-project-path');
      if (activeProjName) activeProjName.textContent = data.current_project;
      if (activeProjPath) activeProjPath.textContent = data.current_path;

      discoveredProjectsList.innerHTML = '';
      const list = data.discovered || [];
      if (list.length === 0) {
        discoveredProjectsList.innerHTML = '<div class="text-xs text-muted p-2">No other CTXbank projects found in sibling folders. Enter any folder path below to open or scan it.</div>';
        return;
      }

      list.forEach(p => {
        const isCurrent = (p.path === data.current_path);
        const item = document.createElement('div');
        item.className = `project-card-item ${isCurrent ? 'is-current' : ''}`;
        item.innerHTML = `
          <div>
            <div class="font-bold text-sm flex items-center gap-2">
              <span>${p.name}</span>
              ${isCurrent ? '<span class="chip clean text-xs">Current</span>' : ''}
            </div>
            <div class="text-xs text-muted font-mono" style="margin-top: 2px;">${p.path}</div>
          </div>
          <div>
            ${isCurrent ? '<span class="text-xs text-muted">Active</span>' : `<button class="btn btn-sm btn-secondary" onclick="window.ctxSwitchProject('${encodeURIComponent(p.path)}')">Switch</button>`}
          </div>
        `;
        discoveredProjectsList.appendChild(item);
      });
    } catch (e) {
      discoveredProjectsList.innerHTML = '<div class="text-xs text-muted">Failed to discover projects.</div>';
    }
  }

  // Global window handler for inline onclick
  window.ctxSwitchProject = async function(encodedPath) {
    const path = decodeURIComponent(encodedPath);
    await executeSwitchProject(path);
  };

  async function executeSwitchProject(path) {
    showToast(`Switching workspace to ${path}...`);
    try {
      const res = await fetch('/api/project/switch', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: path })
      });
      const data = await res.json();
      if (data.success) {
        if (projectModal) projectModal.classList.remove('active');
        showToast(`Workspace switched to ${data.name}!`);
        loadStatus();
        loadFiles();
        loadGraph();
        loadCheckpoints();
      } else {
        showToast('Failed to switch: ' + (data.error || 'unknown'));
      }
    } catch (e) {
      showToast('Error switching workspace.');
    }
  }

  if (btnInspectCustomFolder && inputCustomFolderPath && folderInspectResult) {
    btnInspectCustomFolder.addEventListener('click', async () => {
      const path = inputCustomFolderPath.value.trim();
      if (!path) {
        showToast('Please enter a folder path.');
        return;
      }

      btnInspectCustomFolder.disabled = true;
      btnInspectCustomFolder.textContent = 'Inspecting...';
      folderInspectResult.style.display = 'block';
      folderInspectResult.innerHTML = '<div class="text-xs text-muted">Inspecting directory structure and CTXbank state...</div>';

      try {
        const res = await fetch(`/api/project/inspect?path=${encodeURIComponent(path)}`);
        const data = await res.json();

        if (data.error) {
          folderInspectResult.innerHTML = `<div class="text-xs" style="color: var(--negative);">Error: ${data.error}</div>`;
          return;
        }

        if (data.has_bank) {
          folderInspectResult.innerHTML = `
            <div class="flex-between">
              <div>
                <div class="font-bold text-sm">${data.name}</div>
                <div class="text-xs text-muted font-mono">${data.path}</div>
                <div class="text-xs text-muted" style="margin-top: 4px;">Git Branch: <strong>${data.branch || 'unknown'}</strong> | Status: <strong>${data.dirty_count} modified</strong></div>
              </div>
              <button class="btn btn-primary btn-sm" id="btn-open-inspected-project">Open This Project</button>
            </div>
          `;
          document.getElementById('btn-open-inspected-project').addEventListener('click', () => {
            executeSwitchProject(data.path);
          });
        } else {
          folderInspectResult.innerHTML = `
            <div class="flex-between" style="align-items: flex-start; gap: 16px;">
              <div>
                <div class="font-bold text-sm flex items-center gap-2">
                  <span>${data.name}</span>
                  <span class="chip" style="background: var(--attention-soft); color: var(--attention);">Uninitialized Codebase</span>
                </div>
                <div class="text-xs text-muted font-mono" style="margin-top: 2px;">${data.path}</div>
                <p class="text-xs text-muted" style="margin-top: 6px; line-height: 1.5;">
                  This project has not been initialized with CTXbank yet. Clicking below will create its memory bank and automatically run AST reconnaissance to detect its tech stack, frameworks, and architecture rules.
                </p>
              </div>
              <button class="btn btn-primary" id="btn-init-inspected-project" style="white-space: nowrap;">✨ Scan & Initialize</button>
            </div>
          `;
          document.getElementById('btn-init-inspected-project').addEventListener('click', async () => {
            const btn = document.getElementById('btn-init-inspected-project');
            btn.disabled = true;
            btn.textContent = 'Scanning Codebase...';
            showToast('Initializing CTXbank and scanning code...');

            try {
              const initRes = await fetch('/api/project/init', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ path: data.path })
              });
              const initData = await initRes.json();
              if (initData.success) {
                if (projectModal) projectModal.classList.remove('active');
                showToast(`Project ${initData.name} initialized and scanned!`);
                loadStatus();
                loadFiles();
                loadGraph();
                loadCheckpoints();
              } else {
                showToast('Initialization failed: ' + (initData.error || 'unknown'));
              }
            } catch (err) {
              showToast('Error during initialization.');
            }
          });
        }

      } catch (err) {
        folderInspectResult.innerHTML = '<div class="text-xs" style="color: var(--negative);">Failed to inspect folder. Make sure the path exists.</div>';
      } finally {
        btnInspectCustomFolder.disabled = false;
        btnInspectCustomFolder.textContent = 'Inspect Folder';
      }
    });
  }

  // 9. Diagnostics Buttons
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
    if (!container) return;
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
