import {Component, Inject} from "@angular/core";
import {MAT_DIALOG_DATA} from "@angular/material/dialog";

@Component({
  selector: 'app-invite-dialog',
  template: `
    <h2 mat-dialog-title>Invite Others</h2>
    <mat-dialog-content>
      <p>Share the following URL with others to join the session:</p>
      <p>http://localhost:4200/session/{{ data.sessionId }}/join</p>
    </mat-dialog-content>
    <mat-dialog-actions>
      <button mat-button mat-dialog-close>Close</button>
    </mat-dialog-actions>
  `,
})
export class InviteDialogComponent {
  constructor(@Inject(MAT_DIALOG_DATA) public data: { sessionId: string }) {
  }
}
