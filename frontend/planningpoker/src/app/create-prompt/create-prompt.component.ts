import { Component } from '@angular/core';
import {PlanningPokerService} from "../planning-poker.service";
import {switchMap} from "rxjs";

@Component({
  selector: 'app-create-prompt',
  templateUrl: './create-prompt.component.html',
  styleUrls: ['./create-prompt.component.scss']
})
export class CreatePromptComponent {
  promptText = '';

  constructor(private service: PlanningPokerService) {
  }

  createPrompt() {
    this.service.session$.pipe(
      switchMap(session => this.service.createPrompt(session, this.promptText)),
    ).subscribe();
  }
}
