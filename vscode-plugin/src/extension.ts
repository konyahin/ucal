import * as vscode from 'vscode';
import { calculate, openConnection } from './ucal';
import { activateDecorator } from './decorator';

export function activate(context: vscode.ExtensionContext) {
    openConnection(context);
    activateDecorator(context);

    const calculateCommand = vscode.commands.registerCommand('ucal.calculate', async () => {
        const expression = await vscode.window.showInputBox({ title: 'insert expression' });
        if (!expression) {
            return;
        }

        try {
            const res = await calculate(expression);
            vscode.window.showInformationMessage(`${res.low.toFixed(2)} ~ ${res.high.toFixed(2)}`);
        } catch (e) {
            const msg = e instanceof Error ? e.message : String(e);
            vscode.window.showErrorMessage(`ucal: ${msg}`);
        }
    });

    context.subscriptions.push(calculateCommand);
}

export function deactivate() { }
