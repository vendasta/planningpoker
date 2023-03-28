import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CreateSessionComponent} from "./create-session/create-session.component";
import {CreatePromptComponent} from "./create-prompt/create-prompt.component";

const routes: Routes = [
  { path: 'create-session', component: CreateSessionComponent },
  { path: 'create-prompt', component: CreatePromptComponent },
  // add any other routes here
];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
