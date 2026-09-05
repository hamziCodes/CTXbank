import * as vscode from 'vscode';
import * as path from 'path';

export function activate(context: vscode.ExtensionContext) {
  const workspaceRoot = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
  if (!workspaceRoot) return;

  // 1. Register Tree Data Providers
  const focusProvider = new ActiveFocusProvider(workspaceRoot);
  const filesProvider = new MemoryBankFilesProvider(workspaceRoot);
  const ckptProvider = new CheckpointsProvider(workspaceRoot);

  vscode.window.registerTreeDataProvider('ctxbank-focus', focusProvider);
  vscode.window.registerTreeDataProvider('ctxbank-files', filesProvider);
  vscode.window.registerTreeDataProvider('ctxbank-checkpoints', ckptProvider);

  // 2. Status Bar Integration
  const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 10);
  statusBarItem.command = 'ctxbank.openDashboard';
  context.subscriptions.push(statusBarItem);
  updateStatusBar(workspaceRoot, statusBarItem);

  // 3. Register Commands
  context.subscriptions.push(
    vscode.commands.registerCommand('ctxbank.refresh', () => {
      focusProvider.refresh();
      filesProvider.refresh();
      ckptProvider.refresh();
      updateStatusBar(workspaceRoot, statusBarItem);
      vscode.window.showInformationMessage('CTXbank state refreshed.');
    }),

    vscode.commands.registerCommand('ctxbank.openDashboard', () => {
      openVisualDashboard(context, workspaceRoot);
    }),

    vscode.commands.registerCommand('ctxbank.createCheckpoint', async () => {
      const focus = await vscode.window.showInputBox({
        prompt: 'Enter active focus for this checkpoint snapshot:',
        placeHolder: 'e.g. Completed auth middleware and unit tests'
      });
      if (!focus) return;

      const notes = await vscode.window.showInputBox({
        prompt: '(Optional) Out-of-band notes or manual fixes:',
        placeHolder: 'e.g. Updated go.mod dependencies manually'
      });

      runCtxCommand(['pause'], focus, notes || '');
      focusProvider.refresh();
      ckptProvider.refresh();
      updateStatusBar(workspaceRoot, statusBarItem);
    }),

    vscode.commands.registerCommand('ctxbank.pause', () => {
      vscode.commands.executeCommand('ctxbank.createCheckpoint');
    }),

    vscode.commands.registerCommand('ctxbank.resume', () => {
      const activePath = path.join(workspaceRoot, 'memory-bank', 'activeContext.md');
      vscode.workspace.openTextDocument(activePath).then(doc => {
        vscode.window.showTextDocument(doc, { preview: true });
      });
    }),

    vscode.commands.registerCommand('ctxbank.lintMemory', () => {
      const terminal = vscode.window.createTerminal('CTXbank');
      terminal.show();
      terminal.sendText('ctx lint-memory');
    }),

    vscode.commands.registerCommand('ctxbank.runAudit', () => {
      const terminal = vscode.window.createTerminal('CTXbank');
      terminal.show();
      terminal.sendText('ctx audit');
    })
  );

  // Watch memory bank files for live refresh
  const watcher = vscode.workspace.createFileSystemWatcher('**/memory-bank/**');
  watcher.onDidChange(() => {
    focusProvider.refresh();
    filesProvider.refresh();
    ckptProvider.refresh();
    updateStatusBar(workspaceRoot, statusBarItem);
  });
  context.subscriptions.push(watcher);
}

// Tree Item Models
class ActiveFocusProvider implements vscode.TreeDataProvider<vscode.TreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<vscode.TreeItem | undefined>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private workspaceRoot: string) {}

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: vscode.TreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(): Promise<vscode.TreeItem[]> {
    const activePath = path.join(this.workspaceRoot, 'memory-bank', 'activeContext.md');
    try {
      const data = await vscode.workspace.fs.readFile(vscode.Uri.file(activePath));
      const content = Buffer.from(data).toString('utf8');
      const lines = content.split('\n');

      let focus = 'No active focus declared';
      let inFocus = false;
      for (const line of lines) {
        const trimmed = line.trim();
        if (trimmed.startsWith('## Focus')) {
          inFocus = true;
          continue;
        } else if (trimmed.startsWith('## ')) {
          inFocus = false;
        }
        if (inFocus && trimmed !== '' && !trimmed.startsWith('#')) {
          focus = trimmed;
          break;
        }
      }

      const focusItem = new vscode.TreeItem(`Focus: ${focus}`, vscode.TreeItemCollapsibleState.None);
      focusItem.iconPath = new vscode.ThemeIcon('target');
      focusItem.tooltip = focus;

      const budgetItem = new vscode.TreeItem(`Budget: ${lines.length} / 150 lines`, vscode.TreeItemCollapsibleState.None);
      budgetItem.iconPath = lines.length >= 150 ? new vscode.ThemeIcon('error') : (lines.length > 120 ? new vscode.ThemeIcon('warning') : new vscode.ThemeIcon('pass'));

      const actionItem = new vscode.TreeItem('Save Checkpoint Snapshot', vscode.TreeItemCollapsibleState.None);
      actionItem.command = { command: 'ctxbank.createCheckpoint', title: 'Save Checkpoint' };
      actionItem.iconPath = new vscode.ThemeIcon('diff-added');

      return [focusItem, budgetItem, actionItem];
    } catch {
      const emptyItem = new vscode.TreeItem('Memory bank not initialized', vscode.TreeItemCollapsibleState.None);
      emptyItem.description = 'Run "ctx init" to start';
      return [emptyItem];
    }
  }
}

class MemoryBankFilesProvider implements vscode.TreeDataProvider<vscode.TreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<vscode.TreeItem | undefined>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private workspaceRoot: string) {}

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: vscode.TreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(): Promise<vscode.TreeItem[]> {
    const files = [
      'projectbrief.md',
      'productContext.md',
      'systemPatterns.md',
      'techContext.md',
      'activeContext.md',
      'progress.md',
      'decisionLog.md'
    ];

    const items: vscode.TreeItem[] = [];
    for (const f of files) {
      const filePath = path.join(this.workspaceRoot, 'memory-bank', f);
      const uri = vscode.Uri.file(filePath);
      try {
        const stat = await vscode.workspace.fs.stat(uri);
        const data = await vscode.workspace.fs.readFile(uri);
        const lines = Buffer.from(data).toString('utf8').split('\n').length;

        const item = new vscode.TreeItem(f, vscode.TreeItemCollapsibleState.None);
        item.description = `${lines} lines (${stat.size}b)`;
        item.command = {
          command: 'vscode.open',
          title: 'Open File',
          arguments: [uri]
        };
        item.iconPath = f === 'activeContext.md' ? new vscode.ThemeIcon('flame') : new vscode.ThemeIcon('file-text');
        items.push(item);
      } catch {
        // File not created yet
      }
    }
    return items;
  }
}

class CheckpointsProvider implements vscode.TreeDataProvider<vscode.TreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<vscode.TreeItem | undefined>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private workspaceRoot: string) {}

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: vscode.TreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(): Promise<vscode.TreeItem[]> {
    const ckptDir = path.join(this.workspaceRoot, 'memory-bank', '.state', 'checkpoints');
    try {
      const entries = await vscode.workspace.fs.readDirectory(vscode.Uri.file(ckptDir));
      const items: vscode.TreeItem[] = [];

      for (const [name, type] of entries.reverse()) {
        if (type === vscode.FileType.File && name.endsWith('.json')) {
          const filePath = path.join(ckptDir, name);
          const uri = vscode.Uri.file(filePath);
          const data = await vscode.workspace.fs.readFile(uri);
          const parsed = JSON.parse(Buffer.from(data).toString('utf8'));

          const item = new vscode.TreeItem(parsed.id || name, vscode.TreeItemCollapsibleState.None);
          item.description = parsed.active_focus || 'snapshot';
          item.tooltip = `${parsed.timestamp} • branch: ${parsed.branch}`;
          item.iconPath = new vscode.ThemeIcon('archive');
          item.command = {
            command: 'vscode.open',
            title: 'Open Snapshot JSON',
            arguments: [uri]
          };
          items.push(item);
        }
      }
      return items.length > 0 ? items : [new vscode.TreeItem('No checkpoints recorded yet', vscode.TreeItemCollapsibleState.None)];
    } catch {
      return [new vscode.TreeItem('No checkpoints found', vscode.TreeItemCollapsibleState.None)];
    }
  }
}

// Visual Webview Dashboard (VERTEX Design)
function openVisualDashboard(context: vscode.ExtensionContext, workspaceRoot: string) {
  const panel = vscode.window.createWebviewPanel(
    'ctxbankDashboard',
    'CTXbank: Architecture Tree',
    vscode.ViewColumn.One,
    { enableScripts: true }
  );

  panel.webview.html = getWebviewContent(workspaceRoot);
}

function getWebviewContent(workspaceRoot: string): string {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <style>
    :root {
      --brand-primary: #042940;
      --accent-text: #609abe;
      --border-width: 1.5px;
      --shadow-slab: 2px 2px 0 var(--vscode-widget-border, #30363d);
    }
    body {
      font-family: var(--vscode-font-family);
      background: var(--vscode-editor-background);
      color: var(--vscode-editor-foreground);
      padding: 24px;
      margin: 0;
    }
    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      border-bottom: var(--border-width) solid var(--vscode-widget-border, #30363d);
      padding-bottom: 16px;
      margin-bottom: 24px;
    }
    .title {
      font-size: 18px;
      font-weight: 700;
    }
    .tree-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 16px;
    }
    .node-card {
      background: var(--vscode-sideBar-background, #161b22);
      border: var(--border-width) solid var(--vscode-widget-border, #30363d);
      border-radius: 6px;
      padding: 16px;
      box-shadow: var(--shadow-slab);
      transition: transform 120ms ease;
    }
    .node-card:hover {
      transform: translate(-1px, -1px);
      border-color: var(--accent-text);
    }
    .node-title {
      font-weight: 700;
      font-size: 14px;
      margin-bottom: 4px;
    }
    .node-sub {
      font-size: 12px;
      color: var(--vscode-descriptionForeground, #8b949e);
    }
    .hot-badge {
      display: inline-block;
      font-size: 10px;
      font-weight: 700;
      padding: 2px 6px;
      background: #609abe;
      color: #000000;
      border-radius: 3px;
      margin-top: 8px;
    }
  </style>
</head>
<body>
  <div class="header">
    <div class="title">CTXbank Architecture & Memory Flow</div>
    <div style="font-size: 12px; color: var(--vscode-descriptionForeground);">Interactive Node-Link Map</div>
  </div>

  <div class="tree-grid">
    <div class="node-card">
      <div class="node-title">projectbrief.md</div>
      <div class="node-sub">Core Mission & Scope</div>
    </div>
    <div class="node-card">
      <div class="node-title">productContext.md</div>
      <div class="node-sub">User Experience & Ingested Notes</div>
    </div>
    <div class="node-card">
      <div class="node-title">systemPatterns.md</div>
      <div class="node-sub">Architectural Decisions & Rules</div>
    </div>
    <div class="node-card">
      <div class="node-title">techContext.md</div>
      <div class="node-sub">Tech Stack & Dependencies</div>
    </div>
    <div class="node-card" style="border-color: var(--accent-text);">
      <div class="node-title">activeContext.md</div>
      <div class="node-sub">Active Focus & Next Steps (&lt;150 lines)</div>
      <span class="hot-badge">HOT CONTEXT</span>
    </div>
    <div class="node-card">
      <div class="node-title">progress.md</div>
      <div class="node-sub">Milestone Completion Ledger</div>
    </div>
    <div class="node-card">
      <div class="node-title">decisionLog.md</div>
      <div class="node-sub">Architectural Decision Records</div>
    </div>
  </div>
</body>
</html>`;
}

async function updateStatusBar(workspaceRoot: string, item: vscode.StatusBarItem) {
  const activePath = path.join(workspaceRoot, 'memory-bank', 'activeContext.md');
  try {
    const data = await vscode.workspace.fs.readFile(vscode.Uri.file(activePath));
    const lines = Buffer.from(data).toString('utf8').split('\n').length;
    item.text = `$(archive) CTX: ${lines}/150 lines`;
    item.tooltip = `CTXbank Active Context Budget: ${lines}/150 lines (Click to open visual tree)`;
    item.show();
  } catch {
    item.hide();
  }
}

function runCtxCommand(args: string[], focus: string, notes: string) {
  const terminal = vscode.window.createTerminal('CTXbank');
  terminal.show();
  terminal.sendText(`ctx pause`);
}

export function deactivate() {}
