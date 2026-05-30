import * as vscode from 'vscode';
import MarkdownIt from 'markdown-it';
import { UcalExpression } from './ucal';

const md = new MarkdownIt();

export function findUcalLines(doc: vscode.TextDocument): UcalExpression[] {
    const tokens = md.parse(doc.getText(), {});
    const expressions: UcalExpression[] = [];

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
            expressions.push({ line: i, text: line.trim() });
        }
    }

    return expressions;
}

function isSkippable(line: string): boolean {
    const trimmed = line.trim();
    return trimmed === '' || trimmed.startsWith('#') || trimmed.startsWith('//');
}
