import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {BehaviorSubject, interval, map, Observable, switchMap, tap} from 'rxjs';
import {Session, Vote, WaitForVotesResponse} from "./objects";
import {WaitForPromptResponse} from "./objects";


@Injectable({
  providedIn: 'root'
})
export class GuestService {
  private sessionData$$: BehaviorSubject<Session> = new BehaviorSubject({session_id: '', token: ''});
  private sessionData$: Observable<Session> = this.sessionData$$.asObservable();

  constructor(private http: HttpClient) {
  }

  joinSession(sessionId: string, participantId: string): Observable<Session> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
    });

    const body = {
      participant_id: participantId
    };

    return this.http.post<any>(`http://localhost:8080/session/${sessionId}/join`, body, {
      headers,
      withCredentials: true
    }).pipe(
      map(response => {
        return {session_id: sessionId, token: response.token};
      }),
      tap(sessionData => {
        this.sessionData$$.next(sessionData);
      }),
    );
  }

  waitForPrompt(session: Session, lastPromptId: string): Observable<WaitForPromptResponse> {
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${session.token}`
    };

    const body = {
      last_prompt_id: lastPromptId
    };

    return this.http.post<any>(`http://localhost:8080/session/${session.session_id}/prompt/wait`, body,{
        headers,
        withCredentials: true
    });
  }

  submitVote(session: Session, promptId: string, vote: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${session.token}`
    });

    const body = {
      vote: vote
    };

    return this.http.post<any>(`http://localhost:8080/session/${session.session_id}/prompt/${promptId}/vote`, body, { headers, withCredentials: true });
  }

  getVotes(session: Session, promptId: string): Observable<WaitForVotesResponse> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${session.token}`
    });

    return interval(1000).pipe(
      switchMap(() => this.http.get<any>(`http://localhost:8080/session/${session.session_id}/prompt/${promptId}/watch`, { headers, withCredentials: true })),
    )
  }

  get session$(): Observable<Session> {
    return this.sessionData$;
  }
}
