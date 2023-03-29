import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {BehaviorSubject, Observable, tap} from 'rxjs';
import {Session} from "./objects";

@Injectable({
  providedIn: 'root'
})
export class HostService {
  private sessionData$$: BehaviorSubject<Session> = new BehaviorSubject({session_id: '', token: ''});
  private sessionData$: Observable<Session> = this.sessionData$$.asObservable();

  constructor(private http: HttpClient) {
  }

  createSession(participantId: string): Observable<Session> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
    });

    const body = {
      participant_id: participantId
    };

    return this.http.post<any>('http://localhost:8080/session/create', body, {headers, withCredentials: true}).pipe(
      tap(sessionData => {
        this.sessionData$$.next(sessionData);
      }),
    );
  }

  createPrompt(session: Session, promptText: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${session.token}`
    });

    const body = {
      prompt: promptText
    };

    return this.http.post<any>(`http://localhost:8080/session/${session.session_id}/prompt/create`, body, {
      headers,
      withCredentials: true
    });
  }

  get session$(): Observable<Session> {
    return this.sessionData$;
  }
}
