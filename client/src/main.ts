document.addEventListener("DOMContentLoaded", main)

import { C, E, gID } from "./DOM";
import { GuessAnnotated, GuessStatus_Allowed, KeyboardDataset, KeyboardKey, Params } from "./Types";
import { API } from "./Comms";

interface GameInterface {
    params: Params,
    tries: GuessAnnotated[],
    resetTimeout: number,
    clockInterval: number,
}

let Game: GameInterface;

function halt() {
    alert("Fatal error happened!")
}

async function getParams() {
    try {
        Game = {
            params: await API.params(),
            tries: [],
            resetTimeout: -1,
            clockInterval: -1,
        };
    } catch (e) {
        console.error("Error during getting params", e);
        halt(); // todo make retriable
    }
}

function populateGuessPrototype() {
    for (let i = 0; i < Game.params.length; i++) {
        E.tablePrototypeRow().appendChild(C.tablePrototypeCellClone())
    }
}

function populateGuessTable() {
    for (let i = 0; i < Game.params.maxTries; i++) {
        E.tableGuesses().appendChild(C.tablePrototypeRowClone())
    }
}

async function getTries() {
    try {
        Game.tries = await API.tries();
        await updateAll();
    } catch (e) {
        console.error("Error during getting tries", e);
        halt(); // todo make retriable
    }
}

function clearTable() {
    E.tableGuesses().innerHTML = "";
    populateGuessTable();
}

function updateGuessTable() {
    clearTable()
    const rows = E.tableGuesses().getElementsByClassName("grid-row");
    for (let i = 0; i < Game.tries.length; i++) {
        const cells = rows[i].getElementsByClassName("grid-cell");
        const guess = Game.tries[i];
        for (let j = 0; j < guess.length; j++) {
            const cell = cells[j] as HTMLElement;
            const {letter, correctness} = guess[j];
            cell.innerHTML = letter;
            cell.classList.add(correctness);
        }
    }
}

async function submitGuess(event: Event) {
    event.preventDefault();
    const fields = E.formFields() as HTMLFieldSetElement;
    fields.disabled = true;
    const guess = (E.inputText() as HTMLInputElement).value;
    try {
        const out = await API.guess(guess);
        Game.tries = out;
        await updateAll();
        (E.inputText() as HTMLInputElement).value = "";
    } catch (e) {
        console.error("Error submitting guess", e);
        alert("Ошибка!\n"+e);
    }
    fields.disabled = false;
}

function keyboardHandler(event: MouseEvent) {
    const cell = event.target as KeyboardKey;
    console.log("Clicked cell with params", cell.dataset);
    switch (cell.dataset.function) {
        case "enter":
            E.inputButton().click();
            break;
        case "backspace":
            const input = (E.inputText() as HTMLInputElement);
            input.value = input.value.slice(0,-1);
            break;
        case "letter":
            (E.inputText() as HTMLInputElement).value += cell.dataset.letter!;
            break;
        default:
            console.error("No action defined for cell", cell);
            break;
    }
}

function generateKeyboard() {
    E.tableKeyboard().innerHTML = "";
    const rows = 3;
    // const rowLimit = Math.ceil(Game.params.letters.length / rows);
    const rowStarts = ["ф", "я"];
    let currentRow = C.tablePrototypeRowClone(false) as HTMLElement;
    for (let i = 0, c = 0; i < Game.params.letters.length; i++, c++) {
        // if (c >= rowLimit)
        const identity = Game.params.letters[i];
        if (identity.some(e=>rowStarts.includes(e)))
        {
            c = 0;
            E.tableKeyboard().appendChild(currentRow);
            currentRow = C.tablePrototypeRowClone(false) as HTMLElement;
        }
        const cell = C.tablePrototypeCellClone(identity.join("/")) as KeyboardKey;
        cell.dataset.function = "letter";
        cell.dataset.letter = identity[0];
        cell.addEventListener("click", keyboardHandler);
        currentRow.append(cell);
    }

    // special keys
    const cellEnter = C.tablePrototypeCellClone("Enter") as KeyboardKey;
    cellEnter.dataset.function = "enter";
    const cellBackspace = C.tablePrototypeCellClone("&lt;-") as KeyboardKey;
    cellBackspace.dataset.function = "backspace";
    [cellEnter, cellBackspace].forEach(c => c.addEventListener("click", keyboardHandler));
    // currentRow.insertAdjacentElement("afterbegin", cellEnter);
    currentRow.prepend(cellEnter);
    currentRow.append(cellBackspace);

    E.tableKeyboard().appendChild(currentRow);
}

function updateKeyboard() {
    const cells = [...E.tableKeyboard().getElementsByClassName("grid-cell")] as KeyboardKey[];
    for (let cell of cells) {
        if (cell.dataset.function !== "letter") continue;
        const letter = cell.dataset.letter!;
        const classes = Game.tries.flat(1).filter(la=>la.letter===letter[0]).map(la=>la.correctness)//.filter(c=>c!=="wrong");
        cell.classList.remove(...GuessStatus_Allowed);
        cell.classList.add(...classes);
    }
}

async function updateAll() {
    updateGuessTable()
    updateKeyboard();
    if (hasFailed() && !hasWon()) { // todo move ailleurs
        // refetch the word
        try {
            const word = await API.word();
            E.wordDetails().hidden = false;
            E.word().innerHTML = word;
        } catch (e) {
            console.error("Error getting word", e);
        }
    }
}

function initTimers() {
    let secondsLeft = Math.floor(Game.params.resetTime-Date.now()/1000);
    Game.resetTimeout = setTimeout(()=>{
        alert("Время вышло!")
        location.reload();
    }, (secondsLeft+2)*1000);
    const formatter = new Intl.DurationFormat("ru", {style: "digital"});
    const setter = ()=>{
        E.resetTime().innerHTML = formatter.format(Temporal.Duration.from({
            seconds: Math.max(0, secondsLeft)
        }).round({largestUnit: "hour"}));
    }
    Game.clockInterval = setInterval(()=>{
        setter();
        secondsLeft--;
    }, 1000);
    setter();
}

function hasWon() {
    return Game.tries.some(ga =>
        !ga.some(la =>
            la.correctness !== "correct"
        )
    );
}

function hasFailed() {
    return Game.tries.length >= Game.params.maxTries;
}

function updateParams() {
    initTimers();
    E.wordNumber().innerHTML = ""+Game.params.wordNumber;
    E.wordLength().innerHTML = ""+Game.params.length;
    E.maxTries().innerHTML = ""+Game.params.maxTries;
}

async function main() {
    // add handlers
    E.inputButton().addEventListener("click", submitGuess)
    // get params of the game
    await getParams(); // todo consider expanding here
    updateParams();
    // set up the interface
    populateGuessPrototype();
    generateKeyboard();
    // get current status of the game
    await getTries();
    // reflect changes
    await updateAll();
}
