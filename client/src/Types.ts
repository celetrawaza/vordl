// Game

type Letter = string;
type Identity = Letter[];
type Alphabet = Identity[];

export interface Params {
    length: number,
    maxTries: number,
    letters: Alphabet,
    resetTime: number,
    wordNumber: number,
};

export type Guess = string;

export const GuessStatus_Allowed = ["correct", "present", "wrong"] as const;
export type GuessStatus = (typeof GuessStatus_Allowed)[number];

export interface LetterAnnotated {
    letter: Letter,
    correctness: GuessStatus,
}

export type GuessAnnotated = LetterAnnotated[];

// UI

export const KeyFunction_Allowed = ["letter", "enter", "backspace"] as const;
export type KeyFunction = (typeof KeyFunction_Allowed)[number];

export interface KeyboardDataset extends DOMStringMap {
    letter?: string,
    function: KeyFunction,
}

export interface KeyboardKey extends HTMLElement {
    dataset: KeyboardDataset,
}
