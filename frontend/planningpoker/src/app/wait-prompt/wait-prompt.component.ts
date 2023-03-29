import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {GuestService} from "../guest.service";
import {interval, switchMap, tap} from "rxjs";

@Component({
  selector: 'app-wait-prompt',
  templateUrl: './wait-prompt.component.html',
  styleUrls: ['./wait-prompt.component.scss']
})
export class WaitPromptComponent implements OnInit {
  promptId = '';
  sessionId = '';
  prompt = '';
  vote = '';

  constructor(private route: ActivatedRoute, private service: GuestService, private router: Router) { }

  ngOnInit(): void {
    const sessionId = this.route.snapshot.paramMap.get('sessionId');
    if (sessionId) {
      this.sessionId = sessionId;
    }
    this.pollForPrompt();
  }

  submitVote() {
    this.service.session$.pipe(
      switchMap(session => this.service.submitVote(session, this.promptId, this.vote)),
    ).subscribe(() => {
      this.router.navigateByUrl(`/session/${this.sessionId}/prompt/${this.promptId}/watch`)
    })
  }

  pollForPrompt() {
    this.service.session$.pipe(
      switchMap(session => this.service.waitForPrompt(session, this.promptId)),
    ).subscribe(response => {
      if (response.prompt_id && response.prompt) {
        this.promptId = response.prompt_id;
        this.prompt = response.prompt;
      } else {
        this.pollForPrompt();
      }
    });
  }
}
