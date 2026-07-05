import { GuessAnnotated, GuessStatus, LetterAnnotated, Params } from "./Types";

export const PhonyParams: Params = {
    length: 5,
    maxTries: 7,
    letters: [["а"],["б"],["в"],["г"],["д"],["е","ё"],["ж"],["з"],["и"],["й"],["к"],["л"],["м"],["н"],["о"],["п"],["р"],["с"],["т"],["у"],["ф"],["х"],["ц"],["ч"],["ш"],["щ"],["ъ","ь"],["ы"],["э"],["ю"],["я"]],
    resetTime: Date.now()/1000+1000,
    wordNumber: 10,
}

function makeLetterAnnotated(): LetterAnnotated {
    return {
        letter: PhonyParams.letters[Math.floor(Math.random()*PhonyParams.letters.length)][0],
        correctness: (["wrong", "present", "correct"] as GuessStatus[])[Math.floor(Math.random()*3)],
    };
}

function makeGuessAnnotated(): GuessAnnotated {
    const out: GuessAnnotated = [];
    for (let i = 0; i < PhonyParams.length; i++) {
        out.push(makeLetterAnnotated());
    }
    return out;
}

export const PhonyTries: GuessAnnotated[] = []
for (let i = 0; i < PhonyParams.maxTries/2; i++) {
    PhonyTries.push(makeGuessAnnotated());
}
