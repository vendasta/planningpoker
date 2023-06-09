import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {CreateSessionComponent} from './create-session/create-session.component';
import {FormsModule} from "@angular/forms";
import {HttpClientModule} from "@angular/common/http";
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatCardModule} from "@angular/material/card";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatDividerModule} from "@angular/material/divider";
import {MatListModule} from "@angular/material/list";
import {MatInputModule} from "@angular/material/input";
import {MatButtonModule} from "@angular/material/button";
import {MatLineModule} from "@angular/material/core";
import {InviteDialogComponent} from "./create-session/invite-dialog.component";
import {MatDialogModule} from "@angular/material/dialog";
import {CreatePromptComponent} from './create-prompt/create-prompt.component';
import {JoinSessionComponent} from './join-session/join-session.component';
import { WaitPromptComponent } from './wait-prompt/wait-prompt.component';
import {MatRadioModule} from "@angular/material/radio";
import { WatchVotesComponent } from './watch-votes/watch-votes.component';
import {MatTableModule} from "@angular/material/table";

@NgModule({
  declarations: [
    AppComponent,
    CreateSessionComponent,
    InviteDialogComponent,
    CreatePromptComponent,
    JoinSessionComponent,
    WaitPromptComponent,
    WatchVotesComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatCardModule,
    MatFormFieldModule,
    MatDividerModule,
    MatListModule,
    MatInputModule,
    MatButtonModule,
    MatLineModule,
    MatDialogModule,
    MatRadioModule,
    MatTableModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
