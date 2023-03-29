import {Component} from '@angular/core';
import {ActivatedRoute, Router} from "@angular/router";
import {GuestService} from "../guest.service";

@Component({
  selector: 'app-join-session',
  templateUrl: './join-session.component.html',
  styleUrls: ['./join-session.component.scss']
})
export class JoinSessionComponent {
  participantId = '';
  sessionId = '';

  constructor(private route: ActivatedRoute, private service: GuestService, private router: Router) {
    const sessionId = this.route.snapshot.paramMap.get('sessionId');
    if (sessionId) {
      this.sessionId = sessionId;
    }
  }

  joinSession() {
    this.service.joinSession(this.sessionId, this.participantId).subscribe(response => {
      this.router.navigateByUrl(`session/${response.session_id}/wait`)
    });
  }
}
