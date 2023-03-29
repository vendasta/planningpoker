
export interface Session {
  session_id: string;
  token: string;
}

export interface WaitForPromptResponse {
  prompt_id: string;
  prompt: string;
}

export interface JoinSessionResponse {
  token: string;
}

export interface WaitForVotesResponse {
  votes: Vote[];
}

export interface Vote {
  participant_id: string;
  vote: string;
}
