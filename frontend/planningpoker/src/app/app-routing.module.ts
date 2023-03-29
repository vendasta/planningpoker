import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {CreateSessionComponent} from "./create-session/create-session.component";
import {CreatePromptComponent} from "./create-prompt/create-prompt.component";
import {JoinSessionComponent} from "./join-session/join-session.component";
import {WaitPromptComponent} from "./wait-prompt/wait-prompt.component";
import {WatchVotesComponent} from "./watch-votes/watch-votes.component";

const routes: Routes = [
  {path: 'create-session', component: CreateSessionComponent},
  {path: 'create-prompt', component: CreatePromptComponent},
  {path: 'session/:sessionId/join', component: JoinSessionComponent},
  {path: 'session/:sessionId/wait', component: WaitPromptComponent},
  {path: 'session/:sessionId/prompt/:promptId/watch', component: WatchVotesComponent},
  // add any other routes here
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
