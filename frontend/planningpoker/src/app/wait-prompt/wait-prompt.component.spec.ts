import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WaitPromptComponent } from './wait-prompt.component';

describe('WaitPromptComponent', () => {
  let component: WaitPromptComponent;
  let fixture: ComponentFixture<WaitPromptComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ WaitPromptComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(WaitPromptComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
