import * as vscode from 'vscode';
import MarkdownIt from 'markdown-it';
import { calculate } from './ucal';

type State = {
    timer?: NodeJS.Timeout;
    source?: vscode.CancellationTokenSource;
};

type LineDecoration = { line: number; ok: boolean; text: string };

const states = new Map<string, State>();

const DEBOUNCE_MS = 300;
const md = new MarkdownIt();

const resultType = vscode.window.createTextEditorDecorationType({
    after: {
        color: new vscode.ThemeColor('editorCodeLens.foreground'),
        margin: '0 0 0 1rem',
    },
});

const errorType = vscode.window.createTextEditorDecorationType({
    after: {
        color: new vscode.ThemeColor('errorForeground'),
        margin: '0 0 0 1rem',
    },
});

function getState(uri: vscode.Uri): State {
    const key = uri.toString();
    let state = states.get(key);
    if (!state) {
        state = {};
        states.set(key, state);
    }
    return state;
}

function cancelPending(state: State): void {
    if (state.timer) {
        clearTimeout(state.timer);
        state.timer = undefined;
    }
    if (state.source) {
        state.source.cancel();
        state.source.dispose();
        state.source = undefined;
    }
}

function scheduleUpdate(doc: vscode.TextDocument): void {
    if (doc.languageId !== 'markdown') {
        return;
    }

    const state = getState(doc.uri);
    cancelPending(state);

    state.timer = setTimeout(() => {
        state.timer = undefined;
        state.source = new vscode.CancellationTokenSource();
        void updateDecorations(doc, state.source.token);
    }, DEBOUNCE_MS);
}

function isSkippable(line: string): boolean {
    const trimmed = line.trim();
    return trimmed === '' || trimmed.startsWith('#') || trimmed.startsWith('//');
}

async function evaluateLine(
    line: number,
    expression: string,
    token: vscode.CancellationToken,
): Promise<LineDecoration | null> {
    try {
        const res = await calculate(expression, token);
        return {
            line,
            ok: true,
            text: `= ${res.low.toFixed(2)} ~ ${res.high.toFixed(2)}`,
        };
    } catch (err) {
        if (token.isCancellationRequested) {
            return null;
        }
        const msg = err instanceof Error ? err.message : String(err);
        return { line, ok: false, text: `error: ${msg}` };
    }
}

async function updateDecorations(
    doc: vscode.TextDocument,
    token: vscode.CancellationToken,
): Promise<void> {
    const docVersion = doc.version;
    const tokens = md.parse(doc.getText(), {});

    const pending: Promise<LineDecoration | null>[] = [];

    for (const t of tokens) {
        if (t.type !== 'fence' || !t.map) {
            continue;
        }
        const lang = t.info.trim().split(/\s+/)[0];
        if (lang !== 'ucal') {
            continue;
        }
        const [start, end] = t.map;
        for (let i = start + 1; i < end - 1; i++) {
            if (i >= doc.lineCount) {
                break;
            }
            const line = doc.lineAt(i).text;
            if (isSkippable(line)) {
                continue;
            }
            pending.push(evaluateLine(i, line.trim(), token));
        }
    }

    const results = await Promise.all(pending);
    if (token.isCancellationRequested || doc.version !== docVersion) {
        return;
    }

    const okOptions: vscode.DecorationOptions[] = [];
    const errOptions: vscode.DecorationOptions[] = [];

    for (const r of results) {
        if (!r) {
            continue;
        }
        if (r.line >= doc.lineCount) {
            continue;
        }
        const lineLength = doc.lineAt(r.line).text.length;
        const range = new vscode.Range(r.line, lineLength, r.line, lineLength);
        const opt: vscode.DecorationOptions = {
            range,
            renderOptions: { after: { contentText: r.text } },
        };
        (r.ok ? okOptions : errOptions).push(opt);
    }

    for (const editor of vscode.window.visibleTextEditors) {
        if (editor.document !== doc) {
            continue;
        }
        editor.setDecorations(resultType, okOptions);
        editor.setDecorations(errorType, errOptions);
    }
}

export function activateDecorator(ctx: vscode.ExtensionContext): void {
    ctx.subscriptions.push(resultType, errorType);

    ctx.subscriptions.push(
        vscode.workspace.onDidOpenTextDocument(doc => scheduleUpdate(doc)),
        vscode.workspace.onDidChangeTextDocument(e => scheduleUpdate(e.document)),
        vscode.window.onDidChangeVisibleTextEditors(editors => {
            for (const editor of editors) {
                scheduleUpdate(editor.document);
            }
        }),
        vscode.workspace.onDidCloseTextDocument(doc => {
            const key = doc.uri.toString();
            const state = states.get(key);
            if (state) {
                cancelPending(state);
                states.delete(key);
            }
        }),
        {
            dispose: () => {
                for (const state of states.values()) {
                    cancelPending(state);
                }
                states.clear();
            },
        },
    );

    for (const editor of vscode.window.visibleTextEditors) {
        scheduleUpdate(editor.document);
    }
}
