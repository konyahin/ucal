import * as vscode from 'vscode';
import { spawn, ChildProcess } from 'child_process';
import {
    CancellationTokenSource,
    createMessageConnection,
    MessageConnection,
    RequestType,
    StreamMessageReader,
    StreamMessageWriter,
} from 'vscode-jsonrpc/node';


export interface EvaluateResult { low: number; high: number; }
interface EvaluateParams { expression: string; }
const Evaluate = new RequestType<EvaluateParams, EvaluateResult, void>('evaluate');

let context: vscode.ExtensionContext;
let proc: ChildProcess;
let conn: MessageConnection;
let output: vscode.OutputChannel;

function getConnection(): MessageConnection {
    if (conn) {
        return conn;
    }

    const bin = vscode.workspace.getConfiguration('ucal').get<string>('binaryPath') ?? 'ucal';
    proc = spawn(bin, ['-serve'], { stdio: ['pipe', 'pipe', 'pipe'] });
    context.subscriptions.push({ dispose: () => proc.kill() });

    proc.stderr?.on('data', d => output.append(d.toString()));
    proc.on('error', e => {
        const message = `spawn error: ${e.message}`;
        output.appendLine(message);
        vscode.window.showErrorMessage(message);
    });

    conn = createMessageConnection(
        new StreamMessageReader(proc.stdout!),
        new StreamMessageWriter(proc.stdin!),
    );
    context.subscriptions.push(conn);

    conn.listen();
    return conn;
}

export function openConnection(ctxt: vscode.ExtensionContext) {
    context = ctxt;

    output = vscode.window.createOutputChannel('ucal');
    context.subscriptions.push(output);
}

export async function calculate(expression: string): Promise<EvaluateResult> {
    const source = new CancellationTokenSource();
    const timer = setTimeout(() => {
        source.cancel();
        output.appendLine(`calculation of ${expression} cancelled by timeout`);
    }, 5000);

    try {
        return await getConnection().sendRequest(Evaluate, { expression }, source.token);
    } finally {
        clearTimeout(timer);
        source.dispose();
    }
}