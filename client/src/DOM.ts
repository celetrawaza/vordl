export const [gID, gClass, cE,] = [
    (id: string) => document.getElementById(id)!,
    (c: string) => document.getElementsByClassName(c),
    (name: string) => document.createElement(name),
]

export const E = {
    wordNumber: ()=>gID("wordNumber"),
    tablePrototypeRow: ()=>gID("tablePrototypes").children[0],
    tablePrototypeCell: ()=>gID("tablePrototypes").children[1],
    tableGuesses: ()=>gID("tableGuesses"),
    inputButton: ()=>gID("inputButton"),
    inputText: ()=>gID("inputText"),
    tableKeyboard: ()=>gID("tableKeyboard"),
    resetTime: ()=>gID("resetTime"),
    wordDetails: ()=>gID("wordDetails"),
    word: ()=>gID("word"),
    formFields: ()=>gID("formFields"),
    wordLength: ()=>gID("wordLength"),
    maxTries: ()=>gID("maxTries"),
}

export const C = {
    tablePrototypeCellClone: (content?: string) => {
        const out = E.tablePrototypeCell().cloneNode(true);
        (out as HTMLElement).innerHTML = content ?? "";
        return out;
    },
    tablePrototypeRowClone: (deep = true) => {
        return E.tablePrototypeRow().cloneNode(deep);
    },
}

