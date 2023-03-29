import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WatchVotesComponent } from './watch-votes.component';

describe('WatchVotesComponent', () => {
  let component: WatchVotesComponent;
  let fixture: ComponentFixture<WatchVotesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ WatchVotesComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(WatchVotesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
