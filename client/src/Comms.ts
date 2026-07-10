import { Guess, GuessAnnotated, Params } from "./Types";

interface APIError {
    message?: string,
}

async function extractError(response: Response) {
    const text = await response.text().catch(e=>""+(e??""))
    return text || `HTTP error ${response.status}`;
}

export const API = {
    params: async () => fetch("/api/params").then(r=>r.json()),
    tries: async (): Promise<GuessAnnotated[]> => fetch("/api/tries").then(r=>r.json()),
    guess: async (guess: Guess): Promise<GuessAnnotated[]> => {
        const response = await fetch("/api/guess", {
            method: "POST",
            body: JSON.stringify(guess),
        });
        if (!response.ok) {
            throw new Error(await extractError(response));
        }
        return response.json();
    },
    word: async (): Promise<string> => {
        const response = await fetch("/api/word");
        if (!response.ok) {
            throw new Error(await extractError(response));
        }
        return response.json();
    }
};
