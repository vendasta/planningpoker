import {Component} from '@angular/core';
import {HostService} from "../host.service";
import {switchMap} from "rxjs";

@Component({
  selector: 'app-create-prompt',
  templateUrl: './create-prompt.component.html',
  styleUrls: ['./create-prompt.component.scss']
})
export class CreatePromptComponent {
  promptText = '';

  constructor(private service: HostService) {
  }

  createPrompt() {
    this.service.session$.pipe(
      switchMap(session => this.service.createPrompt(session, this.promptText)),
    ).subscribe();
  }
}
