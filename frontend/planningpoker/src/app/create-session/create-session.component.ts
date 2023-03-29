import {Component} from '@angular/core';
import {HostService} from '../host.service';
import {MatDialog} from "@angular/material/dialog";
import {InviteDialogComponent} from "./invite-dialog.component";
import {Router} from "@angular/router";

@Component({
  selector: 'app-create-session',
  templateUrl: './create-session.component.html',
  styleUrls: ['./create-session.component.scss']
})
export class CreateSessionComponent {
  participantId = '';

  constructor(
    private sessionService: HostService,
    private router: Router,
    private dialog: MatDialog,
  ) {
  }

  createSession() {
    this.sessionService.createSession(this.participantId)
      .subscribe(response => {
        const dialogRef = this.dialog.open(InviteDialogComponent, {
          data: {sessionId: response.session_id}
        });

        dialogRef.afterClosed().subscribe(result => {
          this.router.navigateByUrl('/create-prompt');
        });
      });
  }
}

