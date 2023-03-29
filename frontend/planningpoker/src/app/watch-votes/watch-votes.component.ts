import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import {Vote} from "../objects";
import {GuestService} from "../guest.service";
import {switchMap} from "rxjs";

@Component({
  selector: 'app-watch-votes',
  templateUrl: './watch-votes.component.html',
  styleUrls: ['./watch-votes.component.scss']
})
export class WatchVotesComponent implements OnInit {
  sessionId = '';
  promptId = '';
  votes: Vote[] = []

  constructor(
    private route: ActivatedRoute,
    private service: GuestService
  ) { }

  ngOnInit(): void {
    const sessionId = this.route.snapshot.paramMap.get('sessionId');
    if (sessionId) {
      this.sessionId = sessionId;
    }
    const promptId = this.route.snapshot.paramMap.get('promptId');
    if (promptId) {
      this.promptId = promptId;
    }

    this.getVotes();
  }

  getVotes(): void {
    this.service.session$.pipe(
      switchMap(session => this.service.getVotes(session, this.promptId)),
    ).subscribe(response => {
      this.votes = response.votes;
    });
  }
}
